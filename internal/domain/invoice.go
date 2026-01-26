package domain

import (
	"time"

	"gorm.io/gorm"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusPending       InvoiceStatus = "PENDING"
	InvoiceStatusPaid          InvoiceStatus = "PAID"
	InvoiceStatusPartiallyPaid InvoiceStatus = "PARTIALLY_PAID"
	InvoiceStatusCancelled     InvoiceStatus = "CANCELLED"
)

// Invoice represents a billing invoice
type Invoice struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Invoice Code (auto-generated: INV-YYYYMMDD-XXXX)
	InvoiceCode string `gorm:"uniqueIndex;size:20;not null" json:"invoice_code"`

	// Foreign Keys
	VisitID *uint  `gorm:"index" json:"visit_id,omitempty"` // Nullable
	Visit   *Visit `gorm:"foreignKey:VisitID" json:"visit,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	// Invoice Details
	InvoiceDate    time.Time     `gorm:"not null;index" json:"invoice_date"`
	DueDate        time.Time     `gorm:"not null" json:"due_date"`
	Subtotal       float64       `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	TaxAmount      float64       `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	DiscountAmount float64       `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TotalAmount    float64       `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status         InvoiceStatus `gorm:"size:20;not null;index;default:'PENDING'" json:"status"`
	Notes          string        `gorm:"type:text" json:"notes"`

	// Relationships
	Items           []*InvoiceItem    `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	Payments        []*Payment        `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
	InsuranceClaims []*InsuranceClaim `gorm:"foreignKey:InvoiceID" json:"insurance_claims,omitempty"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TableName specifies the table name for Invoice model
func (Invoice) TableName() string {
	return "invoices"
}
