package main

import (
	"loan/config"
	"loan/routes"
	"loan/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Start server
	utils.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		utils.Fatal("Failed to start server: " + err.Error())
	}
}
