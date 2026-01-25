package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=dto.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	authResp, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid username or password")
			return
		}
		if errors.Is(err, service.ErrUserInactive) {
			response.Forbidden(c, "User account is inactive")
			return
		}
		response.InternalServerError(c, "Failed to login")
		return
	}

	response.Success(c, "Login successful", authResp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	response.Success(c, "Token refreshed successfully", map[string]string{
		"access_token": accessToken,
	})
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get user profile")
		return
	}

	response.Success(c, "Profile retrieved successfully", user)
}
