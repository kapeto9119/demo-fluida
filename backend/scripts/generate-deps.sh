#!/bin/bash

# Generate go.sum file if it doesn't exist
if [ ! -f "../go.sum" ]; then
  echo "Generating go.sum file for Docker builds..."
  cd ..
  go mod tidy
  cd scripts
fi

echo "go.sum file is ready for Docker builds." 