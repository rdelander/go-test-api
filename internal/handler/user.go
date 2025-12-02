package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-test-api/internal/user"
	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// UserHandler handles user-related requests
type UserHandler struct {
	validator *validator.Validator
	repo      user.Repo
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(v *validator.Validator, repo user.Repo) *UserHandler {
	return &UserHandler{
		validator: v,
		repo:      repo,
	}
}

// Create handles POST /users (idempotent - creates or updates by email)
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest

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

	// Upsert user in database (idempotent operation)
	user, err := h.repo.Upsert(r.Context(), &req)
	if err != nil {
		// Log full error for debugging (in production, use proper logging)
		// log.Printf("Failed to upsert user: %v", err)
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	// Build response - always return 200 OK for idempotent operations
	// (could be 201 for new, 200 for updated, but consistent 200 is simpler)
	response.JSON(w, http.StatusOK, user)
}

// List handles GET /users with optional email filter
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check for email query parameter
	email := r.URL.Query().Get("email")

	var users []*user.UserResponse
	var err error

	if email != "" {
		// Filter by email (repository will add wildcards)
		users, err = h.repo.ListByEmail(r.Context(), email)
	} else {
		// List all users
		users, err = h.repo.List(r.Context())
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	response.JSON(w, http.StatusOK, users)
}
