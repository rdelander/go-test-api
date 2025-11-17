package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func validateHelloRequest(req *HelloRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required and cannot be empty")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name must be 100 characters or less")
	}
	return nil
}

func getHelloWorld(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Message: "Hello, World!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func postHelloWorld(w http.ResponseWriter, r *http.Request) {
	var req HelloRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}
	defer r.Body.Close()

	// Validate the request
	if err := validateHelloRequest(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", strings.TrimSpace(req.Name)),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("GET /hello_world", getHelloWorld)
	http.HandleFunc("POST /hello_world", postHelloWorld)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
