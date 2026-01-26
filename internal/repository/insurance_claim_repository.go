package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// InsuranceClaimRepository handles insurance claim data operations
type InsuranceClaimRepository struct {
	db *gorm.DB
}

// NewInsuranceClaimRepository creates a new insurance claim repository
func NewInsuranceClaimRepository(db *gorm.DB) *InsuranceClaimRepository {
	return &InsuranceClaimRepository{db: db}
}

// Create creates insurance claim
func (r *InsuranceClaimRepository) Create(claim *domain.InsuranceClaim) error {
	return r.db.Create(claim).Error
}

// FindByInvoiceID finds claims for an invoice
func (r *InsuranceClaimRepository) FindByInvoiceID(invoiceID uint) ([]*domain.InsuranceClaim, error) {
	var claims []*domain.InsuranceClaim
	err := r.db.Where("invoice_id = ?", invoiceID).
		Order("claim_date DESC").
		Find(&claims).Error
	return claims, err
}

// FindByPatientID finds claims for a patient
func (r *InsuranceClaimRepository) FindByPatientID(patientID uint) ([]*domain.InsuranceClaim, error) {
	var claims []*domain.InsuranceClaim
	err := r.db.Preload("Invoice").
		Where("patient_id = ?", patientID).
		Order("claim_date DESC").
		Find(&claims).Error
	return claims, err
}

// FindByID finds claim by ID
func (r *InsuranceClaimRepository) FindByID(id uint) (*domain.InsuranceClaim, error) {
	var claim domain.InsuranceClaim
	err := r.db.Preload("Invoice").
		Preload("Patient").
		First(&claim, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &claim, nil
}

// Update updates claim
func (r *InsuranceClaimRepository) Update(claim *domain.InsuranceClaim) error {
	return r.db.Save(claim).Error
}

// GenerateClaimCode generates a unique claim code
func (r *InsuranceClaimRepository) GenerateClaimCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("CLM-%s-", today)

	var lastClaim domain.InsuranceClaim
	err := r.db.Where("claim_code LIKE ?", prefix+"%").
		Order("claim_code DESC").
		First(&lastClaim).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastClaim.ClaimCode != "" {
		var lastSeq int
		fmt.Sscanf(lastClaim.ClaimCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
