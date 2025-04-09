import Link from 'next/link'

/**
 * Home page component for the Fluida Invoice Generator
 * Provides navigation to create new invoices or view existing ones
 */
export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-white to-gray-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <div className="flex justify-center">
            <div className="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                className="h-8 w-8 text-primary-600" 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path 
                  strokeLinecap="round" 
                  strokeLinejoin="round" 
                  strokeWidth={2} 
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" 
                />
              </svg>
            </div>
          </div>
          <h1 className="mt-4 text-3xl font-extrabold text-gray-900 sm:text-4xl">
            Fluida Invoice Generator
          </h1>
          <p className="mt-3 text-gray-600 sm:mt-4">
            Create and manage invoices for USDC payments on Solana
          </p>
        </div>
        
        <div className="mt-10 space-y-4">
          <Link 
            href="/create-invoice" 
            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
          >
            Create New Invoice
          </Link>
          
          <Link 
            href="/invoices" 
            className="w-full flex justify-center py-3 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
          >
            View All Invoices
          </Link>
        </div>
        
        <div className="mt-8 text-center">
          <p className="text-xs text-gray-500">
            Running on Solana Testnet
          </p>
        </div>
      </div>
    </div>
  )
} 