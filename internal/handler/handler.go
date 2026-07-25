package handler

import (
	"errors"
	"log/slog"
	"net/http"
	appErrors "url-audit-service/internal/errors"
	"url-audit-service/internal/model"
	"url-audit-service/internal/response"
	"url-audit-service/internal/service"
	"url-audit-service/internal/validator"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService service.AuditService
}

func NewAuditHandler(s service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: s,
	}
}

func (h *AuditHandler) Audit(c *gin.Context) {
	var req model.AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request payload")
		return
	}

	if err := validator.ValidateSafetyURL(req.URL); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_URL", err.Error())
		return
	}

	result, err := h.auditService.AuditURL(c.Request.Context(), req.URL)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidURL) {
			response.Error(c, http.StatusBadRequest, "INVALID_URL", err.Error())
			return
		}
		if errors.Is(err, appErrors.ErrURLTimeout) {
			response.Error(c, http.StatusGatewayTimeout, "TIMEOUT", err.Error())
			return
		}
		if errors.Is(err, appErrors.ErrConnectionFailed) {
			response.Error(c, http.StatusBadGateway, "CONNECTION_FAILED", err.Error())
			return
		}
		slog.Error("audit service error", slog.Any("error", err))
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuditHandler) GetHistory(c *gin.Context) {
	// TODO: Implement audit history logic
	response.JSON(c, http.StatusOK, gin.H{"message": "History endpoint placeholder"})
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{
		"status": "UP",
	})
}
