package service

import (
	"errors"
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/pkg/jwt"
	"github.com/minhtran/his/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwt.Manager
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Login authenticates a user
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		User: &dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			IsActive:    user.IsActive,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken generates a new access token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	accessToken, err := s.jwtManager.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	return accessToken, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
	}, nil
}
