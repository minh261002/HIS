package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// InvoiceRepository handles invoice data operations
type InvoiceRepository struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new invoice repository
func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

// Create creates invoice with items
func (r *InvoiceRepository) Create(invoice *domain.Invoice) error {
	return r.db.Create(invoice).Error
}

// FindByID finds invoice by ID with relationships
func (r *InvoiceRepository) FindByID(id uint) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.Preload("Patient").
		Preload("Visit").
		Preload("Items").
		Preload("Payments").
		Preload("InsuranceClaims").
		First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// FindByCode finds invoice by code
func (r *InvoiceRepository) FindByCode(code string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.Preload("Patient").
		Preload("Visit").
		Preload("Items").
		Preload("Payments").
		Preload("InsuranceClaims").
		Where("invoice_code = ?", code).
		First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// FindByPatientID finds invoices for a patient
func (r *InvoiceRepository) FindByPatientID(patientID uint) ([]*domain.Invoice, error) {
	var invoices []*domain.Invoice
	err := r.db.Preload("Items").
		Where("patient_id = ?", patientID).
		Order("invoice_date DESC").
		Find(&invoices).Error
	return invoices, err
}

// Update updates invoice
func (r *InvoiceRepository) Update(invoice *domain.Invoice) error {
	return r.db.Save(invoice).Error
}

// GenerateInvoiceCode generates a unique invoice code
func (r *InvoiceRepository) GenerateInvoiceCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("INV-%s-", today)

	var lastInvoice domain.Invoice
	err := r.db.Where("invoice_code LIKE ?", prefix+"%").
		Order("invoice_code DESC").
		First(&lastInvoice).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastInvoice.InvoiceCode != "" {
		var lastSeq int
		fmt.Sscanf(lastInvoice.InvoiceCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
