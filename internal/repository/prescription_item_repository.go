package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// PrescriptionItemRepository handles prescription item data operations
type PrescriptionItemRepository struct {
	db *gorm.DB
}

// NewPrescriptionItemRepository creates a new prescription item repository
func NewPrescriptionItemRepository(db *gorm.DB) *PrescriptionItemRepository {
	return &PrescriptionItemRepository{db: db}
}

// Create creates a new prescription item
func (r *PrescriptionItemRepository) Create(item *domain.PrescriptionItem) error {
	return r.db.Create(item).Error
}

// FindByID finds a prescription item by ID
func (r *PrescriptionItemRepository) FindByID(id uint) (*domain.PrescriptionItem, error) {
	var item domain.PrescriptionItem
	err := r.db.Preload("Medication").First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// FindByPrescriptionID finds all items for a prescription
func (r *PrescriptionItemRepository) FindByPrescriptionID(prescriptionID uint) ([]*domain.PrescriptionItem, error) {
	var items []*domain.PrescriptionItem
	err := r.db.Preload("Medication").
		Where("prescription_id = ?", prescriptionID).
		Find(&items).Error
	return items, err
}

// Update updates a prescription item
func (r *PrescriptionItemRepository) Update(item *domain.PrescriptionItem) error {
	return r.db.Save(item).Error
}

// Delete deletes a prescription item
func (r *PrescriptionItemRepository) Delete(id uint) error {
	return r.db.Delete(&domain.PrescriptionItem{}, id).Error
}
