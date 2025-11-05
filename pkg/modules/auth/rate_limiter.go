package auth

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter controls how frequently sensitive auth operations can run.
type RateLimiter interface {
	Allow(action, key string) bool
}

type rateLimitRule struct {
	limit rate.Limit
	burst int
}

func defaultRateLimitRules() map[string]rateLimitRule {
	return map[string]rateLimitRule{
		"default":                {limit: rate.Every(500 * time.Millisecond), burst: 20}, // ~120 req/min
		"signin":                 {limit: rate.Every(time.Minute / 8), burst: 5},         // 8 attempts/min
		"signup":                 {limit: rate.Every(time.Minute / 3), burst: 3},         // 3 attempts/min
		"password_reset":         {limit: rate.Every(time.Minute / 4), burst: 3},         // 4 attempts/min
		"password_reset_confirm": {limit: rate.Every(time.Minute / 4), burst: 5},         // 4 attempts/min
		"mfa_verify":             {limit: rate.Every(time.Minute / 6), burst: 6},         // 10 attempts/min approx
	}
}

type memoryRateLimiter struct {
	mu          sync.Mutex
	rules       map[string]rateLimitRule
	buckets     map[string]*rate.Limiter
	lastSeen    map[string]time.Time
	ttl         time.Duration
	nextCleanup time.Time
}

func newMemoryRateLimiter(rules map[string]rateLimitRule) RateLimiter {
	if len(rules) == 0 {
		rules = defaultRateLimitRules()
	}
	return &memoryRateLimiter{
		rules:    rules,
		buckets:  make(map[string]*rate.Limiter),
		lastSeen: make(map[string]time.Time),
		ttl:      30 * time.Minute,
	}
}

func (m *memoryRateLimiter) Allow(action, key string) bool {
	if key == "" {
		key = "anonymous"
	}

	now := time.Now()
	composite := action + ":" + key

	m.mu.Lock()
	defer m.mu.Unlock()

	rule, ok := m.rules[action]
	if !ok {
		rule = m.rules["default"]
	}
	if rule.limit <= 0 {
		return true
	}
	if rule.burst <= 0 {
		rule.burst = 1
	}

	limiter, ok := m.buckets[composite]
	if !ok {
		limiter = rate.NewLimiter(rule.limit, rule.burst)
		m.buckets[composite] = limiter
	}
	m.lastSeen[composite] = now

	allowed := limiter.Allow()
	if allowed {
		m.maybeCleanupLocked(now)
		return true
	}

	m.maybeCleanupLocked(now)
	return false
}

func (m *memoryRateLimiter) maybeCleanupLocked(now time.Time) {
	if m.ttl <= 0 {
		return
	}
	if !m.nextCleanup.IsZero() && now.Before(m.nextCleanup) {
		return
	}
	cutoff := now.Add(-m.ttl)
	for key, seen := range m.lastSeen {
		if seen.Before(cutoff) {
			delete(m.lastSeen, key)
			delete(m.buckets, key)
		}
	}
	m.nextCleanup = now.Add(m.ttl)
}
