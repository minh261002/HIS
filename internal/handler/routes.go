package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/jwt"
	"github.com/minhtran/his/internal/pkg/response"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	r *gin.Engine,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	patientHandler *PatientHandler,
	allergyHandler *PatientAllergyHandler,
	historyHandler *PatientMedicalHistoryHandler,
	jwtManager *jwt.Manager,
	rbacMiddleware *middleware.RBACMiddleware,
	allowedOrigins []string,
) {
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
		// Public auth routes (only login and refresh)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth protected routes
			protected.GET("/auth/profile", authHandler.GetProfile)

			// User management routes (admin only)
			users := protected.Group("/users")
			users.Use(rbacMiddleware.RequirePermission("users.manage"))
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.POST("/:id/roles", userHandler.AssignRoles)
			}

			// Patient management routes
			patients := protected.Group("/patients")
			{
				// Statistics (requires view permission)
				patients.GET("/stats", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatientStats)

				// Search (requires view permission)
				patients.GET("/search", rbacMiddleware.RequirePermission("patients.view"), patientHandler.SearchPatients)

				// Get by code (requires view permission)
				patients.GET("/code/:code", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatientByCode)

				// Register patient (requires create permission)
				patients.POST("", rbacMiddleware.RequirePermission("patients.create"), patientHandler.RegisterPatient)

				// List patients (requires view permission)
				patients.GET("", rbacMiddleware.RequirePermission("patients.view"), patientHandler.ListPatients)

				// Get patient details (requires view permission)
				patients.GET("/:id", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatient)

				// Update patient (requires update permission)
				patients.PUT("/:id", rbacMiddleware.RequirePermission("patients.update"), patientHandler.UpdatePatient)

				// Delete patient (requires delete permission - admin only)
				patients.DELETE("/:id", rbacMiddleware.RequirePermission("patients.delete"), patientHandler.DeletePatient)

				// Patient allergies sub-routes
				patients.POST("/:id/allergies", rbacMiddleware.RequirePermission("patients.create"), allergyHandler.AddAllergy)
				patients.GET("/:id/allergies", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetPatientAllergies)
				patients.GET("/:id/allergies/active", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetActiveAllergies)

				// Patient medical history sub-routes
				patients.POST("/:id/medical-history", rbacMiddleware.RequirePermission("patients.create"), historyHandler.AddMedicalHistory)
				patients.GET("/:id/medical-history", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetPatientHistory)
				patients.GET("/:id/medical-history/active", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetActiveConditions)
			}

			// Allergy routes (standalone)
			allergies := protected.Group("/allergies")
			{
				allergies.GET("/:allergyId", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetAllergy)
				allergies.PUT("/:allergyId", rbacMiddleware.RequirePermission("patients.update"), allergyHandler.UpdateAllergy)
				allergies.DELETE("/:allergyId", rbacMiddleware.RequirePermission("patients.delete"), allergyHandler.DeleteAllergy)
			}

			// Medical history routes (standalone)
			medicalHistory := protected.Group("/medical-history")
			{
				medicalHistory.GET("/:historyId", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetMedicalHistory)
				medicalHistory.PUT("/:historyId", rbacMiddleware.RequirePermission("patients.update"), historyHandler.UpdateMedicalHistory)
				medicalHistory.DELETE("/:historyId", rbacMiddleware.RequirePermission("patients.delete"), historyHandler.DeleteMedicalHistory)
			}
		}
	}
}
