package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// BedRepository handles bed data operations
type BedRepository struct {
	db *gorm.DB
}

// NewBedRepository creates a new bed repository
func NewBedRepository(db *gorm.DB) *BedRepository {
	return &BedRepository{db: db}
}

// FindByID finds a bed by ID
func (r *BedRepository) FindByID(id uint) (*domain.Bed, error) {
	var bed domain.Bed
	err := r.db.Where("is_active = ?", true).First(&bed, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bed, nil
}

// FindByBedNumber finds a bed by bed number
func (r *BedRepository) FindByBedNumber(bedNumber string) (*domain.Bed, error) {
	var bed domain.Bed
	err := r.db.Where("bed_number = ? AND is_active = ?", bedNumber, true).First(&bed).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bed, nil
}

// FindAvailableBeds finds available beds by department and bed type
func (r *BedRepository) FindAvailableBeds(department, bedType string) ([]*domain.Bed, error) {
	var beds []*domain.Bed
	query := r.db.Where("status = ? AND is_active = ?", domain.BedStatusAvailable, true)

	if department != "" {
		query = query.Where("department = ?", department)
	}
	if bedType != "" {
		query = query.Where("bed_type = ?", bedType)
	}

	err := query.Order("bed_number ASC").Find(&beds).Error
	return beds, err
}

// UpdateStatus updates bed status
func (r *BedRepository) UpdateStatus(id uint, status domain.BedStatus) error {
	return r.db.Model(&domain.Bed{}).Where("id = ?", id).Update("status", status).Error
}
