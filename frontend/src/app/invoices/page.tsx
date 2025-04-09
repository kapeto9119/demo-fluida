'use client'

import Link from 'next/link'
import useInvoices from '../../hooks/useInvoices'
import InvoiceTable from '../../components/invoices/InvoiceTable'
import LoadingSpinner from '../../components/ui/LoadingSpinner'
import Button from '../../components/ui/Button'

/**
 * Page component for displaying all invoices
 */
export default function InvoicesList() {
  const { invoices, loading, error, formatDate } = useInvoices()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner text="Loading invoices..." />
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">All Invoices</h1>
        <Link href="/create-invoice">
          <Button ariaLabel="Create new invoice">Create New Invoice</Button>
        </Link>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {invoices.length === 0 ? (
        <div className="bg-white rounded-lg shadow-md p-6 text-center">
          <p className="text-gray-600">No invoices found.</p>
          <Link href="/create-invoice" className="text-primary-600 underline mt-2 inline-block">
            Create your first invoice
          </Link>
        </div>
      ) : (
        <InvoiceTable invoices={invoices} formatDate={formatDate} />
      )}
    </div>
  )
} 