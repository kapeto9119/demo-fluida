'use client'

import { useState, useEffect } from 'react'
import { Invoice } from '../types'
import apiService from '../services/api'

/**
 * Custom hook for handling invoice payment functionality
 */
export const useInvoicePayment = (token: string) => {
  const [invoice, setInvoice] = useState<Invoice | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [paymentStatus, setPaymentStatus] = useState<'pending' | 'processing' | 'success' | 'error'>('pending')

  // Fetch invoice data by token
  useEffect(() => {
    const fetchInvoice = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await apiService.getInvoiceByToken(token)
        setInvoice(data)
        
        // If invoice is already paid, update status
        if (data.status === 'PAID') {
          setPaymentStatus('success')
        }
      } catch (err) {
        console.error('Error fetching invoice:', err)
        setError('Failed to load invoice. Please check the payment link and try again.')
      } finally {
        setLoading(false)
      }
    }

    if (token) {
      fetchInvoice()
    }
  }, [token])

  // Check payment status periodically during processing
  useEffect(() => {
    if (!invoice || paymentStatus !== 'processing') {
      return
    }

    const checkPaymentStatus = async () => {
      try {
        const data = await apiService.getInvoiceByToken(token)
        if (data.status === 'PAID') {
          setPaymentStatus('success')
          setInvoice(data)
        }
      } catch (err) {
        console.error('Error checking payment status:', err)
      }
    }

    // Check every 5 seconds while processing
    const interval = setInterval(checkPaymentStatus, 5000)
    return () => clearInterval(interval)
  }, [invoice, token, paymentStatus])

  // Format date for display
  const formatDate = (dateString: string) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString()
  }

  // Handle payment workflow
  const handlePaymentStart = () => {
    setPaymentStatus('processing')
  }

  const handlePaymentSuccess = async () => {
    try {
      // Update invoice status to PAID
      if (invoice) {
        await apiService.updateInvoiceStatus(invoice.id, 'PAID')
        setInvoice(prev => prev ? { ...prev, status: 'PAID' } : null)
      }
      setPaymentStatus('success')
    } catch (error) {
      console.error('Error updating invoice status:', error)
      setPaymentStatus('error')
    }
  }

  const handlePaymentError = () => {
    setPaymentStatus('error')
  }

  return {
    invoice,
    loading,
    error,
    paymentStatus,
    formatDate,
    handlePaymentStart,
    handlePaymentSuccess,
    handlePaymentError
  }
}

export default useInvoicePayment 