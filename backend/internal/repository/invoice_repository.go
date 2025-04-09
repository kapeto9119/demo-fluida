package repository

import (
	"context"
	"errors"

	"github.com/ncapetillo/demo-fluida/internal/models"
	"gorm.io/gorm"
)

// InvoiceRepository defines methods to interact with invoices in the database
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *models.Invoice) error
	FindByID(ctx context.Context, id int) (*models.Invoice, error)
	FindByLinkToken(ctx context.Context, linkToken string) (*models.Invoice, error)
	FindByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.Invoice, error)
	List(ctx context.Context, page, limit int) ([]models.Invoice, error)
	UpdateStatus(ctx context.Context, id int, status models.InvoiceStatus) error
	FindPendingInvoices(ctx context.Context) ([]models.Invoice, error)
	Update(ctx context.Context, invoice *models.Invoice) error
}

// GORMInvoiceRepository implements InvoiceRepository using GORM
type GORMInvoiceRepository struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new invoice repository
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &GORMInvoiceRepository{db: db}
}

// Create adds a new invoice to the database
func (r *GORMInvoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

// FindByID retrieves an invoice by ID
func (r *GORMInvoiceRepository) FindByID(ctx context.Context, id int) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.WithContext(ctx).First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when record not found
		}
		return nil, err
	}
	return &invoice, nil
}

// FindByLinkToken retrieves an invoice by link token
func (r *GORMInvoiceRepository) FindByLinkToken(ctx context.Context, linkToken string) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.WithContext(ctx).Where("link_token = ?", linkToken).First(&invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// FindByInvoiceNumber retrieves an invoice by invoice number
func (r *GORMInvoiceRepository) FindByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.WithContext(ctx).Where("invoice_number = ?", invoiceNumber).First(&invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// List retrieves invoices with pagination
func (r *GORMInvoiceRepository) List(ctx context.Context, page, limit int) ([]models.Invoice, error) {
	var invoices []models.Invoice
	offset := (page - 1) * limit
	
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&invoices).Error; err != nil {
		return nil, err
	}
	
	return invoices, nil
}

// UpdateStatus updates the status of an invoice
func (r *GORMInvoiceRepository) UpdateStatus(ctx context.Context, id int, status models.InvoiceStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// FindPendingInvoices retrieves all pending invoices
func (r *GORMInvoiceRepository) FindPendingInvoices(ctx context.Context) ([]models.Invoice, error) {
	var invoices []models.Invoice
	
	if err := r.db.WithContext(ctx).
		Where("status = ?", models.StatusPending).
		Find(&invoices).Error; err != nil {
		return nil, err
	}
	
	return invoices, nil
}

// Update updates an invoice
func (r *GORMInvoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	return r.db.WithContext(ctx).Save(invoice).Error
} 