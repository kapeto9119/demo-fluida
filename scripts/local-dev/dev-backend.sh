#!/bin/bash

set -e

echo "========================================"
echo "🚀 Starting Fluida Backend (Dev Mode)"
echo "========================================"

# Check if .env file exists at the root level, otherwise create from example
if [ ! -f .env ]; then
  echo "🔧 Creating root .env file from .env.example..."
  cp .env.example .env
  echo "✅ Root .env file created. You may want to review and update it."
else
  echo "✅ Using existing root .env file"
fi

# Change directory to backend
cd backend

echo "📦 Syncing dependencies..."
go mod tidy

echo "🔄 Checking if database is running..."
if ! docker ps | grep -q fluida-db; then
  echo "ℹ️  Database not running. Starting PostgreSQL database..."
  echo "ℹ️  Running ./dev-db.sh..."
  cd .. && ./dev-db.sh && cd backend
else
  echo "✅ Database is already running"
fi

# Load environment variables from root .env file
echo "🔧 Loading environment variables from root .env file..."
# Export variables from the parent directory's .env file
eval "$(cd .. && grep -v '^#' .env | sed 's/^/export /')"

# Override DB_HOST for local development
export DB_HOST=localhost
export PORT=$BACKEND_PORT

echo "🔥 Starting Go server..."
go run cmd/server/main.go

# Note: Server will keep running until terminated with Ctrl+C 