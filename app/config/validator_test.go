package config

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	v := NewValidator()
	v.Required("field", "")
	if !v.HasErrors() {
		t.Error("Expected error for empty value")
	}

	v = NewValidator()
	v.Required("field", "value")
	if v.HasErrors() {
		t.Error("Expected no error for non-empty value")
	}
}

func TestValidator_NotDefault(t *testing.T) {
	defaults := []string{"secret", "password"}

	v := NewValidator()
	v.NotDefault("field", "secret", defaults)
	if !v.HasErrors() {
		t.Error("Expected error for default value")
	}

	v = NewValidator()
	v.NotDefault("field", "my-secret-key", defaults)
	if v.HasErrors() {
		t.Error("Expected no error for non-default value")
	}
}

func TestValidator_MinLength(t *testing.T) {
	v := NewValidator()
	v.MinLength("field", "ab", 3)
	if !v.HasErrors() {
		t.Error("Expected error for too short value")
	}

	v = NewValidator()
	v.MinLength("field", "abc", 3)
	if v.HasErrors() {
		t.Error("Expected no error for valid length")
	}
}

func TestValidator_OneOf(t *testing.T) {
	allowed := []string{"development", "staging", "production"}

	v := NewValidator()
	v.OneOf("env", "test", allowed)
	if !v.HasErrors() {
		t.Error("Expected error for invalid value")
	}

	v = NewValidator()
	v.OneOf("env", "production", allowed)
	if v.HasErrors() {
		t.Error("Expected no error for valid value")
	}
}

func TestValidateConfig_Production(t *testing.T) {
	cfg := &AppConfig{
		Env:       "production",
		Debug:     true, // Should fail
		DBDsn:     "storage/database.db", // Should fail
		JWTSecret: "short", // Should fail (too short)
		LogLevel:  "info",
		LogFormat: "json",
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("Expected validation errors for production config")
	}

	// Check errors contain expected messages
	errStr := err.Error()
	if !containsStr(errStr, "debug") {
		t.Error("Expected error about debug mode")
	}
}

func TestValidateConfig_Development(t *testing.T) {
	cfg := &AppConfig{
		Env:       "development",
		Debug:     true,
		DBDsn:     "storage/database.db",
		JWTSecret: "", // Should auto-generate
		LogLevel:  "info",
		LogFormat: "json",
	}

	err := ValidateConfig(cfg)
	if err != nil {
		t.Errorf("Unexpected error for development config: %v", err)
	}

	// Check JWT secret was generated
	if cfg.JWTSecret == "" {
		t.Error("JWT secret should be auto-generated")
	}
	if len(cfg.JWTSecret) < 32 {
		t.Error("Generated JWT secret too short")
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
