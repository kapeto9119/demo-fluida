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
    
    try {
      const invoice = await apiService.createInvoice(formData)
      setCreatedInvoice(invoice)
    } catch (error) {
      console.error('Error creating invoice:', error)
      setError('Failed to create invoice. Please try again.')
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
    if (!createdInvoice) return ''
    return `${window.location.origin}/pay/${createdInvoice.linkToken}`
  }

  return {
    formData,
    isLoading,
    createdInvoice,
    error,
    handleChange,
    handleSubmit,
    resetForm,
    getPaymentLink,
  }
}

export default useCreateInvoice 