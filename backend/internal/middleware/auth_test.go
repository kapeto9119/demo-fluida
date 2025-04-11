package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	// Set up environment for testing
	username := "testuser"
	password := "testpass"
	
	originalUsername := os.Getenv("AUTH_USERNAME")
	originalPassword := os.Getenv("AUTH_PASSWORD")
	
	// Set test credentials
	os.Setenv("AUTH_USERNAME", username)
	os.Setenv("AUTH_PASSWORD", password)
	
	// Restore original environment after test
	defer func() {
		os.Setenv("AUTH_USERNAME", originalUsername)
		os.Setenv("AUTH_PASSWORD", originalPassword)
	}()
	
	// Create a test handler that always succeeds
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply auth middleware
	handler := BasicAuth(testHandler)
	
	tests := []struct {
		name       string
		path       string
		auth       string
		wantStatus int
	}{
		{
			name:       "No auth, regular path",
			path:       "/api/invoices",
			auth:       "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Health check bypasses auth",
			path:       "/health",
			auth:       "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "API health check bypasses auth",
			path:       "/api/health",
			auth:       "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Valid auth credentials",
			path:       "/api/invoices",
			auth:       basicAuth(username, password),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid username",
			path:       "/api/invoices",
			auth:       basicAuth("wronguser", password),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid password",
			path:       "/api/invoices",
			auth:       basicAuth(username, "wrongpass"),
			wantStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, _ := http.NewRequest("GET", tt.path, nil)
			
			// Add auth header if provided
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Handle the request
			handler.ServeHTTP(rr, req)
			
			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", 
					status, tt.wantStatus)
			}
		})
	}
}

// basicAuth creates a basic auth header value
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
} 