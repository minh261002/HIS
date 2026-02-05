package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// CreateDepartment handles creating a new department
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID") // Assuming userID is set in context by middleware

	dept, err := h.service.CreateDepartment(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrDepartmentCodeExists) {
			response.BadRequest(c, "Department code already exists", nil)
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "Department created successfully", dept)
}

// GetDepartment handles getting a department
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", nil)
		return
	}

	dept, err := h.service.GetDepartment(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrDepartmentNotFound) {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Department retrieved successfully", dept)
}

// UpdateDepartment handles updating a department
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", nil)
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	dept, err := h.service.UpdateDepartment(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrDepartmentNotFound) {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Department updated successfully", dept)
}

// DeleteDepartment handles deleting a department
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", nil)
		return
	}

	userID := c.GetUint("userID")

	if err := h.service.DeleteDepartment(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrDepartmentNotFound) {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "Department deleted successfully", nil)
}

// ListDepartments handles listing departments
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	depts, total, err := h.service.ListDepartments(page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessPaginated(c, "Departments retrieved successfully", depts, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}
