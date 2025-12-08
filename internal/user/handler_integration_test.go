//go:build integration
// +build integration

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go-test-api/internal/config"
	"go-test-api/internal/database"
	"go-test-api/internal/user/db"
	"go-test-api/internal/validator"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool
var testQueries *db.Queries

func TestMain(m *testing.M) {
	// Setup
	var err error
	ctx := context.Background()

	// Load config (will use development defaults)
	cfg := config.Load()

	testDB, err = database.New(ctx, cfg.Database)
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

// cleanupUsers removes all users from the database
func cleanupUsers(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}

// setupTestUsers creates users via the repository (simulating auth/register)
func setupTestUsers(t *testing.T, repo *Repository, users []struct {
	name  string
	email string
}) {
	t.Helper()
	for _, u := range users {
		userReq := &CreateUserRequest{
			Name:  u.name,
			Email: u.email,
		}
		_, err := repo.Upsert(context.Background(), userReq, "hashedpassword")
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", u.email, err)
		}
	}
}

// setupHandler creates a clean repository and handler for testing
func setupHandler(t *testing.T) (*Repository, *Handler) {
	t.Helper()
	cleanupUsers(t)
	repo := NewRepository(testQueries)
	handler := NewHandler(validator.New(), repo)
	return repo, handler
}

// executeListRequest makes a GET request to /users and returns the parsed response
func executeListRequest(t *testing.T, handler *Handler, queryParams string) []map[string]interface{} {
	t.Helper()
	url := "/users"
	if queryParams != "" {
		url += "?" + queryParams
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	return response
}

func TestUserHandler_List_Integration(t *testing.T) {
	repo, handler := setupHandler(t)

	users := []struct {
		name  string
		email string
	}{
		{"Alice", "alice@example.com"},
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}
	setupTestUsers(t, repo, users)

	response := executeListRequest(t, handler, "")

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
	repo, handler := setupHandler(t)

	users := []struct {
		name  string
		email string
	}{
		{"Alice", "Alice@Example.com"},
		{"Bob", "BOB@example.COM"},
		{"John", "John.Doe@Example.com"},
		{"Johnny", "johnny@EXAMPLE.com"},
	}
	setupTestUsers(t, repo, users)

	response := executeListRequest(t, handler, "email=JoHn")

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
