package domain

import (
	"time"

	"gorm.io/gorm"
)

// DiagnosisType represents the type of diagnosis
type DiagnosisType string

const (
	DiagnosisTypePrimary      DiagnosisType = "PRIMARY"
	DiagnosisTypeSecondary    DiagnosisType = "SECONDARY"
	DiagnosisTypeDifferential DiagnosisType = "DIFFERENTIAL"
	DiagnosisTypeComplication DiagnosisType = "COMPLICATION"
)

// DiagnosisStatus represents the status of a diagnosis
type DiagnosisStatus string

const (
	DiagnosisStatusProvisional DiagnosisStatus = "PROVISIONAL"
	DiagnosisStatusConfirmed   DiagnosisStatus = "CONFIRMED"
	DiagnosisStatusRuledOut    DiagnosisStatus = "RULED_OUT"
)

// Diagnosis represents a patient diagnosis
type Diagnosis struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	VisitID uint   `gorm:"not null;index" json:"visit_id"`
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized for quick access
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	ICD10CodeID uint       `gorm:"not null;index" json:"icd10_code_id"`
	ICD10Code   *ICD10Code `gorm:"foreignKey:ICD10CodeID" json:"icd10_code,omitempty"`

	DiagnosedBy uint  `gorm:"not null;index" json:"diagnosed_by"` // Doctor
	Doctor      *User `gorm:"foreignKey:DiagnosedBy" json:"doctor,omitempty"`

	// Diagnosis Details
	DiagnosisType   DiagnosisType   `gorm:"size:20;not null;index" json:"diagnosis_type"`
	DiagnosisStatus DiagnosisStatus `gorm:"size:20;not null;index" json:"diagnosis_status"`
	ClinicalNotes   string          `gorm:"type:text" json:"clinical_notes"`
	DiagnosedAt     time.Time       `gorm:"not null;index" json:"diagnosed_at"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Diagnosis model
func (Diagnosis) TableName() string {
	return "diagnoses"
}
