package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
)

// InvoiceService handles invoice business logic
type InvoiceService struct {
	invoiceRepo *repository.InvoiceRepository
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(invoiceRepo *repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoiceRepo: invoiceRepo}
}

// CreateInvoice creates invoice with items
func (s *InvoiceService) CreateInvoice(req *dto.CreateInvoiceRequest, createdBy uint) (*dto.InvoiceResponse, error) {
	// Generate invoice code
	code, err := s.invoiceRepo.GenerateInvoiceCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice code: %w", err)
	}

	// Calculate subtotal
	var subtotal float64
	items := make([]*domain.InvoiceItem, len(req.Items))
	for i, itemReq := range req.Items {
		amount := float64(itemReq.Quantity) * itemReq.UnitPrice
		subtotal += amount

		items[i] = &domain.InvoiceItem{
			ItemType:    domain.ItemType(itemReq.ItemType),
			ItemID:      itemReq.ItemID,
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
			Amount:      amount,
		}
	}

	// Calculate total
	totalAmount := subtotal + req.TaxAmount - req.DiscountAmount

	// Create invoice
	now := time.Now()
	invoice := &domain.Invoice{
		InvoiceCode:    code,
		VisitID:        req.VisitID,
		PatientID:      req.PatientID,
		InvoiceDate:    now,
		DueDate:        now.AddDate(0, 0, 30), // 30 days from now
		Subtotal:       subtotal,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    totalAmount,
		Status:         domain.InvoiceStatusPending,
		Notes:          req.Notes,
		Items:          items,
		CreatedBy:      createdBy,
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Reload to get relationships
	invoice, _ = s.invoiceRepo.FindByID(invoice.ID)
	return s.toInvoiceResponse(invoice), nil
}

// GetInvoiceByID gets invoice by ID
func (s *InvoiceService) GetInvoiceByID(id uint) (*dto.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find invoice: %w", err)
	}
	if invoice == nil {
		return nil, ErrInvoiceNotFound
	}
	return s.toInvoiceResponse(invoice), nil
}

// GetInvoiceByCode gets invoice by code
func (s *InvoiceService) GetInvoiceByCode(code string) (*dto.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find invoice: %w", err)
	}
	if invoice == nil {
		return nil, ErrInvoiceNotFound
	}
	return s.toInvoiceResponse(invoice), nil
}

// GetPatientInvoices gets patient's invoices
func (s *InvoiceService) GetPatientInvoices(patientID uint) ([]*dto.InvoiceListItem, error) {
	invoices, err := s.invoiceRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient invoices: %w", err)
	}

	items := make([]*dto.InvoiceListItem, len(invoices))
	for i, inv := range invoices {
		items[i] = s.toInvoiceListItem(inv)
	}
	return items, nil
}

// Helper functions
func (s *InvoiceService) toInvoiceResponse(inv *domain.Invoice) *dto.InvoiceResponse {
	resp := &dto.InvoiceResponse{
		ID:             inv.ID,
		InvoiceCode:    inv.InvoiceCode,
		VisitID:        inv.VisitID,
		PatientID:      inv.PatientID,
		InvoiceDate:    inv.InvoiceDate,
		DueDate:        inv.DueDate,
		Subtotal:       inv.Subtotal,
		TaxAmount:      inv.TaxAmount,
		DiscountAmount: inv.DiscountAmount,
		TotalAmount:    inv.TotalAmount,
		Status:         string(inv.Status),
		Notes:          inv.Notes,
		CreatedAt:      inv.CreatedAt,
	}

	if inv.Patient != nil {
		resp.PatientName = inv.Patient.FullName
	}

	// Add items
	resp.Items = make([]*dto.InvoiceItemResponse, len(inv.Items))
	for i, item := range inv.Items {
		resp.Items[i] = &dto.InvoiceItemResponse{
			ID:          item.ID,
			ItemType:    string(item.ItemType),
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
		}
	}

	return resp
}

func (s *InvoiceService) toInvoiceListItem(inv *domain.Invoice) *dto.InvoiceListItem {
	item := &dto.InvoiceListItem{
		ID:          inv.ID,
		InvoiceCode: inv.InvoiceCode,
		TotalAmount: inv.TotalAmount,
		Status:      string(inv.Status),
		InvoiceDate: inv.InvoiceDate,
	}

	if inv.Patient != nil {
		item.PatientName = inv.Patient.FullName
	}

	return item
}
