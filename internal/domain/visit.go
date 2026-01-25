package domain

import (
	"time"

	"gorm.io/gorm"
)

// VisitType represents the type of visit
type VisitType string

const (
	VisitTypeScheduled VisitType = "SCHEDULED"
	VisitTypeWalkIn    VisitType = "WALK_IN"
	VisitTypeEmergency VisitType = "EMERGENCY"
	VisitTypeFollowUp  VisitType = "FOLLOW_UP"
)

// VisitStatus represents the status of a visit
type VisitStatus string

const (
	VisitStatusWaiting    VisitStatus = "WAITING"
	VisitStatusInProgress VisitStatus = "IN_PROGRESS"
	VisitStatusCompleted  VisitStatus = "COMPLETED"
	VisitStatusCancelled  VisitStatus = "CANCELLED"
)

// VitalSigns represents patient vital signs
type VitalSigns struct {
	Temperature            float64 `json:"temperature"`              // °C
	BloodPressureSystolic  int     `json:"blood_pressure_systolic"`  // mmHg
	BloodPressureDiastolic int     `json:"blood_pressure_diastolic"` // mmHg
	HeartRate              int     `json:"heart_rate"`               // bpm
	RespiratoryRate        int     `json:"respiratory_rate"`         // breaths/min
	OxygenSaturation       int     `json:"oxygen_saturation"`        // %
	Weight                 float64 `json:"weight"`                   // kg
	Height                 float64 `json:"height"`                   // cm
	BMI                    float64 `json:"bmi"`                      // calculated
}

// Visit represents a patient visit/examination
type Visit struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Visit Code (auto-generated: VST-YYYYMMDD-XXXX)
	VisitCode string `gorm:"uniqueIndex;size:20;not null" json:"visit_code"`

	// Foreign Keys
	AppointmentID *uint        `gorm:"index" json:"appointment_id,omitempty"` // Nullable for walk-ins
	Appointment   *Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	DoctorID uint  `gorm:"not null;index" json:"doctor_id"`
	Doctor   *User `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`

	// Visit Details
	VisitDate time.Time   `gorm:"not null;index" json:"visit_date"`
	VisitTime time.Time   `gorm:"not null" json:"visit_time"`
	VisitType VisitType   `gorm:"size:20;not null" json:"visit_type"`
	Status    VisitStatus `gorm:"size:20;not null;index;default:'WAITING'" json:"status"`

	// Chief Complaint & Symptoms
	ChiefComplaint string `gorm:"type:text;not null" json:"chief_complaint"`
	Symptoms       string `gorm:"type:text" json:"symptoms"`

	// Vital Signs (embedded)
	Temperature            float64 `json:"temperature"`
	BloodPressureSystolic  int     `json:"blood_pressure_systolic"`
	BloodPressureDiastolic int     `json:"blood_pressure_diastolic"`
	HeartRate              int     `json:"heart_rate"`
	RespiratoryRate        int     `json:"respiratory_rate"`
	OxygenSaturation       int     `json:"oxygen_saturation"`
	Weight                 float64 `json:"weight"`
	Height                 float64 `json:"height"`
	BMI                    float64 `json:"bmi"`

	// Clinical Documentation
	PhysicalExamination  string     `gorm:"type:text" json:"physical_examination"`
	ClinicalNotes        string     `gorm:"type:text" json:"clinical_notes"`
	TreatmentPlan        string     `gorm:"type:text" json:"treatment_plan"`
	FollowUpInstructions string     `gorm:"type:text" json:"follow_up_instructions"`
	NextVisitDate        *time.Time `json:"next_visit_date,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Visit model
func (Visit) TableName() string {
	return "visits"
}

// BeforeCreate hook to calculate BMI
func (v *Visit) BeforeCreate(tx *gorm.DB) error {
	if v.Weight > 0 && v.Height > 0 {
		heightInMeters := v.Height / 100
		v.BMI = v.Weight / (heightInMeters * heightInMeters)
	}
	return nil
}

// BeforeUpdate hook to recalculate BMI
func (v *Visit) BeforeUpdate(tx *gorm.DB) error {
	if v.Weight > 0 && v.Height > 0 {
		heightInMeters := v.Height / 100
		v.BMI = v.Weight / (heightInMeters * heightInMeters)
	}
	return nil
}
