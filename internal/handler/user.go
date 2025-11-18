package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-test-api/internal/model"
	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// UserHandler handles user-related requests
type UserHandler struct {
	validator *validator.Validator
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(v *validator.Validator) *UserHandler {
	return &UserHandler{
		validator: v,
	}
}

// Create handles POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest

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

	// TODO: Save to database (for now, just mock it)
	// In a real app, you'd inject a database/repository here
	userID := "usr_" + fmt.Sprintf("%d", len(req.Name)) // Mock ID generation

	// Build response
	resp := model.UserResponse{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}
	response.JSON(w, http.StatusCreated, resp)
}
