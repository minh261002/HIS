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

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment handles creating payment
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	payment, err := h.paymentService.CreatePayment(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvoiceNotFound) {
			response.NotFound(c, "Invoice not found")
			return
		}
		if errors.Is(err, service.ErrPaymentExceedsBalance) {
			response.BadRequest(c, "Payment amount exceeds invoice balance", nil)
			return
		}
		response.InternalServerError(c, "Failed to create payment")
		return
	}

	response.Created(c, "Payment created successfully", payment)
}

// GetInvoicePayments handles getting invoice's payments
func (h *PaymentHandler) GetInvoicePayments(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid invoice ID", nil)
		return
	}

	payments, err := h.paymentService.GetInvoicePayments(uint(invoiceID))
	if err != nil {
		response.InternalServerError(c, "Failed to get invoice payments")
		return
	}

	response.Success(c, "Invoice payments retrieved successfully", payments)
}
