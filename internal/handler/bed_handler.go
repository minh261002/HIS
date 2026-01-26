package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// BedHandler handles bed HTTP requests
type BedHandler struct {
	bedService *service.BedService
}

// NewBedHandler creates a new bed handler
func NewBedHandler(bedService *service.BedService) *BedHandler {
	return &BedHandler{bedService: bedService}
}

// GetAvailableBeds handles getting available beds
func (h *BedHandler) GetAvailableBeds(c *gin.Context) {
	department := c.Query("department")
	bedType := c.Query("bed_type")

	beds, err := h.bedService.GetAvailableBeds(department, bedType)
	if err != nil {
		response.InternalServerError(c, "Failed to get available beds")
		return
	}

	response.Success(c, "Available beds retrieved successfully", beds)
}

// GetBedByNumber handles getting bed by bed number
func (h *BedHandler) GetBedByNumber(c *gin.Context) {
	bedNumber := c.Param("number")
	if bedNumber == "" {
		response.BadRequest(c, "Bed number is required", nil)
		return
	}

	bed, err := h.bedService.GetBedByNumber(bedNumber)
	if err != nil {
		response.NotFound(c, "Bed not found")
		return
	}

	response.Success(c, "Bed retrieved successfully", bed)
}
