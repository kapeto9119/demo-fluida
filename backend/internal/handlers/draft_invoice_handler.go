package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/models"
	"github.com/ncapetillo/demo-fluida/internal/response"
)

// DraftInvoiceHandler handles HTTP requests related to draft invoices
type DraftInvoiceHandler struct{}

// NewDraftInvoiceHandler creates a new draft invoice handler
func NewDraftInvoiceHandler() *DraftInvoiceHandler {
	return &DraftInvoiceHandler{}
}

// Routes returns a router with all draft invoice-related routes
func (h *DraftInvoiceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	
	r.Post("/", h.CreateDraftInvoice)
	r.Get("/{userId}", h.GetDraftInvoiceByUserID)
	r.Put("/{id}", h.UpdateDraftInvoice)
	r.Delete("/{id}", h.DeleteDraftInvoice)
	r.Get("/check", h.CheckInvoiceNumberExists)
	
	return r
}

// CreateDraftInvoice creates a new draft invoice
func (h *DraftInvoiceHandler) CreateDraftInvoice(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDraftInvoiceRequest
	
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload: "+err.Error())
		return
	}
	
	// Check if the user already has a draft invoice
	existingDraft, err := models.GetDraftInvoiceByUserID(db.DB, req.UserID)
	if err == nil {
		// User already has a draft invoice, update it instead
		updateReq := models.UpdateDraftInvoiceRequest{
			InvoiceData: req.InvoiceData,
		}
		if err := models.UpdateDraftInvoice(db.DB, existingDraft.ID, updateReq); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update existing draft invoice: "+err.Error(), "update_failed")
			return
		}
		
		// Get the updated draft
		updatedDraft, _ := models.GetDraftInvoiceByID(db.DB, existingDraft.ID)
		
		response.JSON(w, http.StatusOK, updatedDraft)
		return
	}
	
	// Create new draft invoice
	draft := models.NewDraftInvoice(req)
	if err := models.CreateDraftInvoice(db.DB, &draft); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create draft invoice: "+err.Error(), "creation_failed")
		return
	}
	
	response.JSON(w, http.StatusCreated, draft)
}

// GetDraftInvoiceByUserID retrieves a draft invoice by user ID
func (h *DraftInvoiceHandler) GetDraftInvoiceByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	
	if userID == "" {
		response.BadRequest(w, "User ID is required")
		return
	}
	
	draft, err := models.GetDraftInvoiceByUserID(db.DB, userID)
	if err != nil {
		response.NotFound(w, "No draft invoice found for this user")
		return
	}
	
	response.JSON(w, http.StatusOK, draft)
}

// UpdateDraftInvoice updates an existing draft invoice
func (h *DraftInvoiceHandler) UpdateDraftInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	if id == "" {
		response.BadRequest(w, "Draft invoice ID is required")
		return
	}
	
	// Check if the draft invoice exists
	_, err := models.GetDraftInvoiceByID(db.DB, id)
	if err != nil {
		response.NotFound(w, "Draft invoice not found")
		return
	}
	
	var req models.UpdateDraftInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload: "+err.Error())
		return
	}
	
	if err := models.UpdateDraftInvoice(db.DB, id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update draft invoice: "+err.Error(), "update_failed")
		return
	}
	
	// Get the updated draft
	updatedDraft, _ := models.GetDraftInvoiceByID(db.DB, id)
	
	response.JSON(w, http.StatusOK, updatedDraft)
}

// DeleteDraftInvoice deletes a draft invoice
func (h *DraftInvoiceHandler) DeleteDraftInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	if id == "" {
		response.BadRequest(w, "Draft invoice ID is required")
		return
	}
	
	// Check if the draft invoice exists
	_, err := models.GetDraftInvoiceByID(db.DB, id)
	if err != nil {
		response.NotFound(w, "Draft invoice not found")
		return
	}
	
	if err := models.DeleteDraftInvoice(db.DB, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete draft invoice: "+err.Error(), "deletion_failed")
		return
	}
	
	response.Success(w, http.StatusOK, "Draft invoice deleted successfully")
}

// CheckInvoiceNumberExists checks if an invoice number already exists
func (h *DraftInvoiceHandler) CheckInvoiceNumberExists(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := r.URL.Query().Get("invoice_number")
	
	if invoiceNumber == "" {
		response.BadRequest(w, "Invoice number is required")
		return
	}
	
	var count int64
	if err := db.DB.Model(&models.Invoice{}).Where("invoice_number = ?", invoiceNumber).Count(&count).Error; err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to check invoice number: "+err.Error(), "database_error")
		return
	}
	
	response.JSON(w, http.StatusOK, map[string]bool{"exists": count > 0})
} 