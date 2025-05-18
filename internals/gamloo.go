package internals

import "fmt"

type GamlooGame struct {
	ID        int    `json:"id"`
	UserA     string `json:"user_a"`
	UserB     string `json:"user_b"`
	Board     [][]string
	LastState [][]string
	turn      string
}

func NewGamlooGame(id int, usera int, userb int) *GamlooGame {
	initialBoard := [][]string{
		{fmt.Sprintf("%d", usera), fmt.Sprintf("%d", usera), fmt.Sprintf("%d", usera)},
		{"0", "0", "0"},
		{fmt.Sprintf("%d", userb), fmt.Sprintf("%d", userb), fmt.Sprintf("%d", userb)},
	}

	return &GamlooGame{
		ID:        id,
		UserA:     fmt.Sprintf("%d", usera),
		UserB:     fmt.Sprintf("%d", userb),
		Board:     initialBoard,
		LastState: initialBoard,
		turn:      fmt.Sprintf("%d", usera),
	}
}

func (g *GamlooGame) PrintBoard() {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "0" {
				print(" . ")
			} else if cell == g.UserA {
				print(" A ")
			} else if cell == g.UserB {
				print(" B ")
			}
		}
		println()
	}
}

func (g *GamlooGame) CheckState(userMessage UserGameResponseDto) (bool, error) {
	g.Board = userMessage.NewState
	g.LastState = userMessage.OldState
	if userMessage.User == g.UserA {
		g.turn = g.UserB
	} else {
		g.turn = g.UserA
	}
	return true, nil
}

func (g *GamlooGame) CheckWin() (bool, error) {
	return false, nil
}
