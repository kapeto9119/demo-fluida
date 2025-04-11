package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/ncapetillo/demo-fluida/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	SQL *sql.DB
)

// Config represents database configuration
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() Config {
	// First check Railway-specific PostgreSQL variables (PGHOST, PGUSER, etc.)
	host := getEnvOrDefault("PGHOST", "")
	port := getEnvOrDefault("PGPORT", "")
	user := getEnvOrDefault("PGUSER", "")
	password := getEnvOrDefault("PGPASSWORD", "")
	dbName := getEnvOrDefault("PGDATABASE", "")
	
	// If Railway variables not found, fall back to our custom DB_* variables
	if host == "" {
		host = getEnvOrDefault("DB_HOST", "db")
	}
	if port == "" {
		port = getEnvOrDefault("DB_PORT", "5432")
	}
	if user == "" {
		user = getEnvOrDefault("DB_USER", "postgres")
	}
	if password == "" {
		password = getEnvOrDefault("DB_PASSWORD", "postgres")
	}
	if dbName == "" {
		dbName = getEnvOrDefault("DB_NAME", "fluida")
	}
	
	log.Printf("Using database configuration: host=%s, port=%s, user=%s, dbname=%s", host, port, user, dbName)
	
	return Config{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		DBName:       dbName,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifetime:  time.Hour,
	}
}

// getEnvOrDefault retrieves environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect establishes a connection to the database
func Connect() (*sql.DB, error) {
	// Check if DATABASE_URL is provided (Railway deployment)
	databaseURL := os.Getenv("DATABASE_URL")
	
	// Check for Railway-specific database URL if DATABASE_URL is not set
	if databaseURL == "" {
		databaseURL = os.Getenv("PGDATABASE_URL")
	}
	
	var dsn string
	if databaseURL != "" {
		// Use the complete connection string from Railway
		dsn = databaseURL
		log.Println("Using database URL for PostgreSQL connection (Railway)")
	} else {
		// Use individual connection parameters (local development)
		config := LoadConfigFromEnv()
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.DBName)
		log.Println("Using individual parameters for PostgreSQL connection")
	}

	// Connect with retry logic
	var err error
	var sqlDB *sql.DB
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		sqlDB, err = sql.Open("postgres", dsn)
		if err == nil {
			err = sqlDB.Ping()
			if err == nil {
				break
			}
		}
		
		// Print more detailed error for diagnostic purposes
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		
		// Print more environment details to help diagnose the issue
		if i == 0 {
			// Only print environment details on first attempt
			config := LoadConfigFromEnv()
			log.Printf("Connection details: host=%s, port=%s, user=%s, dbname=%s", 
				config.Host, config.Port, config.User, config.DBName)
			
			// Print all environment variables with PG prefix for debugging
			log.Println("Environment variables for database connection:")
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "PG") || strings.HasPrefix(env, "DB_") {
					// Don't log passwords
					if !strings.Contains(env, "PASSWORD") && !strings.Contains(env, "PASS") {
						log.Println(env)
					} else {
						parts := strings.SplitN(env, "=", 2)
						if len(parts) == 2 {
							log.Printf("%s=********", parts[0])
						}
					}
				}
			}
		}
		
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	config := LoadConfigFromEnv() // Get default config for connection pool settings
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	// Initialize GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
		PreferSimpleProtocol: true, // Prefer simple protocol for better compatibility
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
		DisableForeignKeyConstraintWhenMigrating: true, // Disable FK constraints during migrations
		PrepareStmt: true, // Cache prepared statements for better performance
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %v", err)
	}

	// Set global DB variable
	DB = gormDB
	SQL = sqlDB

	// Auto-migrate the schema
	if err := migrateSchema(); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	log.Println("Successfully connected to database")
	return sqlDB, nil
}

// migrateSchema automatically creates or updates the database tables
func migrateSchema() error {
	// Create custom enum type for invoice status if it doesn't exist
	// This must be done manually since GORM doesn't handle PostgreSQL enums well
	DB.Exec(`DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'invoice_status') THEN
			CREATE TYPE invoice_status AS ENUM ('PENDING', 'PAID', 'CANCELED');
		END IF;
	END $$;`)

	// Ensure UUID and JSONB support is enabled
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	
	// Run auto migrations for all models
	if err := DB.AutoMigrate(&models.Invoice{}, &models.DraftInvoice{}); err != nil {
		return fmt.Errorf("failed to migrate schema: %v", err)
	}
	
	// Create indexes for better performance
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoice(status);")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_invoices_link_token ON invoice(link_token);")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_draft_invoice_user_id ON draft_invoice(user_id);")
	
	log.Println("Database schema migrated successfully")
	return nil
} 