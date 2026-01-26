package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// ImagingTemplateRepository handles imaging template data operations
type ImagingTemplateRepository struct {
	db *gorm.DB
}

// NewImagingTemplateRepository creates a new imaging template repository
func NewImagingTemplateRepository(db *gorm.DB) *ImagingTemplateRepository {
	return &ImagingTemplateRepository{db: db}
}

// FindByID finds a template by ID
func (r *ImagingTemplateRepository) FindByID(id uint) (*domain.ImagingTemplate, error) {
	var template domain.ImagingTemplate
	err := r.db.Where("is_active = ?", true).First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// FindByCode finds a template by code
func (r *ImagingTemplateRepository) FindByCode(code string) (*domain.ImagingTemplate, error) {
	var template domain.ImagingTemplate
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// Search searches templates by name or code
func (r *ImagingTemplateRepository) Search(query string, limit int) ([]*domain.ImagingTemplate, error) {
	var templates []*domain.ImagingTemplate

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

// FindByModality finds templates by modality
func (r *ImagingTemplateRepository) FindByModality(modality string) ([]*domain.ImagingTemplate, error) {
	var templates []*domain.ImagingTemplate
	err := r.db.Where("modality = ? AND is_active = ?", modality, true).
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}
