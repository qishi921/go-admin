package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	ghttp "github.com/Hlgxz/gai/http"
)

// CSRF token configuration
const (
	csrfTokenLength = 32
	csrfHeaderName  = "X-CSRF-Token"
	csrfCookieName  = "csrf_token"
	csrfFieldName   = "_csrf"
)

// isSecureEnv 检查是否为安全环境（HTTPS）
func isSecureEnv() bool {
	return os.Getenv("APP_ENV") == "production" || os.Getenv("HTTPS") == "true"
}

// generateToken creates a cryptographically secure random token
func generateToken() string {
	bytes := make([]byte, csrfTokenLength)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CSRFMiddleware provides CSRF protection for state-changing requests.
// For API endpoints using JWT authentication, CSRF is typically not needed
// because the token is stored in localStorage (not cookies).
func CSRFMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		// Skip CSRF for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			// Generate and set token for GET requests
			token := generateToken()
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     csrfCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: false, // JavaScript needs to read it
				Secure:   isSecureEnv(),
				SameSite: http.SameSiteStrictMode,
			})
			c.Set("csrf_token", token)
			c.Next()
			return
		}

		// For state-changing methods, validate CSRF token
		token := c.Header(csrfHeaderName)
		if token == "" {
			token = c.Request.FormValue(csrfFieldName)
		}

		// Get token from cookie
		cookie, err := c.Request.Cookie(csrfCookieName)
		if err != nil {
			c.Error(http.StatusForbidden, "CSRF token missing")
			return
		}

		if token == "" || token != cookie.Value {
			c.Error(http.StatusForbidden, "Invalid CSRF token")
			return
		}

		c.Next()
	}
}

// CSRFExemptMiddleware skips CSRF for specific paths.
func CSRFExemptMiddleware(exemptPaths []string) ghttp.HandlerFunc {
	exempt := make(map[string]bool)
	for _, p := range exemptPaths {
		exempt[p] = true
	}

	return func(c *ghttp.Context) {
		path := c.Request.URL.Path

		// Skip CSRF for exempt paths (API endpoints with JWT)
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}

		// Apply CSRF for non-API routes
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
			token := c.Header(csrfHeaderName)
			if token == "" {
				token = c.Request.FormValue(csrfFieldName)
			}

			cookie, err := c.Request.Cookie(csrfCookieName)
			if err != nil || token == "" || token != cookie.Value {
				c.Error(http.StatusForbidden, "Invalid CSRF token")
				return
			}
		}

		c.Next()
	}
}

// GetCSRFToken returns a CSRF token for the current request.
func GetCSRFToken(c *ghttp.Context) string {
	if token, exists := c.Get("csrf_token"); exists {
		if t, ok := token.(string); ok {
			return t
		}
	}
	return ""
}
