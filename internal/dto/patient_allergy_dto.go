package dto

import "time"

// CreateAllergyRequest represents request to add allergy
type CreateAllergyRequest struct {
	Allergen      string `json:"allergen" binding:"required,min=2,max=100"`
	AllergenType  string `json:"allergen_type" binding:"required,oneof=MEDICATION FOOD ENVIRONMENTAL OTHER"`
	Reaction      string `json:"reaction" binding:"omitempty"`
	Severity      string `json:"severity" binding:"required,oneof=MILD MODERATE SEVERE LIFE_THREATENING"`
	DiagnosedDate string `json:"diagnosed_date" binding:"omitempty"` // Format: YYYY-MM-DD
	Notes         string `json:"notes" binding:"omitempty"`
}

// UpdateAllergyRequest represents request to update allergy
type UpdateAllergyRequest struct {
	Allergen      string `json:"allergen" binding:"omitempty,min=2,max=100"`
	AllergenType  string `json:"allergen_type" binding:"omitempty,oneof=MEDICATION FOOD ENVIRONMENTAL OTHER"`
	Reaction      string `json:"reaction" binding:"omitempty"`
	Severity      string `json:"severity" binding:"omitempty,oneof=MILD MODERATE SEVERE LIFE_THREATENING"`
	DiagnosedDate string `json:"diagnosed_date" binding:"omitempty"` // Format: YYYY-MM-DD
	Notes         string `json:"notes" binding:"omitempty"`
	IsActive      *bool  `json:"is_active" binding:"omitempty"`
}

// AllergyResponse represents allergy details
type AllergyResponse struct {
	ID            uint       `json:"id"`
	PatientID     uint       `json:"patient_id"`
	Allergen      string     `json:"allergen"`
	AllergenType  string     `json:"allergen_type"`
	Reaction      string     `json:"reaction"`
	Severity      string     `json:"severity"`
	DiagnosedDate *time.Time `json:"diagnosed_date"`
	Notes         string     `json:"notes"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AllergyListItem represents simplified allergy for list view
type AllergyListItem struct {
	ID           uint   `json:"id"`
	Allergen     string `json:"allergen"`
	AllergenType string `json:"allergen_type"`
	Severity     string `json:"severity"`
	IsActive     bool   `json:"is_active"`
}
