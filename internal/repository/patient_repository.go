package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PatientRepository handles patient data operations
type PatientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

// Create creates a new patient
func (r *PatientRepository) Create(patient *domain.Patient) error {
	return r.db.Create(patient).Error
}

// FindByID finds a patient by ID
func (r *PatientRepository) FindByID(id uint) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.First(&patient, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

// FindByPatientCode finds a patient by patient code
func (r *PatientRepository) FindByPatientCode(code string) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.Where("patient_code = ?", code).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

// FindByNationalID finds a patient by national ID
func (r *PatientRepository) FindByNationalID(nationalID string) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.Where("national_id = ?", nationalID).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

// Update updates a patient
func (r *PatientRepository) Update(patient *domain.Patient) error {
	return r.db.Save(patient).Error
}

// Delete soft deletes a patient
func (r *PatientRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Patient{}, id).Error
}

// List returns a paginated list of patients with optional filters
func (r *PatientRepository) List(page, pageSize int, filters map[string]interface{}) ([]*domain.Patient, int64, error) {
	var patients []*domain.Patient
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&domain.Patient{})

	// Apply filters
	if gender, ok := filters["gender"]; ok && gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if bloodType, ok := filters["blood_type"]; ok && bloodType != "" {
		query = query.Where("blood_type = ?", bloodType)
	}
	if city, ok := filters["city"]; ok && city != "" {
		query = query.Where("city = ?", city)
	}
	if isActive, ok := filters["is_active"]; ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&patients).Error

	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// Search searches for patients by name, phone, email, or patient code
func (r *PatientRepository) Search(query string, page, pageSize int) ([]*domain.Patient, int64, error) {
	var patients []*domain.Patient
	var total int64

	offset := (page - 1) * pageSize

	searchQuery := r.db.Model(&domain.Patient{}).Where(
		"full_name LIKE ? OR phone_number LIKE ? OR email LIKE ? OR patient_code LIKE ? OR national_id LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
	)

	// Count total
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := searchQuery.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&patients).Error

	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// GetPatientStats returns patient statistics
func (r *PatientRepository) GetPatientStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total patients
	var totalPatients int64
	if err := r.db.Model(&domain.Patient{}).Count(&totalPatients).Error; err != nil {
		return nil, err
	}
	stats["total_patients"] = totalPatients

	// Active patients
	var activePatients int64
	if err := r.db.Model(&domain.Patient{}).Where("is_active = ?", true).Count(&activePatients).Error; err != nil {
		return nil, err
	}
	stats["active_patients"] = activePatients

	// New patients today
	today := time.Now().Truncate(24 * time.Hour)
	var newToday int64
	if err := r.db.Model(&domain.Patient{}).Where("DATE(created_at) = ?", today).Count(&newToday).Error; err != nil {
		return nil, err
	}
	stats["new_today"] = newToday

	// New patients this month
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	var newThisMonth int64
	if err := r.db.Model(&domain.Patient{}).Where("created_at >= ?", firstDayOfMonth).Count(&newThisMonth).Error; err != nil {
		return nil, err
	}
	stats["new_this_month"] = newThisMonth

	// Gender distribution
	var genderStats []struct {
		Gender domain.Gender
		Count  int64
	}
	if err := r.db.Model(&domain.Patient{}).
		Select("gender, COUNT(*) as count").
		Group("gender").
		Scan(&genderStats).Error; err != nil {
		return nil, err
	}
	stats["gender_distribution"] = genderStats

	return stats, nil
}

// GeneratePatientCode generates a unique patient code
func (r *PatientRepository) GeneratePatientCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("P-%s-", today)

	// Find the last patient code for today
	var lastPatient domain.Patient
	err := r.db.Where("patient_code LIKE ?", prefix+"%").
		Order("patient_code DESC").
		First(&lastPatient).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Generate next sequence number
	sequence := 1
	if lastPatient.PatientCode != "" {
		// Extract sequence from last code (P-YYYYMMDD-XXXX)
		var lastSeq int
		fmt.Sscanf(lastPatient.PatientCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
