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

// DispensingHandler handles dispensing HTTP requests
type DispensingHandler struct {
	dispensingService *service.DispensingService
}

// NewDispensingHandler creates a new dispensing handler
func NewDispensingHandler(dispensingService *service.DispensingService) *DispensingHandler {
	return &DispensingHandler{dispensingService: dispensingService}
}

// DispensePrescription handles dispensing prescription
func (h *DispensingHandler) DispensePrescription(c *gin.Context) {
	var req dto.DispensePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	dispensings, err := h.dispensingService.DispensePrescription(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			response.BadRequest(c, "Insufficient stock", nil)
			return
		}
		if errors.Is(err, service.ErrExpiredStock) {
			response.BadRequest(c, "Stock has expired", nil)
			return
		}
		response.InternalServerError(c, "Failed to dispense prescription")
		return
	}

	response.Created(c, "Prescription dispensed successfully", dispensings)
}

// GetPrescriptionDispensingRecords handles getting prescription dispensing records
func (h *DispensingHandler) GetPrescriptionDispensingRecords(c *gin.Context) {
	prescriptionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	records, err := h.dispensingService.GetPrescriptionDispensingRecords(uint(prescriptionID))
	if err != nil {
		response.InternalServerError(c, "Failed to get dispensing records")
		return
	}

	response.Success(c, "Dispensing records retrieved successfully", records)
}

// GetPatientDispensingHistory handles getting patient dispensing history
func (h *DispensingHandler) GetPatientDispensingHistory(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	history, err := h.dispensingService.GetPatientDispensingHistory(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get dispensing history")
		return
	}

	response.Success(c, "Dispensing history retrieved successfully", history)
}
