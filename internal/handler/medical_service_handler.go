package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

type MedicalServiceHandler struct {
	service *service.MedicalServiceService
}

func NewMedicalServiceHandler(service *service.MedicalServiceService) *MedicalServiceHandler {
	return &MedicalServiceHandler{service: service}
}

// CreateService handles creating a new service
func (h *MedicalServiceHandler) CreateService(c *gin.Context) {
	var req dto.CreateMedicalServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	srv, err := h.service.CreateService(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrServiceCodeExists) {
			response.BadRequest(c, "Service code already exists", nil)
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "Service created successfully", srv)
}

// GetService handles getting a service
func (h *MedicalServiceHandler) GetService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", nil)
		return
	}

	srv, err := h.service.GetService(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Service retrieved successfully", srv)
}

// UpdateService handles updating a service
func (h *MedicalServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", nil)
		return
	}

	var req dto.UpdateMedicalServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	srv, err := h.service.UpdateService(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Service updated successfully", srv)
}

// ListServices handles listing services
func (h *MedicalServiceHandler) ListServices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var deptID *uint
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		id, _ := strconv.ParseUint(deptIDStr, 10, 32)
		uID := uint(id)
		deptID = &uID
	}

	services, total, err := h.service.ListServices(page, pageSize, deptID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessPaginated(c, "Services retrieved successfully", services, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}
