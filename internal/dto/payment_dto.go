package dto

import "time"

// CreatePaymentRequest represents request to create payment
type CreatePaymentRequest struct {
	InvoiceID     uint    `json:"invoice_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,min=0"`
	Notes         string  `json:"notes" binding:"omitempty"`
}

// PaymentResponse represents payment details
type PaymentResponse struct {
	ID            uint      `json:"id"`
	PaymentCode   string    `json:"payment_code"`
	InvoiceID     uint      `json:"invoice_id"`
	PatientID     uint      `json:"patient_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	PaymentURL    string    `json:"payment_url,omitempty"` // For VNPay/PayOS
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// PaymentListItem represents simplified payment for list view
type PaymentListItem struct {
	ID            uint      `json:"id"`
	PaymentCode   string    `json:"payment_code"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	PaymentDate   time.Time `json:"payment_date"`
}
