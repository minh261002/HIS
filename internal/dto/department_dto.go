package dto

import "time"

// DepartmentResponse represents department data in response
type DepartmentResponse struct {
	ID          uint                `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	IsActive    bool                `json:"is_active"`
	HeadDoctor  *UserDetailResponse `json:"head_doctor,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// CreateDepartmentRequest represents request for creating Department
type CreateDepartmentRequest struct {
	Code         string `json:"code" binding:"required,max=50"`
	Name         string `json:"name" binding:"required,max=100"`
	Description  string `json:"description" binding:"max=255"`
	HeadDoctorID *uint  `json:"head_doctor_id"`
}

// UpdateDepartmentRequest represents request for updating Department
type UpdateDepartmentRequest struct {
	Name         string `json:"name" binding:"max=100"`
	Description  string `json:"description" binding:"max=255"`
	IsActive     bool   `json:"is_active"`
	HeadDoctorID *uint  `json:"head_doctor_id"`
}
