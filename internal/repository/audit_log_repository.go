package repository

import (
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// AuditLogRepository handles audit log data operations
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(log *domain.AuditLog) error {
	return r.db.Create(log).Error
}

// List returns a paginated list of audit logs with filtering
func (r *AuditLogRepository) List(page, pageSize int, userID *uint, resource string, fromDate, toDate *time.Time) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&domain.AuditLog{})

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if fromDate != nil {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", toDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
