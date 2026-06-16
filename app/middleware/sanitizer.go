package middleware

import (
	"html"
	"regexp"
	"strings"
)

// InputSanitizer provides methods for sanitizing user input.
type InputSanitizer struct {
	// AllowedHTMLTags specifies which HTML tags are allowed (empty = no HTML)
	AllowedHTMLTags []string
	// MaxLength limits string length (0 = no limit)
	MaxLength int
	// TrimWhitespace removes leading/trailing whitespace
	TrimWhitespace bool
}

// DefaultSanitizer returns a sanitizer with safe defaults.
func DefaultSanitizer() *InputSanitizer {
	return &InputSanitizer{
		AllowedHTMLTags: []string{}, // No HTML allowed by default
		MaxLength:       0,          // No length limit by default
		TrimWhitespace:  true,
	}
}

// SanitizeString cleans a string for safe storage/display.
func (s *InputSanitizer) SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Trim whitespace
	if s.TrimWhitespace {
		input = strings.TrimSpace(input)
	}

	// Escape HTML entities
	input = html.EscapeString(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters (except newlines and tabs)
	input = removeControlChars(input)

	// Limit length
	if s.MaxLength > 0 && len(input) > s.MaxLength {
		input = input[:s.MaxLength]
	}

	return input
}

// SanitizeMap recursively sanitizes all string values in a map.
func (s *InputSanitizer) SanitizeMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	result := make(map[string]any)
	for key, value := range input {
		// Sanitize key
		cleanKey := s.SanitizeString(key)

		// Sanitize value based on type
		switch v := value.(type) {
		case string:
			result[cleanKey] = s.SanitizeString(v)
		case map[string]any:
			result[cleanKey] = s.SanitizeMap(v)
		case []any:
			result[cleanKey] = s.SanitizeSlice(v)
		default:
			result[cleanKey] = v
		}
	}
	return result
}

// SanitizeSlice sanitizes all elements in a slice.
func (s *InputSanitizer) SanitizeSlice(input []any) []any {
	if input == nil {
		return nil
	}

	result := make([]any, len(input))
	for i, value := range input {
		switch v := value.(type) {
		case string:
			result[i] = s.SanitizeString(v)
		case map[string]any:
			result[i] = s.SanitizeMap(v)
		case []any:
			result[i] = s.SanitizeSlice(v)
		default:
			result[i] = v
		}
	}
	return result
}

// removeControlChars removes control characters except newlines and tabs.
func removeControlChars(input string) string {
	var result strings.Builder
	for _, r := range input {
		// Allow printable chars, newlines, and tabs
		if r >= 32 || r == '\n' || r == '\t' || r == '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// StripTags removes all HTML tags from a string.
func StripTags(input string) string {
	// Simple regex to remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

// SanitizeFilename cleans a filename for safe filesystem use.
func SanitizeFilename(filename string) string {
	if filename == "" {
		return ""
	}

	// Remove path separators
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// Remove parent directory references
	filename = strings.ReplaceAll(filename, "..", "")

	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Trim whitespace
	filename = strings.TrimSpace(filename)

	// Limit length
	if len(filename) > 255 {
		filename = filename[:255]
	}

	return filename
}

// SanitizeEmail normalizes an email address.
func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}

// SanitizePhone normalizes a phone number.
func SanitizePhone(phone string) string {
	// Remove common separators, keep only digits and +
	var result strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' || r == '+' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// IsSafeString checks if a string contains only safe characters.
func IsSafeString(input string) bool {
	for _, r := range input {
		// Check for control characters
		if r < 32 && r != '\n' && r != '\t' && r != '\r' {
			return false
		}
		// Check for null byte
		if r == 0 {
			return false
		}
	}
	return true
}

// SanitizeInput is a convenience function using default sanitizer.
func SanitizeInput(input string) string {
	return DefaultSanitizer().SanitizeString(input)
}

// SanitizeMapInput is a convenience function using default sanitizer.
func SanitizeMapInput(input map[string]any) map[string]any {
	return DefaultSanitizer().SanitizeMap(input)
}
