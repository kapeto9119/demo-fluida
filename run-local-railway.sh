#!/bin/bash
set -e

echo "🚀 Running Fluida locally using the Railway build setup"

# Check if .env file exists, otherwise copy from example
if [ ! -f .env ]; then
    echo "⚠️  .env file not found, copying from .env.example"
    cp .env.example .env
    echo "⚠️  Please check the .env file and update values as needed"
    echo "   Press Enter to continue or Ctrl+C to abort"
    read
fi

# Check if the image exists
if ! docker image inspect fluida-local-image &> /dev/null; then
    echo "❌ Image 'fluida-local-image' not found. Build it first:"
    echo "./local-railway-build.sh"
    exit 1
fi

echo "🐳 Running Docker container with Railway configuration..."
docker run -it --rm \
    -p 8080:8080 \
    -p 3000:3000 \
    --env-file .env \
    fluida-local-image

echo "✅ Application stopped" 