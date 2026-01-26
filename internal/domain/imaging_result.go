package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DICOMFiles represents an array of DICOM file paths
type DICOMFiles []string

// Scan implements the sql.Scanner interface
func (d *DICOMFiles) Scan(value interface{}) error {
	if value == nil {
		*d = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, d)
}

// Value implements the driver.Valuer interface
func (d DICOMFiles) Value() (driver.Value, error) {
	if len(d) == 0 {
		return "[]", nil
	}
	return json.Marshal(d)
}

// ImagingResult represents an imaging result
type ImagingResult struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key (one-to-one with request)
	RequestID uint            `gorm:"uniqueIndex;not null" json:"request_id"`
	Request   *ImagingRequest `gorm:"foreignKey:RequestID" json:"request,omitempty"`

	RadiologistID uint  `gorm:"not null;index" json:"radiologist_id"`
	Radiologist   *User `gorm:"foreignKey:RadiologistID" json:"radiologist,omitempty"`

	// Result Information
	Findings   string     `gorm:"type:text;not null" json:"findings"`
	Impression string     `gorm:"type:text;not null" json:"impression"`
	DICOMFiles DICOMFiles `gorm:"type:json" json:"dicom_files"`
	ReportDate time.Time  `gorm:"not null" json:"report_date"`
	IsCritical bool       `gorm:"default:false" json:"is_critical"`
}

// TableName specifies the table name for ImagingResult model
func (ImagingResult) TableName() string {
	return "imaging_results"
}
