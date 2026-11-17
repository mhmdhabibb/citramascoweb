package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *authHandler {
	return &authHandler{service: service}
}

func (h *authHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username and password are required",
			"error":   err.Error(),
		})
		return
	}

	res, err := h.service.SignIn(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signed in successfully",
		"data":    res,
	})
}

func (h *authHandler) Profile(c *gin.Context) {
	userId := c.GetString("user_id")
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	token := ""
	if len(parts) == 2 {
		token = parts[1]
	}

	profile, err := h.service.GetProfile(userId, token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"user_id": userId,
				"role":    c.GetString("role"),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

