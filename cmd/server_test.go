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
		name           string
		path           string
		wantStatus     int
		wantKey        string
		wantValue      string
	}{
		{
			name:      "healthz returns 200",
			path:      "/healthz",
			wantStatus: http.StatusOK,
			wantKey:   "status",
			wantValue: "ok",
		},
		{
			name:      "readyz returns 200",
			path:      "/readyz",
			wantStatus: http.StatusOK,
			wantKey:   "status",
			wantValue: "ready",
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

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("GET %s: failed to unmarshal JSON: %v", tt.path, err)
			}

			if got, ok := body[tt.wantKey]; !ok || got != tt.wantValue {
				t.Errorf("GET %s: body[%q] = %q, want %q", tt.path, tt.wantKey, got, tt.wantValue)
			}
		})
	}
}

func TestExistingEndpoints(t *testing.T) {
	r := newRouter()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "hello returns 200",
			path:       "/hello",
			wantStatus: http.StatusOK,
		},
		{
			name:       "version returns 200",
			path:       "/version",
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

func TestGreet(t *testing.T) {
	r := newRouter()

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "with name",
			query:      "name=Alice",
			wantStatus: http.StatusOK,
			wantMsg:    "Hello, Alice!",
		},
		{
			name:       "without name",
			query:      "",
			wantStatus: http.StatusOK,
			wantMsg:    "Hello, World!",
		},
		{
			name:       "empty name",
			query:      "name=",
			wantStatus: http.StatusOK,
			wantMsg:    "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/greet?"+tt.query, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GET /greet: got status %d, want %d", rec.Code, tt.wantStatus)
			}

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("GET /greet: failed to unmarshal JSON: %v", err)
			}

			if got := body["message"]; got != tt.wantMsg {
				t.Errorf("GET /greet: message = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestRequestLoggerMiddleware(t *testing.T) {
	r := newRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "GET /hello",
			method:     http.MethodGet,
			path:       "/hello",
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /nonexistent returns 404",
			method:     http.MethodGet,
			path:       "/nonexistent",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("%s %s: got status %d, want %d", tt.method, tt.path, rec.Code, tt.wantStatus)
			}
		})
	}
}
