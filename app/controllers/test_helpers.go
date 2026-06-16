package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
)

// MockContext creates a mock HTTP context for testing.
func MockContext(method, path string, body any) *ghttp.Context {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	recorder := httptest.NewRecorder()
	return &ghttp.Context{
		Request: req,
		Writer:  recorder,
	}
}

// MockDB creates an in-memory SQLite database for testing.
// Note: This is a placeholder - actual implementation depends on Gai's DB setup.
func MockDB() *orm.DB {
	// In a real test, you would:
	// 1. Create an in-memory SQLite database
	// 2. Run migrations
	// 3. Seed test data
	return nil
}

// AssertSuccess checks if response indicates success.
func AssertSuccess(t *testing.T, response map[string]any) {
	code, ok := response["code"].(float64)
	if !ok || code != 0 {
		t.Errorf("Expected success (code=0), got code=%v", response["code"])
	}
}

// AssertError checks if response indicates an error.
func AssertError(t *testing.T, response map[string]any, expectedCode int) {
	code, ok := response["code"].(float64)
	if !ok || int(code) != expectedCode {
		t.Errorf("Expected error code=%d, got code=%v", expectedCode, response["code"])
	}
}

// ParseResponse parses JSON response into map.
func ParseResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	var result map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return result
}
