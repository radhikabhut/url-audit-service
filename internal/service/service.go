package service

import (
	"context"
	"time"
	"url-audit-service/internal/cache"
	"url-audit-service/internal/client"
	"url-audit-service/internal/model"
	"url-audit-service/internal/repository"
)

type AuditService interface {
	AuditURL(ctx context.Context, url string) (*model.AuditResult, error)
	GetAuditHistory(ctx context.Context, url string) (*model.AuditResult, error)
}

type auditService struct {
	client   client.AuditClient
	repo     repository.AuditRepository
	cache    cache.Cache
	cacheTTL time.Duration
	sem      chan struct{}
}

func NewAuditService(client client.AuditClient, repo repository.AuditRepository, c cache.Cache, cacheTTL time.Duration, maxConcurrentAudits int) AuditService {
	// Fallback to 10 if invalid/zero
	if maxConcurrentAudits <= 0 {
		maxConcurrentAudits = 10
	}
	return &auditService{
		client:   client,
		repo:     repo,
		cache:    c,
		cacheTTL: cacheTTL,
		sem:      make(chan struct{}, maxConcurrentAudits),
	}
}

func (s *auditService) AuditURL(ctx context.Context, url string) (*model.AuditResult, error) {
	normURL, err := cache.NormalizeURL(url)
	if err != nil {
		// Fallback to raw url as key if normalization fails
		normURL = url
	}

	// 1. Check Cache
	if val, found, err := s.cache.Get(ctx, normURL); err == nil && found {
		if cachedResult, ok := val.(*model.AuditResult); ok {
			// Clone result to avoid race condition and mutate cached indicator
			cloned := *cachedResult
			cloned.Cached = true
			return &cloned, nil
		}
	}

	// Acquire concurrency slot before executing outbound request
	select {
	case s.sem <- struct{}{}:
		defer func() { <-s.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// 2. Perform Audit
	result, err := s.client.AuditURL(ctx, url)
	if err != nil {
		return nil, err
	}

	// 3. Save to Repository
	_ = s.repo.Save(ctx, result)

	// 4. Save to Cache
	_ = s.cache.Set(ctx, normURL, result, s.cacheTTL)

	return result, nil
}

func (s *auditService) GetAuditHistory(ctx context.Context, url string) (*model.AuditResult, error) {
	normURL, err := cache.NormalizeURL(url)
	if err != nil {
		normURL = url
	}
	return s.repo.GetByURL(ctx, normURL)
}
