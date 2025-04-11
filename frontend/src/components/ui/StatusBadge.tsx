'use client'

type StatusType = 'PAID' | 'PENDING' | 'OVERDUE' | 'CANCELLED'

interface StatusBadgeProps {
  status: StatusType | string
  className?: string
}

/**
 * Reusable StatusBadge component for displaying invoice statuses
 * Automatically applies appropriate colors based on status
 */
export const StatusBadge = ({ status, className = '' }: StatusBadgeProps) => {
  const getStatusStyles = (status: string) => {
    switch (status) {
      case 'PAID':
        return 'bg-green-100 text-green-800'
      case 'PENDING':
        return 'bg-yellow-100 text-yellow-800'
      case 'OVERDUE':
        return 'bg-red-100 text-red-800'
      case 'CANCELLED':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-blue-100 text-blue-800'
    }
  }

  return (
    <span
      className={`status-badge ${getStatusStyles(
        status
      )} ${className}`}
    >
      {status}
    </span>
  )
}

export default StatusBadge 