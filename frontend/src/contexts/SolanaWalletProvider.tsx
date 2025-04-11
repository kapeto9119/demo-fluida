'use client'

import { FC, ReactNode, useMemo, useState, useEffect } from 'react'
import { ConnectionProvider, WalletProvider } from '@solana/wallet-adapter-react'
import { WalletAdapterNetwork } from '@solana/wallet-adapter-base'
import { WalletModalProvider } from '@solana/wallet-adapter-react-ui'
import { clusterApiUrl } from '@solana/web3.js'

// Use our local CSS file that doesn't have @import statements
import '@/styles/wallet-adapter.css'

interface SolanaWalletProviderProps {
  children: ReactNode
}

const SolanaWalletProvider: FC<SolanaWalletProviderProps> = ({ children }) => {
  // State to track if we're in the browser
  const [isMounted, setIsMounted] = useState(false)

  // Set isMounted to true once the component is mounted
  useEffect(() => {
    setIsMounted(true)
    return () => setIsMounted(false)
  }, [])

  // Use devnet for development
  const network = WalletAdapterNetwork.Devnet

  // Set up connection endpoint
  const endpoint = useMemo(() => clusterApiUrl(network), [network])

  // Using empty wallets array - Wallet Standard will auto-detect supported wallets
  const wallets = useMemo(
    () => [],
    []
  )

  // Return null on server-side rendering to prevent hydration issues
  if (!isMounted) {
    return <>{children}</>
  }

  return (
    <ConnectionProvider endpoint={endpoint}>
      <WalletProvider wallets={wallets} autoConnect>
        <WalletModalProvider>
          {children}
        </WalletModalProvider>
      </WalletProvider>
    </ConnectionProvider>
  )
}

export default SolanaWalletProvider