package dto

import "time"

// AuditLogResponse represents audit log data in response
type AuditLogResponse struct {
	ID         uint                `json:"id"`
	User       *UserDetailResponse `json:"user,omitempty"`
	Action     string              `json:"action"`
	Resource   string              `json:"resource"`
	ResourceID string              `json:"resource_id"`
	Details    interface{}         `json:"details,omitempty"`
	IPAddress  string              `json:"ip_address"`
	UserAgent  string              `json:"user_agent"`
	CreatedAt  time.Time           `json:"created_at"`
}
