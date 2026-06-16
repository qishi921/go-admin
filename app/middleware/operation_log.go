package middleware

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/models"
)

// sensitiveFields 需要脱敏的敏感字段名
var sensitiveFields = map[string]bool{
	"password":        true,
	"passwd":          true,
	"pwd":             true,
	"new_password":    true,
	"old_password":    true,
	"confirm_password": true,
	"password_confirm": true,
	"token":          true,
	"access_token":   true,
	"refresh_token":  true,
	"secret":         true,
	"api_key":        true,
	"apikey":         true,
	"authorization":  true,
	"credit_card":    true,
	"card_number":    true,
	"cvv":            true,
	"ssn":            true,
}

// statusWriter wraps http.ResponseWriter to capture the written status code
// while preserving optional interfaces (Flusher, Hijacker).
type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("gai: underlying ResponseWriter does not implement http.Hijacker")
}

// OperationLogMiddleware returns middleware that records operation logs
// asynchronously after request processing. It extracts user information
// from JWT claims and skips logging for login/register endpoints.
func OperationLogMiddleware(db *orm.DB) ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Skip logging for auth endpoints (login/register)
		if shouldSkipLog(path) {
			c.Next()
			return
		}

		// Wrap response writer to capture status code
		sw := &statusWriter{ResponseWriter: c.Writer, code: http.StatusOK}
		c.Writer = sw

		c.Next()

		// Build log entry
		log := &models.OperationLog{
			Method:    method,
			Path:      path,
			Ip:        c.ClientIP(),
			UserAgent: c.Header("User-Agent"),
			Result:    http.StatusText(sw.code),
			Duration:  int(time.Since(start).Milliseconds()),
			Status:    getStatus(sw.code),
			Action:    method + " " + path,
		}

		// Extract user info from JWT claims
		if claims, ok := c.Get("auth_claims"); ok {
			if jwtClaims, ok := claims.(*auth.Claims); ok {
				userID := int(jwtClaims.UserID)
				log.UserId = &userID
				if username, ok := jwtClaims.Extra["username"].(string); ok {
					log.Username = username
				}
			}
		}

		// Capture request params (body + query)
		log.Params = buildParams(c)

		// Save to database synchronously to avoid SQLite concurrency issues
		if _, err := orm.Create[models.OperationLog](db, log); err != nil {
			slog.Error("failed to create operation log", "error", err)
		}
	}
}

// shouldSkipLog returns true for paths that should not be logged.
func shouldSkipLog(path string) bool {
	skipPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/logout",
	}
	for _, skip := range skipPaths {
		if strings.EqualFold(path, skip) {
			return true
		}
	}
	return false
}

// getStatus returns "success" for 2xx status codes, otherwise "error".
func getStatus(code int) string {
	if code >= http.StatusOK && code < http.StatusMultipleChoices {
		return "success"
	}
	return "error"
}

// buildParams combines request body and query parameters as JSON string.
func buildParams(c *ghttp.Context) string {
	params := make(map[string]any)

	// Add query parameters
	if c.Request.URL.RawQuery != "" {
		queryParams := make(map[string]any)
		for key, values := range c.Request.URL.Query() {
			if sensitiveFields[strings.ToLower(key)] {
				queryParams[key] = "***"
			} else if len(values) == 1 {
				queryParams[key] = values[0]
			} else {
				queryParams[key] = values
			}
		}
		params["query"] = queryParams
	}

	// Add body for JSON requests
	if c.IsJSON() && c.Request.Body != nil {
		body, err := c.Body()
		if err == nil && len(body) > 0 {
			var bodyData any
			if err := json.Unmarshal(body, &bodyData); err == nil {
				params["body"] = maskSensitiveFields(bodyData)
			}
		}
	}

	if len(params) == 0 {
		return ""
	}

	data, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	return string(data)
}

// maskSensitiveFields 递归脱敏敏感字段
func maskSensitiveFields(data any) any {
	switch v := data.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, val := range v {
			if sensitiveFields[strings.ToLower(key)] {
				result[key] = "***"
			} else {
				result[key] = maskSensitiveFields(val)
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, val := range v {
			result[i] = maskSensitiveFields(val)
		}
		return result
	default:
		return data
	}
}
