#!/bin/bash

set -e

echo "========================================"
echo "🚀 Starting Fluida Frontend (Dev Mode)"
echo "========================================"

# Change directory to frontend
cd frontend

echo "📦 Installing dependencies..."
npm install || npm install --legacy-peer-deps

echo "🔥 Starting development server..."
npm run dev

# Note: Dev server will keep running until terminated with Ctrl+C 