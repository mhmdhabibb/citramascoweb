package middlewares

import (
	"citramascoweb-backend/config"
	"citramascoweb-backend/pkg/utils"
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
		var tokenString string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized!"})
			return
		}

		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token was expired or invalid!"})
			return
		}

		role := ""
		// Look up user role strictly from Citramasco's user_roles table
		if config.DB != nil && claims.UserId != "" {
			var roleSlug string
			err := config.DB.Table("user_roles").
				Select("roles.slug").
				Joins("JOIN roles ON roles.id = user_roles.role_id").
				Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", claims.UserId).
				Limit(1).
				Pluck("roles.slug", &roleSlug).Error

			if err == nil && roleSlug != "" {
				role = roleSlug
			}
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", role)
		c.Set("claims", claims)

		c.Next()
	}
}
