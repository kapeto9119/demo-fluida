'use client'

import { useAuth } from '@/contexts/AuthContext'
import Link from 'next/link'

export default function Header() {
  const { isAuthenticated, logout } = useAuth()

  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <Link href="/" className="text-xl font-bold text-indigo-600">
          Fluida Invoice
        </Link>
        
        <nav className="flex space-x-4">
          <Link href="/" className="text-gray-600 hover:text-indigo-600">
            Home
          </Link>
          <Link href="/invoices" className="text-gray-600 hover:text-indigo-600">
            Invoices
          </Link>
          <Link href="/create-invoice" className="text-gray-600 hover:text-indigo-600">
            Create Invoice
          </Link>
          
          {isAuthenticated && (
            <button 
              onClick={logout}
              className="text-red-500 hover:text-red-700"
            >
              Logout
            </button>
          )}
        </nav>
      </div>
    </header>
  )
} 