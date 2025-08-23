package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ranggaaprilio/boilerGo/config"
)

// Logger represents the application logger
type Logger struct {
	*slog.Logger
}

// LogLevel represents different log levels
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// New creates a new structured logger instance
func New(conf config.Configurations) *Logger {
	var level slog.Level

	// Parse log level from configuration
	switch strings.ToLower(conf.App.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer = os.Stdout

	// Configure output based on environment
	if conf.IsProduction() {
		// In production, you might want to write to a file or external service
		writer = os.Stdout
	}

	// Create handler based on environment
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: conf.IsDevelopment(),
	}

	if conf.IsDevelopment() {
		// Use text handler for development (more readable)
		handler = slog.NewTextHandler(writer, opts)
	} else {
		// Use JSON handler for production (structured)
		handler = slog.NewJSONHandler(writer, opts)
	}

	// Add service information to all logs
	logger := slog.New(handler).With(
		"service", conf.App.ServiceName,
		"environment", conf.Server.Environment,
		"version", "1.0.0", // This should come from build info
	)

	return &Logger{Logger: logger}
}

// WithContext returns a logger with additional context
func (l *Logger) WithContext(ctx ...interface{}) *Logger {
	if len(ctx)%2 != 0 {
		l.Error("WithContext called with odd number of arguments", "args", ctx)
		return l
	}

	return &Logger{Logger: l.Logger.With(ctx...)}
}

// WithComponent returns a logger with component information
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{Logger: l.Logger.With("component", component)}
}

// WithRequestID returns a logger with request ID
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With("request_id", requestID)}
}

// WithUserID returns a logger with user ID
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{Logger: l.Logger.With("user_id", userID)}
}

// WithError returns a logger with error information
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// LogRequest logs HTTP request information
func (l *Logger) LogRequest(method, path, userAgent, ip string, duration time.Duration, statusCode int) {
	l.Info("HTTP Request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"user_agent", userAgent,
		"client_ip", ip,
	)
}

// LogDatabaseQuery logs database query information
func (l *Logger) LogDatabaseQuery(query string, duration time.Duration, err error) {
	if err != nil {
		l.Error("Database Query Failed",
			"query", query,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
		)
	} else {
		l.Debug("Database Query",
			"query", query,
			"duration_ms", duration.Milliseconds(),
		)
	}
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(port string, config interface{}) {
	l.Info("Application Starting",
		"port", port,
		"config", config,
	)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(reason string) {
	l.Info("Application Shutting Down",
		"reason", reason,
		"timestamp", time.Now(),
	)
}

// Fatal logs a fatal message and exits the application
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}

// Middleware creates an Echo middleware for request logging
func (l *Logger) Middleware() func(next func(c interface{}) error) func(c interface{}) error {
	return func(next func(c interface{}) error) func(c interface{}) error {
		return func(c interface{}) error {
			start := time.Now()

			if err := next(c); err != nil {
				// Handle Echo context here if needed
				// This is a simplified version
				l.WithError(err).Error("Request failed")
				return err
			}

			duration := time.Since(start)
			l.Debug("Request completed", "duration_ms", duration.Milliseconds())

			return nil
		}
	}
}

// Health logs health check information
func (l *Logger) LogHealthCheck(component string, status string, duration time.Duration, message string) {
	l.Info("Health Check",
		"component", component,
		"status", status,
		"duration_ms", duration.Milliseconds(),
		"message", message,
	)
}

// LogConfigLoad logs configuration loading
func (l *Logger) LogConfigLoad(source string, success bool) {
	if success {
		l.Info("Configuration loaded successfully", "source", source)
	} else {
		l.Error("Failed to load configuration", "source", source)
	}
}

// LogDatabaseConnection logs database connection events
func (l *Logger) LogDatabaseConnection(event string, details map[string]interface{}) {
	args := []interface{}{"event", event}
	for k, v := range details {
		args = append(args, k, v)
	}
	l.Info("Database Connection", args...)
}

// Performance logs performance metrics
func (l *Logger) LogPerformance(operation string, duration time.Duration, additionalFields ...interface{}) {
	args := []interface{}{
		"operation", operation,
		"duration_ms", duration.Milliseconds(),
	}
	args = append(args, additionalFields...)

	if duration > 1*time.Second {
		l.Warn("Slow Operation Detected", args...)
	} else {
		l.Debug("Performance Metric", args...)
	}
}

// Security logs security-related events
func (l *Logger) LogSecurity(event string, userID string, ip string, details map[string]interface{}) {
	args := []interface{}{
		"security_event", event,
		"user_id", userID,
		"client_ip", ip,
	}

	for k, v := range details {
		args = append(args, k, v)
	}

	l.Warn("Security Event", args...)
}

// Business logs business logic events
func (l *Logger) LogBusiness(event string, entity string, entityID string, details map[string]interface{}) {
	args := []interface{}{
		"business_event", event,
		"entity", entity,
		"entity_id", entityID,
	}

	for k, v := range details {
		args = append(args, k, v)
	}

	l.Info("Business Event", args...)
}
