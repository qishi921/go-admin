package middleware

import (
	"net/http"
	"sync"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
)

// RateLimitConfig configures the rate limiter.
type RateLimitConfig struct {
	// Requests per window
	Limit int
	// Time window in seconds
	Window int
	// Block duration in seconds after limit exceeded
	BlockDuration int
	// Key extractor function (default: IP)
	KeyExtractor func(c *ghttp.Context) string
}

// DefaultRateLimitConfig returns a default rate limit configuration.
// 100 requests per minute, block for 5 minutes after exceeding.
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Limit:         100,
		Window:        60,
		BlockDuration: 300,
		KeyExtractor:  func(c *ghttp.Context) string { return c.ClientIP() },
	}
}

// StrictRateLimitConfig returns a stricter configuration.
// 30 requests per minute, block for 15 minutes.
func StrictRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Limit:         30,
		Window:        60,
		BlockDuration: 900,
		KeyExtractor:  func(c *ghttp.Context) string { return c.ClientIP() },
	}
}

// rateLimitEntry tracks request counts for a key.
type rateLimitEntry struct {
	Count      int
	WindowStart time.Time
	Blocked    bool
	BlockEnd   time.Time
}

// RateLimiter manages rate limiting for API requests.
type RateLimiter struct {
	mu         sync.RWMutex
	entries    map[string]*rateLimitEntry
	config     *RateLimitConfig
	stopCleanup chan struct{}
}

// NewRateLimiter creates a new rate limiter with the given config.
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimitConfig()
	}
	rl := &RateLimiter{
		entries:     make(map[string]*rateLimitEntry),
		config:      config,
		stopCleanup: make(chan struct{}),
	}
	// 启动定期清理过期条目的 goroutine
	go rl.cleanupExpiredEntries()
	return rl
}

// cleanupExpiredEntries 定期清理过期的条目
func (rl *RateLimiter) cleanupExpiredEntries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			maxAge := time.Duration(rl.config.Window+rl.config.BlockDuration) * time.Second
			for key, entry := range rl.entries {
				// 清理过期的窗口或解封的条目
				if now.Sub(entry.WindowStart) > maxAge {
					delete(rl.entries, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCleanup:
			return
		}
	}
}

// Stop 停止清理 goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// RateLimitMiddleware returns middleware for API rate limiting.
func RateLimitMiddleware(config *RateLimitConfig) ghttp.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *ghttp.Context) {
		key := config.KeyExtractor(c)

		limiter.mu.RLock()
		entry, exists := limiter.entries[key]
		limiter.mu.RUnlock()

		now := time.Now()

		// Check if blocked
		if exists && entry.Blocked {
			if now.Before(entry.BlockEnd) {
				c.Error(http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
			// Block expired, reset
			limiter.mu.Lock()
			delete(limiter.entries, key)
			limiter.mu.Unlock()
			entry = nil
			exists = false
		}

		// Check window
		if exists && now.Sub(entry.WindowStart) > time.Duration(config.Window)*time.Second {
			// Window expired, reset count
			limiter.mu.Lock()
			entry.Count = 1
			entry.WindowStart = now
			entry.Blocked = false
			limiter.mu.Unlock()
		} else if exists {
			// Increment count
			limiter.mu.Lock()
			entry.Count++
			if entry.Count > config.Limit {
				entry.Blocked = true
				entry.BlockEnd = now.Add(time.Duration(config.BlockDuration) * time.Second)
			}
			limiter.mu.Unlock()

			if entry.Count > config.Limit {
				c.Error(http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
		} else {
			// New entry
			limiter.mu.Lock()
			limiter.entries[key] = &rateLimitEntry{
				Count:       1,
				WindowStart: now,
				Blocked:     false,
			}
			limiter.mu.Unlock()
		}

		c.Next()
	}
}

// APIRateLimitMiddleware applies rate limiting to API routes.
// Uses default configuration: 100 requests/minute per IP.
func APIRateLimitMiddleware() ghttp.HandlerFunc {
	return RateLimitMiddleware(DefaultRateLimitConfig())
}

// StrictRateLimitMiddleware applies strict rate limiting.
// Useful for sensitive endpoints: 30 requests/minute per IP.
func StrictRateLimitMiddleware() ghttp.HandlerFunc {
	return RateLimitMiddleware(StrictRateLimitConfig())
}

// UserBasedRateLimit returns rate limiting based on user ID.
// Requires authentication - uses JWT claims for key.
func UserBasedRateLimit(limit int, window int) ghttp.HandlerFunc {
	config := &RateLimitConfig{
		Limit:  limit,
		Window: window,
		BlockDuration: 300,
		KeyExtractor: func(c *ghttp.Context) string {
			// Use user ID if authenticated, otherwise IP
			if userID, ok := c.Get("auth_user_id"); ok {
				return "user:" + userIDToString(userID)
			}
			return "ip:" + c.ClientIP()
		},
	}
	return RateLimitMiddleware(config)
}

// userIDToString converts user ID to string
func userIDToString(id any) string {
	switch v := id.(type) {
	case uint64:
		return uint64ToString(v)
	case int:
		return intToString(v)
	case string:
		return v
	default:
		return ""
	}
}

func uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}