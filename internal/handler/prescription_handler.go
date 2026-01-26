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

// PrescriptionHandler handles prescription HTTP requests
type PrescriptionHandler struct {
	prescriptionService *service.PrescriptionService
}

// NewPrescriptionHandler creates a new prescription handler
func NewPrescriptionHandler(prescriptionService *service.PrescriptionService) *PrescriptionHandler {
	return &PrescriptionHandler{prescriptionService: prescriptionService}
}

// CreatePrescription handles creating a prescription
func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	var req dto.CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	prescription, err := h.prescriptionService.CreatePrescription(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		if errors.Is(err, service.ErrMedicationNotFound) {
			response.NotFound(c, "Medication not found")
			return
		}
		response.InternalServerError(c, "Failed to create prescription")
		return
	}

	response.Created(c, "Prescription created successfully", prescription)
}

// GetPrescription handles getting prescription details
func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	prescription, err := h.prescriptionService.GetPrescriptionByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to get prescription")
		return
	}

	response.Success(c, "Prescription retrieved successfully", prescription)
}

// GetPrescriptionByCode handles getting prescription by code
func (h *PrescriptionHandler) GetPrescriptionByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Prescription code is required", nil)
		return
	}

	prescription, err := h.prescriptionService.GetPrescriptionByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to get prescription")
		return
	}

	response.Success(c, "Prescription retrieved successfully", prescription)
}

// UpdatePrescription handles updating prescription
func (h *PrescriptionHandler) UpdatePrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	var req dto.UpdatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	prescription, err := h.prescriptionService.UpdatePrescription(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to update prescription")
		return
	}

	response.Success(c, "Prescription updated successfully", prescription)
}

// DispensePrescription handles dispensing prescription
func (h *PrescriptionHandler) DispensePrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.prescriptionService.DispensePrescription(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to dispense prescription")
		return
	}

	response.Success(c, "Prescription dispensed successfully", nil)
}

// CompletePrescription handles completing prescription
func (h *PrescriptionHandler) CompletePrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.prescriptionService.CompletePrescription(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to complete prescription")
		return
	}

	response.Success(c, "Prescription completed successfully", nil)
}

// CancelPrescription handles cancelling prescription
func (h *PrescriptionHandler) CancelPrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid prescription ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.prescriptionService.CancelPrescription(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			response.NotFound(c, "Prescription not found")
			return
		}
		response.InternalServerError(c, "Failed to cancel prescription")
		return
	}

	response.Success(c, "Prescription cancelled successfully", nil)
}

// GetVisitPrescriptions handles getting visit prescriptions
func (h *PrescriptionHandler) GetVisitPrescriptions(c *gin.Context) {
	visitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	prescriptions, err := h.prescriptionService.GetVisitPrescriptions(uint(visitID))
	if err != nil {
		response.InternalServerError(c, "Failed to get visit prescriptions")
		return
	}

	response.Success(c, "Visit prescriptions retrieved successfully", prescriptions)
}

// GetPatientPrescriptions handles getting patient prescriptions
func (h *PrescriptionHandler) GetPatientPrescriptions(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	prescriptions, err := h.prescriptionService.GetPatientPrescriptions(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient prescriptions")
		return
	}

	response.Success(c, "Patient prescriptions retrieved successfully", prescriptions)
}
