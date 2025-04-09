'use client'

import { ReactNode } from 'react'

interface ButtonProps {
  children: ReactNode
  onClick?: () => void
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'outline'
  disabled?: boolean
  className?: string
  fullWidth?: boolean
  ariaLabel?: string
}

/**
 * Reusable Button component with different variants
 */
export const Button = ({
  children,
  onClick,
  type = 'button',
  variant = 'primary',
  disabled = false,
  className = '',
  fullWidth = false,
  ariaLabel,
}: ButtonProps) => {
  const getVariantStyles = () => {
    switch (variant) {
      case 'primary':
        return 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500'
      case 'secondary':
        return 'bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500'
      case 'outline':
        return 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 focus:ring-primary-500'
      default:
        return 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500'
    }
  }

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel}
      className={`
        px-4 py-2 rounded-md shadow-sm text-sm font-medium
        focus:outline-none focus:ring-2 focus:ring-offset-2
        transition-colors duration-200
        disabled:opacity-50 disabled:cursor-not-allowed
        ${getVariantStyles()}
        ${fullWidth ? 'w-full' : ''}
        ${className}
      `}
    >
      {children}
    </button>
  )
}

export default Button 