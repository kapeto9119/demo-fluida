'use client'

import { Invoice } from '../../types'
import StatusBadge from '../ui/StatusBadge'

interface InvoiceTableProps {
  invoices: Invoice[]
  formatDate: (dateString: string) => string
}

/**
 * Component to display invoices in a table format
 */
export const InvoiceTable = ({ invoices, formatDate }: InvoiceTableProps) => {
  // Make sure invoices is an array
  const invoiceList = Array.isArray(invoices) ? invoices : [];
  
  return (
    <div className="card">
      <div className="table-responsive">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Invoice #
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Client
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Amount
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Due Date
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {invoiceList.map((invoice) => (
              <tr key={invoice.id} className="hover:bg-gray-50 transition-colors duration-150">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">{invoice.invoiceNumber}</div>
                  <div className="text-xs text-gray-500">{formatDate(invoice.createdAt)}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-gray-900">{invoice.recipientDetails.name}</div>
                  <div className="text-xs text-gray-500">{invoice.recipientDetails.email}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">{invoice.amount} {invoice.currency}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-gray-900">{formatDate(invoice.dueDate)}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <StatusBadge status={invoice.status} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <a
                    href={`/pay/${invoice.linkToken}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary-600 hover:text-primary-900 hover:underline"
                  >
                    View Payment Link
                  </a>
                </td>
              </tr>
            ))}
            {invoiceList.length === 0 && (
              <tr>
                <td colSpan={6} className="px-6 py-4 text-center text-gray-500">
                  No invoices found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      
      {/* Alternative mobile view - only shows below sm breakpoint */}
      <div className="sm:hidden mt-4">
        {invoiceList.length === 0 ? (
          <p className="text-center text-gray-500 py-4">No invoices found</p>
        ) : (
          <div className="space-y-4 px-4 pb-4">
            {invoiceList.map((invoice) => (
              <div key={invoice.id} className="bg-white border rounded-lg p-4 shadow-sm">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <div className="font-medium">{invoice.invoiceNumber}</div>
                    <div className="text-xs text-gray-500">{formatDate(invoice.createdAt)}</div>
                  </div>
                  <StatusBadge status={invoice.status} />
                </div>
                
                <div className="mb-2">
                  <div className="text-sm">{invoice.recipientDetails.name}</div>
                  <div className="text-xs text-gray-500">{invoice.recipientDetails.email}</div>
                </div>
                
                <div className="flex justify-between items-center text-sm mb-3">
                  <div>
                    <span className="text-gray-500">Amount:</span> {invoice.amount} {invoice.currency}
                  </div>
                  <div>
                    <span className="text-gray-500">Due:</span> {formatDate(invoice.dueDate)}
                  </div>
                </div>
                
                <a
                  href={`/pay/${invoice.linkToken}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block text-center text-primary-600 border border-primary-600 rounded-md py-1 px-2 text-sm hover:bg-primary-50"
                >
                  View Payment Link
                </a>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default InvoiceTable 