import './globals.css'
import type { Metadata } from 'next'
import React from 'react'
import SolanaWalletProvider from '@/contexts/SolanaWalletProvider'
import ToastProvider from '@/contexts/ToastProvider'
import { AuthProvider } from '@/contexts/AuthContext'
import Navbar from '@/components/ui/Navbar'
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
            <ToastProvider>
              <Navbar />
              <main className="min-h-screen bg-gray-50 pt-4 pb-8">
                <div className="container-fluid">
                  {children}
                </div>
              </main>
            </ToastProvider>
          </SolanaWalletProvider>
        </AuthProvider>
      </body>
    </html>
  )
} 