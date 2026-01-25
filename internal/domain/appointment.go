package domain

import (
	"time"

	"gorm.io/gorm"
)

// AppointmentType represents the type of appointment
type AppointmentType string

const (
	AppointmentTypeConsultation AppointmentType = "CONSULTATION"
	AppointmentTypeFollowUp     AppointmentType = "FOLLOW_UP"
	AppointmentTypeEmergency    AppointmentType = "EMERGENCY"
	AppointmentTypeCheckup      AppointmentType = "CHECKUP"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusScheduled  AppointmentStatus = "SCHEDULED"
	AppointmentStatusConfirmed  AppointmentStatus = "CONFIRMED"
	AppointmentStatusInProgress AppointmentStatus = "IN_PROGRESS"
	AppointmentStatusCompleted  AppointmentStatus = "COMPLETED"
	AppointmentStatusCancelled  AppointmentStatus = "CANCELLED"
	AppointmentStatusNoShow     AppointmentStatus = "NO_SHOW"
)

// Appointment represents a patient appointment
type Appointment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Appointment Code (auto-generated: APT-YYYYMMDD-XXXX)
	AppointmentCode string `gorm:"uniqueIndex;size:20;not null" json:"appointment_code"`

	// Foreign Keys
	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"`
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	// Appointment Details
	AppointmentDate time.Time         `gorm:"not null;index" json:"appointment_date"`
	AppointmentTime time.Time         `gorm:"not null" json:"appointment_time"`
	DurationMinutes int               `gorm:"not null;default:30" json:"duration_minutes"`
	AppointmentType AppointmentType   `gorm:"size:20;not null" json:"appointment_type"`
	Status          AppointmentStatus `gorm:"size:20;not null;index;default:'SCHEDULED'" json:"status"`

	// Appointment Information
	Reason string `gorm:"type:text" json:"reason"` // Chief complaint
	Notes  string `gorm:"type:text" json:"notes"`

	// Cancellation Details
	CancelledReason string     `gorm:"type:text" json:"cancelled_reason,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy     *uint      `json:"cancelled_by,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Appointment model
func (Appointment) TableName() string {
	return "appointments"
}
