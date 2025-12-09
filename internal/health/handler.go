package health

import (
	"net/http"

	"go-test-api/pkg/response"
)

// Handler handles health check requests
type Handler struct{}

// NewHandler creates a new health Handler
func NewHandler() *Handler {
	return &Handler{}
}

// Check handles GET /health
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
