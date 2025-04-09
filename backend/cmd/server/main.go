package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/handlers"
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
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	invoiceService := services.NewInvoiceService()

	// Initialize handlers
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	// Initialize and start Solana payment watcher
	paymentWatcher, err := solana.NewPaymentWatcher()
	if err != nil {
		log.Printf("Warning: Failed to initialize Solana payment watcher: %v", err)
		log.Println("Continuing without payment watcher - automatic payment detection will not work")
	} else {
		paymentWatcher.Start()
		defer paymentWatcher.Stop()
		log.Println("Solana payment watcher started successfully")
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Mount("/invoices", invoiceHandler.Routes())
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
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
} 