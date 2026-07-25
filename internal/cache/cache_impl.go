package cache

import (
	"context"
	"sync"
	"time"
)

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

type inMemoryCache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

func NewInMemoryCache() Cache {
	return &inMemoryCache{
		store: make(map[string]cacheEntry),
	}
}

func (c *inMemoryCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.store[key]
	if !exists {
		return nil, false, nil
	}

	if time.Now().After(entry.expiration) {
		delete(c.store, key)
		return nil, false, nil
	}

	return entry.value, true, nil
}

func (c *inMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
	return nil
}

func (c *inMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	return nil
}
