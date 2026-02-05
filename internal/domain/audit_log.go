package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AuditAction represents an action type in audit logs
type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionView   AuditAction = "VIEW"
	AuditActionLogin  AuditAction = "LOGIN"
	AuditActionLogout AuditAction = "LOGOUT"
)

// AuditDetails represents JSON details for audit logs
type AuditDetails map[string]interface{}

// Scan implements the sql.Scanner interface
func (a *AuditDetails) Scan(value interface{}) error {
	if value == nil {
		*a = AuditDetails{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the driver.Valuer interface
func (a AuditDetails) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// AuditLog represents a system audit log entry
type AuditLog struct {
	ID         uint         `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time    `json:"created_at"`
	UserID     *uint        `gorm:"index" json:"user_id,omitempty"`
	User       *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     AuditAction  `gorm:"size:50;not null;index" json:"action"`
	Resource   string       `gorm:"size:100;not null;index" json:"resource"`
	ResourceID string       `gorm:"size:100;index" json:"resource_id,omitempty"`
	Details    AuditDetails `gorm:"type:json" json:"details,omitempty"`
	IPAddress  string       `gorm:"size:50" json:"ip_address"`
	UserAgent  string       `gorm:"size:255" json:"user_agent"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
