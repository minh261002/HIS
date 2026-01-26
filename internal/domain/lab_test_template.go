package domain

import (
	"time"

	"gorm.io/gorm"
)

// TestCategory represents the category of lab test
type TestCategory string

const (
	TestCategoryHematology   TestCategory = "HEMATOLOGY"
	TestCategoryBiochemistry TestCategory = "BIOCHEMISTRY"
	TestCategoryMicrobiology TestCategory = "MICROBIOLOGY"
	TestCategoryImmunology   TestCategory = "IMMUNOLOGY"
	TestCategoryUrology      TestCategory = "UROLOGY"
	TestCategorySerology     TestCategory = "SEROLOGY"
)

// SampleType represents the type of sample required
type SampleType string

const (
	SampleTypeBlood  SampleType = "BLOOD"
	SampleTypeUrine  SampleType = "URINE"
	SampleTypeStool  SampleType = "STOOL"
	SampleTypeSputum SampleType = "SPUTUM"
	SampleTypeSwab   SampleType = "SWAB"
	SampleTypeOther  SampleType = "OTHER"
)

// LabTestTemplate represents a lab test template
type LabTestTemplate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Template Information
	Code                    string       `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Name                    string       `gorm:"size:200;not null;index" json:"name"`
	Category                TestCategory `gorm:"size:20;not null;index" json:"category"`
	Description             string       `gorm:"type:text" json:"description"`
	SampleType              SampleType   `gorm:"size:20;not null" json:"sample_type"`
	PreparationInstructions string       `gorm:"type:text" json:"preparation_instructions"`
	TurnaroundTimeHours     int          `json:"turnaround_time_hours"`
	Price                   float64      `gorm:"type:decimal(10,2)" json:"price"`
	IsActive                bool         `gorm:"default:true;index" json:"is_active"`

	// Parameters (one-to-many)
	Parameters []*LabTestTemplateParameter `gorm:"foreignKey:TemplateID" json:"parameters,omitempty"`
}

// TableName specifies the table name for LabTestTemplate model
func (LabTestTemplate) TableName() string {
	return "lab_test_templates"
}
