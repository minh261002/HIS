package dto

import "time"

// CreateLabTestRequestRequest represents request to create a lab test
type CreateLabTestRequestRequest struct {
	VisitID       uint   `json:"visit_id" binding:"required"`
	TemplateID    uint   `json:"template_id" binding:"required"`
	Priority      string `json:"priority" binding:"required,oneof=ROUTINE URGENT STAT"`
	ClinicalNotes string `json:"clinical_notes" binding:"omitempty"`
}

// UpdateLabTestRequestRequest represents request to update a lab test
type UpdateLabTestRequestRequest struct {
	ClinicalNotes string `json:"clinical_notes" binding:"omitempty"`
}

// LabTestResultResponse represents a test result
type LabTestResultResponse struct {
	ID              uint   `json:"id"`
	ParameterName   string `json:"parameter_name"`
	Value           string `json:"value"`
	Unit            string `json:"unit"`
	NormalRangeText string `json:"normal_range_text"`
	IsAbnormal      bool   `json:"is_abnormal"`
	Remarks         string `json:"remarks"`
}

// LabTestRequestResponse represents lab test request details
type LabTestRequestResponse struct {
	ID                uint                     `json:"id"`
	RequestCode       string                   `json:"request_code"`
	VisitID           uint                     `json:"visit_id"`
	PatientID         uint                     `json:"patient_id"`
	PatientName       string                   `json:"patient_name"`
	DoctorID          uint                     `json:"doctor_id"`
	DoctorName        string                   `json:"doctor_name"`
	TemplateID        uint                     `json:"template_id"`
	TemplateName      string                   `json:"template_name"`
	TemplateCode      string                   `json:"template_code"`
	Status            string                   `json:"status"`
	Priority          string                   `json:"priority"`
	RequestedDate     time.Time                `json:"requested_date"`
	SampleCollectedAt *time.Time               `json:"sample_collected_at,omitempty"`
	CompletedAt       *time.Time               `json:"completed_at,omitempty"`
	ClinicalNotes     string                   `json:"clinical_notes"`
	Results           []*LabTestResultResponse `json:"results"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

// LabTestRequestListItem represents simplified request for list view
type LabTestRequestListItem struct {
	ID            uint      `json:"id"`
	RequestCode   string    `json:"request_code"`
	PatientName   string    `json:"patient_name"`
	TemplateName  string    `json:"template_name"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	RequestedDate time.Time `json:"requested_date"`
}
