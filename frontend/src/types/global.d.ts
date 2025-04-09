// This file contains global type definitions to fix TypeScript JSX errors

declare global {
  namespace JSX {
    interface IntrinsicElements {
      [elemName: string]: any;
    }
  }

  // For Solana wallet integration
  interface Window {
    solana: any;
  }
} 