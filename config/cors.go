package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
)

func GetCorsConfig() cors.Config {
	// Get allowed origins from env
	allowedOrigins := strings.Split(getEnvOrDefault("CORS_ALLOWED_ORIGINS", "*"), ",")

	// Get allowed methods from env
	allowedMethods := strings.Split(getEnvOrDefault("CORS_ALLOWED_METHODS",
		"GET,POST,PUT,PATCH,DELETE,OPTIONS"), ",")

	// Get allowed headers from env
	allowedHeaders := strings.Split(getEnvOrDefault("CORS_ALLOWED_HEADERS",
		"Origin,Content-Type,Accept,Authorization"), ",")

	// Get max age from env
	maxAge := getEnvOrDefaultDuration("CORS_MAX_AGE", 12*time.Hour)

	return cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     allowedMethods,
		AllowHeaders:     allowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value + "s"); err == nil {
			return duration
		}
	}
	return defaultValue
}
