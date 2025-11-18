package main

import (
	"go-test-api/internal/server"
)

func main() {
	// Initialize server configuration
	cfg := server.Config{
		Port: "8080",
	}

	// Create and run server
	srv := server.New(cfg)
	srv.Run()
}
