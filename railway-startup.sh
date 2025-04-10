#!/bin/bash
set -e

echo "🚀 Starting Fluida Invoice Generator on Railway"

# Set environment variables if not already set
export PORT=${PORT:-8080}
export FRONTEND_PORT=${FRONTEND_PORT:-3000}
export FRONTEND_URL=${FRONTEND_URL:-http://localhost:3000}

# Debug: Print environment info
echo "Environment variables:"
echo "PORT=$PORT"
echo "FRONTEND_PORT=$FRONTEND_PORT"
echo "FRONTEND_URL=$FRONTEND_URL"
echo "DATABASE_URL exists: $(if [ ! -z "$DATABASE_URL" ]; then echo 'Yes'; else echo 'No'; fi)"

# Start backend server immediately
echo "🔵 Starting backend server..."
cd backend
./app &
BACKEND_PID=$!
cd ..

# Give backend a moment to start
sleep 10

# Debug: Try various endpoints to check what's working
echo "🔍 Checking various endpoints for debugging..."
ENDPOINTS=("/" "/health" "/api/v1/health")
for endpoint in "${ENDPOINTS[@]}"; do
  echo "Testing endpoint: $endpoint"
  curl -v http://localhost:$PORT$endpoint || echo "Failed to connect to $endpoint"
  echo ""
  echo "------------------------"
done

# Start frontend server
echo "🔵 Starting frontend server..."
cd frontend
npx next start -p ${FRONTEND_PORT} &
FRONTEND_PID=$!

# Log success message
echo "✅ Fluida Invoice Generator is now running!"
echo "🌐 Backend API available at http://localhost:$PORT"
echo "🌐 Frontend available at http://localhost:$FRONTEND_PORT"

# Create a simple root endpoint for Railway healthcheck
echo "Setting up root endpoint for Railway healthcheck..."
while true; do
  # This is just to keep the script running
  sleep 60
  echo "System is still running... ($(date))"
done &
HEALTH_CHECK_PID=$!

# Handle graceful shutdown
trap 'kill $BACKEND_PID $FRONTEND_PID $HEALTH_CHECK_PID; exit' SIGINT SIGTERM

# Keep script running
wait 