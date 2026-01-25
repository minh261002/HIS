package domain

import (
	"time"

	"gorm.io/gorm"
)

// ICD10Code represents an ICD-10 disease code
type ICD10Code struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ICD-10 Information
	Code        string `gorm:"uniqueIndex;size:10;not null" json:"code"` // e.g., "E11.9", "I10"
	Description string `gorm:"size:500;not null" json:"description"`
	Category    string `gorm:"size:100;index" json:"category"` // e.g., "Endocrine", "Circulatory"
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
}

// TableName specifies the table name for ICD10Code model
func (ICD10Code) TableName() string {
	return "icd10_codes"
}
