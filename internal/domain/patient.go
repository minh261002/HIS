package domain

import (
	"time"

	"gorm.io/gorm"
)

// Gender represents patient gender
type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

// BloodType represents patient blood type
type BloodType string

const (
	BloodTypeAPositive  BloodType = "A+"
	BloodTypeANegative  BloodType = "A-"
	BloodTypeBPositive  BloodType = "B+"
	BloodTypeBNegative  BloodType = "B-"
	BloodTypeABPositive BloodType = "AB+"
	BloodTypeABNegative BloodType = "AB-"
	BloodTypeOPositive  BloodType = "O+"
	BloodTypeONegative  BloodType = "O-"
)

// Patient represents a patient in the hospital
type Patient struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Patient Code (auto-generated: P-YYYYMMDD-XXXX)
	PatientCode string `gorm:"uniqueIndex;size:20;not null" json:"patient_code"`

	// Personal Information
	FirstName string `gorm:"size:50;not null" json:"first_name"`
	LastName  string `gorm:"size:50;not null" json:"last_name"`
	FullName  string `gorm:"size:100;not null;index" json:"full_name"`

	DateOfBirth time.Time `gorm:"not null" json:"date_of_birth"`
	Age         int       `gorm:"-" json:"age"` // Calculated field
	Gender      Gender    `gorm:"size:10;not null" json:"gender"`
	BloodType   BloodType `gorm:"size:5" json:"blood_type"`

	// Contact Information
	PhoneNumber string `gorm:"size:20;index" json:"phone_number"`
	Email       string `gorm:"size:100;index" json:"email"`
	Address     string `gorm:"size:255" json:"address"`
	City        string `gorm:"size:100" json:"city"`
	State       string `gorm:"size:100" json:"state"`
	PostalCode  string `gorm:"size:20" json:"postal_code"`
	Country     string `gorm:"size:100;default:'Vietnam'" json:"country"`

	// Identification
	NationalID string `gorm:"size:20;uniqueIndex" json:"national_id"` // CCCD/CMND

	// Insurance Information
	InsuranceNumber   string `gorm:"size:50" json:"insurance_number"`
	InsuranceProvider string `gorm:"size:100" json:"insurance_provider"`

	// Emergency Contact
	EmergencyContactName         string `gorm:"size:100" json:"emergency_contact_name"`
	EmergencyContactPhone        string `gorm:"size:20" json:"emergency_contact_phone"`
	EmergencyContactRelationship string `gorm:"size:50" json:"emergency_contact_relationship"`

	// Medical Information
	Allergies         string `gorm:"type:text" json:"allergies"`          // JSON or comma-separated
	ChronicConditions string `gorm:"type:text" json:"chronic_conditions"` // JSON or comma-separated
	Notes             string `gorm:"type:text" json:"notes"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Patient model
func (Patient) TableName() string {
	return "patients"
}

// BeforeCreate hook to calculate age
func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	p.FullName = p.FirstName + " " + p.LastName
	p.Age = calculateAge(p.DateOfBirth)
	return nil
}

// BeforeUpdate hook to update age and full name
func (p *Patient) BeforeUpdate(tx *gorm.DB) error {
	p.FullName = p.FirstName + " " + p.LastName
	p.Age = calculateAge(p.DateOfBirth)
	return nil
}

// AfterFind hook to calculate age when loading from database
func (p *Patient) AfterFind(tx *gorm.DB) error {
	p.Age = calculateAge(p.DateOfBirth)
	return nil
}

// calculateAge calculates age from date of birth
func calculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}
