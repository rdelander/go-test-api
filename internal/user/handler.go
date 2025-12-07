package user

import (
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

// List handles GET /users with optional email filter
// @Summary List all users
// @Description Get list of all users, optionally filtered by email
// @Tags users
// @Produce json
// @Param email query string false "Filter by email"
// @Success 200 {array} UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users [get]
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
