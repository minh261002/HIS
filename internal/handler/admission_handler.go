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

// AdmissionHandler handles admission HTTP requests
type AdmissionHandler struct {
	admissionService *service.AdmissionService
}

// NewAdmissionHandler creates a new admission handler
func NewAdmissionHandler(admissionService *service.AdmissionService) *AdmissionHandler {
	return &AdmissionHandler{admissionService: admissionService}
}

// CreateAdmission handles creating an admission
func (h *AdmissionHandler) CreateAdmission(c *gin.Context) {
	var req dto.CreateAdmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	admission, err := h.admissionService.CreateAdmission(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrVisitNotFound) {
			response.NotFound(c, "Visit not found")
			return
		}
		if errors.Is(err, service.ErrBedNotFound) {
			response.NotFound(c, "Bed not found")
			return
		}
		if errors.Is(err, service.ErrBedNotAvailable) {
			response.BadRequest(c, "Bed is not available", nil)
			return
		}
		response.InternalServerError(c, "Failed to create admission")
		return
	}

	response.Created(c, "Admission created successfully", admission)
}

// GetAdmission handles getting admission details
func (h *AdmissionHandler) GetAdmission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid admission ID", nil)
		return
	}

	admission, err := h.admissionService.GetAdmissionByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrAdmissionNotFound) {
			response.NotFound(c, "Admission not found")
			return
		}
		response.InternalServerError(c, "Failed to get admission")
		return
	}

	response.Success(c, "Admission retrieved successfully", admission)
}

// GetAdmissionByCode handles getting admission by code
func (h *AdmissionHandler) GetAdmissionByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Admission code is required", nil)
		return
	}

	admission, err := h.admissionService.GetAdmissionByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrAdmissionNotFound) {
			response.NotFound(c, "Admission not found")
			return
		}
		response.InternalServerError(c, "Failed to get admission")
		return
	}

	response.Success(c, "Admission retrieved successfully", admission)
}

// DischargeAdmission handles discharging a patient
func (h *AdmissionHandler) DischargeAdmission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid admission ID", nil)
		return
	}

	var req dto.DischargeAdmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.admissionService.DischargeAdmission(uint(id), &req, userID); err != nil {
		if errors.Is(err, service.ErrAdmissionNotFound) {
			response.NotFound(c, "Admission not found")
			return
		}
		response.InternalServerError(c, "Failed to discharge patient")
		return
	}

	response.Success(c, "Patient discharged successfully", nil)
}

// TransferBed handles transferring patient to a new bed
func (h *AdmissionHandler) TransferBed(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid admission ID", nil)
		return
	}

	var req dto.TransferBedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.admissionService.TransferBed(uint(id), &req, userID); err != nil {
		if errors.Is(err, service.ErrBedNotFound) {
			response.NotFound(c, "Bed not found")
			return
		}
		if errors.Is(err, service.ErrBedNotAvailable) {
			response.BadRequest(c, "Bed is not available", nil)
			return
		}
		response.InternalServerError(c, "Failed to transfer bed")
		return
	}

	response.Success(c, "Bed transferred successfully", nil)
}

// GetActiveAdmissions handles getting all active admissions
func (h *AdmissionHandler) GetActiveAdmissions(c *gin.Context) {
	admissions, err := h.admissionService.GetActiveAdmissions()
	if err != nil {
		response.InternalServerError(c, "Failed to get active admissions")
		return
	}

	response.Success(c, "Active admissions retrieved successfully", admissions)
}

// GetPatientAdmissions handles getting patient's admission history
func (h *AdmissionHandler) GetPatientAdmissions(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	admissions, err := h.admissionService.GetPatientAdmissions(uint(patientID))
	if err != nil {
		response.InternalServerError(c, "Failed to get patient admissions")
		return
	}

	response.Success(c, "Patient admissions retrieved successfully", admissions)
}

// CreateNursingNote handles creating a nursing note
func (h *AdmissionHandler) CreateNursingNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid admission ID", nil)
		return
	}

	var req dto.CreateNursingNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.admissionService.CreateNursingNote(uint(id), &req, userID); err != nil {
		response.InternalServerError(c, "Failed to create nursing note")
		return
	}

	response.Created(c, "Nursing note created successfully", nil)
}

// GetAdmissionNursingNotes handles getting admission's nursing notes
func (h *AdmissionHandler) GetAdmissionNursingNotes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid admission ID", nil)
		return
	}

	notes, err := h.admissionService.GetAdmissionNursingNotes(uint(id))
	if err != nil {
		response.InternalServerError(c, "Failed to get nursing notes")
		return
	}

	response.Success(c, "Nursing notes retrieved successfully", notes)
}
