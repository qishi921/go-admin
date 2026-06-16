package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ghttp "github.com/Hlgxz/gai/http"
)

func TestGenerateToken(t *testing.T) {
	token := generateToken()

	// 验证长度（hex 编码后是 64 字符）
	if len(token) != 64 {
		t.Errorf("Expected token length 64, got %d", len(token))
	}

	// 验证两次生成不同的 token
	token2 := generateToken()
	if token == token2 {
		t.Error("Two tokens should be different")
	}
}

func TestCSRFMiddleware_GET(t *testing.T) {
	mw := CSRFMiddleware()

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	mw(c)

	// 应设置 cookie
	cookies := w.Result().Cookies()
	var csrfCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == csrfCookieName {
			csrfCookie = cookie
			break
		}
	}

	if csrfCookie == nil {
		t.Error("CSRF cookie should be set for GET requests")
	}

	// 验证 token 被设置在 context 中
	receivedToken := GetCSRFToken(c)
	if receivedToken == "" {
		t.Error("CSRF token should be set for GET requests")
	}
}

func TestCSRFMiddleware_POST_MissingToken(t *testing.T) {
	mw := CSRFMiddleware()

	req := httptest.NewRequest("POST", "/test", nil)
	// 设置 cookie 但不发送 header
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "testtoken"})
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	mw(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestCSRFMiddleware_POST_ValidToken(t *testing.T) {
	mw := CSRFMiddleware()

	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set(csrfHeaderName, "testtoken")
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "testtoken"})
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	mw(c)

	// 验证请求通过（中间件会调用 Next，但我们无法直接测试）
	// 至少验证没有返回错误状态
	if w.Code == http.StatusForbidden {
		t.Error("Should not return forbidden for valid token")
	}
}

func TestCSRFExemptMiddleware_APIPaths(t *testing.T) {
	mw := CSRFExemptMiddleware([]string{"/public"})

	// API 路径应跳过 CSRF
	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	mw(c)

	// API 路径应直接通过，不会设置 403
	if w.Code == http.StatusForbidden {
		t.Error("API paths should skip CSRF check")
	}
}

func TestGetCSRFToken(t *testing.T) {
	c := ghttp.NewContext(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	// 未设置时返回空
	if GetCSRFToken(c) != "" {
		t.Error("Should return empty string when not set")
	}

	// 设置后返回 token
	c.Set("csrf_token", "testtoken123")
	if GetCSRFToken(c) != "testtoken123" {
		t.Error("Should return the set token")
	}

	// 设置非字符串时返回空
	c.Set("csrf_token", 123)
	if GetCSRFToken(c) != "" {
		t.Error("Should return empty string for non-string value")
	}
}
