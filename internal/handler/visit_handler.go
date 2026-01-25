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

// VisitHandler handles visit HTTP requests
type VisitHandler struct {
	visitService *service.VisitService
}

// NewVisitHandler creates a new visit handler
func NewVisitHandler(visitService *service.VisitService) *VisitHandler {
	return &VisitHandler{visitService: visitService}
}

// CreateVisit handles creating a new visit
func (h *VisitHandler) CreateVisit(c *gin.Context) {
	var req dto.CreateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	visit, err := h.visitService.CreateVisit(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		response.InternalServerError(c, "Failed to create visit")
		return
	}

	response.Created(c, "Visit created successfully", visit)
}

// ListVisits handles listing visits with filters
func (h *VisitHandler) ListVisits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if patientID := c.Query("patient_id"); patientID != "" {
		filters["patient_id"] = patientID
	}
	if doctorID := c.Query("doctor_id"); doctorID != "" {
		filters["doctor_id"] = doctorID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if visitType := c.Query("visit_type"); visitType != "" {
		filters["visit_type"] = visitType
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	visits, total, err := h.visitService.SearchVisits(filters, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list visits")
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.SuccessPaginated(c, "Visits retrieved successfully", visits, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// GetVisit handles getting visit details
func (h *VisitHandler) GetVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	visit, err := h.visitService.GetVisitByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		response.InternalServerError(c, "Failed to get visit")
		return
	}

	response.Success(c, "Visit retrieved successfully", visit)
}

// GetVisitByCode handles getting visit by code
func (h *VisitHandler) GetVisitByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Visit code is required", nil)
		return
	}

	visit, err := h.visitService.GetVisitByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		response.InternalServerError(c, "Failed to get visit")
		return
	}

	response.Success(c, "Visit retrieved successfully", visit)
}

// UpdateVisit handles updating visit details
func (h *VisitHandler) UpdateVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	var req dto.UpdateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	visit, err := h.visitService.UpdateVisit(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		response.InternalServerError(c, "Failed to update visit")
		return
	}

	response.Success(c, "Visit updated successfully", visit)
}

// CompleteVisit handles completing a visit
func (h *VisitHandler) CompleteVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.visitService.CompleteVisit(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		response.InternalServerError(c, "Failed to complete visit")
		return
	}

	response.Success(c, "Visit completed successfully", nil)
}

// CancelVisit handles cancelling a visit
func (h *VisitHandler) CancelVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.visitService.CancelVisit(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		response.InternalServerError(c, "Failed to cancel visit")
		return
	}

	response.Success(c, "Visit cancelled successfully", nil)
}

// GetPatientVisits handles getting patient's visit history
func (h *VisitHandler) GetPatientVisits(c *gin.Context) {
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

	visits, err := h.visitService.GetPatientVisits(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient visits")
		return
	}

	response.Success(c, "Patient visits retrieved successfully", visits)
}

// GetDoctorVisits handles getting doctor's visits
func (h *VisitHandler) GetDoctorVisits(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID", nil)
		return
	}

	date := c.Query("date")
	if date == "" {
		response.BadRequest(c, "Date parameter is required (YYYY-MM-DD)", nil)
		return
	}

	visits, err := h.visitService.GetDoctorVisits(uint(doctorID), date)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to get doctor visits")
		return
	}

	response.Success(c, "Doctor visits retrieved successfully", visits)
}
