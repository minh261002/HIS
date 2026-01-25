package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/config"
	"github.com/minhtran/his/internal/handler"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/jwt"
	"github.com/minhtran/his/internal/pkg/logger"
	"github.com/minhtran/his/internal/repository"
	"github.com/minhtran/his/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Hospital Information System API",
		zap.String("version", "1.0.0"),
		zap.String("mode", cfg.Server.Mode),
	)

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database, cfg.Log.Level)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	allergyRepo := repository.NewPatientAllergyRepository(db)
	historyRepo := repository.NewPatientMedicalHistoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo, db)
	patientService := service.NewPatientService(patientRepo)
	allergyService := service.NewPatientAllergyService(allergyRepo, patientRepo)
	historyService := service.NewPatientMedicalHistoryService(historyRepo, patientRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	patientHandler := handler.NewPatientHandler(patientService)
	allergyHandler := handler.NewPatientAllergyHandler(allergyService)
	historyHandler := handler.NewPatientMedicalHistoryHandler(historyService)

	// Initialize middleware
	rbacMiddleware := middleware.NewRBACMiddleware(userRepo)

	// Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Setup routes
	handler.SetupRoutes(router, authHandler, userHandler, patientHandler, allergyHandler, historyHandler, jwtManager, rbacMiddleware, cfg.Server.AllowedOrigins)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
