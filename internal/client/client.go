package client

import (
	"context"
	"url-audit-service/internal/model"
)

type AuditClient interface {
	AuditURL(ctx context.Context, targetURL string) (*model.AuditResult, error)
}
