package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// DispensingRepository handles dispensing data operations
type DispensingRepository struct {
	db *gorm.DB
}

// NewDispensingRepository creates a new dispensing repository
func NewDispensingRepository(db *gorm.DB) *DispensingRepository {
	return &DispensingRepository{db: db}
}

// Create creates dispensing record
func (r *DispensingRepository) Create(dispensing *domain.Dispensing) error {
	return r.db.Create(dispensing).Error
}

// FindByPrescriptionID finds dispensing records for a prescription
func (r *DispensingRepository) FindByPrescriptionID(prescriptionID uint) ([]*domain.Dispensing, error) {
	var dispensings []*domain.Dispensing
	err := r.db.Preload("Medication").
		Preload("Pharmacist").
		Where("prescription_id = ?", prescriptionID).
		Order("dispensed_date DESC").
		Find(&dispensings).Error
	return dispensings, err
}

// FindByPatientID finds dispensing history for a patient
func (r *DispensingRepository) FindByPatientID(patientID uint) ([]*domain.Dispensing, error) {
	var dispensings []*domain.Dispensing
	err := r.db.Preload("Medication").
		Preload("Pharmacist").
		Where("patient_id = ?", patientID).
		Order("dispensed_date DESC").
		Find(&dispensings).Error
	return dispensings, err
}

// GenerateDispensingCode generates a unique dispensing code
func (r *DispensingRepository) GenerateDispensingCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("DIS-%s-", today)

	var lastDispensing domain.Dispensing
	err := r.db.Where("dispensing_code LIKE ?", prefix+"%").
		Order("dispensing_code DESC").
		First(&lastDispensing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastDispensing.DispensingCode != "" {
		var lastSeq int
		fmt.Sscanf(lastDispensing.DispensingCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
