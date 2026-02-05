package domain

import (
	"time"

	"gorm.io/gorm"
)

// ServiceType represents the type of medical service
type ServiceType string

const (
	ServiceTypeConsultation ServiceType = "CONSULTATION"
	ServiceTypeLabTest      ServiceType = "LAB_TEST"
	ServiceTypeImaging      ServiceType = "IMAGING"
	ServiceTypeProcedure    ServiceType = "PROCEDURE"
	ServiceTypeOther        ServiceType = "OTHER"
)

// MedicalService represents a medical service offered by the hospital
// Mapped to 'services' table
type MedicalService struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Code         string      `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name         string      `gorm:"size:100;not null" json:"name"`
	Description  string      `gorm:"type:text" json:"description"`
	ServiceType  ServiceType `gorm:"size:50;not null" json:"service_type"`
	BasePrice    float64     `gorm:"type:decimal(10,2);not null" json:"base_price"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	DepartmentID *uint       `json:"department_id,omitempty"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// TableName specifies the table name for MedicalService model
func (MedicalService) TableName() string {
	return "services"
}
