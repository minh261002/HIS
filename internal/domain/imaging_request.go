package domain

import (
	"time"

	"gorm.io/gorm"
)

// ImagingRequestStatus represents the status of an imaging request
type ImagingRequestStatus string

const (
	ImagingRequestStatusPending    ImagingRequestStatus = "PENDING"
	ImagingRequestStatusScheduled  ImagingRequestStatus = "SCHEDULED"
	ImagingRequestStatusInProgress ImagingRequestStatus = "IN_PROGRESS"
	ImagingRequestStatusCompleted  ImagingRequestStatus = "COMPLETED"
	ImagingRequestStatusCancelled  ImagingRequestStatus = "CANCELLED"
)

// ImagingPriority represents the priority of an imaging request
type ImagingPriority string

const (
	ImagingPriorityRoutine ImagingPriority = "ROUTINE"
	ImagingPriorityUrgent  ImagingPriority = "URGENT"
	ImagingPriorityStat    ImagingPriority = "STAT"
)

// ImagingRequest represents an imaging request
type ImagingRequest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Request Code (auto-generated: IMG-YYYYMMDD-XXXX)
	RequestCode string `gorm:"uniqueIndex;size:20;not null" json:"request_code"`

	// Foreign Keys
	VisitID uint   `gorm:"not null;index" json:"visit_id"`
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"` // Denormalized
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	TemplateID uint             `gorm:"not null;index" json:"template_id"`
	Template   *ImagingTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`

	// Request Details
	Status              ImagingRequestStatus `gorm:"size:20;not null;index;default:'PENDING'" json:"status"`
	Priority            ImagingPriority      `gorm:"size:20;not null;default:'ROUTINE'" json:"priority"`
	RequestedDate       time.Time            `gorm:"not null;index" json:"requested_date"`
	ScheduledDate       *time.Time           `json:"scheduled_date,omitempty"`
	CompletedAt         *time.Time           `json:"completed_at,omitempty"`
	ClinicalIndication  string               `gorm:"type:text" json:"clinical_indication"`
	SpecialInstructions string               `gorm:"type:text" json:"special_instructions"`

	// Result (one-to-one)
	Result *ImagingResult `gorm:"foreignKey:RequestID" json:"result,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for ImagingRequest model
func (ImagingRequest) TableName() string {
	return "imaging_requests"
}
