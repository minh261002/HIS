package domain

import (
	"time"

	"gorm.io/gorm"
)

// ConditionType represents the type of medical condition
type ConditionType string

const (
	ConditionTypeChronic    ConditionType = "CHRONIC"
	ConditionTypeAcute      ConditionType = "ACUTE"
	ConditionTypeHereditary ConditionType = "HEREDITARY"
	ConditionTypeSurgical   ConditionType = "SURGICAL"
	ConditionTypeOther      ConditionType = "OTHER"
)

// ConditionStatus represents the status of a medical condition
type ConditionStatus string

const (
	ConditionStatusActive     ConditionStatus = "ACTIVE"
	ConditionStatusResolved   ConditionStatus = "RESOLVED"
	ConditionStatusManaged    ConditionStatus = "MANAGED"
	ConditionStatusMonitoring ConditionStatus = "MONITORING"
)

// PatientMedicalHistory represents a patient's medical history record
type PatientMedicalHistory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	// Condition Information
	ConditionName string          `gorm:"size:200;not null" json:"condition_name"`
	ConditionType ConditionType   `gorm:"size:20;not null;index" json:"condition_type"`
	DiagnosisDate *time.Time      `gorm:"index" json:"diagnosis_date"`
	Status        ConditionStatus `gorm:"size:20;not null;index" json:"status"`
	Treatment     string          `gorm:"type:text" json:"treatment"`
	Notes         string          `gorm:"type:text" json:"notes"`

	// Status
	IsActive bool `gorm:"default:true;index" json:"is_active"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for PatientMedicalHistory model
func (PatientMedicalHistory) TableName() string {
	return "patient_medical_history"
}
