package models

import (
	"sort"
)

type MovePool struct {
	Moves       []*Move
	MovesByType map[string][]*Move
}

func NewMovePool() *MovePool {
	return &MovePool{
		Moves:       make([]*Move, 0),
		MovesByType: make(map[string][]*Move),
	}
}

func (mp *MovePool) AddMove(move *Move) {
	mp.Moves = append(mp.Moves, move)
	mp.MovesByType[move.Type] = append(mp.MovesByType[move.Type], move)
}

func SelectSignatureMoves(pokemon *Pokemon, allMoves []*Move, pokemonMoves []*Move) []Move {
	candidates := make([]*Move, 0)

	for _, move := range pokemonMoves {
		for _, poolMove := range allMoves {
			if poolMove.NameEN == move.NameEN {
				candidates = append(candidates, poolMove)
				break
			}
		}
	}

	scores := make(map[int]int)
	for i, move := range candidates {
		score := 0

		hasSTAB := false
		for _, t := range pokemon.Types {
			if move.Type == t {
				hasSTAB = true
				break
			}
		}
		if hasSTAB {
			score += 50
		}

		score += move.Power

		if move.Category == "status" {
			score += 20
		}

		scores[i] = score
	}

	sort.Slice(candidates, func(i, j int) bool {
		return scores[i] > scores[j]
	})

	result := make([]Move, 0, 5)
	for i := 0; i < len(candidates) && i < 5; i++ {
		result = append(result, *candidates[i])
	}

	return result
}
