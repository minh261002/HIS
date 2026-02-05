package dto

import "time"

// MedicalServiceResponse represents medical service data in response
type MedicalServiceResponse struct {
	ID          uint                `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	ServiceType string              `json:"service_type"`
	BasePrice   float64             `json:"base_price"`
	IsActive    bool                `json:"is_active"`
	Department  *DepartmentResponse `json:"department,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// CreateMedicalServiceRequest represents request for creating MedicalService
type CreateMedicalServiceRequest struct {
	Code         string  `json:"code" binding:"required,max=50"`
	Name         string  `json:"name" binding:"required,max=100"`
	Description  string  `json:"description"`
	ServiceType  string  `json:"service_type" binding:"required"`
	BasePrice    float64 `json:"base_price" binding:"required,min=0"`
	DepartmentID *uint   `json:"department_id"`
}

// UpdateMedicalServiceRequest represents request for updating MedicalService
type UpdateMedicalServiceRequest struct {
	Name         string  `json:"name" binding:"max=100"`
	Description  string  `json:"description"`
	ServiceType  string  `json:"service_type"`
	BasePrice    float64 `json:"base_price" binding:"min=0"`
	IsActive     bool    `json:"is_active"`
	DepartmentID *uint   `json:"department_id"`
}
