package repository

import (
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// BedAllocationRepository handles bed allocation data operations
type BedAllocationRepository struct {
	db *gorm.DB
}

// NewBedAllocationRepository creates a new bed allocation repository
func NewBedAllocationRepository(db *gorm.DB) *BedAllocationRepository {
	return &BedAllocationRepository{db: db}
}

// Create creates a bed allocation
func (r *BedAllocationRepository) Create(allocation *domain.BedAllocation) error {
	return r.db.Create(allocation).Error
}

// FindByAdmissionID finds allocations for an admission
func (r *BedAllocationRepository) FindByAdmissionID(admissionID uint) ([]*domain.BedAllocation, error) {
	var allocations []*domain.BedAllocation
	err := r.db.Preload("Bed").
		Where("admission_id = ?", admissionID).
		Order("allocated_date DESC").
		Find(&allocations).Error
	return allocations, err
}

// FindCurrentAllocation finds current bed allocation for an admission
func (r *BedAllocationRepository) FindCurrentAllocation(admissionID uint) (*domain.BedAllocation, error) {
	var allocation domain.BedAllocation
	err := r.db.Preload("Bed").
		Where("admission_id = ? AND is_current = ?", admissionID, true).
		First(&allocation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &allocation, nil
}

// ReleaseAllocation releases a bed allocation
func (r *BedAllocationRepository) ReleaseAllocation(id uint) error {
	now := time.Now()
	return r.db.Model(&domain.BedAllocation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"released_date": now,
			"is_current":    false,
		}).Error
}

// Update updates an allocation
func (r *BedAllocationRepository) Update(allocation *domain.BedAllocation) error {
	return r.db.Save(allocation).Error
}
