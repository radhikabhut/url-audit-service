package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase scheme and host",
			input:    "HTTPS://GOOGLE.COM/Path",
			expected: "https://google.com/Path",
		},
		{
			name:     "Remove default HTTP port",
			input:    "http://example.com:80/path",
			expected: "http://example.com/path",
		},
		{
			name:     "Remove default HTTPS port",
			input:    "https://example.com:443/path",
			expected: "https://example.com/path",
		},
		{
			name:     "Keep non-default port",
			input:    "http://example.com:8080/path",
			expected: "http://example.com:8080/path",
		},
		{
			name:     "Strip trailing slash",
			input:    "https://example.com/path/",
			expected: "https://example.com/path",
		},
		{
			name:     "Sort query parameters",
			input:    "https://example.com/path?b=2&a=1",
			expected: "https://example.com/path?a=1&b=2",
		},
		{
			name:     "Empty path defaults to slash",
			input:    "https://example.com",
			expected: "https://example.com/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeURL(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCacheSetGetDelete(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	key := "test-key"
	val := "test-value"

	// 1. Get empty
	v, found, err := c.Get(ctx, key)
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, v)

	// 2. Set and Get
	err = c.Set(ctx, key, val, 1*time.Second)
	assert.NoError(t, err)

	v, found, err = c.Get(ctx, key)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, val, v)

	// 3. Delete
	err = c.Delete(ctx, key)
	assert.NoError(t, err)

	v, found, err = c.Get(ctx, key)
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestCacheExpiration(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	key := "expire-key"
	val := "expire-val"

	err := c.Set(ctx, key, val, 10*time.Millisecond)
	assert.NoError(t, err)

	// Immediate get
	v, found, err := c.Get(ctx, key)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, val, v)

	// Wait for expiration
	time.Sleep(15 * time.Millisecond)

	v, found, err = c.Get(ctx, key)
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, v)
}
