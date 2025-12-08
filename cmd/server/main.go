package main

import (
	"log/slog"
	"os"

	"go-test-api/internal/config"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func setupLogger(cfg config.Config) {
	var handler slog.Handler

	if cfg.IsProduction() {
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
	cfg := config.Load()

	setupLogger(cfg)

	// Create server
	srv, err := server.New(server.Config{
		Port:      cfg.Port,
		Database:  cfg.Database,
		JWTSecret: cfg.JWTSecret,
		JWTExpiry: cfg.JWTExpiry,
	})
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Run server
	srv.Run()
}
