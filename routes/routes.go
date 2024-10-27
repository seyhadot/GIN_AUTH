package routes

import (
	"loan/config"
	"loan/controllers"
	"loan/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, config *config.Config) {
	// Initialize controllers
	authController := controllers.NewAuthController(config)

	// Public routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			protected.GET("/profile", authController.GetCurrentProfile)
			protected.PUT("/profile", authController.UpdateProfile)
			protected.PUT("/password", authController.UpdatePassword) // New endpoint
			protected.GET("/users", authController.GetAllUsers)
		}
	}
}
