package middlewares

import "github.com/gin-gonic/gin"

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")

		if !exists {
			c.JSON(401, gin.H{"message": "Unathorized: Role Not Found!"})
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole.(string) == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have permission!"})
			c.Abort()
			return
		}

		c.Next()

	}
}
