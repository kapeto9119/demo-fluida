#!/bin/bash
set -e

echo "🚀 Building project locally with Nixpacks (simulating Railway build)"

# Check if Nixpacks is installed
if ! command -v nixpacks &> /dev/null; then
    echo "❌ Nixpacks is not installed. Please install it first:"
    echo "curl -fsSL https://nixpacks.com/install.sh | bash"
    exit 1
fi

# Clean any previous builds
echo "🧹 Cleaning previous builds..."
rm -rf .nixpacks

# Build using nixpacks with the same configuration Railway uses
echo "🔨 Building with Nixpacks..."
nixpacks build . --name fluida-local

# Generate a Docker image from the build
echo "🐳 Creating Docker image from Nixpacks build..."
nixpacks images create --name fluida-local-image --build-id fluida-local

# Success message
echo "✅ Build completed successfully!"
echo ""
echo "To run the application locally, use:"
echo "docker run -p 8080:8080 -p 3000:3000 --env-file .env fluida-local-image"
echo ""
echo "Make sure you have all the required environment variables in a .env file." 