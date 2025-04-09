#!/bin/bash

set -e

echo "=============================================="
echo "🚀 Fluida Invoice Generator Setup"
echo "=============================================="

# Store the current directory
PROJECT_ROOT=$(pwd)

# Function to print status messages
print_status() {
  echo "🔵 $1"
}

# Function to print error messages
print_error() {
  echo "🔴 $1"
}

# Function to print success messages
print_success() {
  echo "✅ $1"
}

# Function to print info messages
print_info() {
  echo "ℹ️  $1"
}

# Go to frontend directory and install dependencies
print_status "[1/4] Installing frontend dependencies..."
cd frontend
npm install || npm install --legacy-peer-deps

# Build the frontend
print_status "[2/4] Building frontend..."
npm run build

# Go back to root directory
cd $PROJECT_ROOT

# Check if Go is installed
if ! command -v go &> /dev/null; then
  print_info "Go is not installed on your system."
  print_info "Please install Go to run the backend in development mode."
  print_info "Visit https://golang.org/doc/install for installation instructions."
fi

# Check if Docker is running
print_status "[3/4] Checking if Docker is running..."
if ! docker info > /dev/null 2>&1; then
  print_error "Docker is not running or not installed."
  print_info "Starting in development mode instead..."
  
  # Start the PostgreSQL database if it's installed locally
  if command -v pg_ctl &> /dev/null; then
    print_info "Starting local PostgreSQL database..."
    # Add PostgreSQL startup commands here if needed
  else
    print_info "PostgreSQL is not installed locally."
    print_info "You will need to set up a database to use the application."
  fi
  
  echo ""
  print_info "Instructions for running in development mode:"
  echo "  1. In this terminal, run the following command to start the frontend:"
  echo "     ./dev-frontend.sh"
  echo ""
  echo "  2. In a new terminal, run the following command to start the backend:"
  echo "     ./dev-backend.sh"
  exit 1
fi

# Prepare dependencies for Docker
print_status "Preparing backend dependencies for Docker..."
cd backend
go mod tidy
cd $PROJECT_ROOT

# Run docker-compose with better error handling
print_status "[4/4] Starting docker containers..."
if ! docker-compose up -d; then
  print_error "Docker Compose failed to start all containers."
  echo ""
  print_info "Let's try to start just the database and run the services manually:"
  
  print_info "Starting PostgreSQL database in Docker..."
  if ! docker-compose up -d db; then
    print_error "Failed to start the database container."
    print_error "Please ensure Docker is running correctly and try again."
    exit 1
  fi
  
  echo ""
  print_info "Instructions for manual startup:"
  echo "  1. In this terminal, run the frontend:"
  echo "     ./dev-frontend.sh"
  echo ""
  echo "  2. In a new terminal, run the backend:"
  echo "     ./dev-backend.sh"
  exit 1
fi

print_success "Application started successfully!"
echo ""
print_info "Access your application at:"
echo "  Frontend: http://localhost:3000"
echo "  Backend: http://localhost:8080"
echo ""
print_info "To view logs:"
echo "  docker-compose logs -f"

print_info "To stop the application:"
echo "  docker-compose down" 