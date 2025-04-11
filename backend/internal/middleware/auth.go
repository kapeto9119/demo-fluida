// Package middleware provides HTTP middleware functions
package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/ncapetillo/demo-fluida/internal/response"
)

// BasicAuth implements a simple Basic Auth middleware
func BasicAuth(next http.Handler) http.Handler {
	// Get credentials from environment, or use defaults if not set
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")
	
	// If credentials aren't set, disable auth
	if username == "" || password == "" {
		return next
	}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check and public APIs
		if r.URL.Path == "/health" || r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Get credentials from request
		user, pass, ok := r.BasicAuth()
		
		// If credentials are invalid, show auth prompt
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || 
		   subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			response.Error(w, http.StatusUnauthorized, "Unauthorized", "authentication_required")
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
} 