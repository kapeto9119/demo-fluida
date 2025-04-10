'use client'

import { useState, useEffect } from 'react'
import { InvoiceFormData, Invoice } from '../types'
import apiService from '../services/api'

/**
 * Initial form data with sensible defaults
 */
const initialFormData: InvoiceFormData = {
  invoiceNumber: '',
  amount: 0,
  currency: 'USDC',
  description: '',
  dueDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  receiverAddr: '',
  senderDetails: {
    name: '',
    email: '',
    address: ''
  },
  recipientDetails: {
    name: '',
    email: '',
    address: ''
  }
}

/**
 * Custom hook for creating invoices
 */
export const useCreateInvoice = () => {
  const [formData, setFormData] = useState<InvoiceFormData>(initialFormData);
  const [isLoading, setIsLoading] = useState(false)
  const [isSavingDraft, setIsSavingDraft] = useState(false)
  const [createdInvoice, setCreatedInvoice] = useState<Invoice | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [draftSaved, setDraftSaved] = useState(false)
  // Field-specific errors
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  // User ID for saving drafts - in a real app, this would come from auth
  const userId = "current_user"; // Placeholder - replace with actual user ID

  /**
   * Validate form data before submission
   * @returns An object with field errors, or empty object if valid
   */
  const validateForm = (): Record<string, string> => {
    const errors: Record<string, string> = {};
    
    // Validate invoice number
    if (!formData.invoiceNumber.trim()) {
      errors.invoiceNumber = 'Invoice number is required';
    } else if (formData.invoiceNumber.length > 50) {
      errors.invoiceNumber = 'Invoice number must be less than 50 characters';
    }
    
    // Validate amount
    if (!formData.amount || formData.amount <= 0) {
      errors.amount = 'Amount must be greater than zero';
    }
    
    // Validate receiver wallet address
    if (!formData.receiverAddr.trim()) {
      errors.receiverAddr = 'Receiver wallet address is required';
    } else if (formData.receiverAddr.length < 32 || formData.receiverAddr.length > 100) {
      errors.receiverAddr = 'Invalid Solana wallet address format';
    }
    
    // Validate due date
    const dueDate = new Date(formData.dueDate);
    const today = new Date();
    today.setHours(0, 0, 0, 0); // Reset time to start of day for fair comparison
    
    if (!formData.dueDate) {
      errors.dueDate = 'Due date is required';
    } else if (dueDate < today) {
      errors.dueDate = 'Due date cannot be in the past';
    }
    
    // Validate sender details
    if (!formData.senderDetails.name.trim()) {
      errors['senderDetails.name'] = 'Sender name is required';
    }
    
    if (!formData.senderDetails.email.trim()) {
      errors['senderDetails.email'] = 'Sender email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.senderDetails.email)) {
      errors['senderDetails.email'] = 'Invalid email format';
    }
    
    // Validate recipient details
    if (!formData.recipientDetails.name.trim()) {
      errors['recipientDetails.name'] = 'Recipient name is required';
    }
    
    if (!formData.recipientDetails.email.trim()) {
      errors['recipientDetails.email'] = 'Recipient email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.recipientDetails.email)) {
      errors['recipientDetails.email'] = 'Invalid email format';
    }
    
    return errors;
  };

  /**
   * Handle input changes, including nested properties
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    
    if (name.includes('.')) {
      // Handle nested properties (senderDetails.name, etc.)
      const [parent, child] = name.split('.')
      setFormData((prev: InvoiceFormData) => ({
        ...prev,
        [parent]: {
          ...(prev[parent as keyof InvoiceFormData] as any),
          [child]: value
        }
      }))
    } else {
      setFormData((prev: InvoiceFormData) => ({ ...prev, [name]: value }))
    }
    
    // Clear the draft saved state when form changes
    if (draftSaved) {
      setDraftSaved(false);
    }
  }

  /**
   * Check if invoice number already exists
   */
  const checkInvoiceNumberExists = async (invoiceNumber: string): Promise<boolean> => {
    try {
      return await apiService.checkInvoiceNumberExists(invoiceNumber);
    } catch (err) {
      console.error('Error checking invoice number:', err);
      return false;
    }
  };

  /**
   * Save the current form data as a draft to the database
   */
  const saveDraftToDatabase = async () => {
    if (!formData.invoiceNumber.trim()) {
      setError('Please enter an invoice number before saving draft');
      return;
    }

    setIsSavingDraft(true);
    try {
      // Store the form data as JSON string
      const invoiceDataJson = JSON.stringify(formData);
      
      // Save to database
      await apiService.saveDraftInvoice(userId, invoiceDataJson);
      
      setDraftSaved(true);
      setError(null);
    } catch (err) {
      console.error('Error saving draft invoice:', err);
      setError('Failed to save draft invoice to database');
    } finally {
      setIsSavingDraft(false);
    }
  };

  /**
   * Submit form data to create a new invoice
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)
    setFieldErrors({})
    
    // Validate form before submission
    const validationErrors = validateForm();
    if (Object.keys(validationErrors).length > 0) {
      setFieldErrors(validationErrors);
      setError('Please fix the errors in the form before submitting');
      setIsLoading(false);
      return;
    }
    
    // Check if invoice number already exists
    if (formData.invoiceNumber) {
      const exists = await checkInvoiceNumberExists(formData.invoiceNumber);
      if (exists) {
        setFieldErrors({
          invoiceNumber: 'This invoice number already exists. Please use a different one.'
        });
        setError('Invoice number already exists. Please use a different invoice number.');
        setIsLoading(false);
        return;
      }
    }
    
    try {
      const invoice = await apiService.createInvoice(formData)
      console.log('Created invoice:', invoice) // Debug info
      
      // Validate the invoice data
      if (!invoice) {
        console.error('Received null or undefined invoice from API')
        setError('Failed to create invoice: Invalid response from server')
        return
      }
      
      // Make sure the invoice has a linkToken
      if (!invoice.linkToken) {
        console.error('Missing linkToken in the invoice response:', invoice)
        
        // Try to look for the token in different properties
        const possibleTokenProps = ['linkToken', 'link_token', 'token', 'paymentToken', 'payment_token'];
        let foundToken = false;
        
        for (const prop of possibleTokenProps) {
          if ((invoice as any)[prop]) {
            console.log(`Found token in alternative property '${prop}':`, (invoice as any)[prop])
            // Create a new invoice object with the token in the right place
            const fixedInvoice = {
              ...invoice,
              linkToken: (invoice as any)[prop]
            };
            setCreatedInvoice(fixedInvoice)
            foundToken = true;
            break;
          }
        }
        
        if (!foundToken) {
          setError('Failed to create payment link: Missing link token in response')
          console.error('Could not find token in any known property')
        }
      } else {
        // Normal case - we have a valid invoice with a linkToken
        setCreatedInvoice(invoice)
      }
    } catch (err: any) {
      console.error('Error creating invoice:', err)
      
      // Check for specific error messages
      const errorResponse = err?.response?.data
      
      if (errorResponse) {
        // Check for duplicate invoice number error
        if (typeof errorResponse === 'string' && errorResponse.includes('already exists')) {
          setFieldErrors({
            invoiceNumber: 'This invoice number already exists. Please use a different one.'
          })
          setError('Invoice number already exists. Please use a different invoice number.')
        } else if (errorResponse.message) {
          setError(errorResponse.message)
        } else {
          setError('Failed to create invoice. Please try again.')
        }
      } else {
        setError('Failed to create invoice. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  /**
   * Reset form to create another invoice
   */
  const resetForm = () => {
    setFormData(initialFormData)
    setCreatedInvoice(null)
    setError(null)
    setFieldErrors({})
    setDraftSaved(false);
  }

  /**
   * Generate payment link from created invoice
   */
  const getPaymentLink = () => {
    if (!createdInvoice) {
      console.error('Cannot generate payment link: No invoice created')
      return `${window.location.origin}/pay/error`;
    }
    
    console.log('Generating payment link from invoice:', createdInvoice)
    
    // Try to get the linkToken directly
    const linkToken = createdInvoice.linkToken;
    
    if (!linkToken) {
      console.error('Missing linkToken in the created invoice:', createdInvoice);
      
      // Check if we can find the token in other properties like in the JSON data
      if (typeof createdInvoice === 'object') {
        // Try to extract from other possible properties
        for (const [key, value] of Object.entries(createdInvoice)) {
          if (key.toLowerCase().includes('token') && typeof value === 'string' && value.length > 10) {
            console.log(`Found possible token in property '${key}':`, value);
            return `${window.location.origin}/pay/${String(value).trim()}`;
          }
        }
        
        // If invoice has ID, we can use that as fallback
        if (createdInvoice.id) {
          console.log('Using invoice ID as fallback for link:', createdInvoice.id);
          return `${window.location.origin}/pay/invoice/${createdInvoice.id}`;
        }
      }
      
      return `${window.location.origin}/pay/error`;
    }
    
    // Ensure token is a string and clean it of any whitespace
    const cleanToken = String(linkToken).trim();
    console.log('Using token for payment link:', cleanToken);
    
    return `${window.location.origin}/pay/${cleanToken}`;
  }

  /**
   * Load a draft invoice from the database
   */
  const loadDraftFromDatabase = async () => {
    setIsLoading(true);
    try {
      const draft = await apiService.getDraftInvoice(userId);
      if (draft && draft.invoiceData) {
        // Parse the JSON data
        const parsedData = JSON.parse(draft.invoiceData) as InvoiceFormData;
        setFormData(parsedData);
      } else {
        setError('No draft invoice found');
      }
    } catch (err) {
      console.error('Error loading draft invoice:', err);
      setError('Failed to load draft invoice');
    } finally {
      setIsLoading(false);
    }
  };

  return {
    formData,
    isLoading,
    isSavingDraft,
    createdInvoice,
    error,
    fieldErrors,
    draftSaved,
    handleChange,
    handleSubmit,
    resetForm,
    getPaymentLink,
    saveDraftToDatabase,
    loadDraftFromDatabase
  }
}

export default useCreateInvoice 