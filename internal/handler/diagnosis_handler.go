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

// DiagnosisHandler handles diagnosis HTTP requests
type DiagnosisHandler struct {
	diagnosisService *service.DiagnosisService
}

// NewDiagnosisHandler creates a new diagnosis handler
func NewDiagnosisHandler(diagnosisService *service.DiagnosisService) *DiagnosisHandler {
	return &DiagnosisHandler{diagnosisService: diagnosisService}
}

// AddDiagnosis handles adding a diagnosis
func (h *DiagnosisHandler) AddDiagnosis(c *gin.Context) {
	var req dto.CreateDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	diagnosis, err := h.diagnosisService.AddDiagnosis(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		if errors.Is(err, service.ErrICD10CodeNotFound) {
			response.NotFound(c, "ICD-10 code not found")
			return
		}
		if errors.Is(err, service.ErrPrimaryDiagnosisExists) {
			response.BadRequest(c, "Primary diagnosis already exists for this visit", nil)
			return
		}
		response.InternalServerError(c, "Failed to add diagnosis")
		return
	}

	response.Created(c, "Diagnosis added successfully", diagnosis)
}

// GetDiagnosis handles getting diagnosis details
func (h *DiagnosisHandler) GetDiagnosis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid diagnosis ID", nil)
		return
	}

	diagnosis, err := h.diagnosisService.GetDiagnosisByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrDiagnosisNotFound) {
			response.NotFound(c, "Diagnosis not found")
			return
		}
		response.InternalServerError(c, "Failed to get diagnosis")
		return
	}

	response.Success(c, "Diagnosis retrieved successfully", diagnosis)
}

// UpdateDiagnosis handles updating a diagnosis
func (h *DiagnosisHandler) UpdateDiagnosis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid diagnosis ID", nil)
		return
	}

	var req dto.UpdateDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	diagnosis, err := h.diagnosisService.UpdateDiagnosis(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrDiagnosisNotFound) {
			response.NotFound(c, "Diagnosis not found")
			return
		}
		response.InternalServerError(c, "Failed to update diagnosis")
		return
	}

	response.Success(c, "Diagnosis updated successfully", diagnosis)
}

// DeleteDiagnosis handles deleting a diagnosis
func (h *DiagnosisHandler) DeleteDiagnosis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid diagnosis ID", nil)
		return
	}

	if err := h.diagnosisService.DeleteDiagnosis(uint(id)); err != nil {
		if errors.Is(err, service.ErrDiagnosisNotFound) {
			response.NotFound(c, "Diagnosis not found")
			return
		}
		response.InternalServerError(c, "Failed to delete diagnosis")
		return
	}

	response.Success(c, "Diagnosis deleted successfully", nil)
}

// GetVisitDiagnoses handles getting visit diagnoses
func (h *DiagnosisHandler) GetVisitDiagnoses(c *gin.Context) {
	visitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	diagnoses, err := h.diagnosisService.GetVisitDiagnoses(uint(visitID))
	if err != nil {
		response.InternalServerError(c, "Failed to get visit diagnoses")
		return
	}

	response.Success(c, "Visit diagnoses retrieved successfully", diagnoses)
}

// GetPatientDiagnoses handles getting patient diagnosis history
func (h *DiagnosisHandler) GetPatientDiagnoses(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if diagnosisType := c.Query("diagnosis_type"); diagnosisType != "" {
		filters["diagnosis_type"] = diagnosisType
	}
	if status := c.Query("diagnosis_status"); status != "" {
		filters["diagnosis_status"] = status
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	diagnoses, err := h.diagnosisService.GetPatientDiagnoses(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient diagnoses")
		return
	}

	response.Success(c, "Patient diagnoses retrieved successfully", diagnoses)
}
