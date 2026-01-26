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

// InvoiceHandler handles invoice HTTP requests
type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

// CreateInvoice handles creating invoice
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	invoice, err := h.invoiceService.CreateInvoice(&req, userID)
	if err != nil {
		response.InternalServerError(c, "Failed to create invoice")
		return
	}

	response.Created(c, "Invoice created successfully", invoice)
}

// GetInvoice handles getting invoice details
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid invoice ID", nil)
		return
	}

	invoice, err := h.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrInvoiceNotFound) {
			response.NotFound(c, "Invoice not found")
			return
		}
		response.InternalServerError(c, "Failed to get invoice")
		return
	}

	response.Success(c, "Invoice retrieved successfully", invoice)
}

// GetInvoiceByCode handles getting invoice by code
func (h *InvoiceHandler) GetInvoiceByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Invoice code is required", nil)
		return
	}

	invoice, err := h.invoiceService.GetInvoiceByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrInvoiceNotFound) {
			response.NotFound(c, "Invoice not found")
			return
		}
		response.InternalServerError(c, "Failed to get invoice")
		return
	}

	response.Success(c, "Invoice retrieved successfully", invoice)
}

// GetPatientInvoices handles getting patient's invoices
func (h *InvoiceHandler) GetPatientInvoices(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	invoices, err := h.invoiceService.GetPatientInvoices(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get patient invoices")
		return
	}

	response.Success(c, "Patient invoices retrieved successfully", invoices)
}
