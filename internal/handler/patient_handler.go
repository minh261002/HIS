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

// PatientHandler handles patient HTTP requests
type PatientHandler struct {
	patientService *service.PatientService
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(patientService *service.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// RegisterPatient handles patient registration
// @Summary Register a new patient
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePatientRequest true "Patient registration request"
// @Success 201 {object} response.Response{data=dto.PatientResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/patients [post]
func (h *PatientHandler) RegisterPatient(c *gin.Context) {
	var req dto.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	patient, err := h.patientService.RegisterPatient(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientExists) {
			response.BadRequest(c, "Patient with this national ID already exists", nil)
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to register patient")
		return
	}

	response.Created(c, "Patient registered successfully", patient)
}

// ListPatients handles listing patients with pagination and filters
// @Summary List patients
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param gender query string false "Filter by gender"
// @Param blood_type query string false "Filter by blood type"
// @Param city query string false "Filter by city"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.PaginatedResponse{data=[]dto.PatientListItem}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/patients [get]
func (h *PatientHandler) ListPatients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if gender := c.Query("gender"); gender != "" {
		filters["gender"] = gender
	}
	if bloodType := c.Query("blood_type"); bloodType != "" {
		filters["blood_type"] = bloodType
	}
	if city := c.Query("city"); city != "" {
		filters["city"] = city
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active, _ := strconv.ParseBool(isActive)
		filters["is_active"] = active
	}

	patients, total, err := h.patientService.ListPatients(page, pageSize, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to list patients")
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.SuccessPaginated(c, "Patients retrieved successfully", patients, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// SearchPatients handles searching patients
// @Summary Search patients
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query (name, phone, email, code, national ID)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]dto.PatientListItem}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/patients/search [get]
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Search query is required", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	patients, total, err := h.patientService.SearchPatients(query, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to search patients")
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.SuccessPaginated(c, "Search results retrieved successfully", patients, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// GetPatient handles getting patient details by ID
// @Summary Get patient details
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response{data=dto.PatientResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		response.InternalServerError(c, "Failed to get patient")
		return
	}

	response.Success(c, "Patient retrieved successfully", patient)
}

// GetPatientByCode handles getting patient by patient code
// @Summary Get patient by code
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Param code path string true "Patient Code"
// @Success 200 {object} response.Response{data=dto.PatientResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/patients/code/{code} [get]
func (h *PatientHandler) GetPatientByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Patient code is required", nil)
		return
	}

	patient, err := h.patientService.GetPatientByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		response.InternalServerError(c, "Failed to get patient")
		return
	}

	response.Success(c, "Patient retrieved successfully", patient)
}

// UpdatePatient handles updating patient information
// @Summary Update patient
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Param request body dto.UpdatePatientRequest true "Update patient request"
// @Success 200 {object} response.Response{data=dto.PatientResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/patients/{id} [put]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	var req dto.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	patient, err := h.patientService.UpdatePatient(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to update patient")
		return
	}

	response.Success(c, "Patient updated successfully", patient)
}

// DeletePatient handles soft deleting a patient
// @Summary Delete patient
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	err = h.patientService.DeletePatient(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		response.InternalServerError(c, "Failed to delete patient")
		return
	}

	response.Success(c, "Patient deleted successfully", nil)
}

// GetPatientStats handles getting patient statistics
// @Summary Get patient statistics
// @Tags patients
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/patients/stats [get]
func (h *PatientHandler) GetPatientStats(c *gin.Context) {
	stats, err := h.patientService.GetPatientStats()
	if err != nil {
		response.InternalServerError(c, "Failed to get patient statistics")
		return
	}

	response.Success(c, "Statistics retrieved successfully", stats)
}
