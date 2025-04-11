import axios from 'axios'
import { Invoice, InvoiceFormData } from '../types'

// TEMPORARY FIX: Hardcode the production API URL
// const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
const API_URL = 'https://serene-radiance-production.up.railway.app'

// Get authentication credentials from localStorage if available, fallback to environment variables
const getAuthCredentials = () => {
  // Check if we're in a browser environment
  if (typeof window !== 'undefined') {
    const storedUsername = localStorage.getItem('auth_username')
    const storedPassword = localStorage.getItem('auth_password')
    
    if (storedUsername && storedPassword) {
      return {
        username: storedUsername,
        password: storedPassword
      }
    }
  }
  
  // Fallback to environment variables
  return {
    username: process.env.NEXT_PUBLIC_AUTH_USERNAME || 'admin',
    password: process.env.NEXT_PUBLIC_AUTH_PASSWORD || 'fluida'
  }
}

// Create basic auth token (in a way that works in both Node.js and browsers)
const createBasicAuthToken = (username: string, password: string) => {
  // For browser environments
  if (typeof window !== 'undefined' && window.btoa) {
    return window.btoa(`${username}:${password}`)
  }
  // For Node.js environments (during SSR)
  return Buffer.from(`${username}:${password}`).toString('base64')
}

const auth = getAuthCredentials()
const basicAuthToken = createBasicAuthToken(auth.username, auth.password)

// Create an axios instance with default settings
const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Basic ${basicAuthToken}`
  },
  withCredentials: true
})

// Add an interceptor to update auth token if credentials change
if (typeof window !== 'undefined') {
  // Check for credential changes every request
  api.interceptors.request.use(config => {
    const currentAuth = getAuthCredentials()
    const currentToken = createBasicAuthToken(currentAuth.username, currentAuth.password)
    
    // Update the Authorization header if token has changed
    if (config.headers && currentToken !== basicAuthToken) {
      config.headers.Authorization = `Basic ${currentToken}`
    }
    
    return config
  })
}

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
    // In case the token might need URL encoding
    const encodedToken = encodeURIComponent(token)
    
    // The API endpoint expects /{linkToken}
    const response = await api.get(`/invoices/${encodedToken}`)
    
    // Handle both wrapped and unwrapped responses
    return response.data.data || response.data
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
    
    // In production, we should not log the entire payload with sensitive data
    const response = await api.post('/invoices', dataToSubmit)
    
    // Return unwrapped data
    return response.data.data || response.data
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
  saveDraftInvoice: async (userId: string, draftData: any): Promise<any> => {
    try {
      const response = await api.post('/invoices/drafts', {
        UserID: userId,
        InvoiceData: draftData
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
        InvoiceData: invoiceData
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
    // In production, we should implement proper error logging
    // but avoid sensitive data exposure
    
    if (process.env.NODE_ENV !== 'production') {
      // Only log detailed errors in development
      if (error.response) {
        console.error('API Error:', {
          status: error.response.status,
          endpoint: error.config?.url
        })
      } else if (error.request) {
        console.error('API Request Error - No Response')
      } else {
        console.error('API Setup Error:', error.message)
      }
    }
    
    return Promise.reject(error)
  }
)

export default apiService 