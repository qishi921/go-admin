package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
)

// LoggerConfig configures the logger behavior.
type LoggerConfig struct {
	// Level: debug, info, warn, error
	Level string
	// Format: json or text
	Format string
	// Output file path (empty = stdout)
	OutputFile string
	// MaxSize in MB for log rotation
	MaxSize int
	// Enable request logging
	LogRequests bool
}

// DefaultLoggerConfig returns default logger configuration.
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:       "info",
		Format:      "json",
		OutputFile:  "",
		LogRequests: true,
	}
}

var (
	logger   *slog.Logger
	logLevel slog.Level
	logFile  *os.File
)

// InitLogger initializes the structured logger.
func InitLogger(config *LoggerConfig) error {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	// Parse log level
	switch config.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Setup output
	var writer = os.Stdout
	if config.OutputFile != "" {
		// Ensure directory exists
		dir := filepath.Dir(config.OutputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		f, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		logFile = f
		writer = f
	}

	// Setup handler
	opts := &slog.HandlerOptions{Level: logLevel}
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

// CloseLogger closes the log file if opened.
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// logStatusWriter wraps ResponseWriter to capture status code for logging
type logStatusWriter struct {
	http.ResponseWriter
	status int
}

func (w *logStatusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// RequestLoggerMiddleware logs HTTP requests with structured logging.
func RequestLoggerMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Wrap response writer to capture status
		sw := &logStatusWriter{ResponseWriter: c.Writer, status: 200}
		c.Writer = sw

		c.Next()

		// Log request
		if logger != nil {
			duration := time.Since(start)
			status := sw.status

			logger.Info("HTTP Request",
				"method", method,
				"path", path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", c.ClientIP(),
				"user_agent", c.Header("User-Agent"),
			)
		}
	}
}

// Logging helper functions
func Debug(msg string, args ...any) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if logger != nil {
		logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if logger != nil {
		logger.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if logger != nil {
		logger.Error(msg, args...)
	}
}

// With returns a logger with preset key-value pairs.
func With(args ...any) *slog.Logger {
	if logger != nil {
		return logger.With(args...)
	}
	return slog.Default()
}
