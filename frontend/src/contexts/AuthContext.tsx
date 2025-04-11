'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useRouter, usePathname } from 'next/navigation'

interface AuthContextType {
  isAuthenticated: boolean
  login: (username: string, password: string) => void
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

  const login = (username: string, password: string) => {
    try {
      // Store credentials in localStorage
      const credentials = btoa(`${username}:${password}`)
      localStorage.setItem('auth_credentials', credentials)
      localStorage.setItem('auth_username', username)
      localStorage.setItem('auth_password', password)
      
      setIsAuthenticated(true)
      router.push('/')
    } catch (error) {
      console.error('Login error:', error)
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