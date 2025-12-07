//go:build unit

package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-test-api/internal/validator"
)

// mockUserRepository is a mock implementation of UserRepository for testing
type mockUserRepository struct {
	upsertFunc      func(ctx context.Context, req *CreateUserRequest, passwordHash string) (*UserResponse, error)
	listFunc        func(ctx context.Context) ([]*UserResponse, error)
	listByEmailFunc func(ctx context.Context, email string) ([]*UserResponse, error)
}

func (m *mockUserRepository) Upsert(ctx context.Context, req *CreateUserRequest, passwordHash string) (*UserResponse, error) {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, req, passwordHash)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) List(ctx context.Context) ([]*UserResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) ListByEmail(ctx context.Context, email string) ([]*UserResponse, error) {
	if m.listByEmailFunc != nil {
		return m.listByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func TestUserHandler_List(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		mockList        func(ctx context.Context) ([]*UserResponse, error)
		mockListByEmail func(ctx context.Context, email string) ([]*UserResponse, error)
		expectedStatus  int
		expectedCount   int
		expectError     bool
	}{
		{
			name: "successful list with users",
			mockList: func(ctx context.Context) ([]*UserResponse, error) {
				return []*UserResponse{
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
			mockList: func(ctx context.Context) ([]*UserResponse, error) {
				return []*UserResponse{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			expectError:    false,
		},
		{
			name: "database error",
			mockList: func(ctx context.Context) ([]*UserResponse, error) {
				return nil, errors.New("database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:  "filter by email",
			query: "john",
			mockListByEmail: func(ctx context.Context, email string) ([]*UserResponse, error) {
				if email != "john" {
					return nil, errors.New("unexpected email")
				}
				return []*UserResponse{{ID: "1", Name: "John Doe", Email: "john@example.com"}}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup
			mockRepo := &mockUserRepository{
				listFunc:        tt.mockList,
				listByEmailFunc: tt.mockListByEmail,
			}
			handler := NewHandler(validator.New(), mockRepo)

			// Create request (include query if provided)
			url := "/users"
			if tt.query != "" {
				url = url + "?email=" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
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
				var users []*UserResponse
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
