package internals

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *UserRepo
}

func NewUserHandler(userRepo *UserRepo) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user UserRequestDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "ok"})
}

func (h *UserHandler) SetOnline(c *gin.Context) {
	var user UserOnlineDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(user.ID)
	if err != nil {
		fmt.Println("Error converting user ID:", err)
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userRepo.SetOnline(userID); err != nil {
		fmt.Println("Error setting user online:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *UserHandler) SetOffline(c *gin.Context) {
	var user UserOnlineDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(user.ID)
	if err != nil {
		fmt.Println("Error converting user ID:", err)
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userRepo.SetOffline(userID); err != nil {
		fmt.Println("Error setting user offline:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
