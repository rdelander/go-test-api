package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// Handler handles user-related HTTP requests
type Handler struct {
	validator *validator.Validator
	repo      Repo
}

// NewHandler creates a new user Handler
func NewHandler(v *validator.Validator, repo Repo) *Handler {
	return &Handler{validator: v, repo: repo}
}

// Create handles POST /users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	defer r.Body.Close()
	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.repo.Upsert(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// List handles GET /users with optional email filter
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	var users []*UserResponse
	var err error
	if email != "" {
		users, err = h.repo.ListByEmail(r.Context(), email)
	} else {
		users, err = h.repo.List(r.Context())
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	response.JSON(w, http.StatusOK, users)
}
