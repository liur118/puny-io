package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liur/puny-io/internal/service"
)

type AuthHandler struct {
	jwtService *service.JWTService
	users      map[string]string
}

func NewAuthHandler(jwtService *service.JWTService, users map[string]string) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		users:      users,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	password, exists := h.users[req.Username]
	if !exists || password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.jwtService.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
	})
}
