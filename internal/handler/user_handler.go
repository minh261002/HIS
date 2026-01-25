package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// UserHandler handles user management HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles creating a new user (admin only)
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 201 {object} response.Response{data=dto.UserDetailResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.BadRequest(c, "User already exists", nil)
			return
		}
		response.InternalServerError(c, "Failed to create user")
		return
	}

	response.Created(c, "User created successfully", user)
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]dto.UserListItem}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list users")
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.SuccessPaginated(c, "Users retrieved successfully", users, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// GetUser handles getting user details by ID
// @Summary Get user details
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=dto.UserDetailResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetUserByIDWithRoles(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get user")
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// UpdateUser handles updating user information
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.Response{data=dto.UserDetailResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to update user")
		return
	}

	response.Success(c, "User updated successfully", user)
}

// DeleteUser handles soft deleting a user
// @Summary Delete user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// AssignRoles handles assigning roles to a user
// @Summary Assign roles to user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body dto.AssignRolesRequest true "Assign roles request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id}/roles [post]
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req dto.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	err = h.userService.AssignRoles(uint(id), req.RoleIDs)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to assign roles")
		return
	}

	response.Success(c, "Roles assigned successfully", nil)
}
