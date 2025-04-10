package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ncapetillo/demo-fluida/internal/models"
	"github.com/ncapetillo/demo-fluida/internal/response"
	"github.com/ncapetillo/demo-fluida/internal/services"
)

// InvoiceHandler handles HTTP requests related to invoices
type InvoiceHandler struct {
	service *services.InvoiceService
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(service *services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		service: service,
	}
}

// Routes returns a router with all invoice-related routes
func (h *InvoiceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	
	r.Get("/", h.GetAllInvoices)
	r.Post("/", h.CreateInvoice)
	r.Get("/{token}", h.GetInvoiceByToken)
	r.Put("/{id}/status", h.UpdateInvoiceStatus)
	
	return r
}

// GetAllInvoices returns all invoices
func (h *InvoiceHandler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	invoices := h.service.GetAllInvoices()
	total := len(invoices)
	
	// Calculate pagination
	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		// Return empty result for out-of-range pages
		response.New().
			WithData([]models.Invoice{}).
			WithPagination(total, page, limit).
			Send(w, http.StatusOK)
		return
	}
	if end > total {
		end = total
	}
	
	// Return paginated result
	response.New().
		WithData(invoices[start:end]).
		WithPagination(total, page, limit).
		Send(w, http.StatusOK)
}

// GetInvoiceByToken retrieves an invoice by token
func (h *InvoiceHandler) GetInvoiceByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	
	if token == "" {
		response.BadRequest(w, "Token is required")
		return
	}
	
	invoice, err := h.service.GetInvoiceByToken(token)
	if err != nil {
		response.NotFound(w, "Invoice not found")
		return
	}
	
	response.JSON(w, http.StatusOK, invoice)
}

// CreateInvoice creates a new invoice
func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req models.CreateInvoiceRequest
	
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload: " + err.Error())
		return
	}
	
	// Log the request for debugging
	reqBytes, _ := json.Marshal(req)
	log.Printf("CreateInvoice request: %s", string(reqBytes))
	
	// Validate the request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		// Convert to response.ValidationError format
		errors := make([]response.ValidationError, 0, len(validationErrors))
		for field, message := range validationErrors {
			errors = append(errors, response.ValidationError{
				Field:   field,
				Message: message,
			})
		}
		response.ValidationErrors(w, errors)
		return
	}
	
	// Create the invoice
	invoice, err := h.service.CreateInvoice(req)
	if err != nil {
		log.Printf("Error creating invoice: %v", err)
		
		// Check if this is a duplicate invoice number error
		if strings.Contains(err.Error(), "already exists") {
			response.Error(w, http.StatusConflict, err.Error(), "duplicate_invoice_number")
			return
		}
		
		// For other errors, don't wrap with "Failed to create invoice:" since that may already be in the message
		response.Error(w, http.StatusBadRequest, err.Error(), "creation_failed")
		return
	}
	
	// For production, we should only log basic info about the created invoice
	log.Printf("Invoice #%s created with token: %s", invoice.InvoiceNumber, invoice.LinkToken)
	
	// Ensure link token is not empty
	if invoice.LinkToken == "" {
		log.Printf("WARNING: LinkToken is empty in the created invoice!")
	}
	
	// Return the created invoice with status 201 Created
	response.JSON(w, http.StatusCreated, invoice)
}

// UpdateInvoiceStatus updates the status of an invoice
func (h *InvoiceHandler) UpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	// Parse the invoice ID
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid invoice ID")
		return
	}
	
	// Parse the request body
	var req models.UpdateInvoiceStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}
	
	// Validate the status
	switch req.Status {
	case models.StatusPending, models.StatusPaid, models.StatusCanceled:
		// Valid status
	default:
		response.BadRequest(w, "Invalid status value")
		return
	}
	
	// Update the invoice status
	invoice, err := h.service.UpdateInvoiceStatus(id, req.Status)
	if err != nil {
		response.NotFound(w, "Invoice not found")
		return
	}
	
	response.JSON(w, http.StatusOK, invoice)
} 