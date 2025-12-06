package main

import (
	"log/slog"
	"os"
	"strconv"

	"go-test-api/internal/database"
	"go-test-api/internal/logging"
	"go-test-api/internal/server"
)

// @title Go Test API
// @version 1.0
// @description A RESTful API for managing users and addresses
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(valueStr); err == nil {
		return intValue
	}
	return defaultValue
}

func setupLogger() {
	var handler slog.Handler
	env := getEnv("ENV", "development")

	if env == "production" {
		// Compact JSON logging for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Pretty-printed JSON for development
		handler = logging.NewPrettyJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	setupLogger()

	// Load configuration from environment variables with defaults
	cfg := server.Config{
		Port: getEnv("PORT", "8080"),
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "gouser"),
			Password: getEnv("DB_PASSWORD", "gopassword"),
			DBName:   getEnv("DB_NAME", "gotestdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	// Create server
	srv, err := server.New(cfg)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Run server
	srv.Run()
}
