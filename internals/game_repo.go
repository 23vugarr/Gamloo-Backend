package internals

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/supabase-community/supabase-go"
)

type GameRepo struct {
	client *supabase.Client
}

func NewGameRepo(client *supabase.Client) *GameRepo {
	return &GameRepo{
		client: client,
	}
}

func (r *GameRepo) CreateGame(game GameResponseDto) error {
	randomNumber := make([]byte, 16)
	_, err := rand.Read(randomNumber)
	if err != nil {
		return err
	}
	id := int(randomNumber[0])
	data, count, err := r.client.From("games").Insert(map[string]interface{}{
		"id":     id,
		"user_a": game.UserA,
		"user_b": game.UserB,
	}, true, "", "", "").Execute()

	if err != nil {
		fmt.Println("Error inserting game:", err)
		return err
	}

	fmt.Println("Inserted game:", data, count)
	return nil
}

func (r *GameRepo) GetUserIDsFromGameID(id string) (GameResponseDto, error) {
	var game []GameResponseDto
	data, _, err := r.client.From("games").Select("user_a, user_b", "1", false).Eq("id", id).Eq("completed", "false").Execute()
	fmt.Println("Data from game:", string(data))
	if err != nil {
		fmt.Println("Error getting game by ID:", err)
		return GameResponseDto{}, err
	}

	err = json.Unmarshal(data, &game)
	if err != nil {
		fmt.Println("Error unmarshalling game data:", err)
		return GameResponseDto{}, err
	}

	if len(game) == 0 {
		fmt.Println("No game found with the given ID")
		return GameResponseDto{}, fmt.Errorf("no game found with ID %s", id)
	}
	return game[0], nil
}

func (r *GameRepo) CompleteGame(id string) error {
	_, _, err := r.client.From("games").Update(map[string]interface{}{
		"completed": "true",
	}, "", "").Eq("id", id).Execute()
	if err != nil {
		fmt.Println("Error completing game:", err)
		return err
	}
	return nil
}
