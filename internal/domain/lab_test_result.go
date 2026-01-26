package domain

import (
	"time"

	"gorm.io/gorm"
)

// LabTestResult represents a lab test result
type LabTestResult struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	RequestID uint            `gorm:"not null;index" json:"request_id"`
	Request   *LabTestRequest `gorm:"foreignKey:RequestID" json:"request,omitempty"`

	// Result Information
	ParameterName   string `gorm:"size:100;not null" json:"parameter_name"`
	Value           string `gorm:"size:100" json:"value"`
	Unit            string `gorm:"size:50" json:"unit"`
	NormalRangeText string `gorm:"size:100" json:"normal_range_text"`
	IsAbnormal      bool   `gorm:"default:false" json:"is_abnormal"`
	Remarks         string `gorm:"type:text" json:"remarks"`
}

// TableName specifies the table name for LabTestResult model
func (LabTestResult) TableName() string {
	return "lab_test_results"
}
