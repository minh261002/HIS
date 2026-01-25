package dto

import "time"

// CreateDiagnosisRequest represents request to add a diagnosis
type CreateDiagnosisRequest struct {
	VisitID         uint   `json:"visit_id" binding:"required"`
	ICD10CodeID     uint   `json:"icd10_code_id" binding:"required"`
	DiagnosisType   string `json:"diagnosis_type" binding:"required,oneof=PRIMARY SECONDARY DIFFERENTIAL COMPLICATION"`
	DiagnosisStatus string `json:"diagnosis_status" binding:"required,oneof=PROVISIONAL CONFIRMED RULED_OUT"`
	ClinicalNotes   string `json:"clinical_notes" binding:"omitempty"`
}

// UpdateDiagnosisRequest represents request to update a diagnosis
type UpdateDiagnosisRequest struct {
	DiagnosisStatus string `json:"diagnosis_status" binding:"omitempty,oneof=PROVISIONAL CONFIRMED RULED_OUT"`
	ClinicalNotes   string `json:"clinical_notes" binding:"omitempty"`
}

// DiagnosisResponse represents diagnosis details
type DiagnosisResponse struct {
	ID               uint      `json:"id"`
	VisitID          uint      `json:"visit_id"`
	PatientID        uint      `json:"patient_id"`
	PatientName      string    `json:"patient_name"`
	ICD10CodeID      uint      `json:"icd10_code_id"`
	ICD10Code        string    `json:"icd10_code"`
	ICD10Description string    `json:"icd10_description"`
	DiagnosedBy      uint      `json:"diagnosed_by"`
	DoctorName       string    `json:"doctor_name"`
	DiagnosisType    string    `json:"diagnosis_type"`
	DiagnosisStatus  string    `json:"diagnosis_status"`
	ClinicalNotes    string    `json:"clinical_notes"`
	DiagnosedAt      time.Time `json:"diagnosed_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DiagnosisListItem represents simplified diagnosis for list view
type DiagnosisListItem struct {
	ID               uint      `json:"id"`
	ICD10Code        string    `json:"icd10_code"`
	ICD10Description string    `json:"icd10_description"`
	DiagnosisType    string    `json:"diagnosis_type"`
	DiagnosisStatus  string    `json:"diagnosis_status"`
	DiagnosedAt      time.Time `json:"diagnosed_at"`
}
