# Fluida Invoice & Payment Link Generator

A web application that allows businesses to generate invoices and payment links for receiving USDC payments on the Solana blockchain (testnet).

## Features

- **Invoice Creation**: Create detailed invoices with sender and recipient information.
- **Draft Saving**: Save invoice drafts and resume work later.
- **Payment Link Generation**: Automatically generate unique payment links for each invoice.
- **Solana Wallet Integration**: Connect to Phantom or other Solana wallets.
- **USDC Payments**: Process payments in USDC on Solana testnet.
- **Payment Detection**: Automatically detect incoming payments and mark invoices as paid.

## Tech Stack

### Backend
- **Language**: Go
- **Framework**: Chi (lightweight HTTP router)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Blockchain**: Solana (using `solana-go` library)

### Frontend
- **Framework**: Next.js with TypeScript 
- **Styling**: TailwindCSS
- **Wallet Integration**: Solana Wallet Adapter
- **HTTP Client**: Axios

### Infrastructure
- **Containerization**: Docker & Docker Compose

## Project Structure

The project follows a clean architecture approach:

```
.
├── backend/                  # Go backend
│   ├── cmd/                  # Command entrypoints
│   │   └── server/           # Main API server
│   ├── internal/             # Internal packages
│   │   ├── db/               # Database connection and utilities
│   │   ├── handlers/         # HTTP request handlers
│   │   ├── models/           # Data models
│   │   ├── services/         # Business logic layer
│   │   └── solana/           # Solana integration
│   ├── Dockerfile            # Backend Docker configuration
│   └── go.mod                # Go module definition
├── frontend/                 # Next.js frontend
│   ├── src/                  # Source code
│   │   ├── app/              # Next.js app directory (pages)
│   │   ├── components/       # React components
│   │   ├── lib/              # Utility functions
│   │   └── types/            # TypeScript type definitions
│   ├── Dockerfile            # Multi-stage Docker build
│   └── package.json          # NPM dependencies
├── docker-compose.yml        # Docker Compose configuration
└── README.md                 # Project documentation
```

### Backend Architecture

The backend follows a layered architecture pattern:

1. **Handlers Layer** (`/internal/handlers`): Handles HTTP requests, validates input, and returns responses.
2. **Services Layer** (`/internal/services`): Contains the business logic of the application.
3. **Models Layer** (`/internal/models`): Defines the data structures and domain objects.
4. **Database Layer** (`/internal/db`): Manages database connections and operations.
5. **Solana Layer** (`/internal/solana`): Handles Solana blockchain integration.

## Getting Started

This project has been reorganized for better maintainability and deployment. The codebase is now structured with:

- `frontend/`: Next.js frontend application
- `backend/`: Go backend API service
- `scripts/`: Helper scripts for development and operations
  - `scripts/local-dev/`: Scripts specifically for local development
  - `scripts/startup.sh`: Main startup script for Docker Compose setup

### Prerequisites

- Node.js 18+ and npm
- Go 1.20+
- Docker and Docker Compose (for local development)
- PostgreSQL (if developing without Docker)

### Local Development

You can run the application in two ways:

#### 1. Using Docker Compose (recommended)

This method spins up all services (frontend, backend, database) in Docker containers:

```bash
# Start all services
docker-compose up

# Or with the helper script (which handles common issues)
./scripts/startup.sh
```

#### 2. Using Individual Development Scripts

For more control during development, you can run each service separately:

```bash
# Start the database
./scripts/local-dev/dev-db.sh

# In a new terminal, start the backend
./scripts/local-dev/dev-backend.sh

# In another terminal, start the frontend
./scripts/local-dev/dev-frontend.sh
```

### Access the Application

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1

## Usage

1. **Create an Invoice**:
   - Navigate to the "Create Invoice" page
   - Fill in the invoice details (sender info, recipient info, amount, etc.)
   - Provide a Solana wallet address to receive the payment
   - Save as draft if you need to complete it later
   - Submit to generate a payment link

2. **Manage Draft Invoices**:
   - Save invoice drafts when you need to pause work
   - Return later to complete and submit the invoice
   - Each user can have one active draft saved in the database

3. **Share the Payment Link**:
   - Copy the generated link and share it with your client
   - Alternatively, email the link directly from the application

4. **Receive Payment**:
   - Client opens the link and sees the invoice details
   - Client connects their Solana wallet (Phantom, etc.)
   - Client approves the USDC transfer
   - The system automatically detects the payment and marks the invoice as paid

## Solana Integration

This application uses the Solana testnet (devnet) for demonstrating payment functionality:

- All transactions occur on the Solana testnet, not the mainnet
- For demo purposes, USDC transfers are simulated with SOL tokens
- In a production environment, real USDC SPL tokens would be used

To test the application, you'll need:
- A Phantom wallet (or other Solana wallet)
- Some devnet SOL (available from faucets)

## Areas for Improvement

With more time, these areas could be enhanced:

1. **Security**:
   - Add authentication and authorization
   - Implement rate limiting
   - Add better error handling and validation

2. **Features**:
   - Email notifications for payments
   - Payment reminders
   - Invoice templates
   - Invoice history and reporting
   - Multi-currency support

3. **Technical**:
   - Proper USDC SPL token implementation
   - WebSocket-based payment detection
   - Comprehensive test suite
   - CI/CD pipeline

## API-First Approach

This project was developed with an API-first mindset, enabling future extensibility and integrations:

### API Design Principles

- **Resource-Oriented**: APIs are organized around resources (invoices, payments)
- **Predictable URLs**: Consistent URL patterns (`/api/invoices/{id}`)
- **Standard HTTP Methods**: Using appropriate HTTP verbs (GET, POST, PUT)
- **Stateless Interactions**: Each request contains all information needed to process it

### Future API Evolution

The API is designed to evolve through:

1. **Versioning Strategy**: 
   - Current implicit version is v1
   - Future versions would use explicit path versioning (`/api/v2/invoices`)
   - Backward compatibility maintained across versions

2. **Extensibility Points**:
   - Webhook support for payment notifications
   - API keys for third-party integrations
   - Pagination and filtering for resource listing
   - Bulk operations for invoice management

3. **Potential Integrations**:
   - Accounting software (QuickBooks, Xero)
   - Payment processors
   - ERP systems
   - CRM platforms

### API Documentation

A full OpenAPI/Swagger specification could be generated to provide:
- Interactive documentation
- SDK generation for clients
- Automated testing of endpoints

## Security Best Practices

The application implements several security best practices, with notes on future enhancements:

### Wallet Key Security

- **Client-Side Only**: Private keys never leave the user's wallet
- **No Key Storage**: The application never stores private keys
- **Public Key Validation**: Implemented address validation before accepting transactions

### Data Protection

- **Sensitive Data Handling**: Personal information is validated and sanitized
- **Future Enhancement**: Implement field-level encryption for sensitive invoice data
- **Future Enhancement**: Add data retention policies and secure deletion

### API Security

- **Input Validation**: All inputs are validated before processing
- **Future Enhancement**: Implement JWT-based authentication
- **Future Enhancement**: Add CORS protection and CSP headers
- **Future Enhancement**: Implement rate limiting to prevent abuse

### Blockchain Security

- **Transaction Verification**: Multiple validations before confirming transactions
- **Future Enhancement**: Implement multi-signature support for high-value transactions
- **Future Enhancement**: Add transaction anomaly detection

### Audit and Compliance

- **Transaction Logging**: All blockchain interactions are logged
- **Future Enhancement**: Implement comprehensive audit logging for security events
- **Future Enhancement**: Add compliance reporting for financial regulations

## Troubleshooting

If you encounter issues:

1. **Docker not running:**
   - The startup script will notify you if Docker isn't available
   - Follow the instructions to run in development mode without Docker

2. **Frontend build issues:**
   - Install dependencies manually: `cd frontend && npm install`
   - Run in development mode: `npm run dev`

3. **Backend connection issues:**
   - Ensure PostgreSQL is running (if using Docker, check `docker ps`)
   - Check backend logs: `docker logs fluida-backend`

4. **Air hot-reload installation error:**
   - If you see an error about `github.com/cosmtrek/air`, it could be one of two issues:
     1. The package has moved to a new path (`github.com/air-verse/air`), OR
     2. The newer versions require Go 1.21+, which conflicts with our Go 1.20 Docker image
   - Solution: Edit the backend/Dockerfile to use a compatible version:
     ```
     RUN go install github.com/cosmtrek/air@v1.29.0
     ```
   - Alternatively, run the backend directly without Docker:
     ```
     cd backend
     go run cmd/server/main.go
     ```

5. **Docker build error: go.sum not found:**
   - This occurs when trying to build the Docker image without first generating the go.sum file
   - Solution 1: Run `go mod tidy` in the backend directory before building Docker containers
   - Solution 2: Use the provided startup script which handles this automatically
   - The error typically looks like: `failed to solve: failed to compute cache key: failed to calculate checksum of ref: "/go.sum": not found`

6. **Frontend Docker build error: Node.js version:**
   - The newer Solana dependencies require Node.js 20+, but our Docker image uses Node.js 18
   - Solutions:
     1. Run without Docker: `cd frontend && npm run dev`
     2. Update the frontend/Dockerfile to use Node.js 20:
        ```
        FROM node:20-alpine
        ```
     3. Install Python in the Docker container for native dependencies:
        ```
        RUN apk add --no-cache python3 make g++
        ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Railway Deployment

This project is configured for deployment on Railway as individual services. This approach provides better scalability, reliability, and maintainability compared to deploying everything as a single service.

### Deployment Architecture

- **Backend API Service**: Deployed as a standalone Go service
- **Frontend Web Application**: Deployed as a standalone Next.js service
- **Database**: Using Railway's managed PostgreSQL database

### How to Deploy to Railway

1. **Set up Railway CLI**:
   ```
   npm install -g @railway/cli
   railway login
   ```

2. **Create a new Railway project**:
   ```
   railway init
   ```

3. **Add a PostgreSQL database to your project**:
   - In the Railway dashboard, click "New"
   - Select "PostgreSQL" to add a database
   - Railway will automatically provide the `DATABASE_URL` environment variable

4. **Deploy the backend service**:
   ```
   cd backend
   railway up
   ```

5. **Deploy the frontend service**:
   ```
   cd frontend
   railway up
   ```

6. **Configure environment variables**:
   - In the Railway dashboard, go to each service
   - Set environment variables according to `.env.example`
   - For the frontend, set `NEXT_PUBLIC_API_URL` to the backend service URL

### Local Development vs Railway Deployment

- **Local Development**: Uses Docker Compose to run all services together (`docker-compose up`)
- **Railway Deployment**: Each service is deployed separately with its own configuration

This approach allows you to develop locally with a setup similar to production while benefiting from Railway's managed services when deployed. 