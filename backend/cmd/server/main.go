package main

import (
	"context"
	"encoding/json"
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
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	
	// Add our custom security middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.ValidateContentType)
	
	// Configure CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fluida Invoice API - v1.0"))
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Health check failed: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database connection unavailable"))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
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
					
					// Return standardized error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusServiceUnavailable)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Database connection unavailable",
						"code":  "database_error",
					})
					return
				}
				
				// Return standardized success response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "OK",
					"version": "1.0.0",
					"dependencies": map[string]string{
						"database": "connected",
					},
				})
			})
			
			r.Mount("/invoices", invoiceHandler.Routes())
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