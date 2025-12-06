package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "go-test-api/docs"
	"go-test-api/internal/address"
	addressdb "go-test-api/internal/address/db"
	"go-test-api/internal/database"
	"go-test-api/internal/middleware"
	"go-test-api/internal/user"
	userdb "go-test-api/internal/user/db"
	"go-test-api/internal/validator"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	port           string
	pool           *pgxpool.Pool
	userHandler    *user.Handler
	addressHandler *address.Handler
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

	return &Server{
		port: cfg.Port,
		pool: pool,
		userHandler: user.NewHandler(
			validator.New(),
			user.NewRepository(
				userdb.New(pool),
			),
		),
		addressHandler: address.NewHandler(
			validator.New(),
			address.NewRepository(
				addressdb.New(pool),
				userdb.New(pool),
			),
		),
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
	// Swagger UI
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	routes := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{"GET", "/users", s.userHandler.List},
		{"POST", "/users", s.userHandler.Create},
		{"GET", "/addresses", s.addressHandler.List},
		{"POST", "/addresses", s.addressHandler.Create},
		{"GET", "/addresses/{id}", s.addressHandler.Get},
		{"PUT", "/addresses/{id}", s.addressHandler.Update},
		{"DELETE", "/addresses/{id}", s.addressHandler.Delete},
	}

	routeList := []string{"GET /swagger/"}
	for _, route := range routes {
		http.HandleFunc(route.method+" "+route.path, route.handler)
		routeList = append(routeList, route.method+" "+route.path)
	}

	log.Printf("Registered routes: %v", routeList)
}

// start starts the HTTP server
func (s *Server) start() error {
	s.setupRoutes()

	// Wrap default mux with logging middleware
	handler := middleware.Logging(http.DefaultServeMux)

	log.Printf("Server starting on port %s...", s.port)
	if err := http.ListenAndServe(":"+s.port, handler); err != nil {
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
