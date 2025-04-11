package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/handlers"
	"github.com/ncapetillo/demo-fluida/internal/middleware"
	"github.com/ncapetillo/demo-fluida/internal/repository"
	"github.com/ncapetillo/demo-fluida/internal/response"
	"github.com/ncapetillo/demo-fluida/internal/services"
	"github.com/ncapetillo/demo-fluida/internal/solana"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	log.Println("Connecting to database...")
	sqlDB, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Ensure database is closed when the program exits
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
	log.Println("Database connected successfully")
	
	// Initialize repository
	invoiceRepo := repository.NewInvoiceRepository(db.DB)
	
	// Initialize services
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// Initialize handlers
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)
	draftInvoiceHandler := handlers.NewDraftInvoiceHandler()

	// Initialize router
	r := chi.NewRouter()

	// Initialize Solana payment watcher
	solanaWatcher, err := solana.NewPaymentWatcher()
	if err != nil {
		log.Printf("Warning: Failed to initialize Solana payment watcher: %v", err)
		log.Println("Automatic payment detection will not work")
	} else {
		// Start watching for payments in a separate goroutine
		go solanaWatcher.WatchForPayments()
		
		// Ensure payment watcher is stopped on shutdown
		defer solanaWatcher.Stop()
	}

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestIDMiddleware)
	
	// Configure CORS - MUST be before BasicAuth for OPTIONS preflight requests
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://demo-fluida-production.up.railway.app", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	
	// Add Basic Authentication
	r.Use(middleware.BasicAuth)
	
	r.Use(middleware.ErrorHandler)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	
	// Configure rate limiter based on environment
	var rateLimit int
	if os.Getenv("ENVIRONMENT") == "production" {
		// Higher limit for production to handle legitimate traffic spikes
		rateLimit = 120 // 120 requests per minute per IP in production
	} else {
		// More restrictive for development/testing
		rateLimit = 60 // 60 requests per minute per IP in development
	}
	
	// Add custom rate limiter
	log.Printf("Configuring rate limiter with %d requests per minute per IP", rateLimit)
	rateLimiter := middleware.NewSimpleRateLimiter(rateLimit)
	r.Use(rateLimiter.RateLimit)
	
	// Add our custom security middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.ValidateContentType)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"name": "Fluida Invoice API",
			"version": "1.0.0",
			"docs": "/api/docs",
		})
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Health check failed: %v", err)
			response.Error(w, http.StatusServiceUnavailable, 
				"Database connection unavailable", "database_error")
			return
		}
		
		response.Success(w, http.StatusOK, "OK")
	})

	// Additional health check endpoint to match Railway.app configuration
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Health check failed: %v", err)
			
			response.Error(w, http.StatusServiceUnavailable, 
				"Database connection unavailable: "+err.Error(), "database_error")
			return
		}
		
		// Add more detailed health check information
		healthResponse := map[string]interface{}{
			"status": "OK",
			"version": "1.0.0",
			"dependencies": map[string]string{
				"database": "connected",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}
		
		response.JSON(w, http.StatusOK, healthResponse)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// API version 1
		r.Route("/v1", func(r chi.Router) {
			// Health check within API namespace
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				// Check database connection
				if err := sqlDB.Ping(); err != nil {
					log.Printf("Health check failed: %v", err)
					
					response.Error(w, http.StatusServiceUnavailable, 
						"Database connection unavailable", "database_error")
					return
				}
				
				response.JSON(w, http.StatusOK, map[string]interface{}{
					"status": "OK",
					"version": "1.0.0",
					"dependencies": map[string]string{
						"database": "connected",
					},
				})
			})
			
			// Authentication verification endpoint - returns 200 if auth is valid
			r.Get("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
				response.JSON(w, http.StatusOK, map[string]string{
					"status": "authenticated",
				})
			})
			
			r.Mount("/invoices", invoiceHandler.Routes())
			
			// Register draft invoice routes
			r.Mount("/invoices/drafts", draftInvoiceHandler.Routes())
			
			// Mount invoice number check endpoint outside of drafts
			r.Get("/invoices/check", draftInvoiceHandler.CheckInvoiceNumberExists)
		})
		
		// Redirect legacy API calls to the versioned API
		r.Mount("/invoices", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/api/v1"+r.URL.Path, http.StatusPermanentRedirect)
		}))
	})

	// Create server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	
	// Listen for interrupt signal
	go func() {
		// Listen for interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		
		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()
		
		log.Println("Server shutting down gracefully...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		serverStopCtx()
	}()
	
	// Start server
	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	
	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Println("Server exited properly")
} 