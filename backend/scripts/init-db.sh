#!/bin/bash

# Script to initialize database for development
set -e

echo "=========================================="
echo "Initializing Fluida Database"
echo "=========================================="

# Database credentials from environment or use defaults
DB_HOST=${DB_HOST:-db}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-fluida}

# Function to execute SQL
execute_sql() {
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "$1"
}

# Function to check if database exists
db_exists() {
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME
}

# Check if database exists
if db_exists; then
    echo "Database '$DB_NAME' already exists."
else
    echo "Creating database '$DB_NAME'..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
fi

# Connect to database and initialize schema
echo "Initializing database schema..."

# Create enum type for invoice status
execute_sql "DO \$\$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'invoice_status') THEN
        CREATE TYPE invoice_status AS ENUM ('PENDING', 'PAID', 'CANCELED');
    END IF;
END \$\$;"

# Enable JSONB support
execute_sql "CREATE EXTENSION IF NOT EXISTS pgcrypto;"

# Create invoices table
execute_sql "CREATE TABLE IF NOT EXISTS invoice (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(255) NOT NULL UNIQUE,
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USDC',
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    due_date TIMESTAMP NOT NULL,
    receiver_addr VARCHAR(255) NOT NULL,
    link_token VARCHAR(255) NOT NULL UNIQUE,
    sender_details JSONB NOT NULL,
    recipient_details JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);"

# Create indexes for performance
execute_sql "CREATE INDEX IF NOT EXISTS idx_invoice_status ON invoice(status);"
execute_sql "CREATE INDEX IF NOT EXISTS idx_invoice_link_token ON invoice(link_token);"
execute_sql "CREATE INDEX IF NOT EXISTS idx_invoice_invoice_number ON invoice(invoice_number);"
execute_sql "CREATE INDEX IF NOT EXISTS idx_invoice_receiver_addr ON invoice(receiver_addr);"

echo "Database initialization completed successfully!" 