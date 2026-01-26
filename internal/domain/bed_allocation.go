package domain

import (
	"time"

	"gorm.io/gorm"
)

// BedAllocation represents a bed allocation to an admission
type BedAllocation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	AdmissionID uint       `gorm:"not null;index" json:"admission_id"`
	Admission   *Admission `gorm:"foreignKey:AdmissionID" json:"admission,omitempty"`

	BedID uint `gorm:"not null;index" json:"bed_id"`
	Bed   *Bed `gorm:"foreignKey:BedID" json:"bed,omitempty"`

	// Allocation Details
	AllocatedDate time.Time  `gorm:"not null" json:"allocated_date"`
	ReleasedDate  *time.Time `json:"released_date,omitempty"`
	IsCurrent     bool       `gorm:"default:true;index" json:"is_current"`
	Notes         string     `gorm:"type:text" json:"notes"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
}

// TableName specifies the table name for BedAllocation model
func (BedAllocation) TableName() string {
	return "bed_allocations"
}
