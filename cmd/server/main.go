package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-audit-service/internal/cache"
	"url-audit-service/internal/client"
	"url-audit-service/internal/config"
	"url-audit-service/internal/handler"
	"url-audit-service/internal/limiter"
	"url-audit-service/internal/logger"
	"url-audit-service/internal/repository"
	"url-audit-service/internal/service"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.LogLevel, cfg.LogFormat)
	slog.Info("Starting URL Audit Service...")

	// 3. Dependency Injection (Constructor-based)
	auditClient := client.NewAuditClient(cfg.RequestTimeout)
	auditRepo := repository.NewInMemoryRepository()
	auditCache := cache.NewInMemoryCache()
	auditSvc := service.NewAuditService(auditClient, auditRepo, auditCache, cfg.CacheTTL, cfg.MaxConcurrentAudits)
	auditHandler := handler.NewAuditHandler(auditSvc)
	healthHandler := handler.NewHealthHandler()
	auditLimiter := limiter.NewInMemoryRateLimiter(cfg.RateLimit, cfg.RateLimitWindow)

	// 4. Setup Router
	r := handler.NewRouter(auditHandler, healthHandler, auditLimiter)

	// 5. Configure HTTP Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.RequestTimeout,
		WriteTimeout: cfg.RequestTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Start Server in Goroutine
	go func() {
		slog.Info("Server is running", slog.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")

	// Create context with a timeout for shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server exited successfully")
}
