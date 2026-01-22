package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MinimalStat represents a minimal stat entry
type MinimalStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

// MinimalType represents a minimal type entry
type MinimalType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

// MinimalPokemon contains only the fields we actually use
type MinimalPokemon struct {
	ID     int           `json:"id"`
	Name   string        `json:"name"`
	Height int           `json:"height"`
	Weight int           `json:"weight"`
	Stats  []MinimalStat `json:"stats"`
	Types  []MinimalType `json:"types"`
}

// MinimalGeneration for generation files
type MinimalGeneration struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PokemonSpecies []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon_species"`
}

func main() {
	inputDir := "assets/api_data"
	outputDir := "assets/embed/api_data"

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Process Pokemon files
	pokemonCount := 0
	totalOrigSize := int64(0)
	totalMinSize := int64(0)

	for i := 1; i <= 1025; i++ {
		inputPath := filepath.Join(inputDir, fmt.Sprintf("pokemon_%d.json", i))
		outputPath := filepath.Join(outputDir, fmt.Sprintf("pokemon_%d.json", i))

		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Printf("Skipping pokemon_%d.json: %v\n", i, err)
			continue
		}
		totalOrigSize += int64(len(data))

		// Parse full JSON
		var fullData map[string]interface{}
		if err := json.Unmarshal(data, &fullData); err != nil {
			fmt.Printf("Error parsing pokemon_%d.json: %v\n", i, err)
			continue
		}

		// Extract only needed fields
		minimal := MinimalPokemon{
			ID:     int(fullData["id"].(float64)),
			Name:   fullData["name"].(string),
			Height: int(fullData["height"].(float64)),
			Weight: int(fullData["weight"].(float64)),
		}

		// Process stats - keep only base_stat and stat.name
		if stats, ok := fullData["stats"].([]interface{}); ok {
			for _, s := range stats {
				stat := s.(map[string]interface{})
				ms := MinimalStat{
					BaseStat: int(stat["base_stat"].(float64)),
				}
				if statInfo, ok := stat["stat"].(map[string]interface{}); ok {
					ms.Stat.Name = statInfo["name"].(string)
				}
				minimal.Stats = append(minimal.Stats, ms)
			}
		}

		// Process types
		if types, ok := fullData["types"].([]interface{}); ok {
			for _, t := range types {
				typ := t.(map[string]interface{})
				mt := MinimalType{
					Slot: int(typ["slot"].(float64)),
				}
				if typeInfo, ok := typ["type"].(map[string]interface{}); ok {
					mt.Type.Name = typeInfo["name"].(string)
				}
				minimal.Types = append(minimal.Types, mt)
			}
		}

		// Write minified JSON (no indentation for smaller size)
		minData, err := json.Marshal(minimal)
		if err != nil {
			fmt.Printf("Error marshaling pokemon_%d.json: %v\n", i, err)
			continue
		}

		if err := os.WriteFile(outputPath, minData, 0644); err != nil {
			fmt.Printf("Error writing pokemon_%d.json: %v\n", i, err)
			continue
		}

		totalMinSize += int64(len(minData))
		pokemonCount++
	}

	// Copy generation files (they're already small and needed)
	for i := 1; i <= 9; i++ {
		inputPath := filepath.Join(inputDir, fmt.Sprintf("generation_%d.json", i))
		outputPath := filepath.Join(outputDir, fmt.Sprintf("generation_%d.json", i))

		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Printf("Skipping generation_%d.json: %v\n", i, err)
			continue
		}
		totalOrigSize += int64(len(data))

		// Parse and re-marshal to remove any extra whitespace
		var genData MinimalGeneration
		if err := json.Unmarshal(data, &genData); err != nil {
			fmt.Printf("Error parsing generation_%d.json: %v\n", i, err)
			continue
		}

		minData, err := json.Marshal(genData)
		if err != nil {
			fmt.Printf("Error marshaling generation_%d.json: %v\n", i, err)
			continue
		}

		if err := os.WriteFile(outputPath, minData, 0644); err != nil {
			fmt.Printf("Error writing generation_%d.json: %v\n", i, err)
			continue
		}
		totalMinSize += int64(len(minData))
	}

	fmt.Printf("\n✅ Minification complete!\n")
	fmt.Printf("   Pokemon processed: %d\n", pokemonCount)
	fmt.Printf("   Original size: %.2f MB\n", float64(totalOrigSize)/1024/1024)
	fmt.Printf("   Minified size: %.2f MB\n", float64(totalMinSize)/1024/1024)
	fmt.Printf("   Reduction: %.1f%%\n", (1-float64(totalMinSize)/float64(totalOrigSize))*100)
}
