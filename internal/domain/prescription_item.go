package domain

import (
	"time"

	"gorm.io/gorm"
)

// PrescriptionItem represents an item in a prescription
type PrescriptionItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	PrescriptionID uint          `gorm:"not null;index" json:"prescription_id"`
	Prescription   *Prescription `gorm:"foreignKey:PrescriptionID" json:"prescription,omitempty"`

	MedicationID uint        `gorm:"not null;index" json:"medication_id"`
	Medication   *Medication `gorm:"foreignKey:MedicationID" json:"medication,omitempty"`

	// Dosage Information
	Quantity     int    `gorm:"not null" json:"quantity"`           // Number of units to dispense
	Dosage       string `gorm:"size:100;not null" json:"dosage"`    // e.g., "1 viên", "5ml"
	Frequency    string `gorm:"size:100;not null" json:"frequency"` // e.g., "2 lần/ngày"
	DurationDays int    `gorm:"not null" json:"duration_days"`      // Number of days
	Instructions string `gorm:"type:text" json:"instructions"`      // How to take
}

// TableName specifies the table name for PrescriptionItem model
func (PrescriptionItem) TableName() string {
	return "prescription_items"
}
