package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-test-api/internal/model"
	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// HelloHandler handles hello world requests
type HelloHandler struct {
	validator *validator.Validator
}

// NewHelloHandler creates a new HelloHandler
func NewHelloHandler(v *validator.Validator) *HelloHandler {
	return &HelloHandler{
		validator: v,
	}
}

// Get handles GET /hello_world
func (h *HelloHandler) Get(w http.ResponseWriter, r *http.Request) {
	resp := model.HelloResponse{
		Message: "Hello, World!",
	}
	response.JSON(w, http.StatusOK, resp)
}

// Post handles POST /hello_world
func (h *HelloHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req model.HelloRequest

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	defer r.Body.Close()

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Build response
	resp := model.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", strings.TrimSpace(req.Name)),
	}
	response.JSON(w, http.StatusOK, resp)
}
