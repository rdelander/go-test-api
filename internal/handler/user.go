package handler

import (
	"encoding/json"
	"net/http"

	"go-test-api/internal/model"
	"go-test-api/internal/repository"
	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// UserHandler handles user-related requests
type UserHandler struct {
	validator *validator.Validator
	repo      *repository.UserRepository
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(v *validator.Validator, repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		validator: v,
		repo:      repo,
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

	// Create user in database
	user, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Build response
	response.JSON(w, http.StatusCreated, user)
}

// List handles GET /users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	response.JSON(w, http.StatusOK, users)
}
