package domain

import (
	"time"

	"gorm.io/gorm"
)

// Department represents a hospital department
type Department struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Code        string `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Orientation/Head of Department (optional)
	HeadDoctorID *uint `json:"head_doctor_id,omitempty"`
	HeadDoctor   *User `gorm:"foreignKey:HeadDoctorID" json:"head_doctor,omitempty"`

	// Audit
	CreatedBy uint `json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Department model
func (Department) TableName() string {
	return "departments"
}
