package repository

import (
	"context"
	"url-audit-service/internal/model"
)

type AuditRepository interface {
	Save(ctx context.Context, result *model.AuditResult) error
	GetByURL(ctx context.Context, url string) (*model.AuditResult, error)
}
