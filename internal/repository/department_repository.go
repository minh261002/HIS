package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// DepartmentRepository handles department data operations
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// Create creates a new department
func (r *DepartmentRepository) Create(dept *domain.Department) error {
	return r.db.Create(dept).Error
}

// FindByID finds a department by ID
func (r *DepartmentRepository) FindByID(id uint) (*domain.Department, error) {
	var dept domain.Department
	err := r.db.Preload("HeadDoctor").First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

// FindByCode finds a department by code
func (r *DepartmentRepository) FindByCode(code string) (*domain.Department, error) {
	var dept domain.Department
	err := r.db.Preload("HeadDoctor").Where("code = ?", code).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

// Update updates a department
func (r *DepartmentRepository) Update(dept *domain.Department) error {
	return r.db.Save(dept).Error
}

// Delete soft deletes a department
func (r *DepartmentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Department{}, id).Error
}

// List returns a paginated list of departments
func (r *DepartmentRepository) List(page, pageSize int) ([]*domain.Department, int64, error) {
	var depts []*domain.Department
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&domain.Department{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("HeadDoctor").
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&depts).Error

	if err != nil {
		return nil, 0, err
	}

	return depts, total, nil
}

// ListActive returns all active departments
func (r *DepartmentRepository) ListActive() ([]*domain.Department, error) {
	var depts []*domain.Department
	err := r.db.Where("is_active = ?", true).Order("name ASC").Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}
