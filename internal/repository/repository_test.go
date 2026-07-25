package repository

import (
	"context"
	"testing"
	"url-audit-service/internal/errors"
	"url-audit-service/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryRepository()

	url := "https://example.com"
	res := &model.AuditResult{
		URL: url,
	}

	// 1. Get not found
	_, err := repo.GetByURL(ctx, url)
	assert.ErrorIs(t, err, errors.ErrNotFound)

	// 2. Save and Get
	err = repo.Save(ctx, res)
	assert.NoError(t, err)

	got, err := repo.GetByURL(ctx, url)
	assert.NoError(t, err)
	assert.Equal(t, res, got)
}
