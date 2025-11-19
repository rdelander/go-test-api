package main

import (
	"log"
	"os"
	"strconv"

	"go-test-api/internal/database"
	"go-test-api/internal/server"
)

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

func main() {
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
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run server
	srv.Run()
}
