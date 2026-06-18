package auth

import (
	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService *authService
}

func NewAuthHandler(authService *authService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	resp, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "User registered successfully!", "data": resp})
}

func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Login Success", "data": resp})
}

func (h *authHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"success": false, "message": "Unauthorized: user not found in context"})
		return
	}

	user, err := h.authService.Profile(userID.(string))
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "User profile retrieved successfully!", "data": user})
}

func (h *authHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"success": false, "message": "Unauthorized: user not found in context"})
		return
	}

	var req dto.UpdateProfileRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = h.authService.UpdateProfile(userID.(string), &req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "User profile updated successfully!"})
}

func (h *authHandler) Logout(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"success": false, "message": "Unauthorized: user not logged in"})
		return
	}

	err := h.authService.Logout()
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "User logged out successfully!"})
}
