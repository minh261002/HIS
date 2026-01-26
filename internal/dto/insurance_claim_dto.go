package dto

import "time"

// CreateInsuranceClaimRequest represents request to create insurance claim
type CreateInsuranceClaimRequest struct {
	InvoiceID         uint    `json:"invoice_id" binding:"required"`
	InsuranceProvider string  `json:"insurance_provider" binding:"required"`
	PolicyNumber      string  `json:"policy_number" binding:"required"`
	ClaimAmount       float64 `json:"claim_amount" binding:"required,min=0"`
	Notes             string  `json:"notes" binding:"omitempty"`
}

// ApproveClaimRequest represents request to approve claim
type ApproveClaimRequest struct {
	ApprovedAmount float64 `json:"approved_amount" binding:"required,min=0"`
}

// RejectClaimRequest represents request to reject claim
type RejectClaimRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required"`
}

// InsuranceClaimResponse represents insurance claim details
type InsuranceClaimResponse struct {
	ID                uint       `json:"id"`
	ClaimCode         string     `json:"claim_code"`
	InvoiceID         uint       `json:"invoice_id"`
	PatientID         uint       `json:"patient_id"`
	PatientName       string     `json:"patient_name"`
	InsuranceProvider string     `json:"insurance_provider"`
	PolicyNumber      string     `json:"policy_number"`
	ClaimAmount       float64    `json:"claim_amount"`
	ApprovedAmount    *float64   `json:"approved_amount,omitempty"`
	ClaimDate         time.Time  `json:"claim_date"`
	ApprovalDate      *time.Time `json:"approval_date,omitempty"`
	Status            string     `json:"status"`
	RejectionReason   string     `json:"rejection_reason"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
}

// InsuranceClaimListItem represents simplified claim for list view
type InsuranceClaimListItem struct {
	ID                uint      `json:"id"`
	ClaimCode         string    `json:"claim_code"`
	InsuranceProvider string    `json:"insurance_provider"`
	ClaimAmount       float64   `json:"claim_amount"`
	Status            string    `json:"status"`
	ClaimDate         time.Time `json:"claim_date"`
}
