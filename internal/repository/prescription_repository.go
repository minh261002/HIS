package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PrescriptionRepository handles prescription data operations
type PrescriptionRepository struct {
	db *gorm.DB
}

// NewPrescriptionRepository creates a new prescription repository
func NewPrescriptionRepository(db *gorm.DB) *PrescriptionRepository {
	return &PrescriptionRepository{db: db}
}

// Create creates a new prescription with items in a transaction
func (r *PrescriptionRepository) Create(prescription *domain.Prescription) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(prescription).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindByID finds a prescription by ID with items
func (r *PrescriptionRepository) FindByID(id uint) (*domain.Prescription, error) {
	var prescription domain.Prescription
	err := r.db.Preload("Items.Medication").
		Preload("Patient").
		Preload("Doctor").
		Preload("Diagnosis.ICD10Code").
		First(&prescription, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prescription, nil
}

// FindByCode finds a prescription by code
func (r *PrescriptionRepository) FindByCode(code string) (*domain.Prescription, error) {
	var prescription domain.Prescription
	err := r.db.Preload("Items.Medication").
		Preload("Patient").
		Preload("Doctor").
		Preload("Diagnosis.ICD10Code").
		Where("prescription_code = ?", code).
		First(&prescription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prescription, nil
}

// FindByVisitID finds prescriptions for a visit
func (r *PrescriptionRepository) FindByVisitID(visitID uint) ([]*domain.Prescription, error) {
	var prescriptions []*domain.Prescription
	err := r.db.Preload("Items.Medication").
		Preload("Doctor").
		Where("visit_id = ?", visitID).
		Order("prescribed_date DESC").
		Find(&prescriptions).Error
	return prescriptions, err
}

// FindByPatientID finds prescriptions for a patient
func (r *PrescriptionRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.Prescription, error) {
	var prescriptions []*domain.Prescription
	query := r.db.Preload("Items.Medication").Preload("Doctor").
		Where("patient_id = ?", patientID)

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("prescribed_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("prescribed_date <= ?", toDate)
	}

	err := query.Order("prescribed_date DESC").Find(&prescriptions).Error
	return prescriptions, err
}

// Update updates a prescription
func (r *PrescriptionRepository) Update(prescription *domain.Prescription) error {
	return r.db.Save(prescription).Error
}

// Delete soft deletes a prescription
func (r *PrescriptionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Prescription{}, id).Error
}

// GeneratePrescriptionCode generates a unique prescription code
func (r *PrescriptionRepository) GeneratePrescriptionCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("PRX-%s-", today)

	var lastPrescription domain.Prescription
	err := r.db.Where("prescription_code LIKE ?", prefix+"%").
		Order("prescription_code DESC").
		First(&lastPrescription).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastPrescription.PrescriptionCode != "" {
		var lastSeq int
		fmt.Sscanf(lastPrescription.PrescriptionCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
