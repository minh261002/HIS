package dto

// CreateUserRequest represents admin creating a new user
type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required,min=2,max=100"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,max=20"`
	RoleIDs     []uint `json:"role_ids" binding:"required,min=1"`
}

// UpdateUserRequest represents updating user information
type UpdateUserRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	FullName    string `json:"full_name" binding:"omitempty,min=2,max=100"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,max=20"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// AssignRolesRequest represents assigning roles to a user
type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
}

// UserDetailResponse represents detailed user information with roles
type UserDetailResponse struct {
	ID          uint           `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	FullName    string         `json:"full_name"`
	PhoneNumber string         `json:"phone_number"`
	IsActive    bool           `json:"is_active"`
	Roles       []RoleResponse `json:"roles"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// RoleResponse represents role information
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// UserListItem represents user in list view
type UserListItem struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	IsActive  bool     `json:"is_active"`
	RoleNames []string `json:"role_names"`
	CreatedAt string   `json:"created_at"`
}
