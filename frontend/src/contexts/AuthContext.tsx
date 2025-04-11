'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useRouter, usePathname } from 'next/navigation'

// Import API URL
const API_URL = 'https://serene-radiance-production.up.railway.app'

interface AuthContextType {
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const router = useRouter()
  const pathname = usePathname()

  // Check for existing auth on mount
  useEffect(() => {
    const storedToken = localStorage.getItem('auth_token')
    if (storedToken) {
      setIsAuthenticated(true)
    } else if (pathname !== '/login' && !isPublicRoute(pathname)) {
      // Redirect to login if not authenticated and not on login page
      router.push('/login')
    }
  }, [pathname, router])

  const login = async (username: string, password: string) => {
    try {
      // Create temporary authentication token
      const tempToken = btoa(`${username}:${password}`)
      
      // Silent authentication verification - using a cleaner approach
      const controller = new AbortController();
      const signal = controller.signal;
      
      try {
        // Use dedicated auth verification endpoint
        const response = await fetch(`${API_URL}/api/v1/auth/verify`, {
          method: 'GET',
          headers: {
            'Authorization': `Basic ${tempToken}`,
            // Prevent OPTIONS preflight for simpler requests
            'Accept': 'application/json',
          },
          // Don't send cookies for this request
          credentials: 'omit',
          signal: signal,
        });
        
        if (response.ok) {
          // Store only the auth token, not the raw credentials
          localStorage.setItem('auth_token', tempToken)
          localStorage.setItem('auth_username', username) // Keep username for UI purposes
          // NEVER store raw password in localStorage
          
          setIsAuthenticated(true)
          router.push('/')
          return;
        }
        
        // For any non-OK response, throw a user-friendly error
        throw new Error('Invalid username or password');
      } catch (fetchError) {
        // If it's a network error or unexpected issue, provide a user-friendly message
        if (fetchError instanceof Error && fetchError.name === 'AbortError') {
          // Request was aborted - stay silent
          return;
        }
        
        // For auth errors and other issues, use a consistent message
        throw new Error('Invalid username or password');
      } finally {
        // Clean up the controller
        controller.abort();
      }
    } catch (error) {
      console.error('Authentication failed');
      throw error; // Rethrow to allow login page to show error
    }
  }

  const logout = () => {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_username')
    localStorage.removeItem('auth_password') // Remove this if it exists (for backward compatibility)
    setIsAuthenticated(false)
    router.push('/login')
  }

  // Check if a route is public (doesn't require authentication)
  const isPublicRoute = (path: string) => {
    // Add any public routes here
    const publicRoutes = ['/login', '/api/health', '/health']
    return publicRoutes.includes(path) || path.startsWith('/pay/')
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
} 