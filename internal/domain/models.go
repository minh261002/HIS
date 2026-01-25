package domain

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User represents a user in the system
type User struct {
	BaseModel
	Username     string  `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string  `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string  `gorm:"size:255;not null" json:"-"`
	FullName     string  `gorm:"size:100" json:"full_name"`
	PhoneNumber  string  `gorm:"size:20" json:"phone_number"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
	Roles        []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// Role represents a role in the system
type Role struct {
	BaseModel
	Name        string        `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string        `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string        `gorm:"size:255" json:"description"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []*User       `gorm:"many2many:user_roles;" json:"-"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// Permission represents a permission in the system
type Permission struct {
	BaseModel
	Name        string  `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string  `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Description string  `gorm:"size:255" json:"description"`
	Module      string  `gorm:"size:50" json:"module"`
	Roles       []*Role `gorm:"many2many:role_permissions;" json:"-"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}
