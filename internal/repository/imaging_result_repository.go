package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// ImagingResultRepository handles imaging result data operations
type ImagingResultRepository struct {
	db *gorm.DB
}

// NewImagingResultRepository creates a new imaging result repository
func NewImagingResultRepository(db *gorm.DB) *ImagingResultRepository {
	return &ImagingResultRepository{db: db}
}

// Create creates a result
func (r *ImagingResultRepository) Create(result *domain.ImagingResult) error {
	return r.db.Create(result).Error
}

// FindByRequestID finds result for a request
func (r *ImagingResultRepository) FindByRequestID(requestID uint) (*domain.ImagingResult, error) {
	var result domain.ImagingResult
	err := r.db.Preload("Radiologist").
		Where("request_id = ?", requestID).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// Update updates a result
func (r *ImagingResultRepository) Update(result *domain.ImagingResult) error {
	return r.db.Save(result).Error
}
