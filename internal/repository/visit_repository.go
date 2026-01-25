package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// VisitRepository handles visit data operations
type VisitRepository struct {
	db *gorm.DB
}

// NewVisitRepository creates a new visit repository
func NewVisitRepository(db *gorm.DB) *VisitRepository {
	return &VisitRepository{db: db}
}

// Create creates a new visit
func (r *VisitRepository) Create(visit *domain.Visit) error {
	return r.db.Create(visit).Error
}

// FindByID finds a visit by ID
func (r *VisitRepository) FindByID(id uint) (*domain.Visit, error) {
	var visit domain.Visit
	err := r.db.Preload("Patient").Preload("Doctor").Preload("Appointment").First(&visit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &visit, nil
}

// FindByCode finds a visit by visit code
func (r *VisitRepository) FindByCode(code string) (*domain.Visit, error) {
	var visit domain.Visit
	err := r.db.Preload("Patient").Preload("Doctor").Preload("Appointment").
		Where("visit_code = ?", code).First(&visit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &visit, nil
}

// FindByAppointmentID finds a visit by appointment ID
func (r *VisitRepository) FindByAppointmentID(appointmentID uint) (*domain.Visit, error) {
	var visit domain.Visit
	err := r.db.Preload("Patient").Preload("Doctor").Preload("Appointment").
		Where("appointment_id = ?", appointmentID).First(&visit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &visit, nil
}

// FindByPatientID finds visits for a patient
func (r *VisitRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.Visit, error) {
	var visits []*domain.Visit
	query := r.db.Preload("Doctor").Where("patient_id = ?", patientID)

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("visit_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("visit_date <= ?", toDate)
	}

	err := query.Order("visit_date DESC, visit_time DESC").Find(&visits).Error
	return visits, err
}

// FindByDoctorID finds visits for a doctor on a specific date
func (r *VisitRepository) FindByDoctorID(doctorID uint, date time.Time) ([]*domain.Visit, error) {
	var visits []*domain.Visit
	err := r.db.Preload("Patient").
		Where("doctor_id = ? AND visit_date = ?", doctorID, date.Format("2006-01-02")).
		Order("visit_time ASC").
		Find(&visits).Error
	return visits, err
}

// Search searches visits with filters
func (r *VisitRepository) Search(filters map[string]interface{}, page, pageSize int) ([]*domain.Visit, int64, error) {
	var visits []*domain.Visit
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&domain.Visit{}).Preload("Patient").Preload("Doctor")

	if patientID, ok := filters["patient_id"]; ok && patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if doctorID, ok := filters["doctor_id"]; ok && doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if visitType, ok := filters["visit_type"]; ok && visitType != "" {
		query = query.Where("visit_type = ?", visitType)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("visit_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("visit_date <= ?", toDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).
		Limit(pageSize).
		Order("visit_date DESC, visit_time DESC").
		Find(&visits).Error

	return visits, total, err
}

// Update updates a visit
func (r *VisitRepository) Update(visit *domain.Visit) error {
	return r.db.Save(visit).Error
}

// Delete soft deletes a visit
func (r *VisitRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Visit{}, id).Error
}

// GenerateVisitCode generates a unique visit code
func (r *VisitRepository) GenerateVisitCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("VST-%s-", today)

	var lastVisit domain.Visit
	err := r.db.Where("visit_code LIKE ?", prefix+"%").
		Order("visit_code DESC").
		First(&lastVisit).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastVisit.VisitCode != "" {
		var lastSeq int
		fmt.Sscanf(lastVisit.VisitCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
