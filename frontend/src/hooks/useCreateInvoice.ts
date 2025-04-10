'use client'

import { useState } from 'react'
import { InvoiceFormData, Invoice } from '../types'
import apiService from '../services/api'

/**
 * Initial form data with sensible defaults
 */
const initialFormData: InvoiceFormData = {
  invoiceNumber: '',
  amount: 0,
  currency: 'USDC',
  description: '',
  dueDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  receiverAddr: '',
  senderDetails: {
    name: '',
    email: '',
    address: ''
  },
  recipientDetails: {
    name: '',
    email: '',
    address: ''
  }
}

/**
 * Custom hook for creating invoices
 */
export const useCreateInvoice = () => {
  const [formData, setFormData] = useState<InvoiceFormData>(initialFormData)
  const [isLoading, setIsLoading] = useState(false)
  const [createdInvoice, setCreatedInvoice] = useState<Invoice | null>(null)
  const [error, setError] = useState<string | null>(null)
  // Field-specific errors
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})

  /**
   * Handle input changes, including nested properties
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    
    if (name.includes('.')) {
      // Handle nested properties (senderDetails.name, etc.)
      const [parent, child] = name.split('.')
      setFormData((prev: InvoiceFormData) => ({
        ...prev,
        [parent]: {
          ...(prev[parent as keyof InvoiceFormData] as any),
          [child]: value
        }
      }))
    } else {
      setFormData((prev: InvoiceFormData) => ({ ...prev, [name]: value }))
    }
  }

  /**
   * Submit form data to create a new invoice
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)
    setFieldErrors({})
    
    try {
      const invoice = await apiService.createInvoice(formData)
      console.log('Created invoice:', invoice) // Debug info
      
      // Validate the invoice data
      if (!invoice) {
        console.error('Received null or undefined invoice from API')
        setError('Failed to create invoice: Invalid response from server')
        return
      }
      
      // Make sure the invoice has a linkToken
      if (!invoice.linkToken) {
        console.error('Missing linkToken in the invoice response:', invoice)
        
        // Try to look for the token in different properties
        const possibleTokenProps = ['linkToken', 'link_token', 'token', 'paymentToken', 'payment_token'];
        let foundToken = false;
        
        for (const prop of possibleTokenProps) {
          if ((invoice as any)[prop]) {
            console.log(`Found token in alternative property '${prop}':`, (invoice as any)[prop])
            // Create a new invoice object with the token in the right place
            const fixedInvoice = {
              ...invoice,
              linkToken: (invoice as any)[prop]
            };
            setCreatedInvoice(fixedInvoice)
            foundToken = true;
            break;
          }
        }
        
        if (!foundToken) {
          setError('Failed to create payment link: Missing link token in response')
          console.error('Could not find token in any known property')
        }
      } else {
        // Normal case - we have a valid invoice with a linkToken
        setCreatedInvoice(invoice)
      }
    } catch (err: any) {
      console.error('Error creating invoice:', err)
      
      // Check for specific error messages
      const errorResponse = err?.response?.data
      
      if (errorResponse) {
        // Check for duplicate invoice number error
        if (typeof errorResponse === 'string' && errorResponse.includes('already exists')) {
          setFieldErrors({
            invoiceNumber: 'This invoice number already exists. Please use a different one.'
          })
          setError('Invoice number already exists. Please use a different invoice number.')
        } else if (errorResponse.message) {
          setError(errorResponse.message)
        } else {
          setError('Failed to create invoice. Please try again.')
        }
      } else {
        setError('Failed to create invoice. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  /**
   * Reset form to create another invoice
   */
  const resetForm = () => {
    setFormData(initialFormData)
    setCreatedInvoice(null)
    setError(null)
  }

  /**
   * Generate payment link from created invoice
   */
  const getPaymentLink = () => {
    if (!createdInvoice) {
      console.error('Cannot generate payment link: No invoice created')
      return `${window.location.origin}/pay/error`;
    }
    
    console.log('Generating payment link from invoice:', createdInvoice)
    
    // Try to get the linkToken directly
    const linkToken = createdInvoice.linkToken;
    
    if (!linkToken) {
      console.error('Missing linkToken in the created invoice:', createdInvoice);
      
      // Check if we can find the token in other properties like in the JSON data
      if (typeof createdInvoice === 'object') {
        // Try to extract from other possible properties
        for (const [key, value] of Object.entries(createdInvoice)) {
          if (key.toLowerCase().includes('token') && typeof value === 'string' && value.length > 10) {
            console.log(`Found possible token in property '${key}':`, value);
            return `${window.location.origin}/pay/${String(value).trim()}`;
          }
        }
        
        // If invoice has ID, we can use that as fallback
        if (createdInvoice.id) {
          console.log('Using invoice ID as fallback for link:', createdInvoice.id);
          return `${window.location.origin}/pay/invoice/${createdInvoice.id}`;
        }
      }
      
      return `${window.location.origin}/pay/error`;
    }
    
    // Ensure token is a string and clean it of any whitespace
    const cleanToken = String(linkToken).trim();
    console.log('Using token for payment link:', cleanToken);
    
    return `${window.location.origin}/pay/${cleanToken}`;
  }

  return {
    formData,
    isLoading,
    createdInvoice,
    error,
    fieldErrors,
    handleChange,
    handleSubmit,
    resetForm,
    getPaymentLink,
  }
}

export default useCreateInvoice 