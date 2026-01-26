package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	ErrInsuranceClaimNotFound = errors.New("insurance claim not found")
)

// InsuranceClaimService handles insurance claim business logic
type InsuranceClaimService struct {
	claimRepo   *repository.InsuranceClaimRepository
	invoiceRepo *repository.InvoiceRepository
}

// NewInsuranceClaimService creates a new insurance claim service
func NewInsuranceClaimService(
	claimRepo *repository.InsuranceClaimRepository,
	invoiceRepo *repository.InvoiceRepository,
) *InsuranceClaimService {
	return &InsuranceClaimService{
		claimRepo:   claimRepo,
		invoiceRepo: invoiceRepo,
	}
}

// CreateInsuranceClaim creates insurance claim
func (s *InsuranceClaimService) CreateInsuranceClaim(req *dto.CreateInsuranceClaimRequest, createdBy uint) (*dto.InsuranceClaimResponse, error) {
	// Validate invoice exists
	invoice, err := s.invoiceRepo.FindByID(req.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find invoice: %w", err)
	}
	if invoice == nil {
		return nil, ErrInvoiceNotFound
	}

	// Generate claim code
	code, err := s.claimRepo.GenerateClaimCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate claim code: %w", err)
	}

	// Create claim
	claim := &domain.InsuranceClaim{
		ClaimCode:         code,
		InvoiceID:         req.InvoiceID,
		PatientID:         invoice.PatientID,
		InsuranceProvider: req.InsuranceProvider,
		PolicyNumber:      req.PolicyNumber,
		ClaimAmount:       req.ClaimAmount,
		ClaimDate:         time.Now(),
		Status:            domain.ClaimStatusSubmitted,
		Notes:             req.Notes,
		CreatedBy:         createdBy,
	}

	if err := s.claimRepo.Create(claim); err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}

	// Reload to get relationships
	claim, _ = s.claimRepo.FindByID(claim.ID)
	return s.toInsuranceClaimResponse(claim), nil
}

// ApproveClaim approves insurance claim
func (s *InsuranceClaimService) ApproveClaim(id uint, req *dto.ApproveClaimRequest, updatedBy uint) error {
	claim, err := s.claimRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find claim: %w", err)
	}
	if claim == nil {
		return ErrInsuranceClaimNotFound
	}

	now := time.Now()
	claim.ApprovedAmount = &req.ApprovedAmount
	claim.ApprovalDate = &now
	claim.Status = domain.ClaimStatusApproved
	claim.UpdatedBy = updatedBy

	return s.claimRepo.Update(claim)
}

// RejectClaim rejects insurance claim
func (s *InsuranceClaimService) RejectClaim(id uint, req *dto.RejectClaimRequest, updatedBy uint) error {
	claim, err := s.claimRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find claim: %w", err)
	}
	if claim == nil {
		return ErrInsuranceClaimNotFound
	}

	claim.RejectionReason = req.RejectionReason
	claim.Status = domain.ClaimStatusRejected
	claim.UpdatedBy = updatedBy

	return s.claimRepo.Update(claim)
}

// GetInvoiceClaims gets invoice's insurance claims
func (s *InsuranceClaimService) GetInvoiceClaims(invoiceID uint) ([]*dto.InsuranceClaimListItem, error) {
	claims, err := s.claimRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice claims: %w", err)
	}

	items := make([]*dto.InsuranceClaimListItem, len(claims))
	for i, c := range claims {
		items[i] = s.toInsuranceClaimListItem(c)
	}
	return items, nil
}

// Helper functions
func (s *InsuranceClaimService) toInsuranceClaimResponse(c *domain.InsuranceClaim) *dto.InsuranceClaimResponse {
	resp := &dto.InsuranceClaimResponse{
		ID:                c.ID,
		ClaimCode:         c.ClaimCode,
		InvoiceID:         c.InvoiceID,
		PatientID:         c.PatientID,
		InsuranceProvider: c.InsuranceProvider,
		PolicyNumber:      c.PolicyNumber,
		ClaimAmount:       c.ClaimAmount,
		ApprovedAmount:    c.ApprovedAmount,
		ClaimDate:         c.ClaimDate,
		ApprovalDate:      c.ApprovalDate,
		Status:            string(c.Status),
		RejectionReason:   c.RejectionReason,
		Notes:             c.Notes,
		CreatedAt:         c.CreatedAt,
	}

	if c.Patient != nil {
		resp.PatientName = c.Patient.FullName
	}

	return resp
}

func (s *InsuranceClaimService) toInsuranceClaimListItem(c *domain.InsuranceClaim) *dto.InsuranceClaimListItem {
	return &dto.InsuranceClaimListItem{
		ID:                c.ID,
		ClaimCode:         c.ClaimCode,
		InsuranceProvider: c.InsuranceProvider,
		ClaimAmount:       c.ClaimAmount,
		Status:            string(c.Status),
		ClaimDate:         c.ClaimDate,
	}
}
