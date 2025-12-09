//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

var baseURL string

func init() {
	baseURL = os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
}

// waitForAPI waits for the API to be ready
func waitForAPI(t *testing.T, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Log("API is ready")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatal("API did not become ready in time")
}

func TestHealthEndpoint(t *testing.T) {
	waitForAPI(t, 30*time.Second)

	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", result["status"])
	}
}

func TestRegisterUser(t *testing.T) {
	waitForAPI(t, 30*time.Second)

	// Generate unique email for this test run
	timestamp := time.Now().Unix()
	email := fmt.Sprintf("e2e-user-%d@test.com", timestamp)

	// Prepare request
	reqBody := map[string]string{
		"name":     "E2E Test User",
		"email":    email,
		"password": "securepassword123",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Make request
	resp, err := http.Post(
		baseURL+"/auth/register",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatalf("Failed to call register endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Verify token fields
	if _, ok := result["token"]; !ok {
		t.Error("Response missing 'token' field")
	}

	if _, ok := result["expires_at"]; !ok {
		t.Error("Response missing 'expires_at' field")
	}

	// Verify user object
	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'user' object")
	}

	if user["email"] != email {
		t.Errorf("Expected email '%s', got '%v'", email, user["email"])
	}

	if user["name"] != "E2E Test User" {
		t.Errorf("Expected name 'E2E Test User', got '%v'", user["name"])
	}

	if _, ok := user["id"]; !ok {
		t.Error("User object missing 'id' field")
	}

	t.Logf("Successfully registered user: %s", email)
}
