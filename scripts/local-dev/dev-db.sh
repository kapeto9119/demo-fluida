#!/bin/bash

set -e

echo "========================================"
echo "🚀 Starting Fluida Database (Dev Mode)"
echo "========================================"

# Check if .env file exists at the root level, otherwise create from example
if [ ! -f .env ]; then
  echo "🔧 Creating .env file from .env.example..."
  cp .env.example .env
  echo "✅ Root .env file created. You may want to review and update it."
else
  echo "✅ Using existing root .env file"
fi

echo "🔄 Starting PostgreSQL database..."
docker-compose up -d db

echo "⏳ Waiting for database to be ready..."
max_retries=20
retry_count=0

while [ $retry_count -lt $max_retries ]; do
  retry_count=$((retry_count+1))
  echo "Checking database connection (attempt $retry_count/$max_retries)..."
  
  if docker-compose exec db pg_isready -U postgres -d fluida > /dev/null 2>&1; then
    echo "✅ Database is ready!"
    break
  else
    if [ $retry_count -eq $max_retries ]; then
      echo "❌ Failed to connect to database after $max_retries attempts"
      exit 1
    fi
    echo "Database not ready yet, retrying in 1 second..."
    sleep 1
  fi
done

echo "🎉 Database is now running and ready for connections!"
echo "   - Host: localhost"
echo "   - Port: 5432"
echo "   - User: postgres"
echo "   - Database: fluida"
echo ""
echo "To connect, use: psql -h localhost -U postgres -d fluida"
echo ""
echo "Database will keep running in the background."
echo "To stop it, run: docker-compose stop db" 