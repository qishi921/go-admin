package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	ghttp "github.com/Hlgxz/gai/http"
)

// SecurityHeadersMiddleware adds essential security headers to all responses.
// These headers protect against common web vulnerabilities:
// - XSS (Cross-Site Scripting)
// - Clickjacking
// - MIME type sniffing
// - Information disclosure
func SecurityHeadersMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		// Prevent clickjacking - don't allow embedding in frames
		c.SetHeader("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing - browser must respect declared content-type
		c.SetHeader("X-Content-Type-Options", "nosniff")

		// Enable XSS filter in browsers
		c.SetHeader("X-XSS-Protection", "1; mode=block")

		// Referrer policy - limit information sent in Referer header
		c.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy - restrict resource loading
		// For admin system, we allow inline scripts (Layui needs them) but restrict external
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.SetHeader("Content-Security-Policy", csp)

		// Permissions Policy (formerly Feature Policy) - restrict browser features
		c.SetHeader("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

		// Cache control for API responses - prevent sensitive data caching
		path := c.Request.URL.Path
		if len(path) >= 4 && path[:4] == "/api" {
			c.SetHeader("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			c.SetHeader("Pragma", "no-cache")
			c.SetHeader("Expires", "0")
			c.SetHeader("Surrogate-Control", "no-store")
		}

		// Remove server identification
		c.SetHeader("Server", "")

		c.Next()
	}
}

// CORSMiddleware provides CORS support for API endpoints.
// Configuration can be adjusted based on environment.
type CORSConfig struct {
	AllowOrigins   []string
	AllowMethods   []string
	AllowHeaders   []string
	ExposeHeaders  []string
	AllowCredentials bool
	MaxAge         int
}

// DefaultCORSConfig returns safe default CORS configuration.
func DefaultCORSConfig() *CORSConfig {
	// 生产环境应通过环境变量指定允许的来源
	allowOrigins := []string{"*"}
	if envOrigins := os.Getenv("CORS_ALLOW_ORIGINS"); envOrigins != "" {
		allowOrigins = strings.Split(envOrigins, ",")
	}

	return &CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token", "X-Requested-With"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORSMiddleware returns middleware with custom CORS configuration.
func CORSMiddleware(config *CORSConfig) ghttp.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *ghttp.Context) {
		origin := c.Header("Origin")
		if origin == "" {
			origin = "*"
		}

		// Check if origin is allowed
		allowed := false
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			c.Next()
			return
		}

		// Set CORS headers
		c.SetHeader("Access-Control-Allow-Origin", origin)
		if len(config.AllowMethods) > 0 {
			c.SetHeader("Access-Control-Allow-Methods", joinHeaders(config.AllowMethods))
		}
		if len(config.AllowHeaders) > 0 {
			c.SetHeader("Access-Control-Allow-Headers", joinHeaders(config.AllowHeaders))
		}
		if len(config.ExposeHeaders) > 0 {
			c.SetHeader("Access-Control-Expose-Headers", joinHeaders(config.ExposeHeaders))
		}
		if config.AllowCredentials {
			c.SetHeader("Access-Control-Allow-Credentials", "true")
		}
		if config.MaxAge > 0 {
			c.SetHeader("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// joinHeaders joins header values with comma
func joinHeaders(headers []string) string {
	result := ""
	for i, h := range headers {
		if i > 0 {
			result += ", "
		}
		result += h
	}
	return result
}