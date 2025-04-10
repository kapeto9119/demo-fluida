#!/bin/bash
set -e

echo "🚀 Starting Fluida Invoice Generator on Railway"

# Get PORT variable from environment
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
echo "🔵 Starting backend server in background..."
cd backend
./app > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# Wait for backend to initialize
echo "Waiting for backend to initialize..."
sleep 5

# Start frontend server
echo "🔵 Starting frontend server..."
cd frontend
npx next start -p ${FRONTEND_PORT} > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Log success message
echo "✅ Fluida Invoice Generator is now running!"
echo "🌐 Backend API available at http://localhost:$PORT"
echo "🌐 Frontend available at http://localhost:$FRONTEND_PORT"

# Keep script running with status updates
while true; do
  sleep 60
  echo "System is still running... ($(date))"
  echo "Backend status: $(if kill -0 $BACKEND_PID 2>/dev/null; then echo 'Running'; else echo 'Stopped'; fi)"
  echo "Frontend status: $(if kill -0 $FRONTEND_PID 2>/dev/null; then echo 'Running'; else echo 'Stopped'; fi)"
done &
STATUS_CHECK_PID=$!

# Handle graceful shutdown
trap 'kill $BACKEND_PID $FRONTEND_PID $STATUS_CHECK_PID; exit' SIGINT SIGTERM

# Keep script running
wait 