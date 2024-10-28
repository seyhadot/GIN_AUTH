package middleware

import (
	"loan/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		utils.Info("Request processed",
			utils.Fields(map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     c.Writer.Status(),
				"duration":   duration,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}),
		)
	}
}
