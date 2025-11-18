package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go-test-api/internal/database"
	"go-test-api/internal/handler"
	"go-test-api/internal/repository"
	"go-test-api/internal/validator"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	port         string
	db           *sql.DB
	helloHandler *handler.HelloHandler
	userHandler  *handler.UserHandler
}

// Config holds server configuration
type Config struct {
	Port     string
	Database database.Config
}

// New creates a new Server instance with all dependencies injected
func New(cfg Config) (*Server, error) {
	// Initialize database connection
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize dependencies
	v := validator.New()
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	helloHandler := handler.NewHelloHandler(v)
	userHandler := handler.NewUserHandler(v, userRepo)

	return &Server{
		port:         cfg.Port,
		db:           db,
		helloHandler: helloHandler,
		userHandler:  userHandler,
	}, nil
}

// Close closes the database connection
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// setupRoutes registers all HTTP routes
func (s *Server) setupRoutes() {
	// Hello endpoints
	http.HandleFunc("GET /hello_world", s.helloHandler.Get)
	http.HandleFunc("POST /hello_world", s.helloHandler.Post)

	// User endpoints
	http.HandleFunc("GET /users", s.userHandler.List)
	http.HandleFunc("POST /users", s.userHandler.Create)
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
	defer s.Close()
	if err := s.start(); err != nil {
		log.Fatal(err)
	}
}
