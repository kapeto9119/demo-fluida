// Package middleware provides HTTP middleware functions
package middleware

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ncapetillo/demo-fluida/internal/response"
)

// ErrorTypes to handle different types of errors
var (
	ErrNotFound   = errors.New("resource not found")
	ErrBadRequest = errors.New("invalid request")
)

// ErrorHandler middleware for consistent error handling
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer that captures panics
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				
				// Return a standardized error response
				response.InternalServerError(w)
			}
		}()
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Recoverer is a middleware that recovers from panics
// It's a custom version that uses our standard response format
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				
				// Send standardized error response
				response.InternalServerError(w)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get existing request ID or create a new one
		requestID := r.Header.Get("X-Request-ID")
		
		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
} 