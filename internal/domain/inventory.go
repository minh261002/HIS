package domain

import (
	"time"

	"gorm.io/gorm"
)

// InventoryUnit represents the unit of inventory
type InventoryUnit string

const (
	InventoryUnitTablet  InventoryUnit = "TABLET"
	InventoryUnitCapsule InventoryUnit = "CAPSULE"
	InventoryUnitBottle  InventoryUnit = "BOTTLE"
	InventoryUnitBox     InventoryUnit = "BOX"
	InventoryUnitVial    InventoryUnit = "VIAL"
	InventoryUnitTube    InventoryUnit = "TUBE"
	InventoryUnitSyringe InventoryUnit = "SYRINGE"
)

// Inventory represents medication inventory
type Inventory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	MedicationID uint        `gorm:"not null;index" json:"medication_id"`
	Medication   *Medication `gorm:"foreignKey:MedicationID" json:"medication,omitempty"`

	// Inventory Details
	BatchNumber  string        `gorm:"size:50;not null;index" json:"batch_number"`
	ExpiryDate   time.Time     `gorm:"not null;index" json:"expiry_date"`
	Quantity     int           `gorm:"not null" json:"quantity"`
	Unit         InventoryUnit `gorm:"size:20;not null" json:"unit"`
	CostPrice    float64       `gorm:"type:decimal(10,2)" json:"cost_price"`
	Supplier     string        `gorm:"size:200" json:"supplier"`
	ReceivedDate time.Time     `gorm:"not null" json:"received_date"`
}

// TableName specifies the table name for Inventory model
func (Inventory) TableName() string {
	return "inventory"
}
