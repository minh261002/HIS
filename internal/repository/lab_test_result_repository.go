package repository

import (
	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// LabTestResultRepository handles lab test result data operations
type LabTestResultRepository struct {
	db *gorm.DB
}

// NewLabTestResultRepository creates a new lab test result repository
func NewLabTestResultRepository(db *gorm.DB) *LabTestResultRepository {
	return &LabTestResultRepository{db: db}
}

// CreateBatch creates multiple results in a transaction
func (r *LabTestResultRepository) CreateBatch(results []*domain.LabTestResult) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, result := range results {
			if err := tx.Create(result).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByRequestID finds all results for a request
func (r *LabTestResultRepository) FindByRequestID(requestID uint) ([]*domain.LabTestResult, error) {
	var results []*domain.LabTestResult
	err := r.db.Where("request_id = ?", requestID).Find(&results).Error
	return results, err
}

// Update updates a result
func (r *LabTestResultRepository) Update(result *domain.LabTestResult) error {
	return r.db.Save(result).Error
}

// UpdateBatch updates multiple results in a transaction
func (r *LabTestResultRepository) UpdateBatch(results []*domain.LabTestResult) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, result := range results {
			if err := tx.Save(result).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
