//go:build integration
// +build integration

package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go-test-api/internal/database"
	"go-test-api/internal/db"
	"go-test-api/internal/validator"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool
var testQueries *db.Queries

func TestMain(m *testing.M) {
	// Setup
	var err error
	ctx := context.Background()

	// Connect to test database
	cfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "gouser"),
		Password: getEnv("DB_PASSWORD", "gopassword"),
		DBName:   getEnv("DB_NAME", "gotestdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	testDB, err = database.New(ctx, cfg)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	testQueries = db.New(testDB)

	// Run tests
	code := m.Run()

	// Teardown
	testDB.Close()

	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var intValue int
	if _, err := fmt.Sscanf(valueStr, "%d", &intValue); err == nil {
		return intValue
	}
	return defaultValue
}

// cleanupUsers removes all users from the database
func cleanupUsers(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}

func TestUserHandler_Create_Integration(t *testing.T) {
	// Setup
	cleanupUsers(t)
	repo := NewRepository(testQueries)
	handler := NewHandler(validator.New(), repo)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		checkFunc      func(t *testing.T, body []byte)
	}{
		{
			name:           "create new user",
			body:           `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["name"] != "John Doe" {
					t.Errorf("Expected name 'John Doe', got '%v'", response["name"])
				}
				if response["email"] != "john@example.com" {
					t.Errorf("Expected email 'john@example.com', got '%v'", response["email"])
				}
				if response["id"] == nil || response["id"] == "" {
					t.Error("Expected ID to be set")
				}
			},
		},
		{
			name:           "upsert existing user by email",
			body:           `{"name":"John Updated","email":"john@example.com"}`,
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["name"] != "John Updated" {
					t.Errorf("Expected name 'John Updated', got '%v'", response["name"])
				}
				if response["email"] != "john@example.com" {
					t.Errorf("Expected email 'john@example.com', got '%v'", response["email"])
				}
				// ID should be "1" since it's an update of the first user
				if response["id"] != "1" {
					t.Errorf("Expected ID '1', got '%v'", response["id"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, w.Body.Bytes())
			}
		})
	}
}

func TestUserHandler_List_Integration(t *testing.T) {
	// Setup
	cleanupUsers(t)
	repo := NewRepository(testQueries)
	handler := NewHandler(validator.New(), repo)

	// Create some test users
	users := []struct {
		name  string
		email string
	}{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}

	for _, u := range users {
		body := fmt.Sprintf(`{"name":"%s","email":"%s"}`, u.name, u.email)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Create(w, req)
	}

	// Test list
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response))
	}

	// Verify users are in order by ID
	expectedNames := []string{"Alice", "Bob", "Charlie"}
	for i, user := range response {
		if user["name"] != expectedNames[i] {
			t.Errorf("Expected user %d to be '%s', got '%v'", i, expectedNames[i], user["name"])
		}
	}
}

func TestUserHandler_ListByEmail_Integration(t *testing.T) {
	// Setup
	cleanupUsers(t)
	repo := NewRepository(testQueries)
	handler := NewHandler(validator.New(), repo)

	// Create some test users (mixed-case emails to verify case-insensitivity)
	users := []struct {
		name  string
		email string
	}{
		{"Alice", "Alice@Example.com"},
		{"Bob", "BOB@example.COM"},
		{"John", "John.Doe@Example.com"},
		{"Johnny", "johnny@EXAMPLE.com"},
	}

	for _, u := range users {
		body := fmt.Sprintf(`{"name":"%s","email":"%s"}`, u.name, u.email)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Create(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("failed to create user %s: status %d body=%s", u.email, w.Code, w.Body.String())
		}
	}

	// Filter by 'JoHn' with mixed case should still match john.doe and johnny
	req := httptest.NewRequest(http.MethodGet, "/users?email=JoHn", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("Expected 2 users matching 'john', got %d", len(response))
	}

	// Collect emails returned (normalize to lowercase for case-insensitive comparison)
	emails := map[string]bool{}
	for _, u := range response {
		if e, ok := u["email"].(string); ok {
			emails[strings.ToLower(e)] = true
		}
	}

	if !emails["john.doe@example.com"] || !emails["johnny@example.com"] {
		t.Fatalf("Filtered results missing expected emails: %v", emails)
	}
}
