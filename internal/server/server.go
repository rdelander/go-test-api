package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go-test-api/internal/database"
	"go-test-api/internal/user"
	"go-test-api/internal/user/db"
	"go-test-api/internal/validator"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	port        string
	pool        *pgxpool.Pool
	userHandler *user.Handler
}

// Config holds server configuration
type Config struct {
	Port     string
	Database database.Config
}

// New creates a new Server instance with all dependencies injected
func New(cfg Config) (*Server, error) {
	// Initialize database connection
	ctx := context.Background()
	pool, err := database.New(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize dependencies
	v := validator.New()
	queries := db.New(pool)
	userRepo := user.NewRepository(queries)

	// Initialize handlers
	userHandler := user.NewHandler(v, userRepo)

	return &Server{
		port:        cfg.Port,
		pool:        pool,
		userHandler: userHandler,
	}, nil
}

// Close closes the database connection
func (s *Server) Close() error {
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}

// setupRoutes registers all HTTP routes
func (s *Server) setupRoutes() {
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
