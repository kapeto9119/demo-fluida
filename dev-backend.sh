#!/bin/bash

set -e

echo "========================================"
echo "🚀 Starting Fluida Backend (Dev Mode)"
echo "========================================"

# Change directory to backend
cd backend

echo "📦 Syncing dependencies..."
go mod tidy

echo "🔄 Checking if database is running..."
if ! docker ps | grep -q fluida-db; then
  echo "ℹ️  Starting PostgreSQL database..."
  docker-compose up -d db
  
  # Wait for database to be ready
  echo "⏳ Waiting for database to be ready..."
  sleep 5
fi

echo "🔥 Starting Go server..."
go run cmd/server/main.go

# Note: Server will keep running until terminated with Ctrl+C 