package main

import (
	"log"

	"go-test-api/internal/database"
	"go-test-api/internal/server"
)

func main() {
	// Initialize server configuration
	cfg := server.Config{
		Port: "8080",
		Database: database.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "gouser",
			Password: "gopassword",
			DBName:   "gotestdb",
			SSLMode:  "disable",
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
