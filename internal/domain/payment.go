package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCash      PaymentMethod = "CASH"
	PaymentMethodVNPay     PaymentMethod = "VNPAY"
	PaymentMethodPayOS     PaymentMethod = "PAYOS"
	PaymentMethodInsurance PaymentMethod = "INSURANCE"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

// GatewayResponse represents payment gateway response data
type GatewayResponse map[string]interface{}

// Scan implements the sql.Scanner interface
func (g *GatewayResponse) Scan(value interface{}) error {
	if value == nil {
		*g = GatewayResponse{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, g)
}

// Value implements the driver.Valuer interface
func (g GatewayResponse) Value() (driver.Value, error) {
	if g == nil {
		return nil, nil
	}
	return json.Marshal(g)
}

// Payment represents a payment transaction
type Payment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Payment Code (auto-generated: PAY-YYYYMMDD-XXXX)
	PaymentCode string `gorm:"uniqueIndex;size:20;not null" json:"payment_code"`

	// Foreign Keys
	InvoiceID uint     `gorm:"not null;index" json:"invoice_id"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`

	PatientID uint     `gorm:"not null;index" json:"patient_id"` // Denormalized
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`

	// Payment Details
	PaymentMethod   PaymentMethod   `gorm:"size:20;not null;index" json:"payment_method"`
	Amount          float64         `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentDate     time.Time       `gorm:"not null;index" json:"payment_date"`
	Status          PaymentStatus   `gorm:"size:20;not null;index;default:'PENDING'" json:"status"`
	TransactionID   string          `gorm:"size:100" json:"transaction_id"`
	GatewayResponse GatewayResponse `gorm:"type:json" json:"gateway_response"`
	Notes           string          `gorm:"type:text" json:"notes"`

	// Audit fields
	CreatedBy uint `gorm:"not null" json:"created_by"`
}

// TableName specifies the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}
