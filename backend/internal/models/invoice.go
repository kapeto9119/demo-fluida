package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceStatus represents the possible states of an invoice
type InvoiceStatus string

const (
	StatusPending  InvoiceStatus = "PENDING"
	StatusPaid     InvoiceStatus = "PAID"
	StatusCanceled InvoiceStatus = "CANCELED"
)

// SenderDetails contains information about the invoice sender
type SenderDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address,omitempty"`
}

// RecipientDetails contains information about the invoice recipient
type RecipientDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address,omitempty"`
}

// Invoice represents a payment invoice in the system
type Invoice struct {
	ID              int       `json:"id"`
	InvoiceNumber   string    `json:"invoiceNumber"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Description     string    `json:"description"`
	DueDate         time.Time `json:"dueDate"`
	Status          InvoiceStatus `json:"status"`
	ReceiverAddr    string    `json:"receiverAddr"`
	LinkToken       string    `json:"linkToken,omitempty"`
	SenderDetails   Person    `json:"senderDetails"`
	RecipientDetails Person   `json:"recipientDetails"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Person represents an entity (sender or recipient) in the system
type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// BeforeCreate hook runs before creating a new invoice
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	// Generate a UUID for the link token if not provided
	if i.LinkToken == "" {
		i.LinkToken = uuid.New().String()
	}
	return nil
}

// CreateInvoice creates a new invoice in the database
func CreateInvoice(db *gorm.DB, invoice *Invoice) error {
	return db.Create(invoice).Error
}

// GetInvoiceByLinkToken fetches an invoice by its link token
func GetInvoiceByLinkToken(db *gorm.DB, linkToken string) (*Invoice, error) {
	var invoice Invoice
	err := db.Where("link_token = ?", linkToken).First(&invoice).Error
	return &invoice, err
}

// GetInvoiceByInvoiceNumber fetches an invoice by its invoice number
func GetInvoiceByInvoiceNumber(db *gorm.DB, invoiceNumber string) (*Invoice, error) {
	var invoice Invoice
	err := db.Where("invoice_number = ?", invoiceNumber).First(&invoice).Error
	return &invoice, err
}

// ListInvoices fetches all invoices with optional pagination
func ListInvoices(db *gorm.DB, page, limit int) ([]Invoice, error) {
	var invoices []Invoice
	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Order("created_at desc").Find(&invoices).Error
	return invoices, err
}

// UpdateInvoiceStatus updates the status of an invoice
func UpdateInvoiceStatus(db *gorm.DB, id uint, status InvoiceStatus) error {
	return db.Model(&Invoice{}).Where("id = ?", id).Update("status", status).Error
}

// CreateInvoiceRequest represents the data required to create a new invoice
type CreateInvoiceRequest struct {
	InvoiceNumber   string    `json:"invoiceNumber"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Description     string    `json:"description"`
	DueDate         time.Time `json:"dueDate"`
	ReceiverAddr    string    `json:"receiverAddr"`
	SenderDetails   Person    `json:"senderDetails"`
	RecipientDetails Person   `json:"recipientDetails"`
}

// UpdateInvoiceStatusRequest represents the data required to update an invoice status
type UpdateInvoiceStatusRequest struct {
	Status InvoiceStatus `json:"status"`
}

// NewInvoice creates a new invoice from a create request
func NewInvoice(req CreateInvoiceRequest) Invoice {
	now := time.Now()
	linkToken := uuid.New().String()
	
	return Invoice{
		InvoiceNumber:   req.InvoiceNumber,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Description:     req.Description,
		DueDate:         req.DueDate,
		Status:          StatusPending,
		ReceiverAddr:    req.ReceiverAddr,
		LinkToken:       linkToken,
		SenderDetails:   req.SenderDetails,
		RecipientDetails: req.RecipientDetails,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
} 