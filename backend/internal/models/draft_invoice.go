package models

import (
	"time"

	"gorm.io/gorm"
)

// DraftInvoice represents a user's draft invoice saved to the database
type DraftInvoice struct {
	ID               string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID           string         `json:"userId" gorm:"not null;index:idx_draft_user_id"`
	InvoiceData      string         `json:"invoiceData" gorm:"type:jsonb"` // Store the entire form data as JSON
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName overrides the table name
func (DraftInvoice) TableName() string {
	return "draft_invoice"
}

// CreateDraftInvoiceRequest represents the data required to create a new draft invoice
type CreateDraftInvoiceRequest struct {
	UserID      string `json:"userId"`
	InvoiceData string `json:"invoiceData"`
}

// UpdateDraftInvoiceRequest represents the data required to update a draft invoice
type UpdateDraftInvoiceRequest struct {
	InvoiceData string `json:"invoiceData"`
}

// NewDraftInvoice creates a new draft invoice from a create request
func NewDraftInvoice(req CreateDraftInvoiceRequest) DraftInvoice {
	now := time.Now()
	
	return DraftInvoice{
		UserID:      req.UserID,
		InvoiceData: req.InvoiceData,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// CreateDraftInvoice creates a new draft invoice in the database
func CreateDraftInvoice(db *gorm.DB, draft *DraftInvoice) error {
	return db.Create(draft).Error
}

// GetDraftInvoiceByUserID fetches a draft invoice by user ID
func GetDraftInvoiceByUserID(db *gorm.DB, userID string) (*DraftInvoice, error) {
	var draft DraftInvoice
	err := db.Where("user_id = ?", userID).First(&draft).Error
	return &draft, err
}

// GetDraftInvoiceByID fetches a draft invoice by its ID
func GetDraftInvoiceByID(db *gorm.DB, id string) (*DraftInvoice, error) {
	var draft DraftInvoice
	err := db.Where("id = ?", id).First(&draft).Error
	return &draft, err
}

// UpdateDraftInvoice updates an existing draft invoice
func UpdateDraftInvoice(db *gorm.DB, id string, req UpdateDraftInvoiceRequest) error {
	return db.Model(&DraftInvoice{}).Where("id = ?", id).Updates(map[string]interface{}{
		"invoice_data": req.InvoiceData,
		"updated_at":   time.Now(),
	}).Error
}

// DeleteDraftInvoice deletes a draft invoice
func DeleteDraftInvoice(db *gorm.DB, id string) error {
	return db.Delete(&DraftInvoice{}, "id = ?", id).Error
} 