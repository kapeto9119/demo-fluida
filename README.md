# Fluida CTO Assignment: Invoice & Payment Link Generator

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
   - ~~Add authentication and authorization~~ ✅ Basic Auth Implemented
   - ~~Implement rate limiting~~ ✅ Implemented
   - ~~Add better error handling and validation~~ ✅ Implemented
   - Additional CSRF protection

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
   - ~~CI/CD pipeline~~ ✅ Implemented

## Recent Improvements

The following enhancements have been implemented:

- **Basic Authentication**: Simple username/password protection for the entire application
- **Standardized Error Handling**: Consistent error responses with proper status codes
- **Rate Limiting**: Protection against API abuse with proxy-aware client detection
- **Response Standardization**: Consistent JSON structure across all endpoints
- **Production Readiness**: Environment-specific configurations and Docker setup for Railway

## API-First Approach

This project follows API-first principles with:

- Resource-oriented design with predictable URLs and standard HTTP methods
- Versioned endpoints and backward compatibility
- Extensibility for future integrations with accounting software and payment processors

## Security Considerations

- **Wallet Security**: Private keys never leave the user's wallet
- **Data Protection**: Input validation and sanitization
- **API Security**: Rate limiting and error handling
- **Blockchain Security**: Transaction verification before confirmation

## Troubleshooting

Common issues:
- **Docker not running**: Use development mode without Docker
- **Connection issues**: Check logs and ensure PostgreSQL is running
- **Build errors**: See full documentation for specific error solutions

## Railway Deployment

The project is configured for Railway deployment as individual services:
- **Backend API**: Standalone Go service 
- **Frontend**: Next.js application
- **Database**: Railway's managed PostgreSQL

Deployment steps and environment variable configuration are provided in the Railway dashboard.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 