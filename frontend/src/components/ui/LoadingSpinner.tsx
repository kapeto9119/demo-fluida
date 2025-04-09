'use client'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg'
  color?: string
  className?: string
  text?: string
}

/**
 * Reusable LoadingSpinner component for indicating loading states
 */
export const LoadingSpinner = ({
  size = 'md',
  color = 'primary-600',
  className = '',
  text,
}: LoadingSpinnerProps) => {
  const getSize = () => {
    switch (size) {
      case 'sm':
        return 'h-6 w-6'
      case 'md':
        return 'h-10 w-10'
      case 'lg':
        return 'h-16 w-16'
      default:
        return 'h-10 w-10'
    }
  }

  return (
    <div className={`flex flex-col items-center justify-center ${className}`}>
      <div className={`animate-spin rounded-full border-b-2 border-${color} ${getSize()}`}></div>
      {text && <p className="mt-4 text-gray-600">{text}</p>}
    </div>
  )
}

export default LoadingSpinner 