package dto

import "time"

// CreateImagingRequestRequest represents request to create an imaging request
type CreateImagingRequestRequest struct {
	VisitID             uint   `json:"visit_id" binding:"required"`
	TemplateID          uint   `json:"template_id" binding:"required"`
	Priority            string `json:"priority" binding:"required,oneof=ROUTINE URGENT STAT"`
	ClinicalIndication  string `json:"clinical_indication" binding:"required"`
	SpecialInstructions string `json:"special_instructions" binding:"omitempty"`
}

// UpdateImagingRequestRequest represents request to update an imaging request
type UpdateImagingRequestRequest struct {
	ClinicalIndication  string `json:"clinical_indication" binding:"omitempty"`
	SpecialInstructions string `json:"special_instructions" binding:"omitempty"`
}

// ScheduleImagingRequest represents request to schedule imaging
type ScheduleImagingRequest struct {
	ScheduledDate string `json:"scheduled_date" binding:"required"` // ISO 8601 format
}

// ImagingResultResponse represents imaging result details
type ImagingResultResponse struct {
	ID              uint      `json:"id"`
	RadiologistID   uint      `json:"radiologist_id"`
	RadiologistName string    `json:"radiologist_name"`
	Findings        string    `json:"findings"`
	Impression      string    `json:"impression"`
	DICOMFiles      []string  `json:"dicom_files"`
	ReportDate      string    `json:"report_date"`
	IsCritical      bool      `json:"is_critical"`
	CreatedAt       time.Time `json:"created_at"`
}

// ImagingRequestResponse represents imaging request details
type ImagingRequestResponse struct {
	ID                  uint                   `json:"id"`
	RequestCode         string                 `json:"request_code"`
	VisitID             uint                   `json:"visit_id"`
	PatientID           uint                   `json:"patient_id"`
	PatientName         string                 `json:"patient_name"`
	DoctorID            uint                   `json:"doctor_id"`
	DoctorName          string                 `json:"doctor_name"`
	TemplateID          uint                   `json:"template_id"`
	TemplateName        string                 `json:"template_name"`
	TemplateCode        string                 `json:"template_code"`
	Modality            string                 `json:"modality"`
	Status              string                 `json:"status"`
	Priority            string                 `json:"priority"`
	RequestedDate       time.Time              `json:"requested_date"`
	ScheduledDate       *time.Time             `json:"scheduled_date,omitempty"`
	CompletedAt         *time.Time             `json:"completed_at,omitempty"`
	ClinicalIndication  string                 `json:"clinical_indication"`
	SpecialInstructions string                 `json:"special_instructions"`
	Result              *ImagingResultResponse `json:"result,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// ImagingRequestListItem represents simplified request for list view
type ImagingRequestListItem struct {
	ID            uint      `json:"id"`
	RequestCode   string    `json:"request_code"`
	PatientName   string    `json:"patient_name"`
	TemplateName  string    `json:"template_name"`
	Modality      string    `json:"modality"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	RequestedDate time.Time `json:"requested_date"`
	HasResult     bool      `json:"has_result"`
}
