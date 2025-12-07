package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go-test-api/internal/user"
	userdb "go-test-api/internal/user/db"
	"go-test-api/internal/validator"
	"go-test-api/pkg/response"

	"github.com/jackc/pgx/v5"
)

// Handler handles authentication HTTP requests
type Handler struct {
	validator   *validator.Validator
	authService *Service
	userRepo    user.Repo
	userQueries *userdb.Queries
}

// NewHandler creates a new auth Handler
func NewHandler(v *validator.Validator, authService *Service, userRepo user.Repo, userQueries *userdb.Queries) *Handler {
	return &Handler{
		validator:   v,
		authService: authService,
		userRepo:    userRepo,
		userQueries: userQueries,
	}
}

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Registration details"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	userReq := &user.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	dbUser, err := h.userRepo.Upsert(r.Context(), userReq, hashedPassword)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(dbUser.ID, dbUser.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	resp := TokenResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(h.authService.jwtExpiry).Unix(),
		User: User{
			ID:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		},
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Login handles POST /auth/login
// @Summary Login user
// @Description Authenticate user with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user by email
	dbUser, err := h.userQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	// Verify password
	if err := h.authService.CheckPassword(dbUser.PasswordHash, req.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(fmt.Sprint(dbUser.ID), dbUser.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	resp := TokenResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(h.authService.jwtExpiry).Unix(),
		User: User{
			ID:    fmt.Sprint(dbUser.ID),
			Name:  dbUser.Name,
			Email: dbUser.Email,
		},
	}

	response.JSON(w, http.StatusOK, resp)
}
