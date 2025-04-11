#!/bin/bash
set -e

echo "=============================================="
echo "🚀 Fluida Invoice Generator"
echo "=============================================="
echo "This is a wrapper script for the reorganized project structure."
echo ""

# Check if the user is asking for help
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  echo "Usage: ./start.sh [option]"
  echo ""
  echo "Options:"
  echo "  docker     Start using Docker Compose (recommended for most users)"
  echo "  frontend   Start only the frontend development server"
  echo "  backend    Start only the backend development server"
  echo "  db         Start only the database"
  echo "  help       Show this help message"
  echo ""
  echo "For more detailed instructions, please see README.md"
  exit 0
fi

# Process command-line arguments
if [[ "$1" == "docker" || "$#" -eq 0 ]]; then
  echo "🐳 Starting using Docker Compose..."
  ./scripts/startup.sh
elif [ "$1" == "frontend" ]; then
  echo "🖥️  Starting frontend development server..."
  ./scripts/local-dev/dev-frontend.sh
elif [ "$1" == "backend" ]; then
  echo "🔧 Starting backend development server..."
  ./scripts/local-dev/dev-backend.sh
elif [ "$1" == "db" ]; then
  echo "🗄️  Starting database..."
  ./scripts/local-dev/dev-db.sh
else
  echo "❌ Unknown option: $1"
  echo "Run './start.sh --help' for usage information."
  exit 1
fi 