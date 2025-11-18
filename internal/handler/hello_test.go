package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-test-api/internal/model"
	"go-test-api/internal/validator"
)

func setupTestHandler() *HelloHandler {
	v := validator.New()
	return NewHelloHandler(v)
}

func TestHelloHandler_Get(t *testing.T) {
	handler := setupTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/hello_world", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response model.HelloResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMessage := "Hello, World!"
	if response.Message != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, response.Message)
	}
}

func TestHelloHandler_Post_Success(t *testing.T) {
	handler := setupTestHandler()
	requestBody := model.HelloRequest{Name: "Alice"}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Post(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response model.HelloResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMessage := "Hello, Alice!"
	if response.Message != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, response.Message)
	}
}

func TestHelloHandler_Post_EmptyName(t *testing.T) {
	handler := setupTestHandler()
	requestBody := model.HelloRequest{Name: ""}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Post(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestHelloHandler_Post_NameTooLong(t *testing.T) {
	handler := setupTestHandler()
	longName := strings.Repeat("a", 101)
	requestBody := model.HelloRequest{Name: longName}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Post(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestHelloHandler_Post_InvalidJSON(t *testing.T) {
	handler := setupTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/hello_world", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Post(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
