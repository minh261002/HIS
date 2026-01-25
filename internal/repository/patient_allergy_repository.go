package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PatientAllergyRepository handles patient allergy data operations
type PatientAllergyRepository struct {
	db *gorm.DB
}

// NewPatientAllergyRepository creates a new patient allergy repository
func NewPatientAllergyRepository(db *gorm.DB) *PatientAllergyRepository {
	return &PatientAllergyRepository{db: db}
}

// Create creates a new allergy record
func (r *PatientAllergyRepository) Create(allergy *domain.PatientAllergy) error {
	return r.db.Create(allergy).Error
}

// FindByID finds an allergy by ID
func (r *PatientAllergyRepository) FindByID(id uint) (*domain.PatientAllergy, error) {
	var allergy domain.PatientAllergy
	err := r.db.First(&allergy, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &allergy, nil
}

// FindByPatientID finds all allergies for a patient
func (r *PatientAllergyRepository) FindByPatientID(patientID uint) ([]*domain.PatientAllergy, error) {
	var allergies []*domain.PatientAllergy
	err := r.db.Where("patient_id = ?", patientID).Order("created_at DESC").Find(&allergies).Error
	if err != nil {
		return nil, err
	}
	return allergies, nil
}

// FindActiveByPatientID finds active allergies for a patient
func (r *PatientAllergyRepository) FindActiveByPatientID(patientID uint) ([]*domain.PatientAllergy, error) {
	var allergies []*domain.PatientAllergy
	err := r.db.Where("patient_id = ? AND is_active = ?", patientID, true).
		Order("severity DESC, created_at DESC").
		Find(&allergies).Error
	if err != nil {
		return nil, err
	}
	return allergies, nil
}

// FindBySeverity finds allergies by severity
func (r *PatientAllergyRepository) FindBySeverity(patientID uint, severity domain.AllergySeverity) ([]*domain.PatientAllergy, error) {
	var allergies []*domain.PatientAllergy
	err := r.db.Where("patient_id = ? AND severity = ?", patientID, severity).
		Order("created_at DESC").
		Find(&allergies).Error
	if err != nil {
		return nil, err
	}
	return allergies, nil
}

// Update updates an allergy record
func (r *PatientAllergyRepository) Update(allergy *domain.PatientAllergy) error {
	return r.db.Save(allergy).Error
}

// Delete soft deletes an allergy record
func (r *PatientAllergyRepository) Delete(id uint) error {
	return r.db.Delete(&domain.PatientAllergy{}, id).Error
}
