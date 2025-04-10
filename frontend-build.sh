#!/bin/bash
set -e

echo "📦 Setting up frontend build environment..."
cd frontend

# Diagnose the environment
echo "🧪 Environment diagnostics:"
echo "Node.js version: $(node --version)"
echo "NPM version: $(npm --version)"
echo "Current directory: $(pwd)"

# Clean the node_modules to ensure a fresh install
echo "🧹 Cleaning previous installations..."
rm -rf node_modules
rm -rf .next

# Ensure all dependencies are installed
echo "📦 Installing dependencies..."
npm ci

# Explicitly install Tailwind CSS and related packages
echo "📦 Ensuring Tailwind CSS is installed..."
npm install --save-dev tailwindcss postcss autoprefixer @tailwindcss/forms

# Verify Tailwind is installed
if [ -d "node_modules/tailwindcss" ]; then
  echo "✅ Tailwind CSS is properly installed"
  ls -la node_modules/tailwindcss
else
  echo "❌ Tailwind CSS installation failed"
  exit 1
fi

# Verify the config files
echo "📝 Checking Tailwind configuration files..."
if [ -f "tailwind.config.js" ]; then
  echo "✅ tailwind.config.js exists"
else
  echo "❌ tailwind.config.js is missing"
  exit 1
fi

if [ -f "postcss.config.js" ]; then
  echo "✅ postcss.config.js exists"
else
  echo "❌ postcss.config.js is missing"
  exit 1
fi

# Run the build with verbose logging
echo "🔨 Building frontend..."
NODE_ENV=production npm run build

echo "✅ Frontend build completed successfully" 