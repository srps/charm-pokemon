package models

import (
	"strconv"
	"strings"
)

type Pokemon struct {
	ID             int
	NamePT         string
	NameEN         string
	Generation     int
	Types          []string
	Height         float64
	Weight         float64
	BaseExperience int
	Stats          PokemonStats
	SignatureMoves []Move
	ArtStandard    string
	ArtShiny       string
	Evolution      *EvolutionChain
	IsFavorite     bool
}

type PokemonStats struct {
	HP      int
	Attack  int
	Defense int
	SpAtk   int
	SpDef   int
	Speed   int
}

type Move struct {
	NamePT   string
	NameEN   string
	Type     string
	Power    int
	Category string
}

type PokemonFilter struct {
	Query      string
	Type       string
	Generation int
}

type Pokedex struct {
	Pokemon       []*Pokemon
	PokemonByID   map[int]*Pokemon
	PokemonByName map[string]*Pokemon
	ByGeneration  map[int][]*Pokemon
	ByType        map[string][]*Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		Pokemon:       make([]*Pokemon, 0),
		PokemonByID:   make(map[int]*Pokemon),
		PokemonByName: make(map[string]*Pokemon),
		ByGeneration:  make(map[int][]*Pokemon),
		ByType:        make(map[string][]*Pokemon),
	}
}

func (p *Pokedex) AddPokemon(pokemon *Pokemon) {
	p.Pokemon = append(p.Pokemon, pokemon)
	p.PokemonByID[pokemon.ID] = pokemon
	p.PokemonByName[pokemon.NameEN] = pokemon
	p.PokemonByName[pokemon.NamePT] = pokemon

	p.ByGeneration[pokemon.Generation] = append(p.ByGeneration[pokemon.Generation], pokemon)
	for _, t := range pokemon.Types {
		p.ByType[t] = append(p.ByType[t], pokemon)
	}
}

func (p *Pokedex) GetByID(id int) *Pokemon {
	return p.PokemonByID[id]
}

func (p *Pokedex) GetByName(name string) *Pokemon {
	return p.PokemonByName[name]
}

func (p *Pokedex) Search(filter PokemonFilter) []*Pokemon {
	results := make([]*Pokemon, 0)

	for _, pokemon := range p.Pokemon {
		if filter.Generation > 0 && pokemon.Generation != filter.Generation {
			continue
		}

		if filter.Type != "" {
			hasType := false
			for _, t := range pokemon.Types {
				if t == filter.Type {
					hasType = true
					break
				}
			}
			if !hasType {
				continue
			}
		}

		if filter.Query == "" {
			results = append(results, pokemon)
			continue
		}

		query := filter.Query
		queryID, parseErr := strconv.Atoi(query)
		if parseErr == nil && pokemon.ID == queryID {
			results = append(results, pokemon)
		} else if containsIgnoreCase(pokemon.NamePT, query) || containsIgnoreCase(pokemon.NameEN, query) {
			results = append(results, pokemon)
		}
	}

	return results
}

func (p *Pokedex) GetNextPokemon(id int) *Pokemon {
	for i, pokemon := range p.Pokemon {
		if pokemon.ID == id && i < len(p.Pokemon)-1 {
			return p.Pokemon[i+1]
		}
	}
	return p.Pokemon[0]
}

func (p *Pokedex) GetPrevPokemon(id int) *Pokemon {
	for i, pokemon := range p.Pokemon {
		if pokemon.ID == id && i > 0 {
			return p.Pokemon[i-1]
		}
	}
	return p.Pokemon[len(p.Pokemon)-1]
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

func (p *Pokedex) GetPokemonByGeneration(gen int) []*Pokemon {
	return p.ByGeneration[gen]
}

func (p *Pokedex) GetPokemonByType(typeName string) []*Pokemon {
	return p.ByType[typeName]
}

func (p *Pokedex) GetCount() int {
	return len(p.Pokemon)
}
