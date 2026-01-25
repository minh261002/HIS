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

// PatientAllergyHandler handles patient allergy HTTP requests
type PatientAllergyHandler struct {
	allergyService *service.PatientAllergyService
}

// NewPatientAllergyHandler creates a new patient allergy handler
func NewPatientAllergyHandler(allergyService *service.PatientAllergyService) *PatientAllergyHandler {
	return &PatientAllergyHandler{
		allergyService: allergyService,
	}
}

// AddAllergy handles adding an allergy to a patient
// @Summary Add patient allergy
// @Tags patient-allergies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Param request body dto.CreateAllergyRequest true "Allergy request"
// @Success 201 {object} response.Response{data=dto.AllergyResponse}
// @Router /api/v1/patients/{id}/allergies [post]
func (h *PatientAllergyHandler) AddAllergy(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	var req dto.CreateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	allergy, err := h.allergyService.AddAllergy(uint(patientID), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to add allergy")
		return
	}

	response.Created(c, "Allergy added successfully", allergy)
}

// GetPatientAllergies handles getting all allergies for a patient
// @Summary Get patient allergies
// @Tags patient-allergies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response{data=[]dto.AllergyListItem}
// @Router /api/v1/patients/{id}/allergies [get]
func (h *PatientAllergyHandler) GetPatientAllergies(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	allergies, err := h.allergyService.GetPatientAllergies(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get allergies")
		return
	}

	response.Success(c, "Allergies retrieved successfully", allergies)
}

// GetActiveAllergies handles getting active allergies for a patient
// @Summary Get active patient allergies
// @Tags patient-allergies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response{data=[]dto.AllergyListItem}
// @Router /api/v1/patients/{id}/allergies/active [get]
func (h *PatientAllergyHandler) GetActiveAllergies(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	allergies, err := h.allergyService.GetActiveAllergies(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get active allergies")
		return
	}

	response.Success(c, "Active allergies retrieved successfully", allergies)
}

// GetAllergy handles getting allergy details
// @Summary Get allergy details
// @Tags patient-allergies
// @Produce json
// @Security BearerAuth
// @Param allergyId path int true "Allergy ID"
// @Success 200 {object} response.Response{data=dto.AllergyResponse}
// @Router /api/v1/allergies/{allergyId} [get]
func (h *PatientAllergyHandler) GetAllergy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("allergyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID", nil)
		return
	}

	allergy, err := h.allergyService.GetAllergyByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrAllergyNotFound) {
			response.NotFound(c, "Allergy not found")
			return
		}
		response.InternalServerError(c, "Failed to get allergy")
		return
	}

	response.Success(c, "Allergy retrieved successfully", allergy)
}

// UpdateAllergy handles updating an allergy
// @Summary Update allergy
// @Tags patient-allergies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param allergyId path int true "Allergy ID"
// @Param request body dto.UpdateAllergyRequest true "Update allergy request"
// @Success 200 {object} response.Response{data=dto.AllergyResponse}
// @Router /api/v1/allergies/{allergyId} [put]
func (h *PatientAllergyHandler) UpdateAllergy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("allergyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID", nil)
		return
	}

	var req dto.UpdateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	allergy, err := h.allergyService.UpdateAllergy(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrAllergyNotFound) {
			response.NotFound(c, "Allergy not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to update allergy")
		return
	}

	response.Success(c, "Allergy updated successfully", allergy)
}

// DeleteAllergy handles deleting an allergy
// @Summary Delete allergy
// @Tags patient-allergies
// @Produce json
// @Security BearerAuth
// @Param allergyId path int true "Allergy ID"
// @Success 200 {object} response.Response
// @Router /api/v1/allergies/{allergyId} [delete]
func (h *PatientAllergyHandler) DeleteAllergy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("allergyId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID", nil)
		return
	}

	err = h.allergyService.DeleteAllergy(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrAllergyNotFound) {
			response.NotFound(c, "Allergy not found")
			return
		}
		response.InternalServerError(c, "Failed to delete allergy")
		return
	}

	response.Success(c, "Allergy deleted successfully", nil)
}
