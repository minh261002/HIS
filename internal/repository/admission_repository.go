package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// AdmissionRepository handles admission data operations
type AdmissionRepository struct {
	db *gorm.DB
}

// NewAdmissionRepository creates a new admission repository
func NewAdmissionRepository(db *gorm.DB) *AdmissionRepository {
	return &AdmissionRepository{db: db}
}

// Create creates a new admission
func (r *AdmissionRepository) Create(admission *domain.Admission) error {
	return r.db.Create(admission).Error
}

// FindByID finds an admission by ID with allocations and notes
func (r *AdmissionRepository) FindByID(id uint) (*domain.Admission, error) {
	var admission domain.Admission
	err := r.db.Preload("BedAllocations.Bed").
		Preload("NursingNotes.Nurse").
		Preload("Patient").
		Preload("Doctor").
		First(&admission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admission, nil
}

// FindByCode finds an admission by code
func (r *AdmissionRepository) FindByCode(code string) (*domain.Admission, error) {
	var admission domain.Admission
	err := r.db.Preload("BedAllocations.Bed").
		Preload("NursingNotes.Nurse").
		Preload("Patient").
		Preload("Doctor").
		Where("admission_code = ?", code).
		First(&admission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admission, nil
}

// FindByPatientID finds admissions for a patient
func (r *AdmissionRepository) FindByPatientID(patientID uint) ([]*domain.Admission, error) {
	var admissions []*domain.Admission
	err := r.db.Preload("BedAllocations.Bed").
		Where("patient_id = ?", patientID).
		Order("admission_date DESC").
		Find(&admissions).Error
	return admissions, err
}

// FindActiveAdmissions finds all active admissions
func (r *AdmissionRepository) FindActiveAdmissions() ([]*domain.Admission, error) {
	var admissions []*domain.Admission
	err := r.db.Preload("Patient").
		Preload("BedAllocations.Bed", "is_current = ?", true).
		Where("status = ?", domain.AdmissionStatusAdmitted).
		Order("admission_date DESC").
		Find(&admissions).Error
	return admissions, err
}

// Update updates an admission
func (r *AdmissionRepository) Update(admission *domain.Admission) error {
	return r.db.Save(admission).Error
}

// GenerateAdmissionCode generates a unique admission code
func (r *AdmissionRepository) GenerateAdmissionCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("ADM-%s-", today)

	var lastAdmission domain.Admission
	err := r.db.Where("admission_code LIKE ?", prefix+"%").
		Order("admission_code DESC").
		First(&lastAdmission).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastAdmission.AdmissionCode != "" {
		var lastSeq int
		fmt.Sscanf(lastAdmission.AdmissionCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
