'use client'

import { usePathname } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import Header from './Header'

export default function HeaderClientWrapper() {
  const pathname = usePathname()
  const { isAuthenticated } = useAuth()
  
  // Don't show header on login page or payment pages
  const isLoginPage = pathname === '/login'
  const isPaymentPage = pathname?.startsWith('/pay/')
  
  if (isLoginPage || isPaymentPage) {
    return null
  }
  
  // Show header on all other pages when authenticated
  return isAuthenticated ? <Header /> : null
} 