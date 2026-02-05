package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// MedicalServiceRepository handles medical service data operations
type MedicalServiceRepository struct {
	db *gorm.DB
}

// NewMedicalServiceRepository creates a new medical service repository
func NewMedicalServiceRepository(db *gorm.DB) *MedicalServiceRepository {
	return &MedicalServiceRepository{db: db}
}

// Create creates a new medical service
func (r *MedicalServiceRepository) Create(service *domain.MedicalService) error {
	return r.db.Create(service).Error
}

// FindByID finds a medical service by ID
func (r *MedicalServiceRepository) FindByID(id uint) (*domain.MedicalService, error) {
	var service domain.MedicalService
	err := r.db.Preload("Department").First(&service, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &service, nil
}

// FindByCode finds a medical service by code
func (r *MedicalServiceRepository) FindByCode(code string) (*domain.MedicalService, error) {
	var service domain.MedicalService
	err := r.db.Preload("Department").Where("code = ?", code).First(&service).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &service, nil
}

// Update updates a medical service
func (r *MedicalServiceRepository) Update(service *domain.MedicalService) error {
	return r.db.Save(service).Error
}

// Delete soft deletes a medical service
func (r *MedicalServiceRepository) Delete(id uint) error {
	return r.db.Delete(&domain.MedicalService{}, id).Error
}

// List returns a paginated list of medical services
func (r *MedicalServiceRepository) List(page, pageSize int, departmentID *uint) ([]*domain.MedicalService, int64, error) {
	var services []*domain.MedicalService
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&domain.MedicalService{})

	if departmentID != nil {
		query = query.Where("department_id = ?", departmentID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Department").
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&services).Error

	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}
