package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/pkg/logger"
	"github.com/minhtran/his/internal/pkg/response"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return error response
				response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}
