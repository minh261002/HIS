package domain

import (
	"time"

	"gorm.io/gorm"
)

// LabTestTemplateParameter represents a parameter in a lab test template
type LabTestTemplateParameter struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	TemplateID uint             `gorm:"not null;index" json:"template_id"`
	Template   *LabTestTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`

	// Parameter Information
	ParameterName   string  `gorm:"size:100;not null" json:"parameter_name"`
	Unit            string  `gorm:"size:50" json:"unit"`
	NormalRangeMin  float64 `gorm:"type:decimal(10,2)" json:"normal_range_min"`
	NormalRangeMax  float64 `gorm:"type:decimal(10,2)" json:"normal_range_max"`
	NormalRangeText string  `gorm:"size:100" json:"normal_range_text"` // For display
	DisplayOrder    int     `gorm:"default:0" json:"display_order"`
}

// TableName specifies the table name for LabTestTemplateParameter model
func (LabTestTemplateParameter) TableName() string {
	return "lab_test_template_parameters"
}
