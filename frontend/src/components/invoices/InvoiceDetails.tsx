'use client'

import { Invoice } from '../../types'
import StatusBadge from '../ui/StatusBadge'

interface InvoiceDetailsProps {
  invoice: Invoice
}

/**
 * Component to display detailed invoice information
 */
export const InvoiceDetails = ({ invoice }: InvoiceDetailsProps) => {
  // Format date for display
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString()
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <div className="flex justify-between items-start mb-6">
        <div>
          <h2 className="text-xl font-semibold">Invoice #{invoice.invoiceNumber}</h2>
          <p className="text-gray-500 text-sm">Created on {formatDate(invoice.createdAt)}</p>
        </div>
        <StatusBadge status={invoice.status} />
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        {/* Sender Details */}
        <div className="border rounded-md p-4">
          <h3 className="text-sm font-medium text-gray-500 mb-2">From</h3>
          <p className="font-medium">{invoice.senderDetails.name}</p>
          <p className="text-gray-600">{invoice.senderDetails.email}</p>
          <p className="text-gray-600 whitespace-pre-line">{invoice.senderDetails.address}</p>
        </div>
        
        {/* Recipient Details */}
        <div className="border rounded-md p-4">
          <h3 className="text-sm font-medium text-gray-500 mb-2">To</h3>
          <p className="font-medium">{invoice.recipientDetails.name}</p>
          <p className="text-gray-600">{invoice.recipientDetails.email}</p>
          <p className="text-gray-600 whitespace-pre-line">{invoice.recipientDetails.address}</p>
        </div>
      </div>
      
      {/* Invoice Amount and Details */}
      <div className="border rounded-md p-4 mb-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-medium">Amount Due</h3>
          <p className="text-2xl font-bold">{invoice.amount} {invoice.currency}</p>
        </div>
        
        <div className="grid grid-cols-2 gap-4">
          <div>
            <h4 className="text-sm font-medium text-gray-500">Invoice Number</h4>
            <p>{invoice.invoiceNumber}</p>
          </div>
          <div>
            <h4 className="text-sm font-medium text-gray-500">Due Date</h4>
            <p>{formatDate(invoice.dueDate)}</p>
          </div>
          <div className="col-span-2">
            <h4 className="text-sm font-medium text-gray-500">Payment Address</h4>
            <p className="font-mono text-sm break-all">{invoice.receiverAddr}</p>
          </div>
        </div>
      </div>
      
      {/* Description */}
      {invoice.description && (
        <div className="border rounded-md p-4">
          <h3 className="text-sm font-medium text-gray-500 mb-2">Description</h3>
          <p className="whitespace-pre-line">{invoice.description}</p>
        </div>
      )}
    </div>
  )
}

export default InvoiceDetails 