'use client'

import { InvoiceFormData } from '../../types'
import TextField from '../ui/TextField'
import TextArea from '../ui/TextArea'
import Button from '../ui/Button'

interface InvoiceFormProps {
  formData: InvoiceFormData
  handleChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void
  handleSubmit: (e: React.FormEvent) => void
  isLoading: boolean
}

/**
 * Form component for creating new invoices
 * Handles all input fields and form submission
 */
export const InvoiceForm = ({
  formData,
  handleChange,
  handleSubmit,
  isLoading
}: InvoiceFormProps) => {
  return (
    <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-md space-y-6">
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
        {/* Invoice Details */}
        <div className="space-y-4 col-span-2">
          <h2 className="text-lg font-medium">Invoice Details</h2>
          
          <TextField
            id="invoiceNumber"
            name="invoiceNumber"
            label="Invoice Number"
            value={formData.invoiceNumber}
            onChange={handleChange}
            required
          />
          
          <div className="grid grid-cols-2 gap-4">
            <TextField
              id="amount"
              name="amount"
              label="Amount"
              value={formData.amount}
              onChange={handleChange}
              type="number"
              min={0}
              step="0.01"
              required
            />
            
            <TextField
              id="currency"
              name="currency"
              label="Currency"
              value={formData.currency}
              onChange={handleChange}
              readOnly
            />
          </div>
          
          <TextArea
            id="description"
            name="description"
            label="Description"
            value={formData.description}
            onChange={handleChange}
            rows={3}
          />
          
          <div className="grid grid-cols-2 gap-4">
            <TextField
              id="dueDate"
              name="dueDate"
              label="Due Date"
              value={formData.dueDate}
              onChange={handleChange}
              type="date"
              required
            />
            
            <TextField
              id="receiverAddr"
              name="receiverAddr"
              label="Receiver Wallet Address"
              value={formData.receiverAddr}
              onChange={handleChange}
              placeholder="Solana wallet address to receive payment"
              required
            />
          </div>
        </div>
        
        {/* Sender Details */}
        <div className="space-y-4">
          <h2 className="text-lg font-medium">Sender Details</h2>
          
          <TextField
            id="senderDetails.name"
            name="senderDetails.name"
            label="Your Name"
            value={formData.senderDetails.name}
            onChange={handleChange}
            required
          />
          
          <TextField
            id="senderDetails.email"
            name="senderDetails.email"
            label="Your Email"
            value={formData.senderDetails.email}
            onChange={handleChange}
            type="email"
            required
          />
          
          <TextArea
            id="senderDetails.address"
            name="senderDetails.address"
            label="Your Address"
            value={formData.senderDetails.address}
            onChange={handleChange}
            rows={3}
          />
        </div>
        
        {/* Recipient Details */}
        <div className="space-y-4">
          <h2 className="text-lg font-medium">Recipient Details</h2>
          
          <TextField
            id="recipientDetails.name"
            name="recipientDetails.name"
            label="Recipient Name"
            value={formData.recipientDetails.name}
            onChange={handleChange}
            required
          />
          
          <TextField
            id="recipientDetails.email"
            name="recipientDetails.email"
            label="Recipient Email"
            value={formData.recipientDetails.email}
            onChange={handleChange}
            type="email"
            required
          />
          
          <TextArea
            id="recipientDetails.address"
            name="recipientDetails.address"
            label="Recipient Address"
            value={formData.recipientDetails.address}
            onChange={handleChange}
            rows={3}
          />
        </div>
      </div>
      
      <div className="flex justify-end">
        <Button
          type="submit"
          disabled={isLoading}
          ariaLabel="Create invoice"
        >
          {isLoading ? 'Creating...' : 'Create Invoice'}
        </Button>
      </div>
    </form>
  )
}

export default InvoiceForm 