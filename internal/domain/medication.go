package domain

import (
	"time"

	"gorm.io/gorm"
)

// DosageForm represents the form of medication
type DosageForm string

const (
	DosageFormTablet    DosageForm = "TABLET"
	DosageFormCapsule   DosageForm = "CAPSULE"
	DosageFormSyrup     DosageForm = "SYRUP"
	DosageFormInjection DosageForm = "INJECTION"
	DosageFormCream     DosageForm = "CREAM"
	DosageFormOintment  DosageForm = "OINTMENT"
	DosageFormDrops     DosageForm = "DROPS"
	DosageFormInhaler   DosageForm = "INHALER"
)

// Medication represents a medication/drug
type Medication struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Medication Information
	Name         string     `gorm:"size:200;not null;index" json:"name"`
	GenericName  string     `gorm:"size:200;index" json:"generic_name"`
	DosageForm   DosageForm `gorm:"size:20;not null;index" json:"dosage_form"`
	Strength     string     `gorm:"size:50" json:"strength"` // e.g., "500mg", "10ml"
	Unit         string     `gorm:"size:20" json:"unit"`     // mg, ml, g, etc.
	Manufacturer string     `gorm:"size:200" json:"manufacturer"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`
}

// TableName specifies the table name for Medication model
func (Medication) TableName() string {
	return "medications"
}
