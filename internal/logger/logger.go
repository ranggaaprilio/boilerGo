package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LogrusLogger wraps logrus.Logger with additional functionality
type LogrusLogger struct {
	*logrus.Logger
	component string
}

// NewLogrusLogger creates a new logrus-based logger
func NewLogrusLogger(logLevel string, isDevelopment bool, serviceName string) *LogrusLogger {
	logger := logrus.New()

	// Set log level
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set formatter based on environment
	if isDevelopment {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
	}

	// Set output to stdout
	logger.SetOutput(os.Stdout)

	// Add service name to all logs
	logger = logger.WithField("service", serviceName).Logger

	return &LogrusLogger{
		Logger: logger,
	}
}

// WithComponent returns a logger with component information
func (l *LogrusLogger) WithComponent(component string) *LogrusLogger {
	return &LogrusLogger{
		Logger:    l.Logger.WithField("component", component).Logger,
		component: component,
	}
}

// WithFields returns a logger with additional fields
func (l *LogrusLogger) WithFields(fields logrus.Fields) *LogrusLogger {
	return &LogrusLogger{
		Logger:    l.Logger.WithFields(fields).Logger,
		component: l.component,
	}
}

// WithField returns a logger with an additional field
func (l *LogrusLogger) WithField(key string, value interface{}) *LogrusLogger {
	return &LogrusLogger{
		Logger:    l.Logger.WithField(key, value).Logger,
		component: l.component,
	}
}

// WithError returns a logger with error information
func (l *LogrusLogger) WithError(err error) *LogrusLogger {
	return &LogrusLogger{
		Logger:    l.Logger.WithError(err).Logger,
		component: l.component,
	}
}

// Fatal logs a fatal message and exits
func (l *LogrusLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.WithFields(argsToFields(args...)).Fatal(msg)
	} else {
		l.Logger.Fatal(msg)
	}
}

// Error logs an error message
func (l *LogrusLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.WithFields(argsToFields(args...)).Error(msg)
	} else {
		l.Logger.Error(msg)
	}
}

// Warn logs a warning message
func (l *LogrusLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.WithFields(argsToFields(args...)).Warn(msg)
	} else {
		l.Logger.Warn(msg)
	}
}

// Info logs an info message
func (l *LogrusLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.WithFields(argsToFields(args...)).Info(msg)
	} else {
		l.Logger.Info(msg)
	}
}

// Debug logs a debug message
func (l *LogrusLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.WithFields(argsToFields(args...)).Debug(msg)
	} else {
		l.Logger.Debug(msg)
	}
}

// argsToFields converts key-value pairs to logrus.Fields
func argsToFields(args ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		fields[key] = args[i+1]
	}
	return fields
}

// SimpleLogger creates a basic logrus logger without config dependency
func SimpleLogger(component string) *LogrusLogger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logger.SetOutput(os.Stdout)

	return &LogrusLogger{
		Logger:    logger.WithField("component", component).Logger,
		component: component,
	}
}
