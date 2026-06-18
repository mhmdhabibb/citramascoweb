package user

import (
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService *userService
}

func NewUserHandler(userService *userService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetAllCustomer(c *gin.Context) {
	customers, err := h.userService.GetAllCustomer()
	if err != nil {
		c.JSON(400, gin.H{"message": "Error getting customers data!", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "List of customers", "data": customers})
}

func (h *userHandler) GetAllByRole(c *gin.Context) {
	role := c.Query("role")
	users, err := h.userService.GetAllByRole(role)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error getting users data!", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "List of users", "data": users})
}

func (h *userHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error deleting user!", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
