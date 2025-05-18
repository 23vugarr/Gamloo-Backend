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

type Position struct {
	Row int
	Col int
}

var validityMap map[Position][]Position = map[Position][]Position{
	{0, 0}: {{0, 1}, {1, 0}, {1, 1}},
	{0, 1}: {{0, 0}, {0, 2}, {1, 1}},
	{0, 2}: {{0, 1}, {1, 1}, {1, 2}},
	{1, 0}: {{0, 0}, {1, 1}, {2, 0}},
	{1, 1}: {{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
	{1, 2}: {{0, 2}, {1, 1}, {2, 2}},
	{2, 0}: {{1, 0}, {1, 1}, {2, 1}},
	{2, 1}: {{2, 0}, {1, 1}, {2, 2}},
	{2, 2}: {{1, 1}, {1, 2}, {2, 1}},
}

func NewGamlooGame(id int, usera int, userb int) *GamlooGame {
	initialBoard := [][]string{
		{fmt.Sprintf("%d", usera), fmt.Sprintf("%d", usera), fmt.Sprintf("%d", usera)},
		{".", ".", "."},
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
			if cell == "." {
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

	valid, err := g.CheckMoveValidity(userMessage)
	if !valid {
		return false, err
	}

	if userMessage.User == g.UserA {
		g.turn = g.UserB
	} else {
		g.turn = g.UserA
	}
	return true, nil
}

func isInSlice(slice []Position, item Position) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (g *GamlooGame) CheckMoveValidity(userMessage UserGameResponseDto) (bool, error) {
	changes := g.CaptureChange(userMessage.OldState, userMessage.NewState)
	if len(changes) != 2 {
		return false, fmt.Errorf("invalid move: more than one change detected")
	}

	initialPos := Position{Row: changes[0][0], Col: changes[0][1]}
	lastPos := Position{Row: changes[1][0], Col: changes[1][1]}

	if g.Board[initialPos.Row][initialPos.Col] != "." {
		return false, fmt.Errorf("invalid move: cell %s is already occupied", changes[0][0])
	}

	if !isInSlice(validityMap[initialPos], lastPos) {
		return false, fmt.Errorf("invalid move: cell %s is not a valid move from cell %s", lastPos, initialPos)
	}

	return true, nil
}

func (g *GamlooGame) CaptureChange(old, new [][]string) [][]int {
	var changes [][]int
	for i := 0; i < len(old); i++ {
		for j := 0; j < len(old[i]); j++ {
			if old[i][j] != new[i][j] {
				changes = append(changes, []int{i, j})
			}
		}
	}
	var finalOrder [][]int
	if old[changes[0][0]][changes[0][1]] != "." {
		finalOrder = append(finalOrder, changes[0])
		finalOrder = append(finalOrder, changes[1])
	} else {
		finalOrder = append(finalOrder, changes[1])
		finalOrder = append(finalOrder, changes[0])
	}

	fmt.Println("Changes detected:", finalOrder)
	return finalOrder
}

func (g *GamlooGame) CheckWin() (bool, error) {
	if g.Board[0][0] == g.Board[0][1] && g.Board[0][1] == g.Board[0][2] {
		if g.Board[0][0] == g.UserB {
			fmt.Println("User A wins")
			return true, nil
		}
	} else if g.Board[2][0] == g.Board[2][1] && g.Board[2][1] == g.Board[2][2] {
		if g.Board[2][0] == g.UserA {
			fmt.Println("User B wins")
			return true, nil
		}
	}
	return false, nil
}
