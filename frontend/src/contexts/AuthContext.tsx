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
    const storedCredentials = localStorage.getItem('auth_credentials')
    if (storedCredentials) {
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
      
      // Verify credentials with backend before storing
      const response = await fetch(`${API_URL}/api/v1/health`, {
        headers: {
          'Authorization': `Basic ${tempToken}`
        }
      })
      
      if (!response.ok) {
        throw new Error('Invalid credentials')
      }
      
      // Store verified credentials in localStorage
      localStorage.setItem('auth_credentials', tempToken)
      localStorage.setItem('auth_username', username)
      localStorage.setItem('auth_password', password)
      
      setIsAuthenticated(true)
      router.push('/')
    } catch (error) {
      console.error('Login error:', error)
      throw error // Rethrow to allow login page to show error
    }
  }

  const logout = () => {
    localStorage.removeItem('auth_credentials')
    localStorage.removeItem('auth_username')
    localStorage.removeItem('auth_password')
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