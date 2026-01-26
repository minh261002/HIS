package dto

import "time"

// CreateInventoryRequest represents request to add inventory
type CreateInventoryRequest struct {
	MedicationID uint    `json:"medication_id" binding:"required"`
	BatchNumber  string  `json:"batch_number" binding:"required"`
	ExpiryDate   string  `json:"expiry_date" binding:"required"` // YYYY-MM-DD format
	Quantity     int     `json:"quantity" binding:"required,min=1"`
	Unit         string  `json:"unit" binding:"required"`
	CostPrice    float64 `json:"cost_price" binding:"omitempty"`
	Supplier     string  `json:"supplier" binding:"omitempty"`
	ReceivedDate string  `json:"received_date" binding:"required"` // YYYY-MM-DD format
}

// InventoryResponse represents inventory details
type InventoryResponse struct {
	ID             uint      `json:"id"`
	MedicationID   uint      `json:"medication_id"`
	MedicationName string    `json:"medication_name"`
	BatchNumber    string    `json:"batch_number"`
	ExpiryDate     string    `json:"expiry_date"`
	Quantity       int       `json:"quantity"`
	Unit           string    `json:"unit"`
	CostPrice      float64   `json:"cost_price"`
	Supplier       string    `json:"supplier"`
	ReceivedDate   string    `json:"received_date"`
	CreatedAt      time.Time `json:"created_at"`
}

// InventoryListItem represents simplified inventory for list view
type InventoryListItem struct {
	ID             uint   `json:"id"`
	MedicationName string `json:"medication_name"`
	BatchNumber    string `json:"batch_number"`
	ExpiryDate     string `json:"expiry_date"`
	Quantity       int    `json:"quantity"`
	Unit           string `json:"unit"`
}

// LowStockAlert represents low stock alert
type LowStockAlert struct {
	MedicationID   uint   `json:"medication_id"`
	MedicationName string `json:"medication_name"`
	TotalQuantity  int    `json:"total_quantity"`
	Unit           string `json:"unit"`
}
