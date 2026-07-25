package limiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type inMemoryRateLimiter struct {
	mu       sync.Mutex
	ips      map[string]*ipLimiterEntry
	rate     rate.Limit
	burst    int
	cleanupD time.Duration
}

func NewInMemoryRateLimiter(rateLimit float64, window time.Duration) RateLimiter {
	// Calculate rate limit as events per second
	eventsPerSec := rateLimit / window.Seconds()

	lim := &inMemoryRateLimiter{
		ips:      make(map[string]*ipLimiterEntry),
		rate:     rate.Limit(eventsPerSec),
		burst:    int(rateLimit),
		cleanupD: 10 * time.Minute,
	}

	if lim.burst < 1 {
		lim.burst = 1
	}

	// Periodically clean up old IPs
	go lim.startCleanupLoop()

	return lim
}

func (l *inMemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, exists := l.ips[key]
	if !exists {
		entry = &ipLimiterEntry{
			limiter: rate.NewLimiter(l.rate, l.burst),
		}
		l.ips[key] = entry
	}

	entry.lastSeen = time.Now()
	return entry.limiter.Allow(), nil
}

func (l *inMemoryRateLimiter) startCleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, entry := range l.ips {
			if now.Sub(entry.lastSeen) > l.cleanupD {
				delete(l.ips, ip)
			}
		}
		l.mu.Unlock()
	}
}
