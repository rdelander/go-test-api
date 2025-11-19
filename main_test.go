package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func setupTestServer() *Server {
	v := validator.New()
	return NewServer(v)
}

func TestGetHelloWorld(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/hello_world", nil)
	w := httptest.NewRecorder()

	server.getHelloWorld(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response HelloResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMessage := "Hello, World!"
	if response.Message != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, response.Message)
	}
}

func TestPostHelloWorld_Success(t *testing.T) {
	server := setupTestServer()
	requestBody := HelloRequest{Name: "Alice"}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.postHelloWorld(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response HelloResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMessage := "Hello, Alice!"
	if response.Message != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, response.Message)
	}
}

func TestPostHelloWorld_EmptyName(t *testing.T) {
	server := setupTestServer()
	requestBody := HelloRequest{Name: ""}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.postHelloWorld(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var response ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestPostHelloWorld_NameTooLong(t *testing.T) {
	server := setupTestServer()
	longName := strings.Repeat("a", 101)
	requestBody := HelloRequest{Name: longName}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.postHelloWorld(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestPostHelloWorld_InvalidJSON(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.postHelloWorld(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var response ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON' error, got '%s'", response.Error)
	}
}
