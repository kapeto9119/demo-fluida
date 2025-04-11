'use client'

import { ReactNode } from 'react'

interface ButtonProps {
  children: ReactNode
  onClick?: () => void
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'outline' | 'danger' | 'success'
  size?: 'small' | 'medium' | 'large'
  disabled?: boolean
  className?: string
  fullWidth?: boolean
  ariaLabel?: string
  isLoading?: boolean
  icon?: ReactNode
  iconPosition?: 'left' | 'right'
}

/**
 * Reusable Button component with different variants and loading state
 */
export const Button = ({
  children,
  onClick,
  type = 'button',
  variant = 'primary',
  size = 'medium',
  disabled = false,
  className = '',
  fullWidth = false,
  ariaLabel,
  isLoading = false,
  icon,
  iconPosition = 'left'
}: ButtonProps) => {
  // Get variant-specific styles
  const getVariantStyles = () => {
    switch (variant) {
      case 'primary':
        return 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500 active:bg-primary-800'
      case 'secondary':
        return 'bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500 active:bg-gray-800'
      case 'outline':
        return 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 focus:ring-primary-500 active:bg-gray-100'
      case 'danger':
        return 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 active:bg-red-800'
      case 'success':
        return 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 active:bg-green-800'
      default:
        return 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500 active:bg-primary-800'
    }
  }
  
  // Get size-specific styles
  const getSizeStyles = () => {
    switch (size) {
      case 'small':
        return 'px-2.5 py-1.5 text-xs'
      case 'medium':
        return 'px-4 py-2 text-sm'
      case 'large':
        return 'px-6 py-3 text-base'
      default:
        return 'px-4 py-2 text-sm'
    }
  }

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || isLoading}
      aria-label={ariaLabel}
      className={`
        ${getSizeStyles()}
        rounded-md shadow-sm font-medium
        focus:outline-none focus:ring-2 focus:ring-offset-2
        transition-all duration-200
        disabled:opacity-50 disabled:cursor-not-allowed
        ${getVariantStyles()}
        ${fullWidth ? 'w-full' : ''}
        ${isLoading ? 'relative' : ''}
        flex items-center justify-center
        ${className}
      `}
    >
      {isLoading && (
        <span className="absolute inset-0 flex items-center justify-center">
          <svg className="animate-spin h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </span>
      )}
      
      <span className={`${isLoading ? 'invisible' : ''} flex items-center`}>
        {icon && iconPosition === 'left' && <span className="mr-2">{icon}</span>}
        {children}
        {icon && iconPosition === 'right' && <span className="ml-2">{icon}</span>}
      </span>
    </button>
  )
}

export default Button 