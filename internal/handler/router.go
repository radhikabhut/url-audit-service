package handler

import (
	"url-audit-service/internal/limiter"
	"url-audit-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(auditHandler *AuditHandler, healthHandler *HealthHandler, l limiter.RateLimiter) *gin.Engine {
	// Set gin mode based on environment or default to release
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Register global middlewares
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Health check route
	r.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := r.Group("/api/v1", middleware.RateLimiter(l))
	{
		v1.POST("/audit", auditHandler.Audit)
		v1.GET("/audit/history", auditHandler.GetHistory)
	}

	return r
}
