'use client'

import useCreateInvoice from '../../hooks/useCreateInvoice'
import InvoiceForm from '../../components/forms/InvoiceForm'
import Button from '../../components/ui/Button'

/**
 * Page component for creating new invoices
 */
export default function CreateInvoice() {
  const {
    formData,
    isLoading,
    createdInvoice,
    error,
    fieldErrors,
    hasUnsavedChanges,
    handleChange,
    handleSubmit,
    resetForm,
    getPaymentLink,
    clearSavedDraft
  } = useCreateInvoice()

  return (
    <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <h1 className="text-2xl font-bold mb-6">Create New Invoice</h1>
      
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}
      
      {/* Show draft notification */}
      {!createdInvoice && hasUnsavedChanges && (
        <div className="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded mb-6 flex justify-between items-center">
          <p>You have an unsaved draft invoice. Your changes are being saved automatically.</p>
          <button 
            onClick={clearSavedDraft} 
            className="text-sm underline hover:no-underline"
            aria-label="Clear draft"
          >
            Clear Draft
          </button>
        </div>
      )}
      
      {createdInvoice ? (
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">Invoice Created Successfully!</h2>
          <p className="mb-2"><strong>Invoice Number:</strong> {createdInvoice.invoiceNumber}</p>
          <p className="mb-2"><strong>Amount:</strong> {createdInvoice.amount} {createdInvoice.currency}</p>
          
          <div className="mt-6">
            <h3 className="text-lg font-medium mb-2">Payment Link:</h3>
            <div className="flex items-center">
              <input
                type="text"
                value={getPaymentLink()}
                readOnly
                className="flex-1 p-2 border rounded-l-md"
              />
              <button
                onClick={() => {
                  navigator.clipboard.writeText(getPaymentLink())
                  alert('Payment link copied to clipboard!')
                }}
                className="bg-primary-600 text-white px-4 py-2 rounded-r-md"
              >
                Copy
              </button>
            </div>
            
            {/* Debugging information */}
            <div className="mt-2 p-2 bg-gray-50 rounded text-xs text-gray-500">
              <p><strong>Debug - Link Token:</strong> {createdInvoice.linkToken || 'Not available'}</p>
              <p><strong>Debug - Generated URL:</strong> {getPaymentLink()}</p>
              <p><strong>Debug - Invoice ID:</strong> {createdInvoice.id}</p>
              <details>
                <summary className="cursor-pointer">Show raw invoice data</summary>
                <pre className="mt-1 overflow-auto max-h-32">
                  {JSON.stringify(createdInvoice, null, 2)}
                </pre>
              </details>
            </div>
            
            <div className="mt-4">
              <a
                href={getPaymentLink()}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-block bg-primary-600 text-white px-4 py-2 rounded-md"
              >
                View Payment Page
              </a>
            </div>
          </div>
          
          <div className="mt-6">
            <Button variant="outline" onClick={resetForm} ariaLabel="Create another invoice">
              Create Another Invoice
            </Button>
          </div>
        </div>
      ) : (
        <InvoiceForm
          formData={formData}
          handleChange={handleChange}
          handleSubmit={handleSubmit}
          isLoading={isLoading}
          fieldErrors={fieldErrors}
        />
      )}
    </div>
  )
} 