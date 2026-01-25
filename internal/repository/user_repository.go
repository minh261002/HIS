package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// List returns a paginated list of users
func (r *UserRepository) List(page, pageSize int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Preload("Roles").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserWithRoles finds a user by ID with roles preloaded
func (r *UserRepository) GetUserWithRoles(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UserHasPermission checks if a user has a specific permission
func (r *UserRepository) UserHasPermission(userID uint, permissionCode string) (bool, error) {
	var count int64
	err := r.db.Table("users").
		Select("COUNT(DISTINCT permissions.id)").
		Joins("INNER JOIN user_roles ON users.id = user_roles.user_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Joins("INNER JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.id = ? AND permissions.code = ? AND users.deleted_at IS NULL", userID, permissionCode).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
