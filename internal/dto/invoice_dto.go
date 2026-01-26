package dto

import "time"

// InvoiceItemRequest represents an invoice line item
type InvoiceItemRequest struct {
	ItemType    string  `json:"item_type" binding:"required"`
	ItemID      *uint   `json:"item_id" binding:"omitempty"`
	Description string  `json:"description" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required,min=0"`
}

// CreateInvoiceRequest represents request to create invoice
type CreateInvoiceRequest struct {
	VisitID        *uint                 `json:"visit_id" binding:"omitempty"`
	PatientID      uint                  `json:"patient_id" binding:"required"`
	Items          []*InvoiceItemRequest `json:"items" binding:"required,min=1,dive"`
	TaxAmount      float64               `json:"tax_amount" binding:"omitempty"`
	DiscountAmount float64               `json:"discount_amount" binding:"omitempty"`
	Notes          string                `json:"notes" binding:"omitempty"`
}

// InvoiceItemResponse represents invoice item details
type InvoiceItemResponse struct {
	ID          uint    `json:"id"`
	ItemType    string  `json:"item_type"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}

// InvoiceResponse represents invoice details
type InvoiceResponse struct {
	ID             uint                   `json:"id"`
	InvoiceCode    string                 `json:"invoice_code"`
	VisitID        *uint                  `json:"visit_id,omitempty"`
	PatientID      uint                   `json:"patient_id"`
	PatientName    string                 `json:"patient_name"`
	InvoiceDate    time.Time              `json:"invoice_date"`
	DueDate        time.Time              `json:"due_date"`
	Subtotal       float64                `json:"subtotal"`
	TaxAmount      float64                `json:"tax_amount"`
	DiscountAmount float64                `json:"discount_amount"`
	TotalAmount    float64                `json:"total_amount"`
	Status         string                 `json:"status"`
	Notes          string                 `json:"notes"`
	Items          []*InvoiceItemResponse `json:"items"`
	CreatedAt      time.Time              `json:"created_at"`
}

// InvoiceListItem represents simplified invoice for list view
type InvoiceListItem struct {
	ID          uint      `json:"id"`
	InvoiceCode string    `json:"invoice_code"`
	PatientName string    `json:"patient_name"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	InvoiceDate time.Time `json:"invoice_date"`
}
