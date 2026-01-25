package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/repository"
)

// RBACMiddleware creates middleware for role-based access control
type RBACMiddleware struct {
	userRepo *repository.UserRepository
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(userRepo *repository.UserRepository) *RBACMiddleware {
	return &RBACMiddleware{
		userRepo: userRepo,
	}
}

// RequireRole checks if user has any of the specified roles
func (m *RBACMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		// Get user with roles
		user, err := m.userRepo.GetUserWithRoles(userID)
		if err != nil {
			response.InternalServerError(c, "Failed to get user roles")
			c.Abort()
			return
		}

		if user == nil {
			response.Unauthorized(c, "User not found")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, userRole := range user.Roles {
			for _, requiredRole := range roles {
				if userRole.Code == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission checks if user has a specific permission
func (m *RBACMiddleware) RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		// Check if user has permission
		hasPermission, err := m.userRepo.UserHasPermission(userID, permissionCode)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
