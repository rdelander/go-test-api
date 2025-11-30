//go:build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-test-api/internal/model"
	"go-test-api/internal/validator"
)

// mockUserRepository is a mock implementation of UserRepository for testing
type mockUserRepository struct {
	upsertFunc func(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	listFunc   func(ctx context.Context) ([]*model.UserResponse, error)
}

func (m *mockUserRepository) Upsert(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) List(ctx context.Context) ([]*model.UserResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockUpsert     func(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "valid user creation",
			body: `{"name":"John Doe","email":"john@example.com"}`,
			mockUpsert: func(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
				return &model.UserResponse{
					ID:    "1",
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    "1",
				"name":  "John Doe",
				"email": "john@example.com",
			},
		},
		{
			name:           "invalid JSON",
			body:           `{invalid json}`,
			mockUpsert:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid JSON",
			},
		},
		{
			name:           "missing required name",
			body:           `{"email":"john@example.com"}`,
			mockUpsert:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Field 'Name' failed validation 'required'",
			},
		},
		{
			name:           "invalid email format",
			body:           `{"name":"John Doe","email":"invalid-email"}`,
			mockUpsert:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Field 'Email' failed validation 'email'",
			},
		},
		{
			name: "database error",
			body: `{"name":"John Doe","email":"john@example.com"}`,
			mockUpsert: func(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
				return nil, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "failed to create user: database connection failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{
				upsertFunc: tt.mockUpsert,
			}
			handler := NewUserHandler(validator.New(), mockRepo)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handler.Create(w, req)

			// Assert status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Assert response body
			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			for key, expectedValue := range tt.expectedBody {
				if response[key] != expectedValue {
					t.Errorf("expected %s to be %v, got %v", key, expectedValue, response[key])
				}
			}
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockList       func(ctx context.Context) ([]*model.UserResponse, error)
		expectedStatus int
		expectedCount  int
		expectError    bool
	}{
		{
			name: "successful list with users",
			mockList: func(ctx context.Context) ([]*model.UserResponse, error) {
				return []*model.UserResponse{
					{ID: "1", Name: "John Doe", Email: "john@example.com"},
					{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectError:    false,
		},
		{
			name: "successful list with no users",
			mockList: func(ctx context.Context) ([]*model.UserResponse, error) {
				return []*model.UserResponse{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			expectError:    false,
		},
		{
			name: "database error",
			mockList: func(ctx context.Context) ([]*model.UserResponse, error) {
				return nil, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{
				listFunc: tt.mockList,
			}
			handler := NewUserHandler(validator.New(), mockRepo)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()

			// Execute
			handler.List(w, req)

			// Assert status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Assert response
			if tt.expectError {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if _, ok := response["error"]; !ok {
					t.Error("expected error field in response")
				}
			} else {
				var users []*model.UserResponse
				if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(users) != tt.expectedCount {
					t.Errorf("expected %d users, got %d", tt.expectedCount, len(users))
				}
			}
		})
	}
}
