/**
 * Shared type definitions for the Fluida application
 */

// Invoice type used across all components
export interface Invoice {
  id: number
  invoiceNumber: string
  amount: number
  currency: string
  description: string
  dueDate: string
  receiverAddr: string
  status: string
  linkToken: string
  senderDetails: {
    name: string
    email: string
    address: string
  }
  recipientDetails: {
    name: string
    email: string
    address: string
  }
  createdAt: string
}

// Form data for creating a new invoice
export interface InvoiceFormData {
  invoiceNumber: string
  amount: number
  currency: string
  description: string
  dueDate: string
  receiverAddr: string
  senderDetails: {
    name: string
    email: string
    address: string
  }
  recipientDetails: {
    name: string
    email: string
    address: string
  }
} 