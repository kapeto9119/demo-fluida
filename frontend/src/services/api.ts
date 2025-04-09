import axios from 'axios'
import { Invoice, InvoiceFormData } from '../types'

// Create an axios instance with default settings
const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

/**
 * API service to handle all API calls in a centralized location
 */
export const apiService = {
  /**
   * Fetch all invoices
   */
  getInvoices: async (): Promise<Invoice[]> => {
    const response = await api.get('/invoices')
    return response.data
  },

  /**
   * Get a single invoice by ID
   */
  getInvoice: async (id: number): Promise<Invoice> => {
    const response = await api.get(`/invoices/${id}`)
    return response.data
  },

  /**
   * Get an invoice by token
   */
  getInvoiceByToken: async (token: string): Promise<Invoice> => {
    const response = await api.get(`/invoices/token/${token}`)
    return response.data
  },

  /**
   * Create a new invoice
   */
  createInvoice: async (invoiceData: InvoiceFormData): Promise<Invoice> => {
    // Format the data before sending
    const dataToSubmit = {
      ...invoiceData,
      amount: parseFloat(invoiceData.amount.toString()),
      dueDate: new Date(invoiceData.dueDate).toISOString(),
    }
    
    const response = await api.post('/invoices', dataToSubmit)
    return response.data
  },

  /**
   * Update invoice status
   */
  updateInvoiceStatus: async (id: number, status: string): Promise<Invoice> => {
    const response = await api.patch(`/invoices/${id}/status`, { status })
    return response.data
  },
}

export default apiService 