package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type HelloRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper function to send JSON responses
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error responses
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

// Helper function to decode and validate JSON request
func decodeAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON")
	}
	if err := validate.Struct(v); err != nil {
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

func getHelloWorld(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Message: "Hello, World!",
	}
	sendJSON(w, http.StatusOK, response)
}

func postHelloWorld(w http.ResponseWriter, r *http.Request) {
	var req HelloRequest
	if err := decodeAndValidate(r, &req); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	response := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", strings.TrimSpace(req.Name)),
	}
	sendJSON(w, http.StatusOK, response)
}

func main() {
	// Initialize validator
	validate = validator.New()

	http.HandleFunc("GET /hello_world", getHelloWorld)
	http.HandleFunc("POST /hello_world", postHelloWorld)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
