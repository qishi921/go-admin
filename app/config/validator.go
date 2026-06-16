package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Validator validates application configuration.
type Validator struct {
	errors []string
}

// NewValidator creates a new config validator.
func NewValidator() *Validator {
	return &Validator{errors: make([]string, 0)}
}

// AddError adds a validation error.
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, fmt.Sprintf("%s: %s", field, message))
}

// HasErrors returns true if there are validation errors.
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors.
func (v *Validator) Errors() []string {
	return v.errors
}

// Error returns a formatted error message.
func (v *Validator) Error() string {
	if !v.HasErrors() {
		return ""
	}
	return "Configuration errors:\n  - " + strings.Join(v.errors, "\n  - ")
}

// Required checks that a value is not empty.
func (v *Validator) Required(field, value string) {
	if value == "" {
		v.AddError(field, "is required")
	}
}

// NotDefault checks that a value is not a known weak default.
func (v *Validator) NotDefault(field, value string, defaults []string) {
	for _, d := range defaults {
		if value == d {
			v.AddError(field, fmt.Sprintf("must not be default value '%s'", d))
			return
		}
	}
}

// MinLength checks minimum string length.
func (v *Validator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters", min))
	}
}

// OneOf checks that value is one of allowed values.
func (v *Validator) OneOf(field, value string, allowed []string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
}

// AppConfig represents validated application configuration.
type AppConfig struct {
	// Server
	Port    int
	Env     string
	Debug   bool

	// Database
	DBDriver string
	DBDsn    string

	// Auth
	JWTSecret string
	JWTTL     int

	// Logging
	LogLevel  string
	LogFormat string
	LogFile   string
}

// ValidateConfig validates the configuration and returns errors if any.
func ValidateConfig(cfg *AppConfig) error {
	v := NewValidator()

	// Environment
	v.OneOf("env", cfg.Env, []string{"development", "staging", "production", "test"})

	// JWT Secret - critical security check
	if cfg.JWTSecret == "" {
		// Generate a random secret for development
		if cfg.Env == "development" {
			cfg.JWTSecret = generateRandomSecret(32)
			fmt.Println("Warning: JWT_SECRET not set, generated random secret for development")
		} else {
			v.Required("jwt_secret", cfg.JWTSecret)
		}
	} else {
		// Check for weak defaults
		v.NotDefault("jwt_secret", cfg.JWTSecret, []string{
			"your-secret-key-change-me",
			"change-me-to-a-random-string",
			"secret",
			"jwt-secret",
			"123456",
		})
		// Minimum length for security
		v.MinLength("jwt_secret", cfg.JWTSecret, 16)
	}

	// In production, require stronger settings
	if cfg.Env == "production" {
		if cfg.Debug {
			v.AddError("debug", "must be false in production")
		}
		if !strings.HasPrefix(cfg.DBDsn, "mysql://") && !strings.HasPrefix(cfg.DBDsn, "postgres://") {
			v.AddError("database.dsn", "production should use MySQL or PostgreSQL")
		}
	}

	// Log level
	v.OneOf("log.level", cfg.LogLevel, []string{"debug", "info", "warn", "error"})

	// Log format
	v.OneOf("log.format", cfg.LogFormat, []string{"json", "text"})

	if v.HasErrors() {
		return v
	}
	return nil
}

// generateRandomSecret creates a cryptographically secure random string.
func generateRandomSecret(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetEnvOrDefault gets an environment variable or returns default.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if running in development mode.
func IsDevelopment(env string) bool {
	return env == "development" || env == "dev" || env == ""
}

// IsProduction returns true if running in production mode.
func IsProduction(env string) bool {
	return env == "production" || env == "prod"
}
