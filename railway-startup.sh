#!/bin/bash
set -e

echo "🚀 Starting Fluida Invoice Generator on Railway"

# Set environment variables if not already set
export PORT=${PORT:-8080}
export FRONTEND_PORT=${FRONTEND_PORT:-3000}
export FRONTEND_URL=${FRONTEND_URL:-http://localhost:3000}

# Wait for PostgreSQL to be available
echo "🔄 Checking database connection..."
MAX_RETRIES=5
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  # Use Go app's health check to verify DB connection
  echo "Attempt $((RETRY_COUNT+1))/$MAX_RETRIES: Checking database connection..."
  if curl -s http://localhost:$PORT/health 2>&1 | grep -q "OK"; then
    echo "✅ Database connection successful!"
    break
  fi
  
  # Start the backend if this is the first attempt
  if [ $RETRY_COUNT -eq 0 ]; then
    echo "🔵 Starting backend server..."
    cd backend
    ./app &
    BACKEND_PID=$!
    sleep 5 # Give backend a chance to start
  fi
  
  RETRY_COUNT=$((RETRY_COUNT+1))
  if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "⏳ Waiting for database connection... (${RETRY_COUNT}/${MAX_RETRIES})"
    sleep 10
  fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
  echo "❌ Failed to connect to database after $MAX_RETRIES attempts"
  echo "🔍 Backend logs:"
  # Show logs from backend to help debugging
  if [ ! -z "$BACKEND_PID" ]; then
    kill $BACKEND_PID
  fi
  exit 1
fi

# If backend wasn't started in the DB check loop, start it now
if [ -z "$BACKEND_PID" ]; then
  echo "🔵 Starting backend server..."
  cd backend
  ./app &
  BACKEND_PID=$!
  cd ..
else
  cd ..
fi

# Start frontend server
echo "🔵 Starting frontend server..."
cd frontend
npx next start -p ${FRONTEND_PORT:-3000} &
FRONTEND_PID=$!

# Log success message
echo "✅ Fluida Invoice Generator is now running!"
echo "🌐 Backend API available at http://localhost:$PORT"
echo "🌐 Frontend available at http://localhost:$FRONTEND_PORT"

# Handle graceful shutdown
trap 'kill $BACKEND_PID $FRONTEND_PID; exit' SIGINT SIGTERM

# Keep script running
wait 