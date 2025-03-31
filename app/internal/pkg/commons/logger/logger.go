// Package logger wraps slog lib to define a simpler DSL to generate structured logs
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"userservice/internal/pkg/errors"
)

// Logger provides the entry point to the custom logger wrapping slog
type Logger struct {
	logger *slog.Logger
}

// LoggerEntry entry wraps slog to provide structured logging
type LoggerEntry struct {
	logger  *slog.Logger
	level   slog.Level
	context context.Context
}

func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // TODO: make this configurable through the config YAML
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &Logger{
		logger: slog.New(handler),
	}
}

func (l *Logger) Info(ctx context.Context) *LoggerEntry {
	return &LoggerEntry{
		logger:  l.logger,
		level:   slog.LevelInfo,
		context: ctx,
	}
}

func (l *Logger) Error(ctx context.Context) *LoggerEntry {
	return &LoggerEntry{
		logger:  l.logger,
		level:   slog.LevelError,
		context: ctx,
	}
}

func (l *Logger) Debug(ctx context.Context) *LoggerEntry {
	return &LoggerEntry{
		logger:  l.logger,
		level:   slog.LevelDebug,
		context: ctx,
	}
}

// ServiceError logs a ServiceError, capturing its structured details for better traceability
func (l *LoggerEntry) ServiceError(serviceError *errors.ServiceError) {
	attrs := []slog.Attr{
		slog.String("category", serviceError.Category().String()),
		slog.String("info", serviceError.Error()),
		slog.String("errorID", serviceError.ErrorID()),
	}

	// Since our ServiceErrors do not always wrap an error, we can omit it
	if innerErr := serviceError.InnerError(); innerErr != nil {
		attrs = append(attrs, slog.Any("innerError", innerErr))
	}

	l.logger.LogAttrs(l.context, l.level, serviceError.InternalMessage(), attrs...)
}

// UserServiceError logs a ServiceError alongside the user entity unique identifier for better traceability of the actions on the resource
func (l *LoggerEntry) UserServiceError(serviceError *errors.ServiceError, userID string) {
	attrs := []slog.Attr{
		slog.String("userId", userID),
		slog.String("category", serviceError.Category().String()),
		slog.String("info", serviceError.Error()),
		slog.String("errorID", serviceError.ErrorID()),
	}

	// Since our ServiceErrors do not always wrap an error, we can omit it
	if innerErr := serviceError.InnerError(); innerErr != nil {
		attrs = append(attrs, slog.Any("innerError", innerErr))
	}

	l.logger.LogAttrs(l.context, l.level, serviceError.InternalMessage(), attrs...)
}

// Attributes logs the attributes alongside the provided message at the previously indicated log level
// i.e:
// logger.Debug(ctx).Attributes(map[string]any{"foo":"bar"}) will log "foo":"bar" structured at the DEBUG level
func (l *LoggerEntry) Attributes(message string, attrs map[string]any) {
	convertedAttrs := make([]slog.Attr, 0, len(attrs))
	for key, value := range attrs {
		convertedAttrs = append(convertedAttrs, slog.Any(key, value))
	}
	l.logger.LogAttrs(l.context, l.level, message, convertedAttrs...)
}

// Message logs the message at the previously indicated log level
// i.e:
// logger.Info(ctx).Message("foo") will log "foo" at the INFO level
func (l *LoggerEntry) Message(format string, args ...interface{}) {
	l.logger.Log(l.context, l.level, fmt.Sprintf(format, args...))
}

// NativeErr logs a native go error at the previously indicated log level
func (l *LoggerEntry) NativeErr(nativeErr error) {
	l.logger.Log(l.context, l.level, nativeErr.Error())
}
