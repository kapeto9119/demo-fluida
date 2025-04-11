# Base image with Alpine for smaller size
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache python3 make g++ linux-headers libc6-compat eudev-dev udev

# Create necessary directories
RUN mkdir -p public .next

# Copy package files first for better caching
COPY package.json package-lock.json* ./

# Install dependencies with explicit TailwindCSS installation
RUN npm ci
RUN npm install --no-save tailwindcss postcss autoprefixer @tailwindcss/forms
RUN ls -la node_modules/tailwindcss || echo "TailwindCSS not found"

# Copy source code
COPY . .

# Build the application
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build

# Set production environment
ENV NODE_ENV production

# Set default authentication (can be overridden)
ENV NEXT_PUBLIC_AUTH_USERNAME admin
ENV NEXT_PUBLIC_AUTH_PASSWORD fluida

# Expose port
EXPOSE 3000

# Start the application
CMD ["npm", "run", "start"] 