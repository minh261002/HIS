package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// ImagingRequestHandler handles imaging request HTTP requests
type ImagingRequestHandler struct {
	requestService *service.ImagingRequestService
}

// NewImagingRequestHandler creates a new imaging request handler
func NewImagingRequestHandler(requestService *service.ImagingRequestService) *ImagingRequestHandler {
	return &ImagingRequestHandler{requestService: requestService}
}

// CreateImagingRequest handles creating an imaging request
func (h *ImagingRequestHandler) CreateImagingRequest(c *gin.Context) {
	var req dto.CreateImagingRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	imaging, err := h.requestService.CreateImagingRequest(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		if errors.Is(err, service.ErrImagingTemplateNotFound) {
			response.NotFound(c, "Imaging template not found")
			return
		}
		response.InternalServerError(c, "Failed to create imaging request")
		return
	}

	response.Created(c, "Imaging request created successfully", imaging)
}

// GetImagingRequest handles getting imaging request details
func (h *ImagingRequestHandler) GetImagingRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	imaging, err := h.requestService.GetImagingRequestByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to get imaging request")
		return
	}

	response.Success(c, "Imaging request retrieved successfully", imaging)
}

// GetImagingRequestByCode handles getting imaging request by code
func (h *ImagingRequestHandler) GetImagingRequestByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Request code is required", nil)
		return
	}

	imaging, err := h.requestService.GetImagingRequestByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to get imaging request")
		return
	}

	response.Success(c, "Imaging request retrieved successfully", imaging)
}

// ScheduleImaging handles scheduling imaging
func (h *ImagingRequestHandler) ScheduleImaging(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	var req dto.ScheduleImagingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	scheduledDate, err := time.Parse(time.RFC3339, req.ScheduledDate)
	if err != nil {
		response.BadRequest(c, "Invalid scheduled date format", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.ScheduleImaging(uint(id), scheduledDate, userID); err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to schedule imaging")
		return
	}

	response.Success(c, "Imaging scheduled successfully", nil)
}

// StartImaging handles starting imaging
func (h *ImagingRequestHandler) StartImaging(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.StartImaging(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to start imaging")
		return
	}

	response.Success(c, "Imaging started successfully", nil)
}

// CompleteImaging handles completing imaging
func (h *ImagingRequestHandler) CompleteImaging(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CompleteImaging(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to complete imaging")
		return
	}

	response.Success(c, "Imaging completed successfully", nil)
}

// CancelImaging handles cancelling imaging
func (h *ImagingRequestHandler) CancelImaging(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CancelImaging(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to cancel imaging")
		return
	}

	response.Success(c, "Imaging cancelled successfully", nil)
}

// CreateOrUpdateResult handles creating/updating imaging result
func (h *ImagingRequestHandler) CreateOrUpdateResult(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid imaging request ID", nil)
		return
	}

	var req dto.CreateImagingResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CreateOrUpdateResult(uint(id), &req, userID); err != nil {
		if errors.Is(err, service.ErrImagingRequestNotFound) {
			response.NotFound(c, "Imaging request not found")
			return
		}
		response.InternalServerError(c, "Failed to save result")
		return
	}

	response.Success(c, "Result saved successfully", nil)
}

// GetVisitImagingRequests handles getting visit's imaging requests
func (h *ImagingRequestHandler) GetVisitImagingRequests(c *gin.Context) {
	visitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	imagingRequests, err := h.requestService.GetVisitImagingRequests(uint(visitID))
	if err != nil {
		response.InternalServerError(c, "Failed to get visit imaging requests")
		return
	}

	response.Success(c, "Visit imaging requests retrieved successfully", imagingRequests)
}

// GetPatientImagingRequests handles getting patient's imaging request history
func (h *ImagingRequestHandler) GetPatientImagingRequests(c *gin.Context) {
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

	imagingRequests, err := h.requestService.GetPatientImagingRequests(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient imaging requests")
		return
	}

	response.Success(c, "Patient imaging requests retrieved successfully", imagingRequests)
}
