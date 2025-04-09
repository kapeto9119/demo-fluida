'use client'

import { FC, useState } from 'react'
import { Invoice } from '../types'
import Button from './ui/Button'

// Interface for component props
interface WalletComponentsProps {
  invoice: Invoice
  onPaymentStart: () => void
  onPaymentSuccess: () => void
  onPaymentError: () => void
}

/**
 * WalletComponents provides wallet connection and payment functionality
 * This is a simplified implementation for development purposes
 */
const WalletComponents: FC<WalletComponentsProps> = ({ 
  invoice, 
  onPaymentStart, 
  onPaymentSuccess, 
  onPaymentError 
}) => {
  const [isProcessing, setIsProcessing] = useState(false)
  const [isWalletConnected, setIsWalletConnected] = useState(false)
  
  /**
   * Simulate connecting to a wallet
   */
  const handleConnectWallet = () => {
    // Simulating wallet connection
    setIsWalletConnected(true)
  }
  
  /**
   * Disconnect from wallet
   */
  const handleDisconnectWallet = () => {
    setIsWalletConnected(false)
  }
  
  /**
   * Process payment when user clicks Pay button
   */
  const handlePayment = async () => {
    if (!isWalletConnected) {
      return
    }
    
    try {
      setIsProcessing(true)
      onPaymentStart()
      
      // Simulate payment processing with a delay
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      onPaymentSuccess()
    } catch (error) {
      console.error('Payment simulation error:', error)
      onPaymentError()
    } finally {
      setIsProcessing(false)
    }
  }
  
  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
        {!isWalletConnected ? (
          <Button
            onClick={handleConnectWallet}
            variant="secondary"
            ariaLabel="Connect wallet"
          >
            Connect Wallet
          </Button>
        ) : (
          <div className="flex flex-col sm:flex-row gap-4 items-center">
            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              <span className="w-2 h-2 mr-1 bg-green-400 rounded-full"></span>
              Wallet Connected
            </span>
            
            <button
              onClick={handleDisconnectWallet}
              className="text-sm text-gray-600 hover:text-gray-900"
              aria-label="Disconnect wallet"
            >
              Disconnect
            </button>
          </div>
        )}
        
        <Button
          onClick={handlePayment}
          disabled={isProcessing || !isWalletConnected}
          ariaLabel={`Pay ${invoice.amount} ${invoice.currency}`}
        >
          {isProcessing ? 'Processing...' : `Pay ${invoice.amount} ${invoice.currency}`}
        </Button>
      </div>
      
      <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-md">
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3">
            <h3 className="text-sm font-medium text-yellow-800">Development Mode</h3>
            <div className="mt-1 text-sm text-yellow-700">
              <p>
                This is a simplified implementation for development. In production, this component would integrate with real Solana wallets and process actual transactions on the blockchain.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default WalletComponents 