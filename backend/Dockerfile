FROM golang:1.20-alpine

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git bash

# Install a compatible version of air for hot reloading
RUN go install github.com/cosmtrek/air@v1.29.0

# Copy go.mod first to leverage Docker cache
COPY go.mod ./
# If go.sum exists, copy it too, but don't fail if it doesn't
COPY go.sum* ./
COPY scripts/generate-deps.sh ./scripts/

# Make the script executable and run it
RUN chmod +x ./scripts/generate-deps.sh && ./scripts/generate-deps.sh

# Initialize the module and download dependencies
RUN go mod tidy && go mod download

# Copy the rest of the application
COPY . .

# Expose port
EXPOSE 8080

# Run air for hot reloading in development
CMD ["air", "-c", ".air.toml"] 