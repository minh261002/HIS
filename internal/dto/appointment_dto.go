package dto

import "time"

// CreateAppointmentRequest represents request to schedule an appointment
type CreateAppointmentRequest struct {
	PatientID       uint   `json:"patient_id" binding:"required"`
	DoctorID        uint   `json:"doctor_id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"` // YYYY-MM-DD
	AppointmentTime string `json:"appointment_time" binding:"required"` // HH:MM
	DurationMinutes int    `json:"duration_minutes" binding:"omitempty,oneof=15 30 45 60"`
	AppointmentType string `json:"appointment_type" binding:"required,oneof=CONSULTATION FOLLOW_UP EMERGENCY CHECKUP"`
	Reason          string `json:"reason" binding:"required,min=5"`
	Notes           string `json:"notes" binding:"omitempty"`
}

// UpdateAppointmentRequest represents request to update/reschedule appointment
type UpdateAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date" binding:"omitempty"`
	AppointmentTime string `json:"appointment_time" binding:"omitempty"`
	DurationMinutes int    `json:"duration_minutes" binding:"omitempty,oneof=15 30 45 60"`
	Reason          string `json:"reason" binding:"omitempty,min=5"`
	Notes           string `json:"notes" binding:"omitempty"`
}

// CancelAppointmentRequest represents request to cancel appointment
type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required,min=5"`
}

// AppointmentResponse represents appointment details
type AppointmentResponse struct {
	ID              uint       `json:"id"`
	AppointmentCode string     `json:"appointment_code"`
	PatientID       uint       `json:"patient_id"`
	PatientName     string     `json:"patient_name"`
	DoctorID        uint       `json:"doctor_id"`
	DoctorName      string     `json:"doctor_name"`
	AppointmentDate string     `json:"appointment_date"`
	AppointmentTime string     `json:"appointment_time"`
	DurationMinutes int        `json:"duration_minutes"`
	AppointmentType string     `json:"appointment_type"`
	Status          string     `json:"status"`
	Reason          string     `json:"reason"`
	Notes           string     `json:"notes"`
	CancelledReason string     `json:"cancelled_reason,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AppointmentListItem represents simplified appointment for list view
type AppointmentListItem struct {
	ID              uint   `json:"id"`
	AppointmentCode string `json:"appointment_code"`
	PatientName     string `json:"patient_name"`
	DoctorName      string `json:"doctor_name"`
	AppointmentDate string `json:"appointment_date"`
	AppointmentTime string `json:"appointment_time"`
	AppointmentType string `json:"appointment_type"`
	Status          string `json:"status"`
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	Time      string `json:"time"`
	Available bool   `json:"available"`
}
