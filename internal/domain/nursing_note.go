package domain

import (
	"time"

	"gorm.io/gorm"
)

// NursingNote represents a nursing note
type NursingNote struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	AdmissionID uint       `gorm:"not null;index" json:"admission_id"`
	Admission   *Admission `gorm:"foreignKey:AdmissionID" json:"admission,omitempty"`

	NurseID uint  `gorm:"not null;index" json:"nurse_id"`
	Nurse   *User `gorm:"foreignKey:NurseID" json:"nurse,omitempty"`

	// Note Details
	NoteDate      time.Time  `gorm:"not null;index" json:"note_date"`
	VitalSigns    VitalSigns `gorm:"type:json" json:"vital_signs"`
	Observations  string     `gorm:"type:text;not null" json:"observations"`
	Interventions string     `gorm:"type:text" json:"interventions"`
}

// TableName specifies the table name for NursingNote model
func (NursingNote) TableName() string {
	return "nursing_notes"
}
