package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// ImagingRequestRepository handles imaging request data operations
type ImagingRequestRepository struct {
	db *gorm.DB
}

// NewImagingRequestRepository creates a new imaging request repository
func NewImagingRequestRepository(db *gorm.DB) *ImagingRequestRepository {
	return &ImagingRequestRepository{db: db}
}

// Create creates a new imaging request
func (r *ImagingRequestRepository) Create(request *domain.ImagingRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a request by ID with result
func (r *ImagingRequestRepository) FindByID(id uint) (*domain.ImagingRequest, error) {
	var request domain.ImagingRequest
	err := r.db.Preload("Template").
		Preload("Result.Radiologist").
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
func (r *ImagingRequestRepository) FindByCode(code string) (*domain.ImagingRequest, error) {
	var request domain.ImagingRequest
	err := r.db.Preload("Template").
		Preload("Result.Radiologist").
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
func (r *ImagingRequestRepository) FindByVisitID(visitID uint) ([]*domain.ImagingRequest, error) {
	var requests []*domain.ImagingRequest
	err := r.db.Preload("Template").Preload("Result").
		Where("visit_id = ?", visitID).
		Order("requested_date DESC").
		Find(&requests).Error
	return requests, err
}

// FindByPatientID finds requests for a patient
func (r *ImagingRequestRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.ImagingRequest, error) {
	var requests []*domain.ImagingRequest
	query := r.db.Preload("Template").Preload("Result").
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
func (r *ImagingRequestRepository) Update(request *domain.ImagingRequest) error {
	return r.db.Save(request).Error
}

// Delete soft deletes a request
func (r *ImagingRequestRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ImagingRequest{}, id).Error
}

// GenerateRequestCode generates a unique request code
func (r *ImagingRequestRepository) GenerateRequestCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("IMG-%s-", today)

	var lastRequest domain.ImagingRequest
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
