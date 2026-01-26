package dto

import "time"

// PrescriptionItemRequest represents a prescription item
type PrescriptionItemRequest struct {
	MedicationID uint   `json:"medication_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	DurationDays int    `json:"duration_days" binding:"required,min=1"`
	Instructions string `json:"instructions" binding:"omitempty"`
}

// CreatePrescriptionRequest represents request to create a prescription
type CreatePrescriptionRequest struct {
	VisitID     uint                       `json:"visit_id" binding:"required"`
	DiagnosisID *uint                      `json:"diagnosis_id" binding:"omitempty"`
	Notes       string                     `json:"notes" binding:"omitempty"`
	Items       []*PrescriptionItemRequest `json:"items" binding:"required,min=1,dive"`
}

// UpdatePrescriptionRequest represents request to update prescription
type UpdatePrescriptionRequest struct {
	Notes string `json:"notes" binding:"omitempty"`
}

// AddPrescriptionItemRequest represents request to add an item
type AddPrescriptionItemRequest struct {
	MedicationID uint   `json:"medication_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	DurationDays int    `json:"duration_days" binding:"required,min=1"`
	Instructions string `json:"instructions" binding:"omitempty"`
}

// UpdatePrescriptionItemRequest represents request to update an item
type UpdatePrescriptionItemRequest struct {
	Quantity     int    `json:"quantity" binding:"omitempty,min=1"`
	Dosage       string `json:"dosage" binding:"omitempty"`
	Frequency    string `json:"frequency" binding:"omitempty"`
	DurationDays int    `json:"duration_days" binding:"omitempty,min=1"`
	Instructions string `json:"instructions" binding:"omitempty"`
}

// PrescriptionItemResponse represents prescription item details
type PrescriptionItemResponse struct {
	ID             uint   `json:"id"`
	MedicationID   uint   `json:"medication_id"`
	MedicationName string `json:"medication_name"`
	DosageForm     string `json:"dosage_form"`
	Strength       string `json:"strength"`
	Quantity       int    `json:"quantity"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	DurationDays   int    `json:"duration_days"`
	Instructions   string `json:"instructions"`
}

// PrescriptionResponse represents prescription details
type PrescriptionResponse struct {
	ID               uint                        `json:"id"`
	PrescriptionCode string                      `json:"prescription_code"`
	VisitID          uint                        `json:"visit_id"`
	PatientID        uint                        `json:"patient_id"`
	PatientName      string                      `json:"patient_name"`
	DoctorID         uint                        `json:"doctor_id"`
	DoctorName       string                      `json:"doctor_name"`
	DiagnosisID      *uint                       `json:"diagnosis_id,omitempty"`
	DiagnosisCode    string                      `json:"diagnosis_code,omitempty"`
	DiagnosisDesc    string                      `json:"diagnosis_description,omitempty"`
	Status           string                      `json:"status"`
	PrescribedDate   string                      `json:"prescribed_date"`
	Notes            string                      `json:"notes"`
	Items            []*PrescriptionItemResponse `json:"items"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
}

// PrescriptionListItem represents simplified prescription for list view
type PrescriptionListItem struct {
	ID               uint      `json:"id"`
	PrescriptionCode string    `json:"prescription_code"`
	PatientName      string    `json:"patient_name"`
	DoctorName       string    `json:"doctor_name"`
	Status           string    `json:"status"`
	PrescribedDate   string    `json:"prescribed_date"`
	ItemCount        int       `json:"item_count"`
	CreatedAt        time.Time `json:"created_at"`
}
