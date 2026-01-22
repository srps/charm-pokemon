package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FavoritesManager struct {
	Favorites    map[int]bool
	FavoritesDir string
	FilePath     string
}

func NewFavoritesManager() *FavoritesManager {
	execDir, _ := os.Executable()
	favoritesDir := filepath.Join(filepath.Dir(execDir), "assets")

	if _, err := os.Stat(favoritesDir); os.IsNotExist(err) {
		os.MkdirAll(favoritesDir, 0755)
	}

	filePath := filepath.Join(favoritesDir, "favorites.json")

	fm := &FavoritesManager{
		Favorites:    make(map[int]bool),
		FavoritesDir: favoritesDir,
		FilePath:     filePath,
	}

	fm.load()

	return fm
}

func (fm *FavoritesManager) load() error {
	data, err := os.ReadFile(fm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &fm.Favorites)
}

func (fm *FavoritesManager) save() error {
	data, err := json.MarshalIndent(fm.Favorites, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fm.FilePath, data, 0644)
}

func (fm *FavoritesManager) IsFavorite(id int) bool {
	return fm.Favorites[id]
}

func (fm *FavoritesManager) AddFavorite(id int) error {
	fm.Favorites[id] = true
	return fm.save()
}

func (fm *FavoritesManager) RemoveFavorite(id int) error {
	delete(fm.Favorites, id)
	return fm.save()
}

func (fm *FavoritesManager) ToggleFavorite(id int) (bool, error) {
	if fm.Favorites[id] {
		delete(fm.Favorites, id)
		err := fm.save()
		return false, err
	} else {
		fm.Favorites[id] = true
		err := fm.save()
		return true, err
	}
}

func (fm *FavoritesManager) GetAllFavorites() []int {
	ids := make([]int, 0, len(fm.Favorites))
	for id := range fm.Favorites {
		ids = append(ids, id)
	}
	return ids
}

func (fm *FavoritesManager) GetCount() int {
	return len(fm.Favorites)
}
