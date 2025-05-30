package internals

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GameRoom struct {
	game  *GamlooGame
	conns []*websocket.Conn
}

var (
	gameConnections  = make(map[string]GameRoom)
	connectionsMutex sync.RWMutex
)

func broadcastToGame(gameID string, message []byte) {
	connectionsMutex.RLock()
	defer connectionsMutex.RUnlock()

	room, exists := gameConnections[gameID]
	if !exists {
		return
	}

	for _, conn := range room.conns {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Printf("Error broadcasting to connection: %v\n", err)
			// remove failed connections
		}
	}
}

type GameHandler struct {
	gameRepo *GameRepo
	userRepo *UserRepo
	ws       *websocket.Upgrader
}

func NewGameHandler(gameRepo *GameRepo, userRepo *UserRepo) *GameHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &GameHandler{
		gameRepo: gameRepo,
		userRepo: userRepo,
		ws:       &upgrader,
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var game GameRequestDto
	if err := c.ShouldBindJSON(&game); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, err := h.userRepo.CheckUserInGame(game.UserA)
	if err != nil {
		fmt.Println("Error checking user in game:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User exists in game:", exists)

	if exists {
		c.JSON(400, gin.H{"error": "user is already in a game"})
		return
	}

	matchedUser, err := h.userRepo.GetRandomOnlineUser(game.UserA)
	fmt.Println("Matched user:", matchedUser)
	if err != nil {
		fmt.Println("Error getting random online user:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if matchedUser == 0 {
		c.JSON(404, gin.H{"error": "no online users found"})
		return
	}

	UserAInt, err := strconv.Atoi(game.UserA)
	if err != nil {
		fmt.Println("Error converting matched user ID:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	gameResponse := GameResponseDto{
		UserA: UserAInt,
		UserB: int(matchedUser),
	}
	if err := h.gameRepo.CreateGame(gameResponse); err != nil {
		fmt.Println("Error creating game:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "ok", "matched_user": fmt.Sprintf("%d", int(matchedUser))})
}

func (h *GameHandler) WebSocketGame(c *gin.Context) {
	gameID := c.Param("id")
	fmt.Println("Game ID:", gameID)

	gameUsers, err := h.gameRepo.GetUserIDsFromGameID(gameID)
	if err != nil {
		fmt.Println("Error getting user IDs from game ID:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		fmt.Println("Error converting game ID:", err)
		c.JSON(400, gin.H{"error": "invalid game ID"})
		return
	}

	conn, err := h.ws.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	connectionsMutex.Lock()

	room, exists := gameConnections[gameID]
	if !exists {
		gamloo := NewGamlooGame(gameIDInt, int(gameUsers.UserA), int(gameUsers.UserB))
		gamloo.PrintBoard()

		room = GameRoom{
			game:  gamloo,
			conns: []*websocket.Conn{},
		}
		room.conns = append(room.conns, conn)
		gameConnections[gameID] = room

		conn.WriteMessage(websocket.TextMessage, []byte("waiting for opponent to join"))
	} else {
		room.conns = append(room.conns, conn)
		gameConnections[gameID] = room
	}
	connectionsMutex.Unlock()

	defer func() {
		connectionsMutex.Lock()
		// for i, c := range gameConnections[gameID].conns {
		// 	if c == conn {
		// 		fmt.Println("Removing connection from game:", gameID)
		// 		// gameConnections[gameID].conns = append(gameConnections[gameID].conns[:i], gameConnections[gameID].conns[i+1:]...)
		// 		gameCons := gameConnections[gameID]
		// 		gameCons.conns = append(gameCons.conns[:i], gameCons.conns[i+1:]...)
		// 		break
		// 	}
		// }
		if len(gameConnections[gameID].conns) == 2 {
			for _, c := range gameConnections[gameID].conns {
				c.WriteMessage(websocket.TextMessage, []byte("game completed"))
				c.Close()
			}
			delete(gameConnections, gameID)

			h.gameRepo.CompleteGame(gameID)
			fmt.Println("Game completed and removed from database")
		}
		connectionsMutex.Unlock()

		conn.Close()
		fmt.Println("WebSocket connection closed and removed from game", gameID, gameConnections[gameID].conns)
	}()

	currentGame := gameConnections[gameID].game

	for {
		var userGameResponseDto UserGameResponseDto
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Println("Received message type:", messageType)
		fmt.Println("Received message:", string(msg))

		if err := json.Unmarshal(msg, &userGameResponseDto); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			break
		}
		fmt.Println("User game response DTO:", userGameResponseDto)

		if userGameResponseDto.User != currentGame.UserA && userGameResponseDto.User != currentGame.UserB {
			fmt.Println("Invalid user in message")
			if err := conn.WriteMessage(websocket.TextMessage, []byte("invalid user")); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
			continue
		}

		if userGameResponseDto.User != currentGame.turn {
			fmt.Println("Not your turn")
			broadcastToGame(gameID, []byte("{\"error\":\"not your turn\", \"turn\":\""+currentGame.turn+"\"}"))
			continue
		}

		legal, err := currentGame.CheckState(userGameResponseDto)
		if !legal {
			fmt.Println("Illegal move")
			if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"error\":\""+err.Error()+"\"}")); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
			continue
		}
		checkWin, _ := currentGame.CheckWin()
		if checkWin {
			fmt.Println("User", userGameResponseDto.User, "wins!")
			broadcastToGame(gameID, []byte("win"))
			break
		}

		data := map[string]interface{}{
			"board": currentGame.Board,
			"turn":  currentGame.turn,
		}
		newState, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling new state:", err)
			break
		}
		broadcastToGame(gameID, newState)
		fmt.Printf("Received message: %s\n", msg)
	}
	fmt.Println("WebSocket connection closed")
}
