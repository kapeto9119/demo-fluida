// Package middleware provides HTTP middleware functions
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		// Set Strict-Transport-Security only in production
		// TODO: Use environment variable to detect production
		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		next.ServeHTTP(w, r)
	})
}

// RequestLogger is a customized request logger with better API request logging
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		// Get the start time
		start := time.Now()
		
		// Process the request
		defer func() {
			// Use standard logging for request logging
			log.Printf(
				"Request completed: status=%d duration=%s path=%s method=%s query=%s",
				ww.Status(),
				time.Since(start).String(),
				r.URL.Path,
				r.Method,
				r.URL.RawQuery,
			)
		}()
		
		next.ServeHTTP(ww, r)
	})
}

// RateLimiter implements a simple rate limiting middleware
// TODO: Replace with a more robust solution like github.com/ulule/limiter
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Placeholder for real rate limiting implementation
		next.ServeHTTP(w, r)
	})
}

// ValidateContentType ensures the request content type is application/json for POST/PUT/PATCH
func ValidateContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only check content type for methods that typically include a body
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			
			// Check if Content-Type is application/json
			if contentType != "application/json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(`{"error":{"message":"Content-Type must be application/json","code":"unsupported_media_type"}}`))
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}
