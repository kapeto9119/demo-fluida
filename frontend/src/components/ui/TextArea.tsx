'use client'

import { ChangeEvent } from 'react'

interface TextAreaProps {
  id: string
  name: string
  label: string
  value: string
  onChange: (e: ChangeEvent<HTMLTextAreaElement>) => void
  rows?: number
  placeholder?: string
  required?: boolean
  className?: string
}

/**
 * Reusable TextArea component for multiline text inputs
 */
export const TextArea = ({
  id,
  name,
  label,
  value,
  onChange,
  rows = 3,
  placeholder = '',
  required = false,
  className = '',
}: TextAreaProps) => {
  return (
    <div>
      <label htmlFor={id} className="block text-sm font-medium text-gray-700">
        {label}
      </label>
      <textarea
        id={id}
        name={name}
        value={value}
        onChange={onChange}
        rows={rows}
        placeholder={placeholder}
        required={required}
        className={`mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 ${className}`}
      />
    </div>
  )
}

export default TextArea 