package services

import (
	"fmt"
	"time"

	"github.com/ncapetillo/demo-fluida/internal/models"
)

// InvoiceService handles business logic for invoices
type InvoiceService struct {
	// In a real application, this would have a database repository dependency
	mockInvoices []models.Invoice
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService() *InvoiceService {
	// Initialize with a sample invoice for demonstration purposes
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

	return &InvoiceService{
		mockInvoices: []models.Invoice{sampleInvoice},
	}
}

// GetAllInvoices returns all invoices
func (s *InvoiceService) GetAllInvoices() []models.Invoice {
	return s.mockInvoices
}

// GetInvoiceByToken retrieves an invoice by its payment link token
func (s *InvoiceService) GetInvoiceByToken(token string) (models.Invoice, error) {
	// For demo purposes, just return the first invoice
	// In a real app, this would query the database
	if token == "demo-token" || len(s.mockInvoices) > 0 {
		return s.mockInvoices[0], nil
	}
	
	return models.Invoice{}, fmt.Errorf("invoice not found: %s", token)
}

// CreateInvoice creates a new invoice
func (s *InvoiceService) CreateInvoice(req models.CreateInvoiceRequest) models.Invoice {
	// Create a new invoice from the request
	newInvoice := models.NewInvoice(req)
	
	// Set the ID (in a real app, the database would do this)
	newInvoice.ID = len(s.mockInvoices) + 1
	
	// Add to our mock database
	s.mockInvoices = append(s.mockInvoices, newInvoice)
	
	return newInvoice
}

// UpdateInvoiceStatus updates the status of an invoice
func (s *InvoiceService) UpdateInvoiceStatus(id int, status models.InvoiceStatus) (models.Invoice, error) {
	// Find the invoice
	for i, inv := range s.mockInvoices {
		if inv.ID == id {
			// Update the status
			s.mockInvoices[i].Status = status
			s.mockInvoices[i].UpdatedAt = time.Now()
			
			return s.mockInvoices[i], nil
		}
	}
	
	return models.Invoice{}, fmt.Errorf("invoice not found: %d", id)
} 