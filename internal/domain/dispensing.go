package domain

import (
	"time"

	"gorm.io/gorm"
)

// Dispensing represents a medication dispensing record
type Dispensing struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Dispensing Code (auto-generated: DIS-YYYYMMDD-XXXX)
	DispensingCode string `gorm:"uniqueIndex;size:20;not null" json:"dispensing_code"`

	// Foreign Keys
	PrescriptionID uint          `gorm:"not null;index" json:"prescription_id"`
	Prescription   *Prescription `gorm:"foreignKey:PrescriptionID" json:"prescription,omitempty"`

	PrescriptionItemID uint              `gorm:"not null;index" json:"prescription_item_id"`
	PrescriptionItem   *PrescriptionItem `gorm:"foreignKey:PrescriptionItemID" json:"prescription_item,omitempty"`

	MedicationID uint        `gorm:"not null;index" json:"medication_id"` // Denormalized
	Medication   *Medication `gorm:"foreignKey:MedicationID" json:"medication,omitempty"`

	InventoryID uint       `gorm:"not null;index" json:"inventory_id"`
	Inventory   *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	PharmacistID uint  `gorm:"not null;index" json:"pharmacist_id"`
	Pharmacist   *User `gorm:"foreignKey:PharmacistID" json:"pharmacist,omitempty"`

	// Dispensing Details
	QuantityDispensed int       `gorm:"not null" json:"quantity_dispensed"`
	BatchNumber       string    `gorm:"size:50;not null" json:"batch_number"` // Denormalized from inventory
	DispensedDate     time.Time `gorm:"not null;index" json:"dispensed_date"`
	Notes             string    `gorm:"type:text" json:"notes"`
}

// TableName specifies the table name for Dispensing model
func (Dispensing) TableName() string {
	return "dispensing"
}
