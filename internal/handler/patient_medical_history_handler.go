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

// PatientMedicalHistoryHandler handles patient medical history HTTP requests
type PatientMedicalHistoryHandler struct {
	historyService *service.PatientMedicalHistoryService
}

// NewPatientMedicalHistoryHandler creates a new patient medical history handler
func NewPatientMedicalHistoryHandler(historyService *service.PatientMedicalHistoryService) *PatientMedicalHistoryHandler {
	return &PatientMedicalHistoryHandler{
		historyService: historyService,
	}
}

// AddMedicalHistory handles adding medical history to a patient
// @Summary Add patient medical history
// @Tags patient-medical-history
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Param request body dto.CreateMedicalHistoryRequest true "Medical history request"
// @Success 201 {object} response.Response{data=dto.MedicalHistoryResponse}
// @Router /api/v1/patients/{id}/medical-history [post]
func (h *PatientMedicalHistoryHandler) AddMedicalHistory(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	var req dto.CreateMedicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	history, err := h.historyService.AddMedicalHistory(uint(patientID), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to add medical history")
		return
	}

	response.Created(c, "Medical history added successfully", history)
}

// GetPatientHistory handles getting all medical history for a patient
// @Summary Get patient medical history
// @Tags patient-medical-history
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response{data=[]dto.MedicalHistoryListItem}
// @Router /api/v1/patients/{id}/medical-history [get]
func (h *PatientMedicalHistoryHandler) GetPatientHistory(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	histories, err := h.historyService.GetPatientHistory(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get medical history")
		return
	}

	response.Success(c, "Medical history retrieved successfully", histories)
}

// GetActiveConditions handles getting active medical conditions for a patient
// @Summary Get active patient conditions
// @Tags patient-medical-history
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response{data=[]dto.MedicalHistoryListItem}
// @Router /api/v1/patients/{id}/medical-history/active [get]
func (h *PatientMedicalHistoryHandler) GetActiveConditions(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	histories, err := h.historyService.GetActiveConditions(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get active conditions")
		return
	}

	response.Success(c, "Active conditions retrieved successfully", histories)
}

// GetMedicalHistory handles getting medical history details
// @Summary Get medical history details
// @Tags patient-medical-history
// @Produce json
// @Security BearerAuth
// @Param historyId path int true "Medical History ID"
// @Success 200 {object} response.Response{data=dto.MedicalHistoryResponse}
// @Router /api/v1/medical-history/{historyId} [get]
func (h *PatientMedicalHistoryHandler) GetMedicalHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("historyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid medical history ID", nil)
		return
	}

	history, err := h.historyService.GetMedicalHistoryByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMedicalHistoryNotFound) {
			response.NotFound(c, "Medical history not found")
			return
		}
		response.InternalServerError(c, "Failed to get medical history")
		return
	}

	response.Success(c, "Medical history retrieved successfully", history)
}

// UpdateMedicalHistory handles updating medical history
// @Summary Update medical history
// @Tags patient-medical-history
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param historyId path int true "Medical History ID"
// @Param request body dto.UpdateMedicalHistoryRequest true "Update medical history request"
// @Success 200 {object} response.Response{data=dto.MedicalHistoryResponse}
// @Router /api/v1/medical-history/{historyId} [put]
func (h *PatientMedicalHistoryHandler) UpdateMedicalHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("historyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid medical history ID", nil)
		return
	}

	var req dto.UpdateMedicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	history, err := h.historyService.UpdateMedicalHistory(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrMedicalHistoryNotFound) {
			response.NotFound(c, "Medical history not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to update medical history")
		return
	}

	response.Success(c, "Medical history updated successfully", history)
}

// DeleteMedicalHistory handles deleting medical history
// @Summary Delete medical history
// @Tags patient-medical-history
// @Produce json
// @Security BearerAuth
// @Param historyId path int true "Medical History ID"
// @Success 200 {object} response.Response
// @Router /api/v1/medical-history/{historyId} [delete]
func (h *PatientMedicalHistoryHandler) DeleteMedicalHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("historyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid medical history ID", nil)
		return
	}

	err = h.historyService.DeleteMedicalHistory(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMedicalHistoryNotFound) {
			response.NotFound(c, "Medical history not found")
			return
		}
		response.InternalServerError(c, "Failed to delete medical history")
		return
	}

	response.Success(c, "Medical history deleted successfully", nil)
}
