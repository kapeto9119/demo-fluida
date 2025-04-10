#!/bin/bash
set -e

echo "🚀 Starting Fluida Invoice Generator on Railway"

# Set environment variables if not already set
export PORT=${PORT:-8080}
export DATABASE_URL=${DATABASE_URL:-postgresql://postgres:postgres@localhost:5432/fluida}
export FRONTEND_URL=${FRONTEND_URL:-http://localhost:3000}

# Check if we need to run migrations
if [ "$RUN_MIGRATIONS" = "true" ]; then
  echo "🔄 Running database migrations..."
  # Add migration commands here if needed
fi

# Start the backend server
echo "🔵 Starting backend server..."
cd backend
./app &
BACKEND_PID=$!

# Start frontend server
echo "🔵 Starting frontend server..."
cd ../frontend
npx next start -p ${FRONTEND_PORT:-3000} &
FRONTEND_PID=$!

# Handle graceful shutdown
trap 'kill $BACKEND_PID $FRONTEND_PID; exit' SIGINT SIGTERM

# Keep script running
wait 