package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RawPokemonResponse struct {
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
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`
}

type MinifiedPokemon struct {
	ID             int      `json:"id"`
	NameEN         string   `json:"name_en"`
	NamePT         string   `json:"name_pt"`
	Height         int      `json:"height"`
	Weight         int      `json:"weight"`
	Stats          Stats    `json:"stats"`
	Types          []string `json:"types"`
	SignatureMoves []Move   `json:"signature_moves"`
}

type Stats struct {
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	SpAtk   int `json:"sp_atk"`
	SpDef   int `json:"sp_def"`
	Speed   int `json:"speed"`
}

type Move struct {
	NameEN   string `json:"name_en"`
	NamePT   string `json:"name_pt"`
	Type     string `json:"type"`
	Power    int    `json:"power"`
	Category string `json:"category"`
}

func formatMoveName(name string) string {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.Title(p)
		}
	}
	return strings.Join(parts, " ")
}

func main() {
	inputDir := filepath.Join("..", "..", "assets", "api_data")

	fmt.Println("Cleaning up Pokemon data (Filter & Trim)...")

	for i := 1; i <= 1025; i++ {
		fileName := fmt.Sprintf("pokemon_%d.json", i)
		inputPath := filepath.Join(inputDir, fileName)

		data, err := os.ReadFile(inputPath)
		if err != nil {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			fmt.Printf("Error unmarshaling %s: %v\n", fileName, err)
			continue
		}

		// Keys to keep
		keepKeys := map[string]bool{
			"id":     true,
			"name":   true,
			"height": true,
			"weight": true,
			"stats":  true,
			"types":  true,
			"moves":  true,
		}

		// Trim unused keys
		for k := range raw {
			if !keepKeys[k] {
				delete(raw, k)
			}
		}

		// Filter and Inject Move Metadata
		curatedNames := CuratedMoves[i]
		newMoves := make([]map[string]interface{}, 0)

		rawMoves, ok := raw["moves"].([]interface{})
		if ok && len(curatedNames) > 0 {
			// Create a map of existing moves for this pokemon to bridge to curated
			existingMoves := make(map[string]bool)
			for _, m := range rawMoves {
				mObj, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				moveInfo, ok := mObj["move"].(map[string]interface{})
				if !ok {
					continue
				}
				name, ok := moveInfo["name"].(string)
				if ok {
					existingMoves[name] = true
				}
			}

			for _, targetName := range curatedNames {
				// Validation: check if pokemon actually learns this move
				if !existingMoves[targetName] {
					fmt.Printf("Warning: Pokemon %d does not officially learn %s. Skipping.\n", i, targetName)
					continue
				}

				meta, hasMeta := MoveMetadataMap[targetName]
				if !hasMeta {
					fmt.Printf("Warning: No metadata found for move %s. Skipping.\n", targetName)
					continue
				}

				moveObj := map[string]interface{}{
					"move": map[string]interface{}{
						"name": targetName,
						"url":  fmt.Sprintf("https://pokeapi.co/api/v2/move/%s/", targetName), // Keep structure
					},
					"metadata": map[string]interface{}{
						"type":         meta.Type,
						"power":        meta.Power,
						"damage_class": meta.DamageClass,
					},
				}
				newMoves = append(newMoves, moveObj)
			}
		} else if ok {
			// Fallback: take first 4 moves if no curated list
			for j := 0; j < len(rawMoves) && j < 4; j++ {
				mObj := rawMoves[j].(map[string]interface{})
				moveInfo := mObj["move"].(map[string]interface{})
				name := moveInfo["name"].(string)

				meta, hasMeta := MoveMetadataMap[name]
				injected := mObj
				if hasMeta {
					injected["metadata"] = map[string]interface{}{
						"type":         meta.Type,
						"power":        meta.Power,
						"damage_class": meta.DamageClass,
					}
				}
				newMoves = append(newMoves, injected)
			}
		}

		raw["moves"] = newMoves

		// Minified output
		minData, _ := json.Marshal(raw)
		err = os.WriteFile(inputPath, minData, 0644)
		if err != nil {
			fmt.Printf("Error writing %s: %v\n", fileName, err)
		}

		if i%100 == 0 {
			fmt.Printf("Processed %d/1025 Pokemon\n", i)
		}
	}

	fmt.Println("Cleanup complete!")
}
