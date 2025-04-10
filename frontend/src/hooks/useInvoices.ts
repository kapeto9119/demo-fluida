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
      
      // Ensure we always set an array to state
      if (Array.isArray(data)) {
        setInvoices(data)
      } else {
        console.error('API returned non-array data:', data)
        setInvoices([])
        setError('Invalid response format from server')
      }
    } catch (err) {
      console.error('Error fetching invoices:', err)
      setError('Failed to load invoices. Please try again later.')
      setInvoices([]) // Ensure invoices is an empty array on error
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
    try {
      const date = new Date(dateString)
      return date.toLocaleDateString()
    } catch (err) {
      console.error('Error formatting date:', err)
      return 'Invalid date'
    }
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