package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// contextKey is a type for context keys
type contextKey string

const (
	// RequestIDKey is the context key for the request ID
	RequestIDKey contextKey = "request_id"
	// UserIDKey is the context key for the user ID
	UserIDKey contextKey = "user_id"
)

// Field creates a field for structured logging
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Init initializes the global logger
func Init(logLevel string, isDevelopment bool) error {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		level = zap.InfoLevel
	}

	config := zap.NewProductionConfig()
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var logger *zap.Logger
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// WithRequestID adds request ID to the logger context
func WithRequestID(ctx context.Context) *zap.Logger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return globalLogger.With(zap.String("request_id", requestID))
	}
	return globalLogger
}

// WithUserID adds user ID to the logger context
func WithUserID(ctx context.Context) *zap.Logger {
	if userID, ok := ctx.Value(UserIDKey).(uuid.UUID); ok && userID != uuid.Nil {
		return globalLogger.With(zap.String("user_id", userID.String()))
	}
	return globalLogger
}

// WithContext adds all context values to the logger
func WithContext(ctx context.Context) *zap.Logger {
	logger := globalLogger
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With(zap.String("request_id", requestID))
	}
	if userID, ok := ctx.Value(UserIDKey).(uuid.UUID); ok && userID != uuid.Nil {
		logger = logger.With(zap.String("user_id", userID.String()))
	}
	return logger
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return globalLogger.Sync()
}
