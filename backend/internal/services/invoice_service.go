package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/models"
	"github.com/ncapetillo/demo-fluida/internal/repository"
	"gorm.io/gorm"
)

// InvoiceService handles business logic for invoices
type InvoiceService struct {
	db         *gorm.DB
	repository repository.InvoiceRepository
	mockMode   bool
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(repo repository.InvoiceRepository) *InvoiceService {
	// Check if we're in development mode with mock data
	mockMode := false
	
	service := &InvoiceService{
		mockMode: mockMode,
	}
	
	if mockMode {
		log.Println("Using mock service for testing - database connection disabled")
		service.db = nil
		service.repository = nil
	} else {
		service.db = db.DB
		service.repository = repo
	}
	
	return service
}

// GetAllInvoices returns all invoices
func (s *InvoiceService) GetAllInvoices() []models.Invoice {
	if s.mockMode {
		return createMockInvoices()
	}
	
	// Use context with timeout for database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Default pagination
	invoices, err := s.repository.List(ctx, 1, 100)
	if err != nil {
		log.Printf("Error fetching invoices: %v", err)
		return []models.Invoice{}
	}
	
	return invoices
}

// GetInvoiceByToken retrieves an invoice by its payment link token
func (s *InvoiceService) GetInvoiceByToken(token string) (models.Invoice, error) {
	if s.mockMode {
		mockInvoices := createMockInvoices()
		if token == "demo-token" || len(mockInvoices) > 0 {
			return mockInvoices[0], nil
		}
		return models.Invoice{}, fmt.Errorf("invoice not found: %s", token)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	invoice, err := s.repository.FindByLinkToken(ctx, token)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to get invoice: %w", err)
	}
	
	if invoice == nil {
		return models.Invoice{}, fmt.Errorf("invoice not found: %s", token)
	}
	
	return *invoice, nil
}

// CreateInvoice creates a new invoice
func (s *InvoiceService) CreateInvoice(req models.CreateInvoiceRequest) (models.Invoice, error) {
	// Create a new invoice from the request
	newInvoice := models.NewInvoice(req)
	
	// Validate the invoice
	if err := newInvoice.Validate(); err != nil {
		return models.Invoice{}, fmt.Errorf("invalid invoice data: %w", err)
	}
	
	if s.mockMode {
		mockInvoices := createMockInvoices()
		newInvoice.ID = len(mockInvoices) + 1
		return newInvoice, nil
	}
	
	// Use transaction for safe creation
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Check if invoice number already exists
		txRepo := repository.NewInvoiceRepository(tx)
		ctx := context.Background()
		
		existing, err := txRepo.FindByInvoiceNumber(ctx, newInvoice.InvoiceNumber)
		if err != nil {
			return err
		}
		
		if existing != nil {
			// Return this error directly without additional wrapping
			return fmt.Errorf("invoice number %s already exists. Please use a different invoice number", newInvoice.InvoiceNumber)
		}
		
		// Create the invoice
		return txRepo.Create(ctx, &newInvoice)
	})
	
	if err != nil {
		// Don't wrap errors that already contain specific messages
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "failed to create invoice") {
			return models.Invoice{}, err
		}
		// Only wrap generic errors
		return models.Invoice{}, fmt.Errorf("failed to create invoice: %w", err)
	}
	
	return newInvoice, nil
}

// UpdateInvoiceStatus updates the status of an invoice
func (s *InvoiceService) UpdateInvoiceStatus(id int, status models.InvoiceStatus) (models.Invoice, error) {
	if s.mockMode {
		mockInvoices := createMockInvoices()
		for i, inv := range mockInvoices {
			if inv.ID == id {
				mockInvoices[i].Status = status
				mockInvoices[i].UpdatedAt = time.Now()
				return mockInvoices[i], nil
			}
		}
		return models.Invoice{}, fmt.Errorf("invoice not found: %d", id)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Use transaction for safe update
	var result models.Invoice
	
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewInvoiceRepository(tx)
		
		// Find the invoice
		invoice, err := txRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		
		if invoice == nil {
			return fmt.Errorf("invoice not found: %d", id)
		}
		
		// Update the status
		invoice.Status = status
		invoice.UpdatedAt = time.Now()
		
		if err := txRepo.Update(ctx, invoice); err != nil {
			return err
		}
		
		result = *invoice
		return nil
	})
	
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to update invoice status: %w", err)
	}
	
	return result, nil
}

// GetPendingInvoices returns all pending invoices
func (s *InvoiceService) GetPendingInvoices() ([]models.Invoice, error) {
	if s.mockMode {
		var pending []models.Invoice
		for _, inv := range createMockInvoices() {
			if inv.Status == models.StatusPending {
				pending = append(pending, inv)
			}
		}
		return pending, nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return s.repository.FindPendingInvoices(ctx)
}

// Helper function to create mock invoices
func createMockInvoices() []models.Invoice {
	sampleInvoice := models.Invoice{
		ID:            1,
		InvoiceNumber: "INV-001",
		Amount:        100,
		Currency:      "USDC",
		Description:   "Demo invoice",
		DueDate:       time.Now().AddDate(0, 0, 7),
		Status:        models.StatusPending,
		ReceiverAddr:  "8JQxYTKfELQhAJL4c3jQvnUuNwZWJxsJr7o8G6iTfEV9",
		LinkToken:     "demo-token",
		SenderDetails: models.Person{
			Name:    "Acme Inc",
			Email:   "billing@acme.com",
			Address: "123 Business Ave, Suite 100, San Francisco, CA 94107",
		},
		RecipientDetails: models.Person{
			Name:    "John Smith",
			Email:   "john@example.com",
			Address: "456 Customer Lane, Apt 303, New York, NY 10001",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return []models.Invoice{sampleInvoice}
} 