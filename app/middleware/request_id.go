package middleware

import (
	"github.com/google/uuid"
	ghttp "github.com/Hlgxz/gai/http"
)

// RequestIDMiddleware adds a unique request ID to each request.
// The ID is set in both the response header and the request context.
func RequestIDMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		// 从请求头获取或生成新的请求 ID
		requestID := c.Header("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 设置响应头
		c.SetHeader("X-Request-ID", requestID)

		// 存储到上下文
		c.Set("request_id", requestID)

		c.Next()
	}
}
