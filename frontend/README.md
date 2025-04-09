# Fluida Invoice Generator Frontend

This is the frontend for the Fluida Invoice & Payment Link Generator. It's built with Next.js, TypeScript, Tailwind CSS, and integrates with Solana for blockchain-based payments.

## Key Features

- Create and manage invoices
- Generate payment links
- Connect to Solana wallets
- Process USDC payments on Solana testnet

## Development Setup

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser to see the application.

## Building for Production

```bash
npm run build
npm run start
```

## Key Dependencies

- **Next.js**: React framework for server-rendered applications
- **TypeScript**: Type-safe JavaScript
- **Tailwind CSS**: Utility-first CSS framework
- **Solana Web3.js**: Solana blockchain JavaScript API
- **Solana Wallet Adapter**: Connect to Solana wallets like Phantom
- **Axios**: HTTP client for API requests

## Project Structure

- `/src/app`: Next.js app directory containing page components
- `/src/components`: Reusable React components
- `/src/lib`: Utility functions and helpers
- `/src/types`: TypeScript types and declarations

## TypeScript Support

TypeScript configuration has been set up to provide a good developer experience. Key configuration files:

- `tsconfig.json`: TypeScript compiler configuration
- `src/types/global.d.ts`: Global TypeScript declarations
- `next-env.d.ts`: Next.js TypeScript declarations 