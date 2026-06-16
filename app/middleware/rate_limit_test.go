package middleware

import (
	"testing"
	"time"
)

func TestLoginRateLimiter_IsBlocked(t *testing.T) {
	// 重置状态
	loginLimiter = &LoginRateLimiter{
		attempts:       make(map[string]*LoginAttempt),
		maxAttempts:    5,
		blockDuration:  15 * time.Minute,
		windowDuration: 5 * time.Minute,
	}

	ip := "192.168.1.1"

	// 初始状态不应被阻止
	if IsLoginBlocked(ip) {
		t.Error("IP should not be blocked initially")
	}

	// 记录失败次数
	for i := 0; i < 4; i++ {
		RecordLoginFailure(ip)
	}
	if IsLoginBlocked(ip) {
		t.Error("IP should not be blocked after 4 attempts")
	}

	// 第 5 次失败应被阻止
	RecordLoginFailure(ip)
	if !IsLoginBlocked(ip) {
		t.Error("IP should be blocked after 5 attempts")
	}

	// 成功登录后应清除
	RecordLoginSuccess(ip)
	if IsLoginBlocked(ip) {
		t.Error("IP should not be blocked after successful login")
	}
}

func TestLoginRateLimiter_WindowExpiry(t *testing.T) {
	loginLimiter = &LoginRateLimiter{
		attempts:       make(map[string]*LoginAttempt),
		maxAttempts:    5,
		blockDuration:  15 * time.Minute,
		windowDuration: 100 * time.Millisecond, // 短窗口用于测试
	}

	ip := "192.168.1.2"

	// 记录几次失败
	for i := 0; i < 3; i++ {
		RecordLoginFailure(ip)
	}

	// 等待窗口过期
	time.Sleep(150 * time.Millisecond)

	// 新的失败应重置计数
	RecordLoginFailure(ip)

	loginLimiter.mu.RLock()
	attempt := loginLimiter.attempts[ip]
	loginLimiter.mu.RUnlock()

	if attempt == nil {
		t.Fatal("Attempt should exist")
	}
	if attempt.Count != 1 {
		t.Errorf("Expected count 1 after window expiry, got %d", attempt.Count)
	}
}

func TestRateLimiter_New(t *testing.T) {
	config := &RateLimitConfig{
		Limit:         10,
		Window:        60,
		BlockDuration: 300,
	}

	limiter := NewRateLimiter(config)
	if limiter == nil {
		t.Fatal("Limiter should not be nil")
	}
	if limiter.config.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", limiter.config.Limit)
	}
}

func TestRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRateLimitConfig()
	if config.Limit != 100 {
		t.Errorf("Expected default limit 100, got %d", config.Limit)
	}
	if config.Window != 60 {
		t.Errorf("Expected default window 60, got %d", config.Window)
	}
}

func TestRateLimiter_StrictConfig(t *testing.T) {
	config := StrictRateLimitConfig()
	if config.Limit != 30 {
		t.Errorf("Expected strict limit 30, got %d", config.Limit)
	}
	if config.BlockDuration != 900 {
		t.Errorf("Expected strict block duration 900, got %d", config.BlockDuration)
	}
}

func TestUserIDToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{uint64(123), "123"},
		{int(456), "456"},
		{"789", "789"},
		{0, "0"},
		{int(-10), "-10"},
	}

	for _, tt := range tests {
		result := userIDToString(tt.input)
		if result != tt.expected {
			t.Errorf("userIDToString(%v) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{-456, "-456"},
	}

	for _, tt := range tests {
		result := intToString(tt.input)
		if result != tt.expected {
			t.Errorf("intToString(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestUint64ToString(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{18446744073709551615, "18446744073709551615"}, // max uint64
	}

	for _, tt := range tests {
		result := uint64ToString(tt.input)
		if result != tt.expected {
			t.Errorf("uint64ToString(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
