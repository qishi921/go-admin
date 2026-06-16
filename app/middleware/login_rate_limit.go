package middleware

import (
	"sync"
	"time"
)

// LoginAttempt tracks login attempts for rate limiting.
type LoginAttempt struct {
	Count     int
	FirstTime time.Time
	Blocked   bool
}

// LoginRateLimiter manages login rate limiting.
type LoginRateLimiter struct {
	mu             sync.RWMutex
	attempts       map[string]*LoginAttempt
	maxAttempts    int
	blockDuration  time.Duration
	windowDuration time.Duration
	stopCleanup    chan struct{}
}

var loginLimiter = &LoginRateLimiter{
	attempts:       make(map[string]*LoginAttempt),
	maxAttempts:    5,
	blockDuration:  15 * time.Minute,
	windowDuration: 5 * time.Minute,
	stopCleanup:    make(chan struct{}),
}

// init 启动定期清理
func init() {
	go loginLimiter.cleanupExpiredAttempts()
}

// cleanupExpiredAttempts 定期清理过期的登录尝试记录
func (l *LoginRateLimiter) cleanupExpiredAttempts() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.mu.Lock()
			now := time.Now()
			maxAge := l.blockDuration + l.windowDuration
			for ip, attempt := range l.attempts {
				// 清理超过最大保留时间的记录
				if now.Sub(attempt.FirstTime) > maxAge {
					delete(l.attempts, ip)
				}
			}
			l.mu.Unlock()
		case <-l.stopCleanup:
			return
		}
	}
}

// IsLoginBlocked checks if an IP is blocked from login attempts.
func IsLoginBlocked(ip string) bool {
	loginLimiter.mu.RLock()
	attempt, exists := loginLimiter.attempts[ip]
	loginLimiter.mu.RUnlock()

	if !exists {
		return false
	}

	if !attempt.Blocked {
		return false
	}

	// Check if block duration has passed
	if time.Since(attempt.FirstTime) >= loginLimiter.blockDuration {
		loginLimiter.mu.Lock()
		delete(loginLimiter.attempts, ip)
		loginLimiter.mu.Unlock()
		return false
	}

	return true
}

// RecordLoginFailure records a failed login attempt.
func RecordLoginFailure(ip string) {
	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()

	attempt, exists := loginLimiter.attempts[ip]
	if !exists {
		attempt = &LoginAttempt{
			Count:     1,
			FirstTime: time.Now(),
		}
		loginLimiter.attempts[ip] = attempt
		return
	}

	// Reset if window expired
	if time.Since(attempt.FirstTime) > loginLimiter.windowDuration {
		attempt.Count = 1
		attempt.FirstTime = time.Now()
		attempt.Blocked = false
		return
	}

	attempt.Count++
	if attempt.Count >= loginLimiter.maxAttempts {
		attempt.Blocked = true
	}
}

// RecordLoginSuccess clears the login attempts for an IP.
func RecordLoginSuccess(ip string) {
	loginLimiter.mu.Lock()
	delete(loginLimiter.attempts, ip)
	loginLimiter.mu.Unlock()
}