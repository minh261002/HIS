package domain

import (
	"time"

	"gorm.io/gorm"
)

// Department represents hospital department
type Department string

const (
	DepartmentInternalMedicine Department = "INTERNAL_MEDICINE"
	DepartmentSurgery          Department = "SURGERY"
	DepartmentPediatrics       Department = "PEDIATRICS"
	DepartmentICU              Department = "ICU"
	DepartmentObstetrics       Department = "OBSTETRICS"
	DepartmentEmergency        Department = "EMERGENCY"
)

// BedType represents the type of bed
type BedType string

const (
	BedTypeStandard  BedType = "STANDARD"
	BedTypeICU       BedType = "ICU"
	BedTypeIsolation BedType = "ISOLATION"
	BedTypeVIP       BedType = "VIP"
)

// BedStatus represents the status of a bed
type BedStatus string

const (
	BedStatusAvailable   BedStatus = "AVAILABLE"
	BedStatusOccupied    BedStatus = "OCCUPIED"
	BedStatusMaintenance BedStatus = "MAINTENANCE"
	BedStatusReserved    BedStatus = "RESERVED"
)

// Bed represents a hospital bed
type Bed struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Bed Information
	BedNumber  string     `gorm:"uniqueIndex;size:20;not null" json:"bed_number"`
	Department Department `gorm:"size:30;not null;index" json:"department"`
	Ward       string     `gorm:"size:50;not null" json:"ward"`
	BedType    BedType    `gorm:"size:20;not null;index" json:"bed_type"`
	Status     BedStatus  `gorm:"size:20;not null;index;default:'AVAILABLE'" json:"status"`
	IsActive   bool       `gorm:"default:true;index" json:"is_active"`
}

// TableName specifies the table name for Bed model
func (Bed) TableName() string {
	return "beds"
}
