package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type authController struct{}

type LoginRequest struct {
	Token string `json:"token" binding:"required"`
}

func (a authController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appToken := os.Getenv("APP_TOKEN")
	if appToken == "" {
		// Fallback or error if env not set, for safety.
		// Assuming config is required.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server misconfiguration"})
		return
	}

	if request.Token != appToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

var AuthController = authController{}
