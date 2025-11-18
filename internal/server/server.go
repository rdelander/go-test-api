package server

import (
	"fmt"
	"log"
	"net/http"

	"go-test-api/internal/handler"
	"go-test-api/internal/validator"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	port         string
	helloHandler *handler.HelloHandler
}

// Config holds server configuration
type Config struct {
	Port string
}

// New creates a new Server instance with all dependencies injected
func New(cfg Config) *Server {
	// Initialize dependencies
	v := validator.New()

	// Initialize handlers
	helloHandler := handler.NewHelloHandler(v)

	return &Server{
		port:         cfg.Port,
		helloHandler: helloHandler,
	}
}

// setupRoutes registers all HTTP routes
func (s *Server) setupRoutes() {
	http.HandleFunc("GET /hello_world", s.helloHandler.Get)
	http.HandleFunc("POST /hello_world", s.helloHandler.Post)
}

// start starts the HTTP server
func (s *Server) start() error {
	s.setupRoutes()

	fmt.Printf("Server starting on port %s...\n", s.port)
	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

// Run starts the server and handles any errors
func (s *Server) Run() {
	if err := s.start(); err != nil {
		log.Fatal(err)
	}
}
