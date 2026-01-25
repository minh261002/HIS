package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// ErrorData represents error details
type ErrorData struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessPaginated sends a successful paginated response
func SuccessPaginated(c *gin.Context, message string, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message string, details map[string]interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a bad request error
func BadRequest(c *gin.Context, message string, details map[string]interface{}) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// Unauthorized sends an unauthorized error
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden sends a forbidden error
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFound sends a not found error
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// InternalServerError sends an internal server error
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}

// ValidationError sends a validation error
func ValidationError(c *gin.Context, details map[string]interface{}) {
	Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details)
}
