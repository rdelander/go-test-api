package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error sends an error response with the given status code and message
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}
