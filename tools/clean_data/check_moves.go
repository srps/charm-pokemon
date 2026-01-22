//go:build check_moves

package main

import (
	"fmt"
	"sort"
)

// This will be run manually or via go run

func main() {
	unique := make(map[string]bool)
	for _, moves := range CuratedMoves {
		for _, m := range moves {
			unique[m] = true
		}
	}

	moveList := make([]string, 0, len(unique))
	for m := range unique {
		moveList = append(moveList, m)
	}
	sort.Strings(moveList)

	fmt.Printf("Total unique moves: %d\n", len(moveList))
	for _, m := range moveList {
		fmt.Printf("\"%s\",\n", m)
	}
}
