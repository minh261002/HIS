package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// MedicationHandler handles medication HTTP requests
type MedicationHandler struct {
	medicationService *service.MedicationService
}

// NewMedicationHandler creates a new medication handler
func NewMedicationHandler(medicationService *service.MedicationService) *MedicationHandler {
	return &MedicationHandler{medicationService: medicationService}
}

// SearchMedications handles searching medications
func (h *MedicationHandler) SearchMedications(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	medications, err := h.medicationService.SearchMedications(query, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to search medications")
		return
	}

	response.Success(c, "Medications retrieved successfully", medications)
}

// GetMedication handles getting medication by ID
func (h *MedicationHandler) GetMedication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid medication ID", nil)
		return
	}

	medication, err := h.medicationService.GetMedicationByID(uint(id))
	if err != nil {
		response.NotFound(c, "Medication not found")
		return
	}

	response.Success(c, "Medication retrieved successfully", medication)
}
