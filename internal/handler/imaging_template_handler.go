package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// ImagingTemplateHandler handles imaging template HTTP requests
type ImagingTemplateHandler struct {
	templateService *service.ImagingTemplateService
}

// NewImagingTemplateHandler creates a new imaging template handler
func NewImagingTemplateHandler(templateService *service.ImagingTemplateService) *ImagingTemplateHandler {
	return &ImagingTemplateHandler{templateService: templateService}
}

// SearchTemplates handles searching imaging templates
func (h *ImagingTemplateHandler) SearchTemplates(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	templates, err := h.templateService.SearchTemplates(query, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to search templates")
		return
	}

	response.Success(c, "Templates retrieved successfully", templates)
}

// GetTemplateByCode handles getting template by code
func (h *ImagingTemplateHandler) GetTemplateByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Template code is required", nil)
		return
	}

	template, err := h.templateService.GetTemplateByCode(code)
	if err != nil {
		response.NotFound(c, "Template not found")
		return
	}

	response.Success(c, "Template retrieved successfully", template)
}

// GetTemplatesByModality handles getting templates by modality
func (h *ImagingTemplateHandler) GetTemplatesByModality(c *gin.Context) {
	modality := c.Param("modality")
	if modality == "" {
		response.BadRequest(c, "Modality is required", nil)
		return
	}

	templates, err := h.templateService.GetTemplatesByModality(modality)
	if err != nil {
		response.InternalServerError(c, "Failed to get templates")
		return
	}

	response.Success(c, "Templates retrieved successfully", templates)
}
