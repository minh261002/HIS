package domain

import (
	"time"

	"gorm.io/gorm"
)

// ClaimStatus represents the status of an insurance claim
type ClaimStatus string

const (
	ClaimStatusSubmitted ClaimStatus = "SUBMITTED"
	ClaimStatusApproved  ClaimStatus = "APPROVED"
	ClaimStatusRejected  ClaimStatus = "REJECTED"
	ClaimStatusPending   ClaimStatus = "PENDING"
)

// InsuranceClaim represents an insurance claim
type InsuranceClaim struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Claim Code (auto-generated: CLM-YYYYMMDD-XXXX)
	ClaimCode string `gorm:"uniqueIndex;size:20;not null" json:"claim_code"`

	// Foreign Keys
	InvoiceID uint     `gorm:"not null;index" json:"invoice_id"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	// Claim Details
	InsuranceProvider string      `gorm:"size:200;not null" json:"insurance_provider"`
	PolicyNumber      string      `gorm:"size:100;not null" json:"policy_number"`
	ClaimAmount       float64     `gorm:"type:decimal(10,2);not null" json:"claim_amount"`
	ApprovedAmount    *float64    `gorm:"type:decimal(10,2)" json:"approved_amount,omitempty"`
	ClaimDate         time.Time   `gorm:"not null;index" json:"claim_date"`
	ApprovalDate      *time.Time  `json:"approval_date,omitempty"`
	Status            ClaimStatus `gorm:"size:20;not null;index;default:'SUBMITTED'" json:"status"`
	RejectionReason   string      `gorm:"type:text" json:"rejection_reason"`
	Notes             string      `gorm:"type:text" json:"notes"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for InsuranceClaim model
func (InsuranceClaim) TableName() string {
	return "insurance_claims"
}
