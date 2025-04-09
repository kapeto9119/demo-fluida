'use client'

import { useState, useEffect } from 'react'
import { Invoice } from '../types'
import apiService from '../services/api'

/**
 * Custom hook for fetching and managing invoices
 */
export const useInvoices = () => {
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchInvoices = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await apiService.getInvoices()
      setInvoices(data)
    } catch (err) {
      console.error('Error fetching invoices:', err)
      setError('Failed to load invoices. Please try again later.')
    } finally {
      setLoading(false)
    }
  }

  // Fetch invoices on mount
  useEffect(() => {
    fetchInvoices()
  }, [])

  // Format date for display
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString()
  }

  return {
    invoices,
    loading,
    error,
    refetch: fetchInvoices,
    formatDate,
  }
}

export default useInvoices 