'use client'

import React, { createContext, useState, useContext, ReactNode } from 'react'
import Toast from '../components/ui/Toast'

type ToastType = 'success' | 'error' | 'warning' | 'info'

interface ToastMessage {
  id: string
  message: string
  type: ToastType
  duration?: number
}

interface ToastContextType {
  showToast: (message: string, type: ToastType, duration?: number) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

interface ToastProviderProps {
  children: ReactNode
}

/**
 * Provider component for managing toast notifications throughout the app
 */
export const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  // Show a new toast notification
  const showToast = (message: string, type: ToastType, duration = 3000) => {
    const id = Math.random().toString(36).substring(2, 9)
    
    setToasts((prevToasts) => [
      ...prevToasts,
      { id, message, type, duration }
    ])
  }

  // Remove a toast notification
  const removeToast = (id: string) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id))
  }

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      {toasts.map((toast) => (
        <Toast
          key={toast.id}
          message={toast.message}
          type={toast.type}
          duration={toast.duration}
          onClose={() => removeToast(toast.id)}
        />
      ))}
    </ToastContext.Provider>
  )
}

/**
 * Hook to use the toast functionality in components
 */
export const useToast = (): ToastContextType => {
  const context = useContext(ToastContext)
  
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  
  return context
}

export default ToastProvider 