package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates payment
func (r *PaymentRepository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

// FindByInvoiceID finds payments for an invoice
func (r *PaymentRepository) FindByInvoiceID(invoiceID uint) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.Where("invoice_id = ?", invoiceID).
		Order("payment_date DESC").
		Find(&payments).Error
	return payments, err
}

// FindByPatientID finds payment history for a patient
func (r *PaymentRepository) FindByPatientID(patientID uint) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.Preload("Invoice").
		Where("patient_id = ?", patientID).
		Order("payment_date DESC").
		Find(&payments).Error
	return payments, err
}

// Update updates payment
func (r *PaymentRepository) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

// GeneratePaymentCode generates a unique payment code
func (r *PaymentRepository) GeneratePaymentCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("PAY-%s-", today)

	var lastPayment domain.Payment
	err := r.db.Where("payment_code LIKE ?", prefix+"%").
		Order("payment_code DESC").
		First(&lastPayment).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastPayment.PaymentCode != "" {
		var lastSeq int
		fmt.Sscanf(lastPayment.PaymentCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
