'use client'

import { ChangeEvent } from 'react'

interface TextFieldProps {
  id: string
  name: string
  label: string
  value: string | number
  onChange: (e: ChangeEvent<HTMLInputElement>) => void
  type?: string
  placeholder?: string
  required?: boolean
  min?: number
  step?: string
  readOnly?: boolean
  className?: string
  autoComplete?: string
  error?: string
}

/**
 * Reusable TextField component for consistent form inputs
 * Used for text, number, email, date inputs
 */
export const TextField = ({
  id,
  name,
  label,
  value,
  onChange,
  type = 'text',
  placeholder = '',
  required = false,
  min,
  step,
  readOnly = false,
  className = '',
  autoComplete,
  error,
}: TextFieldProps) => {
  const hasError = Boolean(error);
  
  return (
    <div>
      <label htmlFor={id} className="block text-sm font-medium text-gray-700">
        {label}
      </label>
      <input
        type={type}
        id={id}
        name={name}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        required={required}
        min={min}
        step={step}
        readOnly={readOnly}
        autoComplete={autoComplete}
        className={`mt-1 block w-full rounded-md shadow-sm focus:ring-primary-500 ${
          readOnly ? 'bg-gray-100' : 'border-gray-300'
        } ${
          hasError ? 'border-red-500 focus:border-red-500' : 'focus:border-primary-500'
        } ${className}`}
        aria-invalid={hasError}
        aria-describedby={hasError ? `${id}-error` : undefined}
      />
      {hasError && (
        <p id={`${id}-error`} className="mt-1 text-sm text-red-600">
          {error}
        </p>
      )}
    </div>
  )
}

export default TextField 