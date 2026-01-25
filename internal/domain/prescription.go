package domain

import (
	"time"

	"gorm.io/gorm"
)

// PrescriptionStatus represents the status of a prescription
type PrescriptionStatus string

const (
	PrescriptionStatusPending   PrescriptionStatus = "PENDING"
	PrescriptionStatusDispensed PrescriptionStatus = "DISPENSED"
	PrescriptionStatusCompleted PrescriptionStatus = "COMPLETED"
	PrescriptionStatusCancelled PrescriptionStatus = "CANCELLED"
)

// Prescription represents a medication prescription
type Prescription struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Prescription Code (auto-generated: PRX-YYYYMMDD-XXXX)
	PrescriptionCode string `gorm:"uniqueIndex;size:20;not null" json:"prescription_code"`

	// Foreign Keys
	VisitID uint   `gorm:"not null;index" json:"visit_id"`
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"` // Denormalized
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	DiagnosisID *uint      `gorm:"index" json:"diagnosis_id,omitempty"` // Optional
	Diagnosis   *Diagnosis `gorm:"foreignKey:DiagnosisID" json:"diagnosis,omitempty"`

	// Prescription Details
	Status         PrescriptionStatus `gorm:"size:20;not null;index;default:'PENDING'" json:"status"`
	PrescribedDate time.Time          `gorm:"not null;index" json:"prescribed_date"`
	Notes          string             `gorm:"type:text" json:"notes"`

	// Prescription Items (one-to-many)
	Items []*PrescriptionItem `gorm:"foreignKey:PrescriptionID" json:"items,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Prescription model
func (Prescription) TableName() string {
	return "prescriptions"
}
