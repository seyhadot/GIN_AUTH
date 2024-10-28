package routes

import (
	"loan/config"
	"loan/controllers"
	"loan/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, config *config.Config) {
	// Add request logger middleware
	router.Use(middleware.RequestLogger())

	// Initialize controllers
	authController := controllers.NewAuthController(config)
	companyController := controllers.NewCompanyController(config)
	branchOfficeController := controllers.NewBranchOfficeController(config)
	staffController := controllers.NewStaffController(config)

	// API routes group
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/staff/register", staffController.RegisterStaff)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", authController.GetCurrentProfile)
				users.PUT("/profile", authController.UpdateProfile)
				users.PUT("/password", authController.UpdatePassword)
				users.GET("", authController.GetAllUsers)
			}

			// Company routes
			companies := protected.Group("/companies")
			{
				// Company CRUD operations
				companies.POST("", companyController.CreateCompany)
				companies.GET("", companyController.ListCompanies)
				companies.GET("/detail/:id", companyController.GetCompany)
				companies.PUT("/detail/:id", companyController.UpdateCompany)
				companies.DELETE("/detail/:id", companyController.DeleteCompany)

				// Branch office routes
				companies.POST("/:id/branches", branchOfficeController.CreateBranchOffice)
				companies.GET("/:id/branches", branchOfficeController.ListBranchOffices)
				companies.GET("/:id/branches/:branch_id", branchOfficeController.GetBranchOffice)
				companies.PUT("/:id/branches/:branch_id", branchOfficeController.UpdateBranchOffice)
				companies.DELETE("/:id/branches/:branch_id", branchOfficeController.DeleteBranchOffice)

				// Staff management routes
				companies.POST("/:id/branches/:branch_id/staff", staffController.AssignStaffToBranch)
				companies.GET("/:id/branches/:branch_id/staff", staffController.ListStaffByBranch)
				companies.DELETE("/:id/branches/:branch_id/staff/:user_id", staffController.RemoveStaffFromBranch)
			}
		}
	}
}
