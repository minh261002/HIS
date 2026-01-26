package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// InsuranceClaimHandler handles insurance claim HTTP requests
type InsuranceClaimHandler struct {
	claimService *service.InsuranceClaimService
}

// NewInsuranceClaimHandler creates a new insurance claim handler
func NewInsuranceClaimHandler(claimService *service.InsuranceClaimService) *InsuranceClaimHandler {
	return &InsuranceClaimHandler{claimService: claimService}
}

// CreateInsuranceClaim handles creating insurance claim
func (h *InsuranceClaimHandler) CreateInsuranceClaim(c *gin.Context) {
	var req dto.CreateInsuranceClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	claim, err := h.claimService.CreateInsuranceClaim(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvoiceNotFound) {
			response.NotFound(c, "Invoice not found")
			return
		}
		response.InternalServerError(c, "Failed to create insurance claim")
		return
	}

	response.Created(c, "Insurance claim created successfully", claim)
}

// ApproveClaim handles approving insurance claim
func (h *InsuranceClaimHandler) ApproveClaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid claim ID", nil)
		return
	}

	var req dto.ApproveClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.claimService.ApproveClaim(uint(id), &req, userID); err != nil {
		if errors.Is(err, service.ErrInsuranceClaimNotFound) {
			response.NotFound(c, "Insurance claim not found")
			return
		}
		response.InternalServerError(c, "Failed to approve claim")
		return
	}

	response.Success(c, "Insurance claim approved successfully", nil)
}

// RejectClaim handles rejecting insurance claim
func (h *InsuranceClaimHandler) RejectClaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid claim ID", nil)
		return
	}

	var req dto.RejectClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.claimService.RejectClaim(uint(id), &req, userID); err != nil {
		if errors.Is(err, service.ErrInsuranceClaimNotFound) {
			response.NotFound(c, "Insurance claim not found")
			return
		}
		response.InternalServerError(c, "Failed to reject claim")
		return
	}

	response.Success(c, "Insurance claim rejected successfully", nil)
}

// GetInvoiceClaims handles getting invoice's insurance claims
func (h *InsuranceClaimHandler) GetInvoiceClaims(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid invoice ID", nil)
		return
	}

	claims, err := h.claimService.GetInvoiceClaims(uint(invoiceID))
	if err != nil {
		response.InternalServerError(c, "Failed to get invoice claims")
		return
	}

	response.Success(c, "Invoice claims retrieved successfully", claims)
}
