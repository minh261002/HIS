package domain

import (
	"time"

	"gorm.io/gorm"
)

// LabTestRequestStatus represents the status of a lab test request
type LabTestRequestStatus string

const (
	LabTestRequestStatusPending         LabTestRequestStatus = "PENDING"
	LabTestRequestStatusSampleCollected LabTestRequestStatus = "SAMPLE_COLLECTED"
	LabTestRequestStatusInProgress      LabTestRequestStatus = "IN_PROGRESS"
	LabTestRequestStatusCompleted       LabTestRequestStatus = "COMPLETED"
	LabTestRequestStatusCancelled       LabTestRequestStatus = "CANCELLED"
)

// LabTestPriority represents the priority of a lab test
type LabTestPriority string

const (
	LabTestPriorityRoutine LabTestPriority = "ROUTINE"
	LabTestPriorityUrgent  LabTestPriority = "URGENT"
	LabTestPriorityStat    LabTestPriority = "STAT"
)

// LabTestRequest represents a lab test request
type LabTestRequest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Request Code (auto-generated: LTR-YYYYMMDD-XXXX)
	RequestCode string `gorm:"uniqueIndex;size:20;not null" json:"request_code"`

	// Foreign Keys
	VisitID uint   `gorm:"not null;index" json:"visit_id"`
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"` // Denormalized
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	TemplateID uint             `gorm:"not null;index" json:"template_id"`
	Template   *LabTestTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`

	// Request Details
	Status            LabTestRequestStatus `gorm:"size:20;not null;index;default:'PENDING'" json:"status"`
	Priority          LabTestPriority      `gorm:"size:20;not null;default:'ROUTINE'" json:"priority"`
	RequestedDate     time.Time            `gorm:"not null;index" json:"requested_date"`
	SampleCollectedAt *time.Time           `json:"sample_collected_at,omitempty"`
	CompletedAt       *time.Time           `json:"completed_at,omitempty"`
	ClinicalNotes     string               `gorm:"type:text" json:"clinical_notes"`

	// Results (one-to-many)
	Results []*LabTestResult `gorm:"foreignKey:RequestID" json:"results,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for LabTestRequest model
func (LabTestRequest) TableName() string {
	return "lab_test_requests"
}
