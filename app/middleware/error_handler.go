package middleware

import (
	"log/slog"
	"os"
)

// IsProduction 检查是否为生产环境
func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

// SafeError 返回安全的错误消息
// 生产环境返回通用消息，开发环境返回详细消息
func SafeError(err error, userMessage string) string {
	if err != nil {
		// 记录详细错误到日志
		slog.Error("Internal error", "error", err.Error())
	}

	if IsProduction() {
		return userMessage
	}
	// 开发环境返回详细错误
	if err != nil {
		return userMessage + ": " + err.Error()
	}
	return userMessage
}

// LogError 记录错误到日志
func LogError(message string, err error, args ...any) {
	allArgs := append([]any{"error", err}, args...)
	slog.Error(message, allArgs...)
}
