package main

import (
	"fmt"

	"github.com/23vugarr/gamloo/internals"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func main() {
	router := gin.Default()
	APIURL := "https://vtyrqmskrmeovjyadiqk.supabase.co"
	APIKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InZ0eXJxbXNrcm1lb3ZqeWFkaXFrIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDc1MDExMzIsImV4cCI6MjA2MzA3NzEzMn0.p4msCKeLNnXxDj2vVdYUBZocOGYJi9IbqLNSzQcyoZA"
	supabeClient, err := supabase.NewClient(APIURL, APIKey, &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	userRepo := internals.NewUserRepo(supabeClient)
	gameRepo := internals.NewGameRepo(supabeClient)

	uHandler := internals.NewUserHandler(userRepo)
	gHandler := internals.NewGameHandler(gameRepo, userRepo)

	v1 := router.Group("/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		userGroup := v1.Group("/user")
		{
			userGroup.POST("/", uHandler.CreateUser)
			userGroup.POST("/online", uHandler.SetOnline)
			userGroup.POST("/offline", uHandler.SetOffline)
		}

		gameGroup := v1.Group("/game")
		{
			gameGroup.GET("/:id", gHandler.WebSocketGame)
			gameGroup.POST("/", gHandler.CreateGame)
		}
	}

	go internals.SetUserOfflineEvery5Seconds(userRepo)
	router.Run(":8085")
}
