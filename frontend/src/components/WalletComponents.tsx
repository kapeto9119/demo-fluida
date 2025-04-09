'use client'

import { FC, useState } from 'react'
import { useWallet } from '@solana/wallet-adapter-react'
import { WalletMultiButton } from '@solana/wallet-adapter-react-ui'
import { LAMPORTS_PER_SOL, PublicKey, Transaction, SystemProgram, Connection, clusterApiUrl } from '@solana/web3.js'
import { createTransferCheckedInstruction, getAssociatedTokenAddressSync, getAccount, getMint } from '@solana/spl-token'
import { Invoice } from '../types'
import Button from './ui/Button'

// Interface for component props
interface WalletComponentsProps {
  invoice: Invoice
  onPaymentStart: () => void
  onPaymentSuccess: () => void
  onPaymentError: () => void
}

// USDC mint address on devnet
const USDC_MINT = new PublicKey('4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU')
// USDC decimals
const USDC_DECIMALS = 6

/**
 * WalletComponents provides wallet connection and payment functionality
 */
const WalletComponents: FC<WalletComponentsProps> = ({ 
  invoice, 
  onPaymentStart, 
  onPaymentSuccess, 
  onPaymentError 
}) => {
  const [isProcessing, setIsProcessing] = useState(false)
  
  // Use wallet hook from wallet adapter
  const { publicKey, sendTransaction, connected } = useWallet()
  
  /**
   * Process payment when user clicks Pay button
   */
  const handlePayment = async () => {
    if (!connected || !publicKey) {
      return
    }
    
    try {
      setIsProcessing(true)
      onPaymentStart()
      
      // Create connection to devnet
      const connection = new Connection(clusterApiUrl('devnet'), 'confirmed')
      
      // Get receiver's public key
      const receiverWallet = new PublicKey(invoice.receiverAddr)
      
      // Calculate amount with decimals
      const amount = Math.round(invoice.amount * (10 ** USDC_DECIMALS))
      
      // Get the associated token accounts
      const senderTokenAccount = getAssociatedTokenAddressSync(USDC_MINT, publicKey)
      const receiverTokenAccount = getAssociatedTokenAddressSync(USDC_MINT, receiverWallet)
      
      // Get token mint info
      const mintInfo = await getMint(connection, USDC_MINT)
      
      // Create transfer instruction
      const transferInstruction = createTransferCheckedInstruction(
        senderTokenAccount,
        USDC_MINT,
        receiverTokenAccount,
        publicKey,
        amount,
        mintInfo.decimals
      )
      
      // Add a memo with invoice ID
      transferInstruction.keys.push({
        pubkey: new PublicKey(invoice.invoiceNumber),
        isSigner: false,
        isWritable: false
      })
      
      // Create transaction
      const transaction = new Transaction().add(transferInstruction)
      
      // Set recent blockhash and fee payer
      transaction.feePayer = publicKey
      transaction.recentBlockhash = (await connection.getLatestBlockhash()).blockhash
      
      // Send transaction
      const signature = await sendTransaction(transaction, connection)
      
      // Wait for confirmation
      const confirmation = await connection.confirmTransaction(signature, 'confirmed')
      
      if (confirmation.value.err) {
        throw new Error('Transaction failed')
      }
      
      console.log(`Transaction confirmed: ${signature}`)
      
      // Payment successful
      onPaymentSuccess()
    } catch (error) {
      console.error('Payment error:', error)
      onPaymentError()
    } finally {
      setIsProcessing(false)
    }
  }
  
  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
        <WalletMultiButton className="!bg-indigo-600 hover:!bg-indigo-700" />
        
        <Button
          onClick={handlePayment}
          disabled={isProcessing || !connected}
          ariaLabel={`Pay ${invoice.amount} ${invoice.currency}`}
        >
          {isProcessing ? 'Processing...' : `Pay ${invoice.amount} ${invoice.currency}`}
        </Button>
      </div>
      
      {!connected && (
        <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-md">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">Connect Your Wallet</h3>
              <div className="mt-1 text-sm text-yellow-700">
                <p>
                  Please connect a Solana wallet (such as Phantom or Solflare) to make a payment.
                  Make sure your wallet is connected to Solana devnet.
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
      
      {connected && !isProcessing && (
        <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-blue-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2h-1V9z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">Ready to Pay</h3>
              <div className="mt-1 text-sm text-blue-700">
                <p>
                  Your wallet is connected. Click the Pay button to send {invoice.amount} {invoice.currency} to the recipient.
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default WalletComponents 