package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ncapetillo/demo-fluida/internal/models"
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
	invoices := h.service.GetAllInvoices()
	
	respondWithJSON(w, http.StatusOK, invoices)
}

// GetInvoiceByToken retrieves an invoice by token
func (h *InvoiceHandler) GetInvoiceByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	
	invoice, err := h.service.GetInvoiceByToken(token)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invoice not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, invoice)
}

// CreateInvoice creates a new invoice
func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req models.CreateInvoiceRequest
	
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: " + err.Error())
		return
	}
	
	// Log the request for debugging
	reqBytes, _ := json.Marshal(req)
	log.Printf("CreateInvoice request: %s", string(reqBytes))
	
	// Create the invoice
	invoice, err := h.service.CreateInvoice(req)
	if err != nil {
		log.Printf("Error creating invoice: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed to create invoice: " + err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, invoice)
}

// UpdateInvoiceStatus updates the status of an invoice
func (h *InvoiceHandler) UpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	// Parse the invoice ID
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}
	
	// Parse the request body
	var req models.UpdateInvoiceStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	// Validate the status
	switch req.Status {
	case models.StatusPending, models.StatusPaid, models.StatusCanceled:
		// Valid status
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid status value")
		return
	}
	
	// Update the invoice status
	invoice, err := h.service.UpdateInvoiceStatus(id, req.Status)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invoice not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, invoice)
}

// Helper to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
} 