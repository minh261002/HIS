package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/jwt"
	"github.com/minhtran/his/internal/pkg/response"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, authHandler *AuthHandler, jwtManager *jwt.Manager, allowedOrigins []string) {
	// Apply global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware(allowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Server is healthy", map[string]string{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth protected routes
			protected.GET("/auth/profile", authHandler.GetProfile)

			// Add more protected routes here
		}
	}
}
