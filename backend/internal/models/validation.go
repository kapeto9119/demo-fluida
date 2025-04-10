package models

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

// ValidationErrors maps field names to error messages
type ValidationErrors map[string]string

// Validate checks if a string is not empty
func validateRequired(field, value string, errors ValidationErrors) {
	if strings.TrimSpace(value) == "" {
		errors[field] = "This field is required"
	}
}

// ValidateEmail checks if a string is a valid email address
func validateEmail(field, value string, errors ValidationErrors) {
	if value == "" {
		return // Skip empty values, use validateRequired for required fields
	}
	
	_, err := mail.ParseAddress(value)
	if err != nil {
		errors[field] = "Invalid email address"
	}
}

// ValidateMinLength checks if a string has a minimum length
func validateMinLength(field, value string, minLength int, errors ValidationErrors) {
	if len(value) < minLength {
		errors[field] = fmt.Sprintf("Must be at least %d characters", minLength)
	}
}

// ValidateMaxLength checks if a string has a maximum length
func validateMaxLength(field, value string, maxLength int, errors ValidationErrors) {
	if len(value) > maxLength {
		errors[field] = fmt.Sprintf("Must be no more than %d characters", maxLength)
	}
}

// ValidateMinValue checks if a number is at least a minimum value
func validateMinValue(field string, value, minValue float64, errors ValidationErrors) {
	if value < minValue {
		errors[field] = fmt.Sprintf("Must be at least %g", minValue)
	}
}

// ValidateMaxValue checks if a number is at most a maximum value
func validateMaxValue(field string, value, maxValue float64, errors ValidationErrors) {
	if value > maxValue {
		errors[field] = fmt.Sprintf("Must be at most %g", maxValue)
	}
}

// ValidateFutureDate checks if a date is in the future
func validateFutureDate(field string, date time.Time, errors ValidationErrors) {
	if date.Before(time.Now()) {
		errors[field] = "Date must be in the future"
	}
}

// ValidateRegex checks if a string matches a regex pattern
func validateRegex(field, value, pattern, message string, errors ValidationErrors) {
	if value == "" {
		return // Skip empty values, use validateRequired for required fields
	}
	
	match, _ := regexp.MatchString(pattern, value)
	if !match {
		errors[field] = message
	}
}

// ValidateSolanaAddress checks if a string is a valid Solana address
func validateSolanaAddress(field, value string, errors ValidationErrors) {
	if value == "" {
		return // Skip empty values, use validateRequired for required fields
	}
	
	// Solana addresses are base58 encoded and typically 32-44 characters
	// This is a simplified check - a real implementation would do more thorough validation
	match, _ := regexp.MatchString(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`, value)
	if !match {
		errors[field] = "Invalid Solana wallet address"
	}
}
