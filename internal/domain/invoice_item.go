package domain

import (
	"time"

	"gorm.io/gorm"
)

// ItemType represents the type of invoice item
type ItemType string

const (
	ItemTypeConsultation ItemType = "CONSULTATION"
	ItemTypeMedication   ItemType = "MEDICATION"
	ItemTypeLabTest      ItemType = "LAB_TEST"
	ItemTypeImaging      ItemType = "IMAGING"
	ItemTypeBed          ItemType = "BED"
	ItemTypeOther        ItemType = "OTHER"
)

// InvoiceItem represents a line item in an invoice
type InvoiceItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key
	InvoiceID uint     `gorm:"not null;index" json:"invoice_id"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`

	// Item Details
	ItemType    ItemType `gorm:"size:20;not null" json:"item_type"`
	ItemID      *uint    `gorm:"index" json:"item_id,omitempty"` // Polymorphic reference
	Description string   `gorm:"type:text;not null" json:"description"`
	Quantity    int      `gorm:"not null" json:"quantity"`
	UnitPrice   float64  `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Amount      float64  `gorm:"type:decimal(10,2);not null" json:"amount"`
}

// TableName specifies the table name for InvoiceItem model
func (InvoiceItem) TableName() string {
	return "invoice_items"
}
