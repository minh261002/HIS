package domain

import (
	"time"

	"gorm.io/gorm"
)

// AllergenType represents the type of allergen
type AllergenType string

const (
	AllergenTypeMedication    AllergenType = "MEDICATION"
	AllergenTypeFood          AllergenType = "FOOD"
	AllergenTypeEnvironmental AllergenType = "ENVIRONMENTAL"
	AllergenTypeOther         AllergenType = "OTHER"
)

// AllergySeverity represents the severity of an allergic reaction
type AllergySeverity string

const (
	AllergySeverityMild            AllergySeverity = "MILD"
	AllergySeverityModerate        AllergySeverity = "MODERATE"
	AllergySeveritySevere          AllergySeverity = "SEVERE"
	AllergySeverityLifeThreatening AllergySeverity = "LIFE_THREATENING"
)

// PatientAllergy represents a patient's allergy record
type PatientAllergy struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	// Allergy Information
	Allergen      string          `gorm:"size:100;not null" json:"allergen"`
	AllergenType  AllergenType    `gorm:"size:20;not null;index" json:"allergen_type"`
	Reaction      string          `gorm:"type:text" json:"reaction"`
	Severity      AllergySeverity `gorm:"size:20;not null;index" json:"severity"`
	DiagnosedDate *time.Time      `json:"diagnosed_date"`
	Notes         string          `gorm:"type:text" json:"notes"`

	// Status
	IsActive bool `gorm:"default:true;index" json:"is_active"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for PatientAllergy model
func (PatientAllergy) TableName() string {
	return "patient_allergies"
}
