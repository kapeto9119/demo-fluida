import axios from 'axios'
import { Invoice, InvoiceFormData } from '../types'

// Get the API URL from environment or use localhost in development
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// Create an axios instance with default settings
const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
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
    try {
      const response = await api.get('/invoices')
      // Handle both wrapped and unwrapped responses
      return response.data.data || response.data
    } catch (error) {
      console.error('Error fetching invoices:', error)
      throw error
    }
  },

  /**
   * Get a single invoice by ID
   */
  getInvoice: async (id: number): Promise<Invoice> => {
    try {
      const response = await api.get(`/invoices/${id}`)
      // Handle both wrapped and unwrapped responses
      return response.data.data || response.data
    } catch (error) {
      console.error(`Error fetching invoice #${id}:`, error)
      throw error
    }
  },

  /**
   * Get an invoice by token
   */
  getInvoiceByToken: async (token: string): Promise<Invoice> => {
    try {
      console.log(`Fetching invoice with token: ${token}`)
      // In case the token might need URL encoding
      const encodedToken = encodeURIComponent(token)
      console.log(`Encoded token: ${encodedToken}`)
      
      // The API endpoint expects /{linkToken}
      const response = await api.get(`/invoices/${encodedToken}`)
      
      console.log('API Response for token lookup:', response)
      // Handle both wrapped and unwrapped responses
      return response.data.data || response.data
    } catch (error) {
      console.error(`Error fetching invoice by token ${token}:`, error)
      throw error
    }
  },

  /**
   * Create a new invoice
   */
  createInvoice: async (invoiceData: InvoiceFormData): Promise<Invoice> => {
    try {
      // Format the data before sending
      const dataToSubmit = {
        ...invoiceData,
        amount: parseFloat(invoiceData.amount.toString()),
        dueDate: new Date(invoiceData.dueDate).toISOString(),
      }
      
      console.log('Sending invoice data to backend:', dataToSubmit)
      const response = await api.post('/invoices', dataToSubmit)
      
      // Log the response to see its structure
      console.log('API Response:', response)
      
      // Check if the response is wrapped in a 'data' property (standard backend wrapper)
      const responseData = response.data.data || response.data;
      
      console.log('Unwrapped invoice data:', responseData)
      console.log('Link token from unwrapped data:', responseData.linkToken)
      
      // Return the unwrapped data
      return responseData
    } catch (error: any) {
      console.error('Error creating invoice:', error)
      
      // Check for specific API error messages in the response
      if (error.response) {
        const errorData = error.response.data;
        
        // Log detailed error information
        console.error('API Error Response:', {
          status: error.response.status,
          data: errorData,
          statusText: error.response.statusText
        })
        
        // If the error response contains a message about duplicate invoice
        if (typeof errorData === 'string' && errorData.includes('already exists')) {
          throw new Error(errorData);
        }
      }
      
      throw error;
    }
  },

  /**
   * Update invoice status
   */
  updateInvoiceStatus: async (id: number, status: string): Promise<Invoice> => {
    try {
      const response = await api.put(`/invoices/${id}/status`, { status })
      // Handle both wrapped and unwrapped responses
      return response.data.data || response.data
    } catch (error) {
      console.error(`Error updating invoice #${id} status:`, error)
      throw error
    }
  },

  /**
   * Check if an invoice number already exists
   */
  checkInvoiceNumberExists: async (invoiceNumber: string): Promise<boolean> => {
    try {
      const response = await api.get(`/invoices/check?invoice_number=${encodeURIComponent(invoiceNumber)}`)
      // Return the exists flag from the response
      return response.data.exists || false
    } catch (error) {
      console.error(`Error checking if invoice number ${invoiceNumber} exists:`, error)
      // Default to false in case of error
      return false
    }
  },

  /**
   * Save draft invoice
   */
  saveDraftInvoice: async (userId: string, invoiceData: string): Promise<any> => {
    try {
      const response = await api.post('/invoices/drafts', {
        userId,
        invoiceData
      })
      return response.data.data || response.data
    } catch (error) {
      console.error('Error saving draft invoice:', error)
      throw error
    }
  },

  /**
   * Get draft invoice by user ID
   */
  getDraftInvoice: async (userId: string): Promise<any> => {
    try {
      const response = await api.get(`/invoices/drafts/${userId}`)
      return response.data.data || response.data
    } catch (error: any) {
      // If not found, return null instead of throwing
      if (error.response && error.response.status === 404) {
        return null
      }
      console.error(`Error fetching draft invoice for user ${userId}:`, error)
      throw error
    }
  },

  /**
   * Update draft invoice
   */
  updateDraftInvoice: async (id: string, invoiceData: string): Promise<any> => {
    try {
      const response = await api.put(`/invoices/drafts/${id}`, {
        invoiceData
      })
      return response.data.data || response.data
    } catch (error) {
      console.error(`Error updating draft invoice ${id}:`, error)
      throw error
    }
  },

  /**
   * Delete draft invoice
   */
  deleteDraftInvoice: async (id: string): Promise<any> => {
    try {
      const response = await api.delete(`/invoices/drafts/${id}`)
      return response.data
    } catch (error) {
      console.error(`Error deleting draft invoice ${id}:`, error)
      throw error
    }
  }
}

// Add an interceptor for logging and handling errors
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('API Error Response:', {
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers,
      })
    } else if (error.request) {
      // The request was made but no response was received
      console.error('API Error Request:', error.request)
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('API Error Setup:', error.message)
    }
    return Promise.reject(error)
  }
)

export default apiService 