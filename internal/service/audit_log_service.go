package service

import (
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// AuditLogService handles business logic for audit logs
type AuditLogService struct {
	repo *repository.AuditLogRepository
}

// NewAuditLogService creates a new audit log service
func NewAuditLogService(repo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

// CreateLog creates a new audit log entry
func (s *AuditLogService) CreateLog(action domain.AuditAction, resource string, resourceID string, userID *uint, details domain.AuditDetails, ip string, userAgent string) error {
	log := &domain.AuditLog{
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		Details:    details,
		IPAddress:  ip,
		UserAgent:  userAgent,
	}
	return s.repo.Create(log)
}

// ListLogs returns a list of audit logs as DTOs
func (s *AuditLogService) ListLogs(page, pageSize int, userID *uint, resource string, fromDate, toDate *time.Time) ([]*dto.AuditLogResponse, int64, error) {
	logs, total, err := s.repo.List(page, pageSize, userID, resource, fromDate, toDate)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.AuditLogResponse, len(logs))
	for i, log := range logs {
		var userResp *dto.UserDetailResponse
		if log.User != nil {
			userResp = &dto.UserDetailResponse{
				ID:          log.User.ID,
				Username:    log.User.Username,
				Email:       log.User.Email,
				FullName:    log.User.FullName,
				PhoneNumber: log.User.PhoneNumber,
				IsActive:    log.User.IsActive,
				// Roles and timestamps are omitted for audit listing to avoid extra joins
			}
		}

		responses[i] = &dto.AuditLogResponse{
			ID:         log.ID,
			User:       userResp,
			Action:     string(log.Action),
			Resource:   log.Resource,
			ResourceID: log.ResourceID,
			Details:    log.Details,
			IPAddress:  log.IPAddress,
			UserAgent:  log.UserAgent,
			CreatedAt:  log.CreatedAt,
		}
	}

	return responses, total, nil
}
