package address

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-test-api/internal/validator"
	"go-test-api/pkg/response"
)

// Handler handles HTTP requests for addresses
type Handler struct {
	validator *validator.Validator
	repo      Repo
}

// NewHandler creates a new address Handler
func NewHandler(v *validator.Validator, repo Repo) *Handler {
	return &Handler{validator: v, repo: repo}
}

// Create handles POST /addresses
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	addr, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create address: %v", err))
		return
	}

	response.JSON(w, http.StatusCreated, addr)
}

// Get handles GET /addresses/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	addr, err := h.repo.Get(r.Context(), int32(id))
	if err != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Address not found: %v", err))
		return
	}

	response.JSON(w, http.StatusOK, addr)
}

// List handles GET /addresses?entity_type=user&entity_id=1&address_type=shipping
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("entity_type")
	entityID := r.URL.Query().Get("entity_id")
	addressType := r.URL.Query().Get("address_type")

	if entityType == "" || entityID == "" {
		response.Error(w, http.StatusBadRequest, "entity_type and entity_id are required")
		return
	}

	var addrs []*AddressResponse
	var err error

	if addressType != "" {
		addrs, err = h.repo.ListByEntityAndType(r.Context(), entityType, entityID, addressType)
	} else {
		addrs, err = h.repo.ListByEntity(r.Context(), entityType, entityID)
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list addresses: %v", err))
		return
	}

	response.JSON(w, http.StatusOK, addrs)
}

// Update handles PUT /addresses/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var req UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	addr, err := h.repo.Update(r.Context(), int32(id), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update address: %v", err))
		return
	}

	response.JSON(w, http.StatusOK, addr)
}

// Delete handles DELETE /addresses/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	if err := h.repo.Delete(r.Context(), int32(id)); err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete address: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
