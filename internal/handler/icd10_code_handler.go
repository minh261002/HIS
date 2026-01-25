package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// ICD10CodeHandler handles ICD-10 code HTTP requests
type ICD10CodeHandler struct {
	icd10Service *service.ICD10CodeService
}

// NewICD10CodeHandler creates a new ICD-10 code handler
func NewICD10CodeHandler(icd10Service *service.ICD10CodeService) *ICD10CodeHandler {
	return &ICD10CodeHandler{icd10Service: icd10Service}
}

// SearchICD10Codes handles searching ICD-10 codes
func (h *ICD10CodeHandler) SearchICD10Codes(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	codes, err := h.icd10Service.SearchICD10Codes(query, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to search ICD-10 codes")
		return
	}

	response.Success(c, "ICD-10 codes retrieved successfully", codes)
}

// GetICD10CodeByCode handles getting ICD-10 code by code
func (h *ICD10CodeHandler) GetICD10CodeByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "ICD-10 code is required", nil)
		return
	}

	icd10Code, err := h.icd10Service.GetICD10CodeByCode(code)
	if err != nil {
		response.NotFound(c, "ICD-10 code not found")
		return
	}

	response.Success(c, "ICD-10 code retrieved successfully", icd10Code)
}

// GetICD10CodesByCategory handles getting ICD-10 codes by category
func (h *ICD10CodeHandler) GetICD10CodesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.BadRequest(c, "Category is required", nil)
		return
	}

	codes, err := h.icd10Service.GetICD10CodesByCategory(category)
	if err != nil {
		response.InternalServerError(c, "Failed to get ICD-10 codes")
		return
	}

	response.Success(c, "ICD-10 codes retrieved successfully", codes)
}
