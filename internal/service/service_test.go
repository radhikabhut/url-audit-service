package service

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"url-audit-service/internal/cache"
	"url-audit-service/internal/model"
	"url-audit-service/internal/repository"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	activeCount int32
	maxActive   int32
	delay       time.Duration
}

func (m *mockClient) AuditURL(ctx context.Context, targetURL string) (*model.AuditResult, error) {
	// Increment active calls count
	currentActive := atomic.AddInt32(&m.activeCount, 1)

	// Check if this exceeds the maximum active observed
	for {
		max := atomic.LoadInt32(&m.maxActive)
		if currentActive <= max {
			break
		}
		if atomic.CompareAndSwapInt32(&m.maxActive, max, currentActive) {
			break
		}
	}

	time.Sleep(m.delay)

	atomic.AddInt32(&m.activeCount, -1)
	return &model.AuditResult{URL: targetURL}, nil
}

func TestServiceConcurrencyLimit(t *testing.T) {
	// Setup dependencies
	mc := &mockClient{delay: 20 * time.Millisecond}
	repo := repository.NewInMemoryRepository()
	c := cache.NewInMemoryCache()

	// Limit to max 2 concurrent audits
	maxConcurrency := 2
	svc := NewAuditService(mc, repo, c, 5*time.Minute, maxConcurrency)

	// Launch 5 concurrent audits
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = svc.AuditURL(context.Background(), "https://url-audit-concurrency-test.com")
		}(i)
	}

	wg.Wait()

	// Assert that maximum concurrent active requests observed by client never exceeded 2
	assert.LessOrEqual(t, mc.maxActive, int32(maxConcurrency))
	assert.Greater(t, mc.maxActive, int32(0))
}

func TestAuditService_CacheAndHistory(t *testing.T) {
	ctx := context.Background()
	mc := &mockClient{delay: 0}
	repo := repository.NewInMemoryRepository()
	c := cache.NewInMemoryCache()

	svc := NewAuditService(mc, repo, c, 5*time.Minute, 5)

	targetURL := "https://example.com/cache-test"

	// 1. First audit (Cache miss)
	res1, err := svc.AuditURL(ctx, targetURL)
	assert.NoError(t, err)
	assert.NotNil(t, res1)
	assert.False(t, res1.Cached)

	// 2. Second audit (Cache hit)
	res2, err := svc.AuditURL(ctx, targetURL)
	assert.NoError(t, err)
	assert.NotNil(t, res2)
	assert.True(t, res2.Cached) // served from cache!

	// 3. Get history
	historyRes, err := svc.GetAuditHistory(ctx, targetURL)
	assert.NoError(t, err)
	assert.NotNil(t, historyRes)
	assert.Equal(t, targetURL, historyRes.URL)
}
