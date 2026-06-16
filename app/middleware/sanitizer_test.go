package middleware

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	s := DefaultSanitizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"  hello  ", "hello"},
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"hello\x00world", "helloworld"},
		{"normal text", "normal text"},
		{"  ", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := s.SanitizeString(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeMap(t *testing.T) {
	s := DefaultSanitizer()

	input := map[string]any{
		"name":    "<script>test</script>",
		"email":   "  test@example.com  ",
		"nested":  map[string]any{"value": "<b>bold</b>"},
		"number":  123,
		"nil_val": nil,
	}

	result := s.SanitizeMap(input)

	// Check string sanitization
	if result["name"] != "&lt;script&gt;test&lt;/script&gt;" {
		t.Errorf("Name not sanitized correctly: %v", result["name"])
	}

	if result["email"] != "test@example.com" {
		t.Errorf("Email not trimmed correctly: %v", result["email"])
	}

	// Check nested map
	if nested, ok := result["nested"].(map[string]any); ok {
		if nested["value"] != "&lt;b&gt;bold&lt;/b&gt;" {
			t.Errorf("Nested value not sanitized: %v", nested["value"])
		}
	} else {
		t.Error("Nested map not preserved")
	}

	// Check non-string values preserved
	if result["number"] != 123 {
		t.Errorf("Number not preserved: %v", result["number"])
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		contains string // Check that result does NOT contain this
	}{
		{"../../../etc/passwd", ".."},
		{"file.txt", ""},
		{"/etc/passwd", "/"},
		{"..\\..\\windows\\system32", ".."},
	}

	for _, tt := range tests {
		result := SanitizeFilename(tt.input)
		if tt.contains != "" && containsStr(result, tt.contains) {
			t.Errorf("SanitizeFilename(%q) = %q, should not contain %q", tt.input, result, tt.contains)
		}
	}
}

func TestIsSafeString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"normal text", true},
		{"text with newline\n", true},
		{"text with tab\t", true},
		{"text\x00null", false},
		{"text\x01control", false},
	}

	for _, tt := range tests {
		result := IsSafeString(tt.input)
		if result != tt.expected {
			t.Errorf("IsSafeString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
