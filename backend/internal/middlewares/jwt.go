package middlewares

import (
	"citramascoweb-backend/pkg/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Unatuhorized!"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid Authorization Format!"})
			return
		}

		tokenString := parts[1]

		claims, err := utils.VerifyJWT(tokenString)

		if err != nil {
			c.JSON(401, gin.H{"error": "Token was expired!"})
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)

		fmt.Println(claims.UserId)
		fmt.Println(claims.Role)

		c.Next()

	}
}
