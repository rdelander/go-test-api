package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "go-test-api/docs"
	"go-test-api/internal/address"
	addressdb "go-test-api/internal/address/db"
	"go-test-api/internal/auth"
	"go-test-api/internal/database"
	"go-test-api/internal/health"
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
	authHandler    *auth.Handler
	authService    *auth.Service
	healthHandler  *health.Handler
}

// Config holds server configuration
type Config struct {
	Port      string
	Database  database.Config
	JWTSecret string
	JWTExpiry time.Duration
}

// New creates a new Server instance with all dependencies injected
func New(cfg Config) (*Server, error) {
	// Initialize database connection
	ctx := context.Background()
	pool, err := database.New(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize auth service
	authService := auth.NewService(cfg.JWTSecret, cfg.JWTExpiry)
	userQueries := userdb.New(pool)
	userRepo := user.NewRepository(userQueries)

	return &Server{
		port:          cfg.Port,
		pool:          pool,
		authService:   authService,
		healthHandler: health.NewHandler(),
		userHandler: user.NewHandler(
			validator.New(),
			userRepo,
		),
		addressHandler: address.NewHandler(
			validator.New(),
			address.NewRepository(
				addressdb.New(pool),
				userQueries,
			),
		),
		authHandler: auth.NewHandler(
			validator.New(),
			authService,
			userRepo,
			userQueries,
		),
	}, nil
}

// Close closes the database connection
func (s *Server) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// setupRoutes registers all HTTP routes
func (s *Server) setupRoutes() {
	// Health check (public)
	http.HandleFunc("GET /health", s.healthHandler.Check)

	// Swagger UI (public)
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Auth routes (public)
	http.HandleFunc("POST /auth/register", s.authHandler.Register)
	http.HandleFunc("POST /auth/login", s.authHandler.Login)

	// Protected routes
	authMiddleware := auth.Middleware(s.authService)
	protectedRoutes := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{"GET", "/users", s.userHandler.List},
		{"GET", "/addresses", s.addressHandler.List},
		{"POST", "/addresses", s.addressHandler.Create},
		{"GET", "/addresses/{id}", s.addressHandler.Get},
		{"PUT", "/addresses/{id}", s.addressHandler.Update},
		{"DELETE", "/addresses/{id}", s.addressHandler.Delete},
	}

	routeList := []string{"GET /health", "GET /swagger/", "POST /auth/register", "POST /auth/login"}
	for _, route := range protectedRoutes {
		http.Handle(route.method+" "+route.path, authMiddleware(http.HandlerFunc(route.handler)))
		routeList = append(routeList, route.method+" "+route.path+" (protected)")
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
