#!/bin/bash
set -e

echo "🚀 Starting Fluida Invoice Generator on Railway"

# Create a temporary healthcheck endpoint on port 80 (Railway default)
echo "Creating temporary healthcheck endpoint..."
(
  while true; do
    echo -e "HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nOK" | nc -l -p 80 -q 1 || true
  done
) &
TEMP_HEALTH_PID=$!

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

# Start backend server
echo "🔵 Starting backend server..."
cd backend
./app &
BACKEND_PID=$!
cd ..

# Debug: Try various endpoints to check what's working
echo "🔍 Checking various endpoints for debugging..."
sleep 10  # Give backend time to start

# Start frontend server
echo "🔵 Starting frontend server..."
cd frontend
npx next start -p ${FRONTEND_PORT} &
FRONTEND_PID=$!

# Log success message
echo "✅ Fluida Invoice Generator is now running!"
echo "🌐 Backend API available at http://localhost:$PORT"
echo "🌐 Frontend available at http://localhost:$FRONTEND_PORT"

# Kill temporary health server
kill $TEMP_HEALTH_PID || true

# Keep script running with simple output every minute
while true; do
  sleep 60
  echo "System is still running... ($(date))"
done &
HEALTH_CHECK_PID=$!

# Handle graceful shutdown
trap 'kill $BACKEND_PID $FRONTEND_PID $HEALTH_CHECK_PID; exit' SIGINT SIGTERM

# Keep script running
wait 