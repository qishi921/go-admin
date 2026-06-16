package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	ghttp "github.com/Hlgxz/gai/http"
)

// RecoveryMiddleware recovers from panics in the request chain.
// It logs the error and returns a 500 Internal Server Error response.
func RecoveryMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				Error("Panic recovered",
					"error", fmt.Sprintf("%v", err),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", stack,
				)

				// 返回 500 错误
				c.Error(http.StatusInternalServerError, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}
