package middleware

import (
	"loan/config"
	"loan/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.BadRequest(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.BadRequest(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := config.ValidateToken(parts[1])
		if err != nil {
			utils.BadRequest(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set both user ID and roles in context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
