package models

type Generation struct {
	ID           int
	NamePT       string
	NameEN       string
	Region       string
	PokemonIDs   []int
	PokemonCount int
}

type EvolutionStage struct {
	PokemonID int
	Name      string
	Trigger   string
	MinLevel  int
	Item      string
}

type EvolutionChain struct {
	Base      EvolutionStage
	Evolution []EvolutionStage
}

func (ec *EvolutionChain) GetStageNames() []string {
	names := []string{ec.Base.Name}
	for _, e := range ec.Evolution {
		names = append(names, e.Name)
	}
	return names
}

func (ec *EvolutionChain) GetPokemonIDs() []int {
	ids := []int{ec.Base.PokemonID}
	for _, e := range ec.Evolution {
		ids = append(ids, e.PokemonID)
	}
	return ids
}

func (ec *EvolutionChain) FindStage(pokemonID int) int {
	for i, id := range ec.GetPokemonIDs() {
		if id == pokemonID {
			return i
		}
	}
	return -1
}

func (ec *EvolutionChain) GetNextStage(currentID int) *EvolutionStage {
	stage := ec.FindStage(currentID)
	if stage >= 0 && stage < len(ec.Evolution) {
		return &ec.Evolution[stage]
	}
	return nil
}

func (ec *EvolutionChain) GetPrevStage(currentID int) *EvolutionStage {
	stage := ec.FindStage(currentID)
	if stage > 0 {
		if stage == 1 {
			return &ec.Base
		}
		return &ec.Evolution[stage-2]
	}
	return nil
}
