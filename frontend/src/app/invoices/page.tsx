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
  
  // Ensure invoices is always an array
  const invoicesList = Array.isArray(invoices) ? invoices : [];

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner text="Loading invoices..." />
      </div>
    )
  }

  return (
    <div>
      <div className="flex flex-col sm:flex-row justify-between items-center mb-8 gap-4">
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

      {invoicesList.length === 0 ? (
        <div className="card card-body text-center py-12">
          <p className="text-gray-600 mb-4">No invoices found.</p>
          <Link href="/create-invoice">
            <Button variant="primary" ariaLabel="Create your first invoice">
              Create your first invoice
            </Button>
          </Link>
        </div>
      ) : (
        <InvoiceTable invoices={invoicesList} formatDate={formatDate} />
      )}
    </div>
  )
} 