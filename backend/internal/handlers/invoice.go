package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/models"
)

// InvoiceRoutes returns a router with all invoice-related routes
func InvoiceRoutes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", CreateInvoice)
	r.Get("/", ListInvoices)
	r.Get("/{linkToken}", GetInvoiceByLinkToken)
	r.Get("/number/{invoiceNumber}", GetInvoiceByNumber)
	r.Put("/{id}/status", UpdateInvoiceStatus)
	return r
}

// CreateInvoice handles the creation of a new invoice
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if invoice.Currency == "" {
		invoice.Currency = "USDC"
	}
	if invoice.Status == "" {
		invoice.Status = models.StatusPending
	}
	if invoice.CreatedAt.IsZero() {
		invoice.CreatedAt = time.Now()
	}
	if invoice.UpdatedAt.IsZero() {
		invoice.UpdatedAt = time.Now()
	}

	// Store in database
	err = models.CreateInvoice(db.DB, &invoice)
	if err != nil {
		http.Error(w, "Failed to create invoice: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created invoice with status 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

// ListInvoices handles listing all invoices with optional pagination
func ListInvoices(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	page := 1
	limit := 10

	// Parse page if provided
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse limit if provided
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get invoices from database
	invoices, err := models.ListInvoices(db.DB, page, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve invoices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the invoices
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

// GetInvoiceByLinkToken handles fetching a single invoice by its link token
func GetInvoiceByLinkToken(w http.ResponseWriter, r *http.Request) {
	linkToken := chi.URLParam(r, "linkToken")
	if linkToken == "" {
		http.Error(w, "Link token is required", http.StatusBadRequest)
		return
	}

	// Get invoice from database
	invoice, err := models.GetInvoiceByLinkToken(db.DB, linkToken)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	// Return the invoice
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

// GetInvoiceByNumber handles fetching a single invoice by its invoice number
func GetInvoiceByNumber(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := chi.URLParam(r, "invoiceNumber")
	if invoiceNumber == "" {
		http.Error(w, "Invoice number is required", http.StatusBadRequest)
		return
	}

	// Get invoice from database
	invoice, err := models.GetInvoiceByInvoiceNumber(db.DB, invoiceNumber)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	// Return the invoice
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

// UpdateInvoiceStatus handles updating an invoice status
func UpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	// Get invoice ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// Parse request body to get new status
	var statusUpdate struct {
		Status models.InvoiceStatus `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	var newStatus models.InvoiceStatus
	switch statusUpdate.Status {
	case models.StatusPending, models.StatusPaid, models.StatusCanceled:
		newStatus = statusUpdate.Status
	default:
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Update in database
	err = models.UpdateInvoiceStatus(db.DB, uint(id), newStatus)
	if err != nil {
		http.Error(w, "Failed to update invoice status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invoice status updated successfully"})
} 