package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// ICD10CodeRepository handles ICD-10 code data operations
type ICD10CodeRepository struct {
	db *gorm.DB
}

// NewICD10CodeRepository creates a new ICD-10 code repository
func NewICD10CodeRepository(db *gorm.DB) *ICD10CodeRepository {
	return &ICD10CodeRepository{db: db}
}

// FindByCode finds an ICD-10 code by code
func (r *ICD10CodeRepository) FindByCode(code string) (*domain.ICD10Code, error) {
	var icd10Code domain.ICD10Code
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&icd10Code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &icd10Code, nil
}

// FindByID finds an ICD-10 code by ID
func (r *ICD10CodeRepository) FindByID(id uint) (*domain.ICD10Code, error) {
	var icd10Code domain.ICD10Code
	err := r.db.First(&icd10Code, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &icd10Code, nil
}

// Search searches ICD-10 codes by code or description
func (r *ICD10CodeRepository) Search(query string, limit int) ([]*domain.ICD10Code, error) {
	var codes []*domain.ICD10Code

	if limit == 0 {
		limit = 20
	}

	err := r.db.Where("is_active = ? AND (code LIKE ? OR description LIKE ?)",
		true, "%"+query+"%", "%"+query+"%").
		Order("code ASC").
		Limit(limit).
		Find(&codes).Error

	return codes, err
}

// FindByCategory finds ICD-10 codes by category
func (r *ICD10CodeRepository) FindByCategory(category string) ([]*domain.ICD10Code, error) {
	var codes []*domain.ICD10Code
	err := r.db.Where("category = ? AND is_active = ?", category, true).
		Order("code ASC").
		Find(&codes).Error
	return codes, err
}
