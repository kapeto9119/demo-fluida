package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

// Person represents sender or recipient information
type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address,omitempty"`
}

// Value implements the driver.Valuer interface for Person
func (p Person) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for Person
func (p *Person) Scan(value interface{}) error {
	if value == nil {
		*p = Person{}
		return nil
	}
	
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to scan Person: unexpected type %T", value)
	}
	
	return json.Unmarshal(data, p)
}

// Invoice represents a payment invoice in the system
type Invoice struct {
	ID               int            `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceNumber    string         `json:"invoiceNumber" gorm:"uniqueIndex:idx_invoice_number;not null;type:varchar(50)"`
	Amount           float64        `json:"amount" gorm:"not null;type:decimal(12,2)"`
	Currency         string         `json:"currency" gorm:"not null;default:USDC;type:varchar(10);index:idx_invoice_currency"`
	Description      string         `json:"description" gorm:"type:text"`
	DueDate          time.Time      `json:"dueDate" gorm:"not null;index:idx_invoice_due_date"`
	Status           InvoiceStatus  `json:"status" gorm:"not null;default:PENDING;type:varchar(20);index:idx_invoice_status"`
	ReceiverAddr     string         `json:"receiverAddr" gorm:"not null;type:varchar(100);index:idx_invoice_receiver"`
	LinkToken        string         `json:"linkToken,omitempty" gorm:"uniqueIndex:idx_invoice_link;not null;type:varchar(100)"`
	SenderDetails    Person         `json:"senderDetails" gorm:"type:jsonb;serializer:json"`
	RecipientDetails Person         `json:"recipientDetails" gorm:"type:jsonb;serializer:json"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime;index:idx_invoice_created_at"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName overrides the table name
func (Invoice) TableName() string {
	return "invoice"
}

// BeforeCreate hook runs before creating a new invoice
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	// Generate a UUID for the link token if not provided
	if i.LinkToken == "" {
		i.LinkToken = uuid.New().String()
	}
	
	// Set default status if not provided
	if i.Status == "" {
		i.Status = StatusPending
	}
	
	// Set default currency if not provided
	if i.Currency == "" {
		i.Currency = "USDC"
	}
	
	return nil
}

// Validate validates the invoice data
func (i *Invoice) Validate() error {
	if i.InvoiceNumber == "" {
		return fmt.Errorf("invoice number is required")
	}
	
	if i.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	
	if i.ReceiverAddr == "" {
		return fmt.Errorf("receiver address is required")
	}
	
	if i.SenderDetails.Name == "" || i.SenderDetails.Email == "" {
		return fmt.Errorf("sender name and email are required")
	}
	
	if i.RecipientDetails.Name == "" || i.RecipientDetails.Email == "" {
		return fmt.Errorf("recipient name and email are required")
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
	InvoiceNumber    string    `json:"invoiceNumber"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	Description      string    `json:"description"`
	DueDate          time.Time `json:"dueDate"`
	ReceiverAddr     string    `json:"receiverAddr"`
	SenderDetails    Person    `json:"senderDetails"`
	RecipientDetails Person    `json:"recipientDetails"`
}

// UpdateInvoiceStatusRequest represents the data required to update an invoice status
type UpdateInvoiceStatusRequest struct {
	Status InvoiceStatus `json:"status" binding:"required"`
}

// NewInvoice creates a new invoice from a create request
func NewInvoice(req CreateInvoiceRequest) Invoice {
	now := time.Now()
	linkToken := uuid.New().String()
	
	return Invoice{
		InvoiceNumber:    req.InvoiceNumber,
		Amount:           req.Amount,
		Currency:         req.Currency,
		Description:      req.Description,
		DueDate:          req.DueDate,
		Status:           StatusPending,
		ReceiverAddr:     req.ReceiverAddr,
		LinkToken:        linkToken,
		SenderDetails:    req.SenderDetails,
		RecipientDetails: req.RecipientDetails,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
} 