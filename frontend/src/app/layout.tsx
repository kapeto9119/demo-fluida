import './globals.css'
import type { Metadata } from 'next'
import React from 'react'
import SolanaWalletProvider from '@/contexts/SolanaWalletProvider'
import { AuthProvider } from '@/contexts/AuthContext'
import HeaderClientWrapper from '@/components/HeaderClientWrapper'
import { Inter, DM_Sans } from 'next/font/google'

// Initialize the fonts
const inter = Inter({ 
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
})

const dmSans = DM_Sans({
  weight: ['400', '500', '700'],
  subsets: ['latin'],
  variable: '--font-dm-sans',
  display: 'swap',
})

export const metadata: Metadata = {
  title: 'Fluida - Invoice & Payment Link Generator',
  description: 'Create invoices and payment links for USDC payments on Solana',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className={`${inter.variable} ${dmSans.variable}`} suppressHydrationWarning>
      <body className="font-sans">
        <AuthProvider>
          <SolanaWalletProvider>
            <HeaderClientWrapper />
            <main className="min-h-screen bg-gray-50">
              {children}
            </main>
          </SolanaWalletProvider>
        </AuthProvider>
      </body>
    </html>
  )
} 