package internals

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type Game struct {
	ID        int    `json:"id"`
	UserA     int    `json:"user_a"`
	UserB     int    `json:"user_b"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

type UserOnline struct {
	ID     int  `json:"id"`
	Online bool `json:"online"`
}

type UserLevel struct {
	ID    int `json:"id"` // user id
	Level int `json:"level"`
}
