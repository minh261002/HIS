package dto

import "time"

// VitalSignsRequest represents vital signs data
type VitalSignsRequest struct {
	Temperature            float64 `json:"temperature" binding:"omitempty,min=35,max=42"`
	BloodPressureSystolic  int     `json:"blood_pressure_systolic" binding:"omitempty,min=70,max=200"`
	BloodPressureDiastolic int     `json:"blood_pressure_diastolic" binding:"omitempty,min=40,max=130"`
	HeartRate              int     `json:"heart_rate" binding:"omitempty,min=40,max=200"`
	RespiratoryRate        int     `json:"respiratory_rate" binding:"omitempty,min=8,max=40"`
	OxygenSaturation       int     `json:"oxygen_saturation" binding:"omitempty,min=70,max=100"`
	Weight                 float64 `json:"weight" binding:"omitempty,min=1,max=300"`
	Height                 float64 `json:"height" binding:"omitempty,min=30,max=250"`
}

// CreateVisitRequest represents request to create a new visit
type CreateVisitRequest struct {
	AppointmentID  *uint              `json:"appointment_id" binding:"omitempty"`
	PatientID      uint               `json:"patient_id" binding:"required"`
	DoctorID       uint               `json:"doctor_id" binding:"required"`
	VisitType      string             `json:"visit_type" binding:"required,oneof=SCHEDULED WALK_IN EMERGENCY FOLLOW_UP"`
	ChiefComplaint string             `json:"chief_complaint" binding:"required,min=5"`
	VitalSigns     *VitalSignsRequest `json:"vital_signs" binding:"omitempty"`
}

// UpdateVisitRequest represents request to update visit details
type UpdateVisitRequest struct {
	Symptoms             string             `json:"symptoms" binding:"omitempty"`
	PhysicalExamination  string             `json:"physical_examination" binding:"omitempty"`
	ClinicalNotes        string             `json:"clinical_notes" binding:"omitempty"`
	TreatmentPlan        string             `json:"treatment_plan" binding:"omitempty"`
	FollowUpInstructions string             `json:"follow_up_instructions" binding:"omitempty"`
	NextVisitDate        string             `json:"next_visit_date" binding:"omitempty"` // YYYY-MM-DD
	VitalSigns           *VitalSignsRequest `json:"vital_signs" binding:"omitempty"`
}

// VisitResponse represents visit details
type VisitResponse struct {
	ID             uint   `json:"id"`
	VisitCode      string `json:"visit_code"`
	AppointmentID  *uint  `json:"appointment_id,omitempty"`
	PatientID      uint   `json:"patient_id"`
	PatientName    string `json:"patient_name"`
	DoctorID       uint   `json:"doctor_id"`
	DoctorName     string `json:"doctor_name"`
	VisitDate      string `json:"visit_date"`
	VisitTime      string `json:"visit_time"`
	VisitType      string `json:"visit_type"`
	Status         string `json:"status"`
	ChiefComplaint string `json:"chief_complaint"`
	Symptoms       string `json:"symptoms"`

	// Vital Signs
	Temperature            float64 `json:"temperature"`
	BloodPressureSystolic  int     `json:"blood_pressure_systolic"`
	BloodPressureDiastolic int     `json:"blood_pressure_diastolic"`
	HeartRate              int     `json:"heart_rate"`
	RespiratoryRate        int     `json:"respiratory_rate"`
	OxygenSaturation       int     `json:"oxygen_saturation"`
	Weight                 float64 `json:"weight"`
	Height                 float64 `json:"height"`
	BMI                    float64 `json:"bmi"`

	// Clinical Documentation
	PhysicalExamination  string     `json:"physical_examination"`
	ClinicalNotes        string     `json:"clinical_notes"`
	TreatmentPlan        string     `json:"treatment_plan"`
	FollowUpInstructions string     `json:"follow_up_instructions"`
	NextVisitDate        *time.Time `json:"next_visit_date,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VisitListItem represents simplified visit for list view
type VisitListItem struct {
	ID             uint   `json:"id"`
	VisitCode      string `json:"visit_code"`
	PatientName    string `json:"patient_name"`
	DoctorName     string `json:"doctor_name"`
	VisitDate      string `json:"visit_date"`
	VisitTime      string `json:"visit_time"`
	VisitType      string `json:"visit_type"`
	Status         string `json:"status"`
	ChiefComplaint string `json:"chief_complaint"`
}
