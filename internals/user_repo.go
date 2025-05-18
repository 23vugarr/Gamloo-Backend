package internals

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/supabase-community/supabase-go"
)

type UserRepo struct {
	client *supabase.Client
}

func NewUserRepo(client *supabase.Client) *UserRepo {
	return &UserRepo{
		client: client,
	}
}

func (r *UserRepo) CreateUser(user UserRequestDto) error {
	randomNumber := make([]byte, 16)
	_, err := rand.Read(randomNumber)
	if err != nil {
		return err
	}
	id := int(randomNumber[0])
	data, count, err := r.client.From("users").Insert(map[string]interface{}{
		"username": user.Username,
		"id":       id,
	}, true, "", "", "").Execute()

	if err != nil {
		return err
	}

	onlineData, _, err := r.client.From("online_users").Insert(map[string]interface{}{
		"id":     id,
		"online": true,
	}, true, "", "", "").Execute()
	if err != nil {
		fmt.Println("Error inserting online status:", err)
		return err
	}

	fmt.Println("Inserted user:", data, count)
	fmt.Println("Inserted online status:", onlineData)
	return nil
}

func (r *UserRepo) SetOnline(userID int) error {
	exists, _, err := r.client.From("online_users").Select("id", "1", false).Eq("id", fmt.Sprintf("%d", userID)).Execute()
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return err
	}

	if len(exists) == 0 {
		fmt.Println("User not found in online_users")
		onlineData, _, err := r.client.From("online_users").Insert(map[string]interface{}{
			"id":     userID,
			"online": true,
		}, true, "", "", "").Execute()
		if err != nil {
			fmt.Println("Error inserting online status:", err)
			return err
		}
		fmt.Println("Inserted online status:", onlineData)
	}

	timeNow := time.Now()
	data, count, err := r.client.From("online_users").Update(map[string]interface{}{
		"online":     true,
		"updated_at": timeNow,
	}, "", "").Eq("id", fmt.Sprintf("%d", userID)).Execute()

	if err != nil {
		fmt.Println("Error updating user online status:", err)
		return err
	}

	fmt.Println("Updated user:", data, count)
	return nil
}

func (r *UserRepo) GetAllOnlineUsers() ([]map[string]interface{}, error) {
	lastFiveSeconds := time.Now().Add(-5 * time.Second).Unix()
	lastFiveSecondsTime := time.Unix(lastFiveSeconds, 0).Format("2006-01-02 15:04:05")
	data, _, err := r.client.From("online_users").Select("id", "1", false).Eq("online", "true").Lte("updated_at", lastFiveSecondsTime).Execute()
	if err != nil {
		fmt.Println("Error fetching online users:", err)
		return nil, err
	}

	var users []map[string]interface{}
	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) SetOffline(userID int) error {
	timeNow := time.Now()
	data, count, err := r.client.From("online_users").Update(map[string]interface{}{
		"online":     false,
		"updated_at": timeNow,
	}, "", "").Eq("id", fmt.Sprintf("%d", userID)).Execute()

	if err != nil {
		fmt.Println("Error updating user online status:", err)
		return err
	}

	fmt.Println("Updated user:", data, count)
	return nil
}

func (r *UserRepo) GetRandomOnlineUser(selfId string) (float64, error) {
	data, _, err := r.client.From("online_users").Select("id", "1", false).Eq("online", "true").Neq("id", selfId).Execute()
	if err != nil {
		fmt.Println("Error fetching random online user:", err)
		return 0, err
	}

	fmt.Print("data from online_users: ", string(data))

	if len(data) == 0 {
		return 0, fmt.Errorf("no online users found")
	}

	var userId []map[string]interface{}
	err = json.Unmarshal(data, &userId)
	resListMap := make([]map[string]interface{}, 0)

	for _, user := range userId {
		insideUserId := fmt.Sprintf("%v", user["id"])
		_, count, err := r.client.From("games").Select("user_a", "user_b", false).Neq("user_a", insideUserId).Neq("user_b", insideUserId).Execute()
		if err != nil {
			fmt.Println("Error checking user in games table:", err)
			return 0, err
		}
		if count == 0 {
			resListMap = append(resListMap, user)
		}
	}

	lenOfArray := len(resListMap)
	if lenOfArray == 0 {
		return 0, fmt.Errorf("no online users found")
	}
	if lenOfArray > 1 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(lenOfArray)
		userId = resListMap[randomIndex : randomIndex+1]
		fmt.Println("Random user ID:", userId)
	}

	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return 0, err
	}
	return userId[0]["id"].(float64), nil
}

func (r *UserRepo) CheckUserInGame(userID string) (bool, error) {
	fmt.Print("Checking user in game with ID: \n", userID)
	data1, _, _ := r.client.From("games").Select("user_a", "user_b", false).Eq("user_a", userID).Eq("completed", "false").Execute()
	data2, _, err := r.client.From("games").Select("user_a", "user_b", false).Eq("user_b", userID).Eq("completed", "false").Execute()

	if err != nil {
		fmt.Println("Error checking user in game:", err)
		return false, err
	}
	fmt.Println("Data1:", string(data1))
	fmt.Println("Data2:", string(data2))

	if len(data1) < 3 && len(data2) < 3 {
		return false, nil
	}
	return true, nil
}
