//go:build fetch_metadata

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type MoveMetadata struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Power       int    `json:"power"`
	DamageClass string `json:"damage_class"`
}

type APIResponse struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
	Type  struct {
		Name string `json:"name"`
	} `json:"type"`
	DamageClass struct {
		Name string `json:"name"`
	} `json:"damage_class"`
}

func main() {
	uniqueMoves := make(map[string]bool)
	for _, moves := range CuratedMoves {
		for _, m := range moves {
			uniqueMoves[m] = true
		}
	}

	metadata := make(map[string]MoveMetadata)

	// Load existing metadata if any to resume
	metaPath := "metadata.json"
	if data, err := os.ReadFile(metaPath); err == nil {
		json.Unmarshal(data, &metadata)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	count := 0
	total := len(uniqueMoves)

	fmt.Printf("Fetching metadata for %d unique moves...\n", total)

	for moveName := range uniqueMoves {
		if _, exists := metadata[moveName]; exists {
			continue
		}

		url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", moveName)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", moveName, err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Failed to fetch %s: status %d\n", moveName, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Printf("Error unmarshaling %s: %v\n", moveName, err)
			continue
		}

		metadata[moveName] = MoveMetadata{
			Name:        apiResp.Name,
			Type:        apiResp.Type.Name,
			Power:       apiResp.Power,
			DamageClass: apiResp.DamageClass.Name,
		}

		count++
		if count%10 == 0 {
			fmt.Printf("Fetched %d/%d moves...\n", count, total)
			// Save progress
			saveMetadata(metaPath, metadata)
		}

		// Rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	saveMetadata(metaPath, metadata)
	fmt.Println("Metadata fetch complete!")

	// Generate Go file
	generateGoFile(metadata)
}

func saveMetadata(path string, data map[string]MoveMetadata) {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(path, bytes, 0644)
}

func generateGoFile(metadata map[string]MoveMetadata) {
	f, err := os.Create("metadata.go")
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "var MoveMetadataMap = map[string]MoveMetadata{")
	for name, m := range metadata {
		fmt.Fprintf(f, "\t\"%s\": {Name: \"%s\", Type: \"%s\", Power: %d, DamageClass: \"%s\"},\n",
			name, m.Name, m.Type, m.Power, m.DamageClass)
	}
	fmt.Fprintln(f, "}")

	fmt.Fprintln(f)
	fmt.Fprintln(f, "type MoveMetadata struct {")
	fmt.Fprintln(f, "\tName        string")
	fmt.Fprintln(f, "\tType        string")
	fmt.Fprintln(f, "\tPower       int")
	fmt.Fprintln(f, "\tDamageClass string")
	fmt.Fprintln(f, "}")
}
