package dto

import "time"

// CreatePatientRequest represents patient registration request
type CreatePatientRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=50"`
	LastName    string `json:"last_name" binding:"required,min=2,max=50"`
	DateOfBirth string `json:"date_of_birth" binding:"required"` // Format: YYYY-MM-DD
	Gender      string `json:"gender" binding:"required,oneof=MALE FEMALE OTHER"`
	BloodType   string `json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`

	// Contact
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=20"`
	Email       string `json:"email" binding:"omitempty,email"`
	Address     string `json:"address" binding:"omitempty,max=255"`
	City        string `json:"city" binding:"omitempty,max=100"`
	State       string `json:"state" binding:"omitempty,max=100"`
	PostalCode  string `json:"postal_code" binding:"omitempty,max=20"`
	Country     string `json:"country" binding:"omitempty,max=100"`

	// Identification
	NationalID string `json:"national_id" binding:"omitempty,max=20"`

	// Insurance
	InsuranceNumber   string `json:"insurance_number" binding:"omitempty,max=50"`
	InsuranceProvider string `json:"insurance_provider" binding:"omitempty,max=100"`

	// Emergency Contact
	EmergencyContactName         string `json:"emergency_contact_name" binding:"omitempty,max=100"`
	EmergencyContactPhone        string `json:"emergency_contact_phone" binding:"omitempty,max=20"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship" binding:"omitempty,max=50"`

	// Medical
	Allergies         string `json:"allergies" binding:"omitempty"`
	ChronicConditions string `json:"chronic_conditions" binding:"omitempty"`
	Notes             string `json:"notes" binding:"omitempty"`
}

// UpdatePatientRequest represents patient update request
type UpdatePatientRequest struct {
	FirstName   string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName    string `json:"last_name" binding:"omitempty,min=2,max=50"`
	DateOfBirth string `json:"date_of_birth" binding:"omitempty"` // Format: YYYY-MM-DD
	Gender      string `json:"gender" binding:"omitempty,oneof=MALE FEMALE OTHER"`
	BloodType   string `json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`

	PhoneNumber string `json:"phone_number" binding:"omitempty,min=10,max=20"`
	Email       string `json:"email" binding:"omitempty,email"`
	Address     string `json:"address" binding:"omitempty,max=255"`
	City        string `json:"city" binding:"omitempty,max=100"`
	State       string `json:"state" binding:"omitempty,max=100"`
	PostalCode  string `json:"postal_code" binding:"omitempty,max=20"`
	Country     string `json:"country" binding:"omitempty,max=100"`

	NationalID string `json:"national_id" binding:"omitempty,max=20"`

	InsuranceNumber   string `json:"insurance_number" binding:"omitempty,max=50"`
	InsuranceProvider string `json:"insurance_provider" binding:"omitempty,max=100"`

	EmergencyContactName         string `json:"emergency_contact_name" binding:"omitempty,max=100"`
	EmergencyContactPhone        string `json:"emergency_contact_phone" binding:"omitempty,max=20"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship" binding:"omitempty,max=50"`

	Allergies         string `json:"allergies" binding:"omitempty"`
	ChronicConditions string `json:"chronic_conditions" binding:"omitempty"`
	Notes             string `json:"notes" binding:"omitempty"`
	IsActive          *bool  `json:"is_active" binding:"omitempty"`
}

// PatientResponse represents patient details response
type PatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	BloodType   string `json:"blood_type"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	NationalID  string `json:"national_id"`

	InsuranceNumber   string `json:"insurance_number"`
	InsuranceProvider string `json:"insurance_provider"`

	EmergencyContactName         string `json:"emergency_contact_name"`
	EmergencyContactPhone        string `json:"emergency_contact_phone"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship"`

	Allergies         string `json:"allergies"`
	ChronicConditions string `json:"chronic_conditions"`
	Notes             string `json:"notes"`

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PatientListItem represents patient in list view
type PatientListItem struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

// PatientSearchRequest represents search filters
type PatientSearchRequest struct {
	Query     string `form:"q"`
	Gender    string `form:"gender"`
	BloodType string `form:"blood_type"`
	City      string `form:"city"`
	IsActive  *bool  `form:"is_active"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}
