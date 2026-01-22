import { createAnthropic } from '@ai-sdk/anthropic';
import { createOpenAI } from '@ai-sdk/openai';
import { opencode } from 'ai-sdk-provider-opencode-sdk';
import { generateText, Output } from 'ai';
import { z } from 'zod';
import fs from 'fs';
import path from 'path';
import { parseArgs } from 'util';

const { values: args } = parseArgs({
    options: {
        'raw-dir': { type: 'string', default: 'assets/api_data_raw_backup' },
        'metadata': { type: 'string', default: 'tools/clean_data/metadata.json' },
        'out': { type: 'string', default: 'tools/clean_data/curated_moves.go' },
        'cache': { type: 'string', default: 'tools/curate_moves/curated_moves_cache.json' },
        'start': { type: 'string', default: '1' },
        'end': { type: 'string', default: '1025' },
        'n': { type: 'string', default: '4' },
        'resume': { type: 'boolean', default: true },
        'force': { type: 'boolean', default: false },
        'sleep-ms': { type: 'string', default: '250' },
        'model': { type: 'string', default: 'opencode/glm-4.7-free' },
        'openai-model': { type: 'string', default: 'gpt-5-nano' },
        'base-url': { type: 'string', default: 'https://opencode.ai/zen/v1' },
        'max-candidates': { type: 'string', default: '60' },
        'seed': { type: 'string', default: '0' },
    }
});

const startID = parseInt(args.start!);
const endID = parseInt(args.end!);
const perPokemon = parseInt(args.n!);
const sleepMS = parseInt(args['sleep-ms']!);
const maxCandidates = parseInt(args['max-candidates']!);
const seed = BigInt(args.seed!);

const OPENCODE_API_KEY = process.env.OPENCODE_API_KEY || process.env.ZEN_API_KEY || process.env.OPENAI_API_KEY;
if (!OPENCODE_API_KEY) {
    console.error('Missing API key: set OPENCODE_API_KEY');
    process.exit(1);
}

const openai = createOpenAI({
    name: 'openai',
    baseURL: args['base-url'],
    headers: { Authorization: `Bearer ${OPENCODE_API_KEY}` },
});

const agentModel = opencode(args.model, {
    agent: 'general',
    sessionTitle: 'Curate moves',
    systemPrompt: 'You are a move curator for a Pokemon game. You are given a list of moves and a Pokemon, and you must choose the best moves for that Pokemon.',
});

const ROOT = path.join(import.meta.dir, '../..');
const RAW_DIR = path.isAbsolute(args['raw-dir']!) ? args['raw-dir']! : path.join(ROOT, args['raw-dir']!);
const METADATA_PATH = path.isAbsolute(args.metadata!) ? args.metadata! : path.join(ROOT, args.metadata!);
const OUT_PATH = path.isAbsolute(args.out!) ? args.out! : path.join(ROOT, args.out!);
const CACHE_PATH = path.isAbsolute(args.cache!) ? args.cache! : path.join(ROOT, args.cache!);

async function main() {
    const metadata = JSON.parse(fs.readFileSync(METADATA_PATH, 'utf-8'));
    let cache: Record<number, string[]> = {};
    if (args.resume && !args.force && fs.existsSync(CACHE_PATH)) {
        cache = JSON.parse(fs.readFileSync(CACHE_PATH, 'utf-8'));
    }

    const pokemonNames: Record<number, string> = {};

    for (let id = startID; id <= endID; id++) {
        const pokemonPath = path.join(RAW_DIR, `pokemon_${id}.json`);
        if (!fs.existsSync(pokemonPath)) {
            console.error(`Skipping ${id}: file not found`);
            continue;
        }

        const pokemon = JSON.parse(fs.readFileSync(pokemonPath, 'utf-8'));
        pokemonNames[id] = pokemon.name;

        if (!args.force && cache[id]) {
            continue;
        }

        console.log(`Processing ${id} (${pokemon.name})...`);

        const cand = getCandidates(pokemon, metadata);
        if (cand.length === 0) {
            console.warn(`No candidate moves for ${id} (${pokemon.name}); leaving empty`);
            cache[id] = [];
            saveCache(cache);
            continue;
        }

        let sorted = rankCandidates(pokemon, cand, metadata);
        if (seed !== 0n) {
            sorted = seededShuffle(sorted, seed + BigInt(id));
        }
        if (sorted.length > maxCandidates) {
            sorted = sorted.slice(0, maxCandidates);
        }

        try {
            let moves = await pickSignatureMoves(pokemon, sorted, metadata);
            moves = filterValid(pokemon, moves, metadata);
            if (moves.length > perPokemon) moves = moves.slice(0, perPokemon);
            if (moves.length === 0) moves = fallbackTopN(sorted, metadata, perPokemon);

            cache[id] = moves;
            saveCache(cache);
            console.log(`Success for ${id}: ${moves.join(', ')}`);
        } catch (err: any) {
            console.error(`LLM failed for ${id} (${pokemon.name}):`, err.message || err);
            const moves = fallbackTopN(sorted, metadata, perPokemon);
            cache[id] = moves;
            saveCache(cache);
        }

        if (sleepMS > 0) {
            await new Promise(r => setTimeout(r, sleepMS));
        }
    }

    writeCuratedMovesGo(OUT_PATH, cache, pokemonNames);
}

function getCandidates(pokemon: any, metadata: any) {
    const seen = new Set();
    const cand = [];
    for (const m of pokemon.moves) {
        const name = m.move.name;
        if (seen.has(name)) continue;
        seen.add(name);
        if (metadata[name]) {
            cand.push(name);
        }
    }
    return cand;
}

function rankCandidates(pokemon: any, cand: string[], metadata: any) {
    const types = new Set(pokemon.types.map((t: any) => t.type.name));
    const scored = cand.map(name => {
        const mm = metadata[name];
        let score = 0;
        if (types.has(mm.type)) score += 20;
        if (mm.power > 0) score += Math.min(mm.power, 120) / 5;
        if (mm.damage_class === 'special' || mm.damage_class === 'physical') score += 5;
        if (mm.damage_class === 'status') score += 1;
        return { name, score };
    });

    return scored.sort((a, b) => b.score - a.score || a.name.localeCompare(b.name)).map(s => s.name);
}

function fallbackTopN(ranked: string[], metadata: any, n: number) {
    return ranked.slice(0, n);
}

function filterValid(pokemon: any, moves: string[], metadata: any) {
    const allowed = new Set(pokemon.moves.map((m: any) => m.move.name));
    const seen = new Set();
    const out = [];
    for (const name of moves) {
        const trimmed = name.trim();
        if (!trimmed || seen.has(trimmed)) continue;
        seen.add(trimmed);
        if (!allowed.has(trimmed)) continue;
        if (!metadata[trimmed]) continue;
        out.push(trimmed);
    }
    return out;
}

async function pickSignatureMoves(pokemon: any, candidates: string[], metadata: any) {
    const stats = Object.fromEntries(pokemon.stats.map((s: any) => [s.stat.name, s.base_stat]));
    const types = pokemon.types.map((t: any) => t.type.name).join(', ');

    const candLines = candidates.map(name => {
        const mm = metadata[name];
        return `- ${name} (${mm.type}, power=${mm.power}, class=${mm.damage_class})`;
    }).join('\n');

    const system = "You are a Pokemon move curator for a Pokedex TUI. You must follow constraints exactly.";
    const prompt = `Pokemon: ${pokemon.name} (id=${pokemon.id})
Types: ${types}
Stats: hp=${stats.hp} atk=${stats.attack} def=${stats.defense} sp_atk=${stats['special-attack']} sp_def=${stats['special-defense']} speed=${stats.speed}

Choose signature moves.
Return a JSON object with a 'moves' key containing an array of 4 move names (strings) in kebab-case.
Constraints:
- Every chosen move MUST be from the candidate list below.
- Prefer the most iconic/recognizable moves associated with this Pokemon.
- Prefer STAB and high-impact moves, but include at most 1 purely-status move unless it is iconic.
- Do not include explanation text.

Candidate moves:
${candLines}`;

    const { text: rawOutput } = await generateText({
        model: agentModel,
        system,
        prompt,
    });

    console.log(`--- REASONING FOR ${pokemon.name} ---`);
    console.log(rawOutput);
    console.log('-----------------------------------');

    const { output } = await generateText({
        model: openai.responses(args['openai-model']!),
        prompt: "Extract the JSON object from the text below:\n\n" + rawOutput,
        output: Output.object({
            schema: z.object({
                moves: z.array(z.string()).length(perPokemon).describe(`The ${perPokemon} most iconic moves for this Pokemon`),
            })
        }),
    });

    const candidateSet = new Set(candidates);
    const validMoves = output.moves.filter(m => candidateSet.has(m));

    if (validMoves.length < perPokemon) {
        for (const cand of candidates) {
            if (validMoves.length >= perPokemon) break;
            if (!validMoves.includes(cand)) {
                validMoves.push(cand);
            }
        }
    }

    return validMoves;
}

function seededShuffle(out: string[], seed: bigint): string[] {
    let x = seed;
    if (x === 0n) return out;
    const res = [...out];
    for (let i = res.length - 1; i > 0; i--) {
        x ^= x >> 12n;
        x ^= x << 25n;
        x ^= x >> 27n;
        x *= 2685821657736338717n;
        const j = Number(x % BigInt(i + 1));
        [res[i], res[j]] = [res[j]!, res[i]!];
    }
    return res;
}

function saveCache(cache: Record<number, string[]>) {
    fs.mkdirSync(path.dirname(CACHE_PATH), { recursive: true });
    fs.writeFileSync(CACHE_PATH, JSON.stringify(cache, null, 2));
}

function writeCuratedMovesGo(outPath: string, curated: Record<number, string[]>, names: Record<number, string>) {
    const ids = Object.keys(curated).map(Number).sort((a, b) => a - b);
    let content = "package main\n\n";
    content += "// Code generated by tools/curate_moves; DO NOT EDIT.\n";
    content += "// CuratedMoves maps Pokemon ID to a list of iconic \"signature\" moves.\n";
    content += "var CuratedMoves = map[int][]string{\n";

    for (const id of ids) {
        const moves = curated[id]!.map(m => `"${m}"`).join(', ');
        const name = names[id] ? ` // ${names[id]}` : "";
        content += `\t${id}: {${moves}},${name}\n`;
    }

    content += "}\n";
    fs.writeFileSync(outPath, content);
}

main().catch(console.error);
