package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RawPokemon struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`
}

type MoveMeta struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Power       int    `json:"power"`
	DamageClass string `json:"damage_class"`
}

type chatCompletionsRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionsResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
	Object string `json:"object,omitempty"`
}

var jsonArrayRe = regexp.MustCompile(`(?s)\[[^\]]*\]`)

func main() {
	var (
		rawDir       = flag.String("raw-dir", filepath.FromSlash("assets/api_data_raw_backup"), "Directory containing raw pokemon_*.json files")
		metadataPath = flag.String("metadata", filepath.FromSlash("tools/clean_data/metadata.json"), "Path to move metadata JSON")
		outGo        = flag.String("out", filepath.FromSlash("tools/clean_data/curated_moves.go"), "Output Go file to write")
		cachePath    = flag.String("cache", filepath.FromSlash("tools/curate_moves/curated_moves_cache.json"), "Cache file (id -> curated moves)")
		startID      = flag.Int("start", 1, "Start pokemon ID")
		endID        = flag.Int("end", 1025, "End pokemon ID")
		perPokemon   = flag.Int("n", 4, "How many signature moves to pick (3-6 recommended)")
		resume       = flag.Bool("resume", true, "Resume from cache (if present)")
		force        = flag.Bool("force", false, "Ignore cache and re-generate all requested IDs")
		sleepMS      = flag.Int("sleep-ms", 250, "Sleep between requests (best-effort rate limiting)")
		model        = flag.String("model", "glm-4.7-free", "Model id")
		baseURL      = flag.String("base-url", "https://opencode.ai/zen/v1/chat/completions", "OpenCode Zen chat completions URL")
		maxCand      = flag.Int("max-candidates", 60, "Max candidate moves to include in prompt (sorted by relevance)")
		seed         = flag.Int64("seed", 0, "Optional seed for deterministic ordering (0 = disabled)")
	)
	flag.Parse()

	if *perPokemon < 1 {
		exitf("-n must be >= 1")
	}
	if *startID < 1 || *endID < *startID {
		exitf("invalid range: start=%d end=%d", *startID, *endID)
	}

	apiKey := firstNonEmpty(
		os.Getenv("OPENCODE_API_KEY"),
		os.Getenv("ZEN_API_KEY"),
		os.Getenv("OPENAI_API_KEY"),
	)
	if apiKey == "" {
		exitf("missing API key: set OPENCODE_API_KEY (or ZEN_API_KEY / OPENAI_API_KEY)")
	}

	meta, err := loadMetadata(*metadataPath)
	if err != nil {
		exitf("load metadata: %v", err)
	}

	cache := map[int][]string{}
	if *resume && !*force {
		if c, err := loadCache(*cachePath); err == nil {
			cache = c
		}
	}

	client := &http.Client{Timeout: 60 * time.Second}
	ctx := context.Background()

	pokemonNames := make(map[int]string, (*endID-*startID)+1)

	for id := *startID; id <= *endID; id++ {
		if !*force {
			if _, ok := cache[id]; ok {
				if _, have := pokemonNames[id]; !have {
					if p, err := loadPokemon(*rawDir, id); err == nil {
						pokemonNames[id] = p.Name
					}
				}
				continue
			}
		}

		p, err := loadPokemon(*rawDir, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %d: %v\n", id, err)
			continue
		}
		pokemonNames[id] = p.Name

		cand := candidateMoves(p, meta)
		if len(cand) == 0 {
			fmt.Fprintf(os.Stderr, "No candidate moves for %d (%s); leaving empty\n", id, p.Name)
			cache[id] = []string{}
			saveCacheBestEffort(*cachePath, cache)
			continue
		}

		if *seed != 0 {
			cand = seededShuffle(cand, *seed+int64(id))
		}

		sorted := rankCandidates(p, cand, meta)
		if len(sorted) > *maxCand {
			sorted = sorted[:*maxCand]
		}

		moves, err := pickSignatureMovesLLM(ctx, client, *baseURL, apiKey, *model, p, sorted, meta, *perPokemon)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LLM failed for %d (%s): %v. Using algorithmic fallback.\n", id, p.Name, err)
			moves = fallbackTopN(sorted, meta, *perPokemon)
		}

		moves = filterValid(p, moves, meta)
		if len(moves) > *perPokemon {
			moves = moves[:*perPokemon]
		}
		if len(moves) == 0 {
			moves = fallbackTopN(sorted, meta, *perPokemon)
		}

		cache[id] = moves
		saveCacheBestEffort(*cachePath, cache)

		if *sleepMS > 0 {
			time.Sleep(time.Duration(*sleepMS) * time.Millisecond)
		}
	}

	if err := writeCuratedMovesGo(*outGo, cache, pokemonNames); err != nil {
		exitf("write %s: %v", *outGo, err)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func loadPokemon(rawDir string, id int) (RawPokemon, error) {
	path := filepath.Join(rawDir, fmt.Sprintf("pokemon_%d.json", id))
	b, err := os.ReadFile(path)
	if err != nil {
		return RawPokemon{}, err
	}
	var p RawPokemon
	if err := json.Unmarshal(b, &p); err != nil {
		return RawPokemon{}, err
	}
	if p.ID == 0 {
		p.ID = id
	}
	if p.Name == "" {
		return RawPokemon{}, errors.New("missing name")
	}
	return p, nil
}

func loadMetadata(path string) (map[string]MoveMeta, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := map[string]MoveMeta{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func loadCache(path string) (map[int][]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[int][]string{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func saveCacheBestEffort(path string, cache map[int][]string) {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	b, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, b, 0644)
}

func pokemonTypeSet(p RawPokemon) map[string]bool {
	set := map[string]bool{}
	for _, t := range p.Types {
		name := strings.TrimSpace(t.Type.Name)
		if name != "" {
			set[name] = true
		}
	}
	return set
}

func candidateMoves(p RawPokemon, meta map[string]MoveMeta) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(p.Moves))
	for _, m := range p.Moves {
		name := strings.TrimSpace(m.Move.Name)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		if _, ok := meta[name]; ok {
			out = append(out, name)
		}
	}
	return out
}

func rankCandidates(p RawPokemon, cand []string, meta map[string]MoveMeta) []string {
	types := pokemonTypeSet(p)

	scored := make([]struct {
		Name  string
		Score int
	}, 0, len(cand))

	for _, name := range cand {
		mm, ok := meta[name]
		if !ok {
			continue
		}
		s := 0
		if types[mm.Type] {
			s += 20
		}
		if mm.Power > 0 {
			s += min(mm.Power, 120) / 5
		}
		switch mm.DamageClass {
		case "special", "physical":
			s += 5
		case "status":
			s += 1
		}
		scored = append(scored, struct {
			Name  string
			Score int
		}{Name: name, Score: s})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Name < scored[j].Name
		}
		return scored[i].Score > scored[j].Score
	})

	out := make([]string, 0, len(scored))
	for _, s := range scored {
		out = append(out, s.Name)
	}
	return out
}

func fallbackTopN(ranked []string, meta map[string]MoveMeta, n int) []string {
	out := make([]string, 0, n)
	for _, name := range ranked {
		if len(out) >= n {
			break
		}
		if _, ok := meta[name]; !ok {
			continue
		}
		out = append(out, name)
	}
	return out
}

func filterValid(p RawPokemon, moves []string, meta map[string]MoveMeta) []string {
	allowed := map[string]bool{}
	for _, m := range p.Moves {
		name := strings.TrimSpace(m.Move.Name)
		if name != "" {
			allowed[name] = true
		}
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(moves))
	for _, name := range moves {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		if !allowed[name] {
			continue
		}
		if _, ok := meta[name]; !ok {
			continue
		}
		out = append(out, name)
	}
	return out
}

func pickSignatureMovesLLM(ctx context.Context, client *http.Client, baseURL, apiKey, model string, p RawPokemon, candidates []string, meta map[string]MoveMeta, n int) ([]string, error) {
	candLines := make([]string, 0, len(candidates))
	for _, name := range candidates {
		mm := meta[name]
		candLines = append(candLines, fmt.Sprintf("- %s (%s, power=%d, class=%s)", name, mm.Type, mm.Power, mm.DamageClass))
	}

	types := make([]string, 0, len(p.Types))
	for _, t := range p.Types {
		if t.Type.Name != "" {
			types = append(types, t.Type.Name)
		}
	}

	stats := map[string]int{}
	for _, s := range p.Stats {
		if s.Stat.Name != "" {
			stats[s.Stat.Name] = s.BaseStat
		}
	}

	system := "You are a Pokemon move curator for a Pokedex TUI. You must follow constraints exactly."
	user := strings.Join([]string{
		fmt.Sprintf("Pokemon: %s (id=%d)", p.Name, p.ID),
		fmt.Sprintf("Types: %s", strings.Join(types, ", ")),
		fmt.Sprintf("Stats: hp=%d atk=%d def=%d sp_atk=%d sp_def=%d speed=%d",
			stats["hp"], stats["attack"], stats["defense"], stats["special-attack"], stats["special-defense"], stats["speed"],
		),
		"",
		"Choose signature moves.",
		fmt.Sprintf("Return ONLY a JSON array of %d move names (strings) in kebab-case.", n),
		"Constraints:",
		"- Every chosen move MUST be from the candidate list below.",
		"- Prefer the most iconic/recognizable moves associated with this Pokemon.",
		"- Prefer STAB and high-impact moves, but include at most 1 purely-status move unless it is iconic.",
		"- Do not include explanation text.",
		"Candidate moves:",
		strings.Join(candLines, "\n"),
	}, "\n")

	reqBody := chatCompletionsRequest{
		Model: model,
		Messages: []message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Temperature: 0.2, // Lower temperature for more stable JSON
		MaxTokens:   400,
	}

	// Max retries with exponential backoff
	const maxRetries = 3
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			sleepDur := time.Duration(1<<uint(i)) * time.Second
			fmt.Fprintf(os.Stderr, " (retry %d after %v)...", i, sleepDur)
			select {
			case <-time.After(sleepDur):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		moves, raw, err := tryCallLLM(ctx, client, baseURL, apiKey, reqBody)
		if err == nil {
			if i > 0 {
				fmt.Fprintf(os.Stderr, " success!\n")
			}
			return moves, nil
		}

		// If it's a 500 error from the provider mentioning prompt_tokens or similar infra error, retry.
		// If it's a 4xx error (except 429), it's probably a client error, don't retry.
		var (
			is500   = strings.Contains(err.Error(), "http 5")
			isRate  = strings.Contains(err.Error(), "429")
			isInfra = strings.Contains(raw, "prompt_tokens") || strings.Contains(raw, "undefined") || strings.Contains(raw, "error")
		)

		if !is500 && !isRate && !isInfra && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, " fatal error: %v\n", err)
			return nil, fmt.Errorf("%v (raw=%q)", err, truncate(raw, 240))
		}

		if i == maxRetries {
			fmt.Fprintf(os.Stderr, " all retries failed.\n")
		}
	}

	// One last attempt with combining system/user messages if it failed
	fmt.Fprintf(os.Stderr, " (trying combined prompt fallback)...")
	reqBody.Messages = []message{
		{Role: "user", Content: system + "\n\n" + user},
	}
	moves, raw, err := tryCallLLM(ctx, client, baseURL, apiKey, reqBody)
	if err == nil {
		return moves, nil
	}

	return nil, fmt.Errorf("last failure: %v (raw=%q)", err, truncate(raw, 240))
}

func tryCallLLM(ctx context.Context, client *http.Client, baseURL, apiKey string, reqBody chatCompletionsRequest) ([]string, string, error) {
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(b))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, string(body), fmt.Errorf("http %d", resp.StatusCode)
	}

	var out chatCompletionsResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, string(body), err
	}
	if out.Error != nil {
		return nil, string(body), errors.New(out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return nil, string(body), errors.New("no choices")
	}
	content := strings.TrimSpace(out.Choices[0].Message.Content)
	moves, err := parseMovesArray(content)
	return moves, content, err
}

func parseMovesArray(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	match := jsonArrayRe.FindString(s)
	if match == "" {
		return nil, errors.New("no json array found")
	}
	var arr []string
	if err := json.Unmarshal([]byte(match), &arr); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func seededShuffle(in []string, seed int64) []string {
	out := append([]string(nil), in...)
	n := len(out)
	x := uint64(seed)
	if x == 0 {
		return out
	}
	for i := n - 1; i > 0; i-- {
		x ^= x >> 12
		x ^= x << 25
		x ^= x >> 27
		x *= 2685821657736338717
		j := int(x % uint64(i+1))
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func writeCuratedMovesGo(path string, curated map[int][]string, names map[int]string) error {
	ids := make([]int, 0, len(curated))
	for id := range curated {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	var buf bytes.Buffer
	buf.WriteString("package main\n\n")
	buf.WriteString("// Code generated by tools/curate_moves; DO NOT EDIT.\n")
	buf.WriteString("// CuratedMoves maps Pokemon ID to a list of iconic \"signature\" moves.\n")
	buf.WriteString("var CuratedMoves = map[int][]string{\n")
	for _, id := range ids {
		moves := append([]string(nil), curated[id]...)
		for i := range moves {
			moves[i] = strconv.Quote(moves[i])
		}
		name := names[id]
		if name != "" {
			name = " // " + name
		}
		buf.WriteString(fmt.Sprintf("\t%d: {%s},%s\n", id, strings.Join(moves, ", "), name))
	}
	buf.WriteString("}\n")

	return os.WriteFile(path, buf.Bytes(), 0644)
}
