package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// DiagnosisRepository handles diagnosis data operations
type DiagnosisRepository struct {
	db *gorm.DB
}

// NewDiagnosisRepository creates a new diagnosis repository
func NewDiagnosisRepository(db *gorm.DB) *DiagnosisRepository {
	return &DiagnosisRepository{db: db}
}

// Create creates a new diagnosis
func (r *DiagnosisRepository) Create(diagnosis *domain.Diagnosis) error {
	return r.db.Create(diagnosis).Error
}

// FindByID finds a diagnosis by ID
func (r *DiagnosisRepository) FindByID(id uint) (*domain.Diagnosis, error) {
	var diagnosis domain.Diagnosis
	err := r.db.Preload("ICD10Code").Preload("Doctor").Preload("Patient").
		First(&diagnosis, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &diagnosis, nil
}

// FindByVisitID finds all diagnoses for a visit
func (r *DiagnosisRepository) FindByVisitID(visitID uint) ([]*domain.Diagnosis, error) {
	var diagnoses []*domain.Diagnosis
	err := r.db.Preload("ICD10Code").Preload("Doctor").
		Where("visit_id = ?", visitID).
		Order("diagnosis_type ASC, diagnosed_at DESC").
		Find(&diagnoses).Error
	return diagnoses, err
}

// FindPrimaryDiagnosis finds the primary diagnosis for a visit
func (r *DiagnosisRepository) FindPrimaryDiagnosis(visitID uint) (*domain.Diagnosis, error) {
	var diagnosis domain.Diagnosis
	err := r.db.Preload("ICD10Code").Preload("Doctor").
		Where("visit_id = ? AND diagnosis_type = ?", visitID, domain.DiagnosisTypePrimary).
		First(&diagnosis).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &diagnosis, nil
}

// FindByPatientID finds diagnoses for a patient
func (r *DiagnosisRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.Diagnosis, error) {
	var diagnoses []*domain.Diagnosis
	query := r.db.Preload("ICD10Code").Preload("Doctor").Preload("Visit").
		Where("patient_id = ?", patientID)

	if diagnosisType, ok := filters["diagnosis_type"]; ok && diagnosisType != "" {
		query = query.Where("diagnosis_type = ?", diagnosisType)
	}
	if status, ok := filters["diagnosis_status"]; ok && status != "" {
		query = query.Where("diagnosis_status = ?", status)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("diagnosed_at >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("diagnosed_at <= ?", toDate)
	}

	err := query.Order("diagnosed_at DESC").Find(&diagnoses).Error
	return diagnoses, err
}

// Update updates a diagnosis
func (r *DiagnosisRepository) Update(diagnosis *domain.Diagnosis) error {
	return r.db.Save(diagnosis).Error
}

// Delete soft deletes a diagnosis
func (r *DiagnosisRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Diagnosis{}, id).Error
}
