// Package middleware provides HTTP middleware functions
package middleware

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ncapetillo/demo-fluida/internal/response"
)

// SimpleRateLimiter implements a basic in-memory rate limiter
type SimpleRateLimiter struct {
	// How many requests per minute to allow
	requestsPerMinute int
	
	// Map to store IP addresses and their request counts
	clients map[string]*client
	
	// Mutex to protect the map
	mu sync.Mutex
}

// client stores information about each client
type client struct {
	count      int       // How many requests made
	lastSeen   time.Time // When was the last request
	resetTimer *time.Timer
}

// NewSimpleRateLimiter creates a new rate limiter
func NewSimpleRateLimiter(requestsPerMinute int) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		requestsPerMinute: requestsPerMinute,
		clients:           make(map[string]*client),
	}
}

// RateLimit implements the middleware function
func (rl *SimpleRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for health checks to avoid false positives
		if strings.HasSuffix(r.URL.Path, "/health") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Get client IP address using X-Forwarded-For header if available
		ip := getRealIP(r)
		
		// Add debug logging in production environment
		if os.Getenv("ENVIRONMENT") == "production" && os.Getenv("DEBUG") == "true" {
			log.Printf("Rate limit check: IP=%s, Path=%s, UserAgent=%s", 
				ip, r.URL.Path, r.UserAgent())
		}
		
		// Check if this client has reached the rate limit
		if !rl.allowRequest(ip) {
			// Log rate limit hits
			log.Printf("Rate limit exceeded: IP=%s, Path=%s", ip, r.URL.Path)
			
			// Return a 429 Too Many Requests
			response.Error(w, http.StatusTooManyRequests, 
				"Rate limit exceeded. Please try again later.", "rate_limit_exceeded")
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// getRealIP extracts the real client IP from request headers
// Handles various proxy scenarios including Railway and other cloud platforms
func getRealIP(r *http.Request) string {
	// Check X-Forwarded-For header first (most common in proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs - use the leftmost one
		// as it represents the original client
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			return clientIP
		}
	}
	
	// Check X-Real-IP header (used by some proxies)
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	
	// If no proxy headers are available, use the remote address
	// but strip the port if present
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there's an error (likely no port in the address), use RemoteAddr as is
		return r.RemoteAddr
	}
	
	return ip
}

// allowRequest checks if a request from a client should be allowed
func (rl *SimpleRateLimiter) allowRequest(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Get the client or create a new one
	c, exists := rl.clients[ip]
	if !exists {
		// Create new client
		c = &client{
			count:    1,
			lastSeen: time.Now(),
		}
		
		// Set up a timer to reset this client after 1 minute
		c.resetTimer = time.AfterFunc(time.Minute, func() {
			rl.mu.Lock()
			delete(rl.clients, ip)
			rl.mu.Unlock()
		})
		
		rl.clients[ip] = c
		return true
	}
	
	// Check if client has exceeded rate limit
	if c.count >= rl.requestsPerMinute {
		return false
	}
	
	// Increment count and update last seen
	c.count++
	c.lastSeen = time.Now()
	
	return true
} 