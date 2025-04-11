package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimit(t *testing.T) {
	// Create a new rate limiter with a limit of 3 requests per minute
	limiter := NewSimpleRateLimiter(3)
	
	// Create a simple handler for testing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply rate limiting to the handler
	limitedHandler := limiter.RateLimit(handler)
	
	// Test scenario: Make 5 requests from the same IP
	// Expect first 3 to succeed and last 2 to be rate limited
	for i := 0; i < 5; i++ {
		// Create a test request
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		
		// Set remote address to simulate the same client IP
		req.RemoteAddr = "192.168.1.1:12345"
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Handle the request
		limitedHandler.ServeHTTP(rr, req)
		
		// Check the response
		if i < 3 {
			// First 3 requests should succeed
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %d for request %d, got %d", http.StatusOK, i+1, status)
			}
		} else {
			// Last 2 requests should be rate limited
			if status := rr.Code; status != http.StatusTooManyRequests {
				t.Errorf("Expected status %d for request %d, got %d", http.StatusTooManyRequests, i+1, status)
			}
		}
	}
}

func TestGetRealIP(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		remoteAddr string
		expected string
	}{
		{
			name:     "X-Forwarded-For Header",
			headers:  map[string]string{"X-Forwarded-For": "192.168.1.1"},
			remoteAddr: "10.0.0.1:12345",
			expected: "192.168.1.1",
		},
		{
			name:     "X-Forwarded-For With Multiple IPs",
			headers:  map[string]string{"X-Forwarded-For": "192.168.1.1, 10.0.0.1, 172.16.0.1"},
			remoteAddr: "10.0.0.1:12345",
			expected: "192.168.1.1",
		},
		{
			name:     "X-Real-IP Header",
			headers:  map[string]string{"X-Real-IP": "192.168.1.1"},
			remoteAddr: "10.0.0.1:12345",
			expected: "192.168.1.1",
		},
		{
			name:     "Remote Address Only",
			headers:  map[string]string{},
			remoteAddr: "10.0.0.1:12345",
			expected: "10.0.0.1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			
			actual := getRealIP(req)
			if actual != tt.expected {
				t.Errorf("getRealIP() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestHealthChecksBypass(t *testing.T) {
	// Create a new rate limiter with a limit of 1 request per minute
	limiter := NewSimpleRateLimiter(1)
	
	// Track if handler was called
	handlerCalled := 0
	
	// Create a simple handler for testing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply rate limiting to the handler
	limitedHandler := limiter.RateLimit(handler)
	
	// Create a test request for health check
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.RemoteAddr = "192.168.1.1:12345"
	
	// Make multiple health check requests - they should bypass rate limiting
	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)
		
		// All requests should succeed
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d for health check request %d, got %d", 
				http.StatusOK, i+1, status)
		}
	}
	
	// Verify that handler was called 5 times
	if handlerCalled != 5 {
		t.Errorf("Expected handler to be called 5 times, got %d", handlerCalled)
	}
} 