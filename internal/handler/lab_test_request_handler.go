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

// LabTestRequestHandler handles lab test request HTTP requests
type LabTestRequestHandler struct {
	requestService *service.LabTestRequestService
}

// NewLabTestRequestHandler creates a new lab test request handler
func NewLabTestRequestHandler(requestService *service.LabTestRequestService) *LabTestRequestHandler {
	return &LabTestRequestHandler{requestService: requestService}
}

// CreateLabTestRequest handles creating a lab test request
func (h *LabTestRequestHandler) CreateLabTestRequest(c *gin.Context) {
	var req dto.CreateLabTestRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	labTest, err := h.requestService.CreateLabTestRequest(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		if errors.Is(err, service.ErrLabTestTemplateNotFound) {
			response.NotFound(c, "Lab test template not found")
			return
		}
		response.InternalServerError(c, "Failed to create lab test request")
		return
	}

	response.Created(c, "Lab test request created successfully", labTest)
}

// GetLabTestRequest handles getting lab test request details
func (h *LabTestRequestHandler) GetLabTestRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	labTest, err := h.requestService.GetLabTestRequestByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to get lab test request")
		return
	}

	response.Success(c, "Lab test request retrieved successfully", labTest)
}

// GetLabTestRequestByCode handles getting lab test request by code
func (h *LabTestRequestHandler) GetLabTestRequestByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Request code is required", nil)
		return
	}

	labTest, err := h.requestService.GetLabTestRequestByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to get lab test request")
		return
	}

	response.Success(c, "Lab test request retrieved successfully", labTest)
}

// CollectSample handles marking sample as collected
func (h *LabTestRequestHandler) CollectSample(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CollectSample(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to collect sample")
		return
	}

	response.Success(c, "Sample collected successfully", nil)
}

// StartProcessing handles starting test processing
func (h *LabTestRequestHandler) StartProcessing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.StartProcessing(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to start processing")
		return
	}

	response.Success(c, "Processing started successfully", nil)
}

// CompleteTest handles completing test
func (h *LabTestRequestHandler) CompleteTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CompleteTest(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to complete test")
		return
	}

	response.Success(c, "Test completed successfully", nil)
}

// CancelTest handles cancelling test
func (h *LabTestRequestHandler) CancelTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.requestService.CancelTest(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to cancel test")
		return
	}

	response.Success(c, "Test cancelled successfully", nil)
}

// EnterResults handles entering test results
func (h *LabTestRequestHandler) EnterResults(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid lab test request ID", nil)
		return
	}

	var req dto.EnterLabTestResultsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.requestService.EnterResults(uint(id), &req); err != nil {
		if errors.Is(err, service.ErrLabTestRequestNotFound) {
			response.NotFound(c, "Lab test request not found")
			return
		}
		response.InternalServerError(c, "Failed to enter results")
		return
	}

	response.Success(c, "Results entered successfully", nil)
}

// GetVisitLabTests handles getting visit's lab tests
func (h *LabTestRequestHandler) GetVisitLabTests(c *gin.Context) {
	visitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid visit ID", nil)
		return
	}

	labTests, err := h.requestService.GetVisitLabTests(uint(visitID))
	if err != nil {
		response.InternalServerError(c, "Failed to get visit lab tests")
		return
	}

	response.Success(c, "Visit lab tests retrieved successfully", labTests)
}

// GetPatientLabTests handles getting patient's lab test history
func (h *LabTestRequestHandler) GetPatientLabTests(c *gin.Context) {
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

	labTests, err := h.requestService.GetPatientLabTests(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient lab tests")
		return
	}

	response.Success(c, "Patient lab tests retrieved successfully", labTests)
}
