package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PatientMedicalHistoryRepository handles patient medical history data operations
type PatientMedicalHistoryRepository struct {
	db *gorm.DB
}

// NewPatientMedicalHistoryRepository creates a new patient medical history repository
func NewPatientMedicalHistoryRepository(db *gorm.DB) *PatientMedicalHistoryRepository {
	return &PatientMedicalHistoryRepository{db: db}
}

// Create creates a new medical history record
func (r *PatientMedicalHistoryRepository) Create(history *domain.PatientMedicalHistory) error {
	return r.db.Create(history).Error
}

// FindByID finds a medical history record by ID
func (r *PatientMedicalHistoryRepository) FindByID(id uint) (*domain.PatientMedicalHistory, error) {
	var history domain.PatientMedicalHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

// FindByPatientID finds all medical history for a patient
func (r *PatientMedicalHistoryRepository) FindByPatientID(patientID uint) ([]*domain.PatientMedicalHistory, error) {
	var histories []*domain.PatientMedicalHistory
	err := r.db.Where("patient_id = ?", patientID).Order("diagnosis_date DESC, created_at DESC").Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindActiveByPatientID finds active medical conditions for a patient
func (r *PatientMedicalHistoryRepository) FindActiveByPatientID(patientID uint) ([]*domain.PatientMedicalHistory, error) {
	var histories []*domain.PatientMedicalHistory
	err := r.db.Where("patient_id = ? AND is_active = ?", patientID, true).
		Order("diagnosis_date DESC, created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindByConditionType finds medical history by condition type
func (r *PatientMedicalHistoryRepository) FindByConditionType(patientID uint, conditionType domain.ConditionType) ([]*domain.PatientMedicalHistory, error) {
	var histories []*domain.PatientMedicalHistory
	err := r.db.Where("patient_id = ? AND condition_type = ?", patientID, conditionType).
		Order("diagnosis_date DESC, created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindByStatus finds medical history by status
func (r *PatientMedicalHistoryRepository) FindByStatus(patientID uint, status domain.ConditionStatus) ([]*domain.PatientMedicalHistory, error) {
	var histories []*domain.PatientMedicalHistory
	err := r.db.Where("patient_id = ? AND status = ?", patientID, status).
		Order("diagnosis_date DESC, created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// Update updates a medical history record
func (r *PatientMedicalHistoryRepository) Update(history *domain.PatientMedicalHistory) error {
	return r.db.Save(history).Error
}

// Delete soft deletes a medical history record
func (r *PatientMedicalHistoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.PatientMedicalHistory{}, id).Error
}
