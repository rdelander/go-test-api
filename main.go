package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
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
