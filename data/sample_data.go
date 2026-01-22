//go:build !realdata

package data

import (
	"charm-pokemon/models"
)

var SamplePokemon = []*models.Pokemon{
	{
		ID:             1,
		NamePT:         "Bulbasaur",
		NameEN:         "Bulbasaur",
		Generation:     1,
		Types:          []string{"grama", "veneno"},
		Height:         7.0,
		Weight:         69.0,
		BaseExperience: 64,
		Stats: models.PokemonStats{
			HP:      45,
			Attack:  49,
			Defense: 49,
			SpAtk:   65,
			SpDef:   65,
			Speed:   45,
		},
		SignatureMoves: []models.Move{
			{
				NamePT:   "Razor Leaf",
				NameEN:   "Razor Leaf",
				Type:     "grama",
				Power:    55,
				Category: "physical",
			},
			{
				NamePT:   "Vine Whip",
				NameEN:   "Vine Whip",
				Type:     "grama",
				Power:    45,
				Category: "physical",
			},
		},
		ArtStandard: `
    ⠀⠀⠀⠀⠀⢀⣀⣀⡀⠀⠀⠀⠀⠀
    ⠀⠀⠀⢀⣴⣿⣿⣿⣿⣷⣦⡀⠀⠀
    ⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀
    ⠀⣼⣿⣿⡟⠁⠀⠀⠈⢻⣿⣿⣿⣧
    ⢸⣿⣿⣿⠃⡀⠀⠀⢀⠈⣿⣿⣿⣿
    ⣿⣿⣿⡏⠀⣿⡄⢠⣾⠀⢸⣿⣿⣿
    ⣿⣿⣿⣇⠀⠙⠃⠘⠃⠀⣸⣿⣿⡿
    ⠈⢿⣿⣿⣷⣤⣀⣀⣴⣾⣿⣿⡿⠃
    ⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀
    ⠀⠀⠀⠀⠈⠛⠛⠛⠋⠁⠀⠀⠀⠀
`,
		ArtShiny: `
    ⠀⠀⠀⠀⠀⢀⣀⣀⡀⠀⠀⠀⠀⠀
    ⠀⠀⠀⢀⣴⣿⣿⣿⣿⣷⣦⡀⠀⠀
    ⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀
    ⠀⣼⣿⣿⡟⠁✨⠀⠈⢻⣿⣿⣿⣧
    ⢸⣿⣿⣿⠃⡀⠀⠀⢀⠈⣿⣿⣿⣿
    ⣿⣿⣿⡏⠀⣿⡄⢠⣾⠀⢸⣿⣿⣿
    ⣿⣿⣿⣇⠀⠙⠃⠘⠃⠀⣸⣿⣿⡿
    ⠈⢿⣿⣿⣷⣤⣀⣀⣴⣾⣿⣿⡿⠃
    ⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀
    ⠀⠀⠀⠀⠈⠛⠛⠛⠋⠁⠀⠀⠀⠀
`,
		Evolution: &models.EvolutionChain{
			Base: models.EvolutionStage{PokemonID: 1, Name: "Bulbasaur", Trigger: "level-up", MinLevel: 0, Item: ""},
			Evolution: []models.EvolutionStage{
				{PokemonID: 2, Name: "Ivysaur", Trigger: "level-up", MinLevel: 16, Item: ""},
				{PokemonID: 3, Name: "Venusaur", Trigger: "level-up", MinLevel: 32, Item: ""},
			},
		},
	},
	{
		ID:             4,
		NamePT:         "Charmander",
		NameEN:         "Charmander",
		Generation:     1,
		Types:          []string{"fogo"},
		Height:         6.0,
		Weight:         85.0,
		BaseExperience: 62,
		Stats: models.PokemonStats{
			HP:      39,
			Attack:  52,
			Defense: 43,
			SpAtk:   60,
			SpDef:   50,
			Speed:   65,
		},
		SignatureMoves: []models.Move{
			{
				NamePT:   "Ember",
				NameEN:   "Ember",
				Type:     "fogo",
				Power:    40,
				Category: "special",
			},
			{
				NamePT:   "Flamethrower",
				NameEN:   "Flamethrower",
				Type:     "fogo",
				Power:    90,
				Category: "special",
			},
		},
		ArtStandard: `
    ⠀⠀⠀⠀⠀⠀⢀⣀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⢀⣾⣿⣿⣷⡀⠀⠀⠀⠀
    ⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀
    ⠀⠀⠀⣾⣿⡏⠉⠉⢹⣿⣷⠀⠀⠀
    ⠀⠀⣸⣿⣿⣇⠀⠀⣸⣿⣿⡆⠀⠀
    ⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀
    ⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀🔥⠀⠀
`,
		ArtShiny: `
    ⠀⠀⠀⠀⠀⠀⢀⣀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⢀⣾⣿⣿⣷⡀⠀⠀⠀⠀
    ⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀
    ⠀⠀⠀⣾⣿⡏⠉⠉⢹⣿⣷⠀⠀⠀
    ⠀⠀⣸⣿⣿⣇✨✨⣸⣿⣿⡆⠀⠀
    ⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀
    ⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀💛⠀⠀
`,
		Evolution: &models.EvolutionChain{
			Base: models.EvolutionStage{PokemonID: 4, Name: "Charmander", Trigger: "level-up", MinLevel: 0, Item: ""},
			Evolution: []models.EvolutionStage{
				{PokemonID: 5, Name: "Charmeleon", Trigger: "level-up", MinLevel: 16, Item: ""},
				{PokemonID: 6, Name: "Charizard", Trigger: "level-up", MinLevel: 36, Item: ""},
			},
		},
	},
	{
		ID:             7,
		NamePT:         "Squirtle",
		NameEN:         "Squirtle",
		Generation:     1,
		Types:          []string{"água"},
		Height:         5.0,
		Weight:         90.0,
		BaseExperience: 63,
		Stats: models.PokemonStats{
			HP:      44,
			Attack:  48,
			Defense: 65,
			SpAtk:   50,
			SpDef:   64,
			Speed:   43,
		},
		SignatureMoves: []models.Move{
			{
				NamePT:   "Water Gun",
				NameEN:   "Water Gun",
				Type:     "água",
				Power:    40,
				Category: "special",
			},
			{
				NamePT:   "Hydro Pump",
				NameEN:   "Hydro Pump",
				Type:     "água",
				Power:    110,
				Category: "special",
			},
		},
		ArtStandard: `
    ⠀⠀⠀⠀⢀⣤⣤⣤⣤⡀⠀⠀⠀⠀
    ⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀
    ⠀⠀⣼⣿⢿⡿⢿⡿⢿⣿⣿⣧⠀⠀
    ⠀⣸⣿⣿⣇⠀⠀⠀⣸⣿⣿⣿⡇⠀
    ⠀⣿⣿⣿⣿⠀⠀⠀⣿⣿⣿⣿⣿⠀
    ⠀⣿⣿⣿⣿⣇⠀⣸⣿⣿⣿⣿⡿⠀
    ⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀
    ⠀⠀⠙⠿⣿⣿⣿⣿⣿⠿⠋⠀💧
`,
		ArtShiny: `
    ⠀⠀⠀⠀⢀⣤⣤⣤⣤⡀⠀⠀⠀⠀
    ⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀
    ⠀⠀⣼⣿⢿⡿⢿⡿⢿⣿⣿⣧⠀⠀
    ⠀⣸⣿⣿⣇⠀✨⠀⣸⣿⣿⣿⡇⠀
    ⠀⣿⣿⣿⣿⠀⠀⠀⣿⣿⣿⣿⣿⠀
    ⠀⣿⣿⣿⣿⣇⠀⣸⣿⣿⣿⣿⡿⠀
    ⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀
    ⠀⠀⠙⠿⣿⣿⣿⣿⣿⠿⠋⠀💜
`,
		Evolution: &models.EvolutionChain{
			Base: models.EvolutionStage{PokemonID: 7, Name: "Squirtle", Trigger: "level-up", MinLevel: 0, Item: ""},
			Evolution: []models.EvolutionStage{
				{PokemonID: 8, Name: "Wartortle", Trigger: "level-up", MinLevel: 16, Item: ""},
				{PokemonID: 9, Name: "Blastoise", Trigger: "level-up", MinLevel: 36, Item: ""},
			},
		},
	},
	{
		ID:             25,
		NamePT:         "Pikachu",
		NameEN:         "Pikachu",
		Generation:     1,
		Types:          []string{"elétrico"},
		Height:         4.0,
		Weight:         60.0,
		BaseExperience: 112,
		Stats: models.PokemonStats{
			HP:      35,
			Attack:  55,
			Defense: 40,
			SpAtk:   50,
			SpDef:   50,
			Speed:   90,
		},
		SignatureMoves: []models.Move{
			{
				NamePT:   "Thunderbolt",
				NameEN:   "Thunderbolt",
				Type:     "elétrico",
				Power:    90,
				Category: "special",
			},
			{
				NamePT:   "Quick Attack",
				NameEN:   "Quick Attack",
				Type:     "normal",
				Power:    40,
				Category: "physical",
			},
			{
				NamePT:   "Iron Tail",
				NameEN:   "Iron Tail",
				Type:     "metálico",
				Power:    100,
				Category: "physical",
			},
			{
				NamePT:   "Thunder Wave",
				NameEN:   "Thunder Wave",
				Type:     "elétrico",
				Power:    0,
				Category: "status",
			},
		},
		ArtStandard: `
⣿⠁⣿⣿⣿⣿⣿⣿⣿⣿
⠀⠀⢀⣀⣀⣀⣀⣀⣀⣀⣀⣀⡀⠀⠀
⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀
⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇
⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁
⠀⠹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀
⠀⠀⠈⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠁⠀⠀
`,
		ArtShiny: `
⣿⠁⣿⣿⣿⣿⣿⣿⣿⣿
⠀⠀⢀⣀⣀⣀⣀⣀⣀⣀⣀⣀⡀⠀⠀
⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀
⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇
⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁
⠀⠹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀
⠀⠀⠈⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠁⠀⠀
`,
		Evolution: &models.EvolutionChain{
			Base: models.EvolutionStage{PokemonID: 172, Name: "Pichu", Trigger: "friendship", MinLevel: 0, Item: ""},
			Evolution: []models.EvolutionStage{
				{PokemonID: 25, Name: "Pikachu", Trigger: "stone", MinLevel: 0, Item: "Thunder Stone"},
				{PokemonID: 26, Name: "Raichu", Trigger: "", MinLevel: 0, Item: ""},
			},
		},
	},
	{
		ID:             150,
		NamePT:         "Mewtwo",
		NameEN:         "Mewtwo",
		Generation:     1,
		Types:          []string{"psíquico"},
		Height:         20.0,
		Weight:         1220.0,
		BaseExperience: 340,
		Stats: models.PokemonStats{
			HP:      106,
			Attack:  110,
			Defense: 90,
			SpAtk:   154,
			SpDef:   90,
			Speed:   130,
		},
		SignatureMoves: []models.Move{
			{
				NamePT:   "Psychic",
				NameEN:   "Psychic",
				Type:     "psíquico",
				Power:    90,
				Category: "special",
			},
			{
				NamePT:   "Shadow Ball",
				NameEN:   "Shadow Ball",
				Type:     "fantasma",
				Power:    80,
				Category: "special",
			},
			{
				NamePT:   "Psystrike",
				NameEN:   "Psystrike",
				Type:     "psíquico",
				Power:    100,
				Category: "special",
			},
		},
		ArtStandard: `
⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣶⣶⣶⣶⣤⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀
⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀⠀
⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠀
⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀
⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧
`,
		ArtShiny: `
⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣶⣶⣶⣶⣤⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀
⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀⠀
⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠀
⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀
⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧
`,
		Evolution: nil,
	},
}

func GetSamplePokedex() *models.Pokedex {
	pokedex := models.NewPokedex()
	for _, pokemon := range SamplePokemon {
		pokedex.AddPokemon(pokemon)
	}
	return pokedex
}

func GetPokedex() *models.Pokedex {
	return GetSamplePokedex()
}
