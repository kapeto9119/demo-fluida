// Package response provides standardized API response utilities
package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Response is the standard API response structure
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
}

// ErrorResponse represents an error in the API
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail provides specific information about an error
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ValidationError represents a validation error in a specific field
type ValidationError struct {
	Field   string
	Message string
}

// New creates a new response
func New() *Response {
	return &Response{}
}

// WithData adds data to the response
func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

// WithMeta adds metadata to the response
func (r *Response) WithMeta(meta interface{}) *Response {
	r.Meta = meta
	return r
}

// WithPagination adds pagination metadata to the response
func (r *Response) WithPagination(total, page, limit int) *Response {
	r.Meta = map[string]interface{}{
		"pagination": map[string]int{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}
	return r
}

// WithError adds an error to the response
func (r *Response) WithError(message, code string) *Response {
	r.Error = &ErrorResponse{
		Message: message,
		Code:    code,
	}
	return r
}

// WithValidationErrors adds validation errors to the response
func (r *Response) WithValidationErrors(errors []ValidationError) *Response {
	if r.Error == nil {
		r.Error = &ErrorResponse{
			Message: "Validation error",
			Code:    "validation_error",
		}
	}

	r.Error.Details = make([]ErrorDetail, 0, len(errors))
	for _, err := range errors {
		r.Error.Details = append(r.Error.Details, ErrorDetail{
			Field:   err.Field,
			Message: err.Message,
		})
	}
	return r
}

// Send writes the response as JSON to the HTTP response writer
func (r *Response) Send(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// Helper functions

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Debug the response payload
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling response data: %v", err)
	} else {
		log.Printf("Response payload: %s", string(dataBytes))
	}
	
	New().WithData(data).Send(w, statusCode)
}

// Error sends an error response with the given status code
func Error(w http.ResponseWriter, statusCode int, message, code string) {
	// Clean any potentially duplicated error messages
	// This avoids the "Failed to create invoice: failed to create invoice:" type errors
	cleanedMessage := message
	if strings.Contains(message, "failed to create invoice:") {
		prefix := "failed to create invoice:"
		if strings.HasPrefix(message, prefix) {
			rest := strings.TrimPrefix(message, prefix)
			if strings.Contains(rest, prefix) {
				// If we have duplication, clean it
				cleanedMessage = prefix + strings.TrimPrefix(rest, prefix)
			}
		}
	}
	
	New().WithError(cleanedMessage, code).Send(w, statusCode)
}

// ValidationError sends a validation error response
func ValidationErrors(w http.ResponseWriter, errors []ValidationError) {
	New().WithError("Validation error", "validation_error").
		WithValidationErrors(errors).
		Send(w, http.StatusBadRequest)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(w, http.StatusNotFound, message, "not_found")
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad request"
	}
	Error(w, http.StatusBadRequest, message, "bad_request")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Internal server error", "internal_error")
}

// Success sends a success response with a message
func Success(w http.ResponseWriter, statusCode int, message string) {
	New().WithData(map[string]string{"message": message}).Send(w, statusCode)
}
