package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

type AuditLogHandler struct {
	service *service.AuditLogService
}

func NewAuditLogHandler(service *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

// ListLogs handles listing audit logs
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, _ := strconv.ParseUint(userIDStr, 10, 32)
		uID := uint(id)
		userID = &uID
	}
	resource := c.Query("resource")

	var fromDate, toDate *time.Time
	if fromStr := c.Query("from_date"); fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			fromDate = &t
		} else {
			response.BadRequest(c, "Invalid from_date, expected format YYYY-MM-DD", nil)
			return
		}
	}
	if toStr := c.Query("to_date"); toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			// set to end of day
			endOfDay := t.Add(24*time.Hour - time.Nanosecond)
			toDate = &endOfDay
		} else {
			response.BadRequest(c, "Invalid to_date, expected format YYYY-MM-DD", nil)
			return
		}
	}

	logs, total, err := h.service.ListLogs(page, pageSize, userID, resource, fromDate, toDate)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessPaginated(c, "Audit logs retrieved successfully", logs, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}
