package domain

import (
	"time"

	"gorm.io/gorm"
)

// AdmissionStatus represents the status of an admission
type AdmissionStatus string

const (
	AdmissionStatusAdmitted    AdmissionStatus = "ADMITTED"
	AdmissionStatusDischarged  AdmissionStatus = "DISCHARGED"
	AdmissionStatusTransferred AdmissionStatus = "TRANSFERRED"
)

// Admission represents a patient admission
type Admission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Admission Code (auto-generated: ADM-YYYYMMDD-XXXX)
	AdmissionCode string `gorm:"uniqueIndex;size:20;not null" json:"admission_code"`

	// Foreign Keys
	VisitID uint   `gorm:"not null;index" json:"visit_id"`
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"` // Denormalized
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	// Admission Details
	AdmissionDate      time.Time       `gorm:"not null;index" json:"admission_date"`
	DischargeDate      *time.Time      `json:"discharge_date,omitempty"`
	AdmissionDiagnosis string          `gorm:"type:text;not null" json:"admission_diagnosis"`
	DischargeDiagnosis string          `gorm:"type:text" json:"discharge_diagnosis"`
	DischargeSummary   string          `gorm:"type:text" json:"discharge_summary"`
	Status             AdmissionStatus `gorm:"size:20;not null;index;default:'ADMITTED'" json:"status"`

	// Relationships
	BedAllocations []*BedAllocation `gorm:"foreignKey:AdmissionID" json:"bed_allocations,omitempty"`
	NursingNotes   []*NursingNote   `gorm:"foreignKey:AdmissionID" json:"nursing_notes,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Admission model
func (Admission) TableName() string {
	return "admissions"
}
