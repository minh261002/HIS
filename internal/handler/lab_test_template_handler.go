package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// LabTestTemplateHandler handles lab test template HTTP requests
type LabTestTemplateHandler struct {
	templateService *service.LabTestTemplateService
}

// NewLabTestTemplateHandler creates a new lab test template handler
func NewLabTestTemplateHandler(templateService *service.LabTestTemplateService) *LabTestTemplateHandler {
	return &LabTestTemplateHandler{templateService: templateService}
}

// SearchTemplates handles searching lab test templates
func (h *LabTestTemplateHandler) SearchTemplates(c *gin.Context) {
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
func (h *LabTestTemplateHandler) GetTemplateByCode(c *gin.Context) {
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

// GetTemplatesByCategory handles getting templates by category
func (h *LabTestTemplateHandler) GetTemplatesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.BadRequest(c, "Category is required", nil)
		return
	}

	templates, err := h.templateService.GetTemplatesByCategory(category)
	if err != nil {
		response.InternalServerError(c, "Failed to get templates")
		return
	}

	response.Success(c, "Templates retrieved successfully", templates)
}
