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
	ErrPaymentExceedsBalance = errors.New("payment amount exceeds invoice balance")
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	invoiceRepo *repository.InvoiceRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	invoiceRepo *repository.InvoiceRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
	}
}

// CreatePayment processes payment
func (s *PaymentService) CreatePayment(req *dto.CreatePaymentRequest, createdBy uint) (*dto.PaymentResponse, error) {
	// Validate invoice exists
	invoice, err := s.invoiceRepo.FindByID(req.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find invoice: %w", err)
	}
	if invoice == nil {
		return nil, ErrInvoiceNotFound
	}

	// Calculate paid amount
	var paidAmount float64
	for _, payment := range invoice.Payments {
		if payment.Status == domain.PaymentStatusCompleted {
			paidAmount += payment.Amount
		}
	}

	// Check if payment exceeds balance
	balance := invoice.TotalAmount - paidAmount
	if req.Amount > balance {
		return nil, ErrPaymentExceedsBalance
	}

	// Generate payment code
	code, err := s.paymentRepo.GeneratePaymentCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment code: %w", err)
	}

	// Create payment
	payment := &domain.Payment{
		PaymentCode:   code,
		InvoiceID:     req.InvoiceID,
		PatientID:     invoice.PatientID,
		PaymentMethod: domain.PaymentMethod(req.PaymentMethod),
		Amount:        req.Amount,
		PaymentDate:   time.Now(),
		Notes:         req.Notes,
		CreatedBy:     createdBy,
	}

	// Process based on payment method
	switch payment.PaymentMethod {
	case domain.PaymentMethodCash:
		// Cash payment is immediate
		payment.Status = domain.PaymentStatusCompleted
	case domain.PaymentMethodVNPay, domain.PaymentMethodPayOS:
		// Payment gateway - pending until callback
		payment.Status = domain.PaymentStatusPending
		// TODO: Create payment URL with gateway
		// payment.TransactionID = gatewayTransactionID
	default:
		payment.Status = domain.PaymentStatusPending
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update invoice status if fully paid
	if payment.Status == domain.PaymentStatusCompleted {
		newPaidAmount := paidAmount + payment.Amount
		if newPaidAmount >= invoice.TotalAmount {
			invoice.Status = domain.InvoiceStatusPaid
		} else if newPaidAmount > 0 {
			invoice.Status = domain.InvoiceStatusPartiallyPaid
		}
		s.invoiceRepo.Update(invoice)
	}

	return s.toPaymentResponse(payment), nil
}

// GetInvoicePayments gets invoice's payments
func (s *PaymentService) GetInvoicePayments(invoiceID uint) ([]*dto.PaymentListItem, error) {
	payments, err := s.paymentRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice payments: %w", err)
	}

	items := make([]*dto.PaymentListItem, len(payments))
	for i, p := range payments {
		items[i] = s.toPaymentListItem(p)
	}
	return items, nil
}

// Helper functions
func (s *PaymentService) toPaymentResponse(p *domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:            p.ID,
		PaymentCode:   p.PaymentCode,
		InvoiceID:     p.InvoiceID,
		PatientID:     p.PatientID,
		PaymentMethod: string(p.PaymentMethod),
		Amount:        p.Amount,
		PaymentDate:   p.PaymentDate,
		Status:        string(p.Status),
		TransactionID: p.TransactionID,
		Notes:         p.Notes,
		CreatedAt:     p.CreatedAt,
	}
}

func (s *PaymentService) toPaymentListItem(p *domain.Payment) *dto.PaymentListItem {
	return &dto.PaymentListItem{
		ID:            p.ID,
		PaymentCode:   p.PaymentCode,
		PaymentMethod: string(p.PaymentMethod),
		Amount:        p.Amount,
		Status:        string(p.Status),
		PaymentDate:   p.PaymentDate,
	}
}
