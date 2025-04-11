#!/bin/bash
set -e

echo "🔍 Debugging TailwindCSS installation"

# Go to frontend directory
cd frontend

# Clean previous installations
echo "🧹 Cleaning previous installations..."
rm -rf node_modules
rm -rf .next

# Install dependencies
echo "📦 Installing dependencies..."
npm ci

# Explicitly install TailwindCSS
echo "📦 Installing TailwindCSS..."
npm install --no-save tailwindcss postcss autoprefixer @tailwindcss/forms

# Check if TailwindCSS is properly installed
echo "🔍 Verifying TailwindCSS installation..."
if [ -d "node_modules/tailwindcss" ]; then
  echo "✅ TailwindCSS is properly installed"
  ls -la node_modules/tailwindcss
else
  echo "❌ TailwindCSS installation failed"
  exit 1
fi

# Check for tailwind.config.js
echo "🔍 Checking Tailwind configuration..."
if [ -f "tailwind.config.js" ]; then
  echo "✅ tailwind.config.js exists"
  cat tailwind.config.js
else
  echo "❌ tailwind.config.js is missing"
  exit 1
fi

# Check for postcss.config.js
if [ -f "postcss.config.js" ]; then
  echo "✅ postcss.config.js exists"
  cat postcss.config.js
else
  echo "❌ postcss.config.js is missing"
  exit 1
fi

# Try to build with NODE_ENV=production
echo "🔨 Testing production build..."
NODE_ENV=production npm run build

echo "✅ TailwindCSS debug completed successfully" 