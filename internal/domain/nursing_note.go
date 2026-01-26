package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// VitalSigns represents patient vital signs
type VitalSigns struct {
	Temperature     float64 `json:"temperature"`      // Celsius
	BloodPressure   string  `json:"blood_pressure"`   // e.g., "120/80"
	Pulse           int     `json:"pulse"`            // beats per minute
	RespiratoryRate int     `json:"respiratory_rate"` // breaths per minute
	OxygenSat       int     `json:"oxygen_sat"`       // percentage
}

// Scan implements the sql.Scanner interface
func (v *VitalSigns) Scan(value interface{}) error {
	if value == nil {
		*v = VitalSigns{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, v)
}

// Value implements the driver.Valuer interface
func (v VitalSigns) Value() (driver.Value, error) {
	return json.Marshal(v)
}

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
