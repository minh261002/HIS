package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// LabTestTemplateRepository handles lab test template data operations
type LabTestTemplateRepository struct {
	db *gorm.DB
}

// NewLabTestTemplateRepository creates a new lab test template repository
func NewLabTestTemplateRepository(db *gorm.DB) *LabTestTemplateRepository {
	return &LabTestTemplateRepository{db: db}
}

// FindByID finds a template by ID with parameters
func (r *LabTestTemplateRepository) FindByID(id uint) (*domain.LabTestTemplate, error) {
	var template domain.LabTestTemplate
	err := r.db.Preload("Parameters", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Where("is_active = ?", true).First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// FindByCode finds a template by code
func (r *LabTestTemplateRepository) FindByCode(code string) (*domain.LabTestTemplate, error) {
	var template domain.LabTestTemplate
	err := r.db.Preload("Parameters", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Where("code = ? AND is_active = ?", code, true).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// Search searches templates by name or code
func (r *LabTestTemplateRepository) Search(query string, limit int) ([]*domain.LabTestTemplate, error) {
	var templates []*domain.LabTestTemplate

	if limit == 0 {
		limit = 20
	}

	err := r.db.Where("is_active = ? AND (name LIKE ? OR code LIKE ?)",
		true, "%"+query+"%", "%"+query+"%").
		Order("name ASC").
		Limit(limit).
		Find(&templates).Error

	return templates, err
}

// FindByCategory finds templates by category
func (r *LabTestTemplateRepository) FindByCategory(category string) ([]*domain.LabTestTemplate, error) {
	var templates []*domain.LabTestTemplate
	err := r.db.Where("category = ? AND is_active = ?", category, true).
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}
