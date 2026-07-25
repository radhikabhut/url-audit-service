package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRateLimiter(t *testing.T) {
	ctx := context.Background()
	// Rate limit: 2 requests per 100 milliseconds
	lim := NewInMemoryRateLimiter(2, 100*time.Millisecond)

	ip := "192.168.1.100"

	// Request 1: Allowed
	allowed, err := lim.Allow(ctx, ip)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// Request 2: Allowed
	allowed, err = lim.Allow(ctx, ip)
	assert.NoError(t, err)
	assert.True(t, allowed)

	// Request 3: Blocked (exceeded burst of 2)
	allowed, err = lim.Allow(ctx, ip)
	assert.NoError(t, err)
	assert.False(t, allowed)

	// Wait for refill
	time.Sleep(110 * time.Millisecond)

	// Request 4: Allowed again
	allowed, err = lim.Allow(ctx, ip)
	assert.NoError(t, err)
	assert.True(t, allowed)
}
