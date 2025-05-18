package internals

type UserRequestDto struct {
	Username string `json:"username" binding:"required"`
}

type UserOnlineDto struct {
	ID string `json:"id"`
}

type GameRequestDto struct {
	UserA string `json:"user_a" binding:"required"`
}

type GameResponseDto struct {
	UserA int `json:"user_a"`
	UserB int `json:"user_b"`
}

type GameRoomDto struct {
	ID        int     `json:"id"`
	UserA     string  `json:"user_a"`
	UserB     string  `json:"user_b"`
	Board     [][]int `json:"board"`
	Turn      int     `json:"turn"`
	LastState string  `json:"last_state"`
}

type UserGameResponseDto struct {
	User     string     `json:"user"`
	NewState [][]string `json:"new_state"`
	OldState [][]string `json:"old_state"`
}
