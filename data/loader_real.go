//go:build realdata

package data

import (
	"charm-pokemon/assets"
	"charm-pokemon/models"
	"encoding/json"
	"fmt"
	"strings"
)

type pokeAPIResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type genAPIResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PokemonSpecies []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon_species"`
}

func GetPokedex() *models.Pokedex {
	pokedex := models.NewPokedex()

	// 1. Load all generations to map pokemon to generations
	pokemonToGen := make(map[int]int)
	for i := 1; i <= 9; i++ {
		genData, err := assets.EmbedFS.ReadFile(fmt.Sprintf("embed/api_data/generation_%d.json", i))
		if err != nil {
			continue
		}
		var genResponse genAPIResponse
		if err := json.Unmarshal(genData, &genResponse); err == nil {
			for _, species := range genResponse.PokemonSpecies {
				// Extract ID from URL: https://pokeapi.co/api/v2/pokemon-species/{id}/
				parts := strings.Split(strings.Trim(species.URL, "/"), "/")
				if len(parts) > 0 {
					idStr := parts[len(parts)-1]
					var id int
					fmt.Sscanf(idStr, "%d", &id)
					if id > 0 {
						pokemonToGen[id] = i
					}
				}
			}
		}
	}

	// 2. Load all pokemon data
	for i := 1; i <= 1025; i++ {
		pokemonData, err := assets.EmbedFS.ReadFile(fmt.Sprintf("embed/api_data/pokemon_%d.json", i))
		if err != nil {
			continue
		}

		var resp pokeAPIResponse
		if err := json.Unmarshal(pokemonData, &resp); err != nil {
			continue
		}

		pokemon := &models.Pokemon{
			ID:     resp.ID,
			NameEN: strings.Title(resp.Name),
			NamePT: strings.Title(resp.Name), // Fallback to English
			Height: float64(resp.Height),
			Weight: float64(resp.Weight),
			Stats:  models.PokemonStats{},
		}

		pokemon.Generation = pokemonToGen[pokemon.ID]
		if pokemon.Generation == 0 {
			// Heuristic if generator mapping failed
			if pokemon.ID <= 151 {
				pokemon.Generation = 1
			} else if pokemon.ID <= 251 {
				pokemon.Generation = 2
			} else if pokemon.ID <= 386 {
				pokemon.Generation = 3
			} else if pokemon.ID <= 493 {
				pokemon.Generation = 4
			} else if pokemon.ID <= 649 {
				pokemon.Generation = 5
			} else if pokemon.ID <= 721 {
				pokemon.Generation = 6
			} else if pokemon.ID <= 809 {
				pokemon.Generation = 7
			} else if pokemon.ID <= 905 {
				pokemon.Generation = 8
			} else {
				pokemon.Generation = 9
			}
		}

		for _, s := range resp.Stats {
			val := s.BaseStat
			switch s.Stat.Name {
			case "hp":
				pokemon.Stats.HP = val
			case "attack":
				pokemon.Stats.Attack = val
			case "defense":
				pokemon.Stats.Defense = val
			case "special-attack":
				pokemon.Stats.SpAtk = val
			case "special-defense":
				pokemon.Stats.SpDef = val
			case "speed":
				pokemon.Stats.Speed = val
			}
		}

		for _, t := range resp.Types {
			pokemon.Types = append(pokemon.Types, translateType(t.Type.Name))
		}

		pokedex.AddPokemon(pokemon)
	}

	return pokedex
}

func translateType(t string) string {
	types := map[string]string{
		"normal":   "normal",
		"fire":     "fogo",
		"water":    "água",
		"grass":    "erva",
		"electric": "elétrico",
		"ice":      "gelo",
		"fighting": "lutador",
		"poison":   "veneno",
		"ground":   "terra",
		"flying":   "voador",
		"psychic":  "psíquico",
		"bug":      "inseto",
		"rock":     "pedra",
		"ghost":    "fantasma",
		"dragon":   "dragão",
		"dark":     "sombrio",
		"steel":    "metálico",
		"fairy":    "fada",
	}
	if pt, ok := types[t]; ok {
		return pt
	}
	return t
}
