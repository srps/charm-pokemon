package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	baseAPIURL     = "https://pokeapi.co/api/v2"
	spriteBaseURL  = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork"
	maxPokemonID   = 1025
	outputDir      = "../../assets/api_data"
	spriteDirStd   = "../../assets/sprites/standard"
	spriteDirShiny = "../../assets/sprites/shiny"
)

func main() {
	fmt.Println("Downloading Pokemon data from PokeAPI...")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		downloadPokemonData()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		downloadSprites(spriteDirStd, "")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		downloadSprites(spriteDirShiny, "/shiny")
	}()

	wg.Wait()

	fmt.Println("Download complete!")
}

func downloadPokemonData() {
	fmt.Println("Downloading Pokemon data...")

	os.MkdirAll(outputDir, 0755)

	for i := 1; i <= maxPokemonID; i++ {
		if err := downloadSinglePokemon(i); err != nil {
			fmt.Printf("Error downloading Pokemon %d: %v\n", i, err)
		}

		if i%50 == 0 {
			fmt.Printf("Downloaded %d/%d Pokemon data\n", i, maxPokemonID)
		}
	}

	for i := 1; i <= 9; i++ {
		if err := downloadGeneration(i); err != nil {
			fmt.Printf("Error downloading generation %d: %v\n", i, err)
		}
	}

	fmt.Println("Pokemon data download complete!")
}

func downloadSinglePokemon(id int) error {
	url := fmt.Sprintf("%s/pokemon/%d", baseAPIURL, id)
	filePath := filepath.Join(outputDir, fmt.Sprintf("pokemon_%d.json", id))

	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		return err
	}

	formatted, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, formatted, 0644)
}

func downloadGeneration(id int) error {
	url := fmt.Sprintf("%s/generation/%d", baseAPIURL, id)
	filePath := filepath.Join(outputDir, fmt.Sprintf("generation_%d.json", id))

	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		return err
	}

	formatted, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, formatted, 0644)
}

func downloadSprites(dir string, shinySuffix string) {
	fmt.Printf("Downloading %s sprites...\n", dir)

	os.MkdirAll(dir, 0755)

	for i := 1; i <= maxPokemonID; i++ {
		url := fmt.Sprintf("%s%s/%d.png", spriteBaseURL, shinySuffix, i)
		filePath := filepath.Join(dir, fmt.Sprintf("%d.png", i))

		if _, err := os.Stat(filePath); err == nil {
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error downloading sprite %d: %v\n", i, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Printf("HTTP %d for sprite %d\n", resp.StatusCode, i)
			continue
		}

		file, err := os.Create(filePath)
		if err != nil {
			resp.Body.Close()
			fmt.Printf("Error creating file %d: %v\n", i, err)
			continue
		}

		_, err = io.Copy(file, resp.Body)
		file.Close()
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Error saving sprite %d: %v\n", i, err)
			os.Remove(filePath)
			continue
		}

		if i%100 == 0 {
			fmt.Printf("Downloaded %d/%d sprites\n", i, maxPokemonID)
		}
	}

	fmt.Printf("%s sprites download complete!\n", dir)
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
