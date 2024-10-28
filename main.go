package main

import (
	"loan/config"
	"loan/middleware"
	"loan/routes"
	"loan/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize application
	utils.Debug("Starting application")
	utils.Info("Environment: " + os.Getenv("GIN_MODE"))
	utils.Info("Log Level: " + os.Getenv("LOG_LEVEL"))

	// Load configuration
	cfg := config.LoadConfig()

	// Setup Gin
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(config.GetCorsConfig()))

	// Setup rate limiter - 100 requests per minute
	limiter := middleware.NewIPRateLimiter(rate.Limit(100/60.0), 5)
	router.Use(middleware.RateLimitMiddleware(limiter))

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Start server
	utils.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		utils.Fatal("Failed to start server: " + err.Error())
	}
}
