package dto

import "time"

// CreateMedicalHistoryRequest represents request to add medical history
type CreateMedicalHistoryRequest struct {
	ConditionName string `json:"condition_name" binding:"required,min=2,max=200"`
	ConditionType string `json:"condition_type" binding:"required,oneof=CHRONIC ACUTE HEREDITARY SURGICAL OTHER"`
	DiagnosisDate string `json:"diagnosis_date" binding:"omitempty"` // Format: YYYY-MM-DD
	Status        string `json:"status" binding:"required,oneof=ACTIVE RESOLVED MANAGED MONITORING"`
	Treatment     string `json:"treatment" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
}

// UpdateMedicalHistoryRequest represents request to update medical history
type UpdateMedicalHistoryRequest struct {
	ConditionName string `json:"condition_name" binding:"omitempty,min=2,max=200"`
	ConditionType string `json:"condition_type" binding:"omitempty,oneof=CHRONIC ACUTE HEREDITARY SURGICAL OTHER"`
	DiagnosisDate string `json:"diagnosis_date" binding:"omitempty"` // Format: YYYY-MM-DD
	Status        string `json:"status" binding:"omitempty,oneof=ACTIVE RESOLVED MANAGED MONITORING"`
	Treatment     string `json:"treatment" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
	IsActive      *bool  `json:"is_active" binding:"omitempty"`
}

// MedicalHistoryResponse represents medical history details
type MedicalHistoryResponse struct {
	ID            uint       `json:"id"`
	PatientID     uint       `json:"patient_id"`
	ConditionName string     `json:"condition_name"`
	ConditionType string     `json:"condition_type"`
	DiagnosisDate *time.Time `json:"diagnosis_date"`
	Status        string     `json:"status"`
	Treatment     string     `json:"treatment"`
	Notes         string     `json:"notes"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// MedicalHistoryListItem represents simplified medical history for list view
type MedicalHistoryListItem struct {
	ID            uint   `json:"id"`
	ConditionName string `json:"condition_name"`
	ConditionType string `json:"condition_type"`
	Status        string `json:"status"`
	IsActive      bool   `json:"is_active"`
}
