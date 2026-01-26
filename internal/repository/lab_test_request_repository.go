package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// LabTestRequestRepository handles lab test request data operations
type LabTestRequestRepository struct {
	db *gorm.DB
}

// NewLabTestRequestRepository creates a new lab test request repository
func NewLabTestRequestRepository(db *gorm.DB) *LabTestRequestRepository {
	return &LabTestRequestRepository{db: db}
}

// Create creates a new lab test request
func (r *LabTestRequestRepository) Create(request *domain.LabTestRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a request by ID with results
func (r *LabTestRequestRepository) FindByID(id uint) (*domain.LabTestRequest, error) {
	var request domain.LabTestRequest
	err := r.db.Preload("Template.Parameters", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Preload("Results").
		Preload("Patient").
		Preload("Doctor").
		First(&request, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// FindByCode finds a request by code
func (r *LabTestRequestRepository) FindByCode(code string) (*domain.LabTestRequest, error) {
	var request domain.LabTestRequest
	err := r.db.Preload("Template.Parameters", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Preload("Results").
		Preload("Patient").
		Preload("Doctor").
		Where("request_code = ?", code).
		First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// FindByVisitID finds requests for a visit
func (r *LabTestRequestRepository) FindByVisitID(visitID uint) ([]*domain.LabTestRequest, error) {
	var requests []*domain.LabTestRequest
	err := r.db.Preload("Template").Preload("Results").
		Where("visit_id = ?", visitID).
		Order("requested_date DESC").
		Find(&requests).Error
	return requests, err
}

// FindByPatientID finds requests for a patient
func (r *LabTestRequestRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.LabTestRequest, error) {
	var requests []*domain.LabTestRequest
	query := r.db.Preload("Template").Preload("Results").
		Where("patient_id = ?", patientID)

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("requested_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("requested_date <= ?", toDate)
	}

	err := query.Order("requested_date DESC").Find(&requests).Error
	return requests, err
}

// Update updates a request
func (r *LabTestRequestRepository) Update(request *domain.LabTestRequest) error {
	return r.db.Save(request).Error
}

// Delete soft deletes a request
func (r *LabTestRequestRepository) Delete(id uint) error {
	return r.db.Delete(&domain.LabTestRequest{}, id).Error
}

// GenerateRequestCode generates a unique request code
func (r *LabTestRequestRepository) GenerateRequestCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("LTR-%s-", today)

	var lastRequest domain.LabTestRequest
	err := r.db.Where("request_code LIKE ?", prefix+"%").
		Order("request_code DESC").
		First(&lastRequest).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastRequest.RequestCode != "" {
		var lastSeq int
		fmt.Sscanf(lastRequest.RequestCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
