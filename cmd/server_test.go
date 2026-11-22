package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := newRouter()
	if r == nil {
		t.Fatal("newRouter returned nil")
	}
}

func TestHealthz(t *testing.T) {
	r := newRouter()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "healthz returns 200",
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "readyz returns 200",
			path:       "/readyz",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GET %s: got status %d, want %d", tt.path, rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestHelloEndpoint(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["message"] != "Hello, World!" {
		t.Errorf("expected message 'Hello, World!', got '%s'", body["message"])
	}
}

func TestPingEndpoint(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["message"] != "pong" {
		t.Errorf("expected message 'pong', got '%s'", body["message"])
	}
}

func TestHeyEndpoint(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/hey", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["message"] != "joe" {
		t.Errorf("expected message 'joe', got '%s'", body["message"])
	}
}
