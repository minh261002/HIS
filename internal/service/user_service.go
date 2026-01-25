package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user management business logic
type UserService struct {
	userRepo *repository.UserRepository
	db       *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
	}
}

// CreateUser creates a new user (admin only)
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserDetailResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	existingEmail, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingEmail != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with roles in transaction
	var user *domain.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		user = &domain.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FullName:     req.FullName,
			PhoneNumber:  req.PhoneNumber,
			IsActive:     true,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Assign roles
		if len(req.RoleIDs) > 0 {
			var roles []*domain.Role
			if err := tx.Find(&roles, req.RoleIDs).Error; err != nil {
				return err
			}
			if len(roles) != len(req.RoleIDs) {
				return errors.New("some roles not found")
			}
			user.Roles = roles
			if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Reload user with roles
	user, err = s.userRepo.GetUserWithRoles(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	return s.toUserDetailResponse(user), nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserDetailResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Reload with roles
	user, err = s.userRepo.GetUserWithRoles(id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	return s.toUserDetailResponse(user), nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// AssignRoles assigns roles to a user
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Get roles
	var roles []*domain.Role
	if err := s.db.Find(&roles, roleIDs).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	if len(roles) != len(roleIDs) {
		return errors.New("some roles not found")
	}

	// Assign roles
	if err := s.db.Model(user).Association("Roles").Replace(roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

// GetUserByID gets user details by ID
func (s *UserService) GetUserByIDWithRoles(id uint) (*dto.UserDetailResponse, error) {
	user, err := s.userRepo.GetUserWithRoles(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.toUserDetailResponse(user), nil
}

// ListUsers returns paginated list of users
func (s *UserService) ListUsers(page, pageSize int) ([]*dto.UserListItem, int64, error) {
	users, total, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	items := make([]*dto.UserListItem, len(users))
	for i, user := range users {
		roleNames := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roleNames[j] = role.Name
		}

		items[i] = &dto.UserListItem{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			RoleNames: roleNames,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}
	}

	return items, total, nil
}

// Helper to convert domain user to DTO
func (s *UserService) toUserDetailResponse(user *domain.User) *dto.UserDetailResponse {
	roles := make([]dto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		}
	}

	return &dto.UserDetailResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		Roles:       roles,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}
