package repository

import (
	"context"
	"sync"
	"url-audit-service/internal/errors"
	"url-audit-service/internal/model"
)

type inMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]*model.AuditResult
}

func NewInMemoryRepository() AuditRepository {
	return &inMemoryRepository{
		store: make(map[string]*model.AuditResult),
	}
}

func (r *inMemoryRepository) Save(ctx context.Context, result *model.AuditResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[result.URL] = result
	return nil
}

func (r *inMemoryRepository) GetByURL(ctx context.Context, url string) (*model.AuditResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result, exists := r.store[url]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return result, nil
}
