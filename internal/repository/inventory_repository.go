package repository

import (
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// InventoryRepository handles inventory data operations
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Create creates inventory
func (r *InventoryRepository) Create(inventory *domain.Inventory) error {
	return r.db.Create(inventory).Error
}

// FindByMedicationID finds inventory for a medication
func (r *InventoryRepository) FindByMedicationID(medicationID uint) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	err := r.db.Preload("Medication").
		Where("medication_id = ? AND quantity > 0", medicationID).
		Order("expiry_date ASC, received_date ASC"). // FIFO
		Find(&inventories).Error
	return inventories, err
}

// FindLowStock finds low stock items
func (r *InventoryRepository) FindLowStock(threshold int) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	err := r.db.Preload("Medication").
		Where("quantity <= ? AND quantity > 0", threshold).
		Order("quantity ASC").
		Find(&inventories).Error
	return inventories, err
}

// FindExpiringSoon finds items expiring soon
func (r *InventoryRepository) FindExpiringSoon(days int) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	expiryDate := time.Now().AddDate(0, 0, days)
	err := r.db.Preload("Medication").
		Where("expiry_date <= ? AND expiry_date > ? AND quantity > 0", expiryDate, time.Now()).
		Order("expiry_date ASC").
		Find(&inventories).Error
	return inventories, err
}

// UpdateQuantity updates inventory quantity
func (r *InventoryRepository) UpdateQuantity(id uint, quantity int) error {
	return r.db.Model(&domain.Inventory{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// FindAvailableStock finds available inventory for dispensing (FIFO)
func (r *InventoryRepository) FindAvailableStock(medicationID uint, quantityNeeded int) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.Where("medication_id = ? AND quantity >= ? AND expiry_date > ?",
		medicationID, quantityNeeded, time.Now()).
		Order("expiry_date ASC, received_date ASC"). // FIFO
		First(&inventory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &inventory, nil
}

// FindByID finds inventory by ID
func (r *InventoryRepository) FindByID(id uint) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.Preload("Medication").First(&inventory, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &inventory, nil
}
