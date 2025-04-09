'use client'

import { useParams } from 'next/navigation'
import Link from 'next/link'
import dynamic from 'next/dynamic'
import useInvoicePayment from '../../../hooks/useInvoicePayment'
import LoadingSpinner from '../../../components/ui/LoadingSpinner'
import StatusBadge from '../../../components/ui/StatusBadge'

// Dynamic imports for wallet components to avoid SSR issues
const WalletComponents = dynamic(
  () => import('@/components/WalletComponents'),
  { 
    ssr: false,
    loading: () => <div className="p-4 text-center">Loading wallet components...</div>
  }
)

/**
 * Page component for processing payments for an invoice
 */
export default function PaymentPage() {
  const params = useParams()
  const { token } = params as { token: string }
  
  const {
    invoice,
    loading,
    error,
    paymentStatus,
    formatDate,
    handlePaymentStart,
    handlePaymentSuccess,
    handlePaymentError
  } = useInvoicePayment(token as string)

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner text="Loading invoice..." />
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-lg mx-auto mt-12 p-6 bg-white rounded-lg shadow-md">
        <h1 className="text-xl font-semibold text-red-600 mb-4">Error</h1>
        <p>{error}</p>
        <div className="mt-6">
          <Link href="/" className="text-primary-600 hover:underline">Return to Home</Link>
        </div>
      </div>
    )
  }

  if (!invoice) {
    return (
      <div className="max-w-lg mx-auto mt-12 p-6 bg-white rounded-lg shadow-md">
        <h1 className="text-xl font-semibold text-red-600 mb-4">Invoice Not Found</h1>
        <p>The invoice you're looking for doesn't exist or has been removed.</p>
        <div className="mt-6">
          <Link href="/" className="text-primary-600 hover:underline">Return to Home</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {/* Invoice Header */}
        <div className="border-b border-gray-200 bg-gray-50 px-6 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-xl font-semibold text-gray-900">Invoice #{invoice.invoiceNumber}</h1>
            <StatusBadge status={invoice.status} />
          </div>
        </div>

        {/* Invoice Details */}
        <div className="px-6 py-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Left Column */}
            <div>
              <h2 className="text-lg font-medium">From</h2>
              <div className="mt-2">
                <p className="text-sm text-gray-800">{invoice.senderDetails.name}</p>
                <p className="text-sm text-gray-600">{invoice.senderDetails.email}</p>
                <p className="text-sm text-gray-600 whitespace-pre-line">{invoice.senderDetails.address}</p>
              </div>

              <h2 className="text-lg font-medium mt-6">To</h2>
              <div className="mt-2">
                <p className="text-sm text-gray-800">{invoice.recipientDetails.name}</p>
                <p className="text-sm text-gray-600">{invoice.recipientDetails.email}</p>
                <p className="text-sm text-gray-600 whitespace-pre-line">{invoice.recipientDetails.address}</p>
              </div>
            </div>

            {/* Right Column */}
            <div>
              <div className="bg-gray-50 p-4 rounded-md">
                <h2 className="text-lg font-medium">Payment Details</h2>
                <div className="mt-4 space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount Due:</span>
                    <span className="font-medium">{invoice.amount} {invoice.currency}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Due Date:</span>
                    <span>{formatDate(invoice.dueDate)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Payment Status:</span>
                    <span className={`font-medium ${
                      paymentStatus === 'success' ? 'text-green-600' : 
                      paymentStatus === 'error' ? 'text-red-600' : 'text-yellow-600'
                    }`}>
                      {paymentStatus === 'pending' && 'Awaiting Payment'}
                      {paymentStatus === 'processing' && 'Processing...'}
                      {paymentStatus === 'success' && 'Paid'}
                      {paymentStatus === 'error' && 'Failed'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Description */}
              {invoice.description && (
                <div className="mt-6">
                  <h2 className="text-lg font-medium">Description</h2>
                  <p className="mt-2 text-sm text-gray-600 whitespace-pre-line">{invoice.description}</p>
                </div>
              )}
            </div>
          </div>

          {/* Payment Section */}
          {paymentStatus !== 'success' && (
            <div className="mt-8 border-t border-gray-200 pt-8">
              <h2 className="text-lg font-medium">Pay with Solana</h2>
              <p className="mt-2 text-sm text-gray-600">
                Connect your wallet and pay {invoice.amount} {invoice.currency} to address:
              </p>
              <div className="mt-2 bg-gray-50 p-3 rounded-md text-sm font-mono break-all">
                {invoice.receiverAddr}
              </div>

              {/* This is where the Solana wallet components will be rendered */}
              <div className="mt-6">
                <WalletComponents 
                  invoice={invoice}
                  onPaymentStart={handlePaymentStart}
                  onPaymentSuccess={handlePaymentSuccess}
                  onPaymentError={handlePaymentError}
                />
              </div>
            </div>
          )}

          {/* Success Message */}
          {paymentStatus === 'success' && (
            <div className="mt-8 border-t border-gray-200 pt-8 text-center">
              <div className="rounded-full bg-green-100 p-3 mx-auto w-16 h-16 flex items-center justify-center">
                <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h2 className="mt-4 text-xl font-medium text-green-600">Payment Complete!</h2>
              <p className="mt-2 text-gray-600">
                Thank you for your payment. The invoice has been marked as paid.
              </p>
              <div className="mt-6">
                <Link href="/" className="text-primary-600 hover:underline">Return to Home</Link>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
} 