package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type HelloRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Server holds all dependencies
type Server struct {
	validator *validator.Validate
}

// NewServer creates a new server with all dependencies injected
func NewServer(v *validator.Validate) *Server {
	return &Server{
		validator: v,
	}
}

// Helper methods on Server for dependency injection
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, ErrorResponse{Error: message})
}

func (s *Server) decodeAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON")
	}
	if err := s.validator.Struct(v); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		firstError := validationErrors[0]
		errorMsg := fmt.Sprintf("Field '%s' failed validation '%s'",
			firstError.Field(),
			firstError.Tag())
		if firstError.Param() != "" {
			errorMsg += fmt.Sprintf(" (expected: %s)", firstError.Param())
		}
		return fmt.Errorf("%s", errorMsg)
	}
	return nil
}

// Handler methods
func (s *Server) getHelloWorld(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Message: "Hello, World!",
	}
	s.sendJSON(w, http.StatusOK, response)
}

func (s *Server) postHelloWorld(w http.ResponseWriter, r *http.Request) {
	var req HelloRequest
	if err := s.decodeAndValidate(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	response := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", strings.TrimSpace(req.Name)),
	}
	s.sendJSON(w, http.StatusOK, response)
}

func main() {
	// Initialize dependencies
	validator := validator.New()

	// Create server with dependencies
	server := NewServer(validator)

	// Register routes using server methods
	http.HandleFunc("GET /hello_world", server.getHelloWorld)
	http.HandleFunc("POST /hello_world", server.postHelloWorld)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
