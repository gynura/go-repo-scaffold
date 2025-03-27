// Package logger wraps slog lib to define a simpler DSL to generate structured logs
package logger

import (
	"log/slog"
	"os"
	"userservice/internal/pkg/errors"
)

// Logger wraps slog to provide structured logging
type Logger struct {
	logger *slog.Logger
	level  slog.Level
	msg    string
	attrs  []slog.Attr
}

func NewLogger() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &Logger{
		logger: slog.New(handler),
	}
}

func (l *Logger) Info(msg string) *Logger {
	return &Logger{
		logger: l.logger,
		level:  slog.LevelInfo,
		msg:    msg,
	}
}

func (l *Logger) Error(msg string) *Logger {
	return &Logger{
		logger: l.logger,
		level:  slog.LevelError,
		msg:    msg,
	}
}

// ServiceError logs a ServiceError, capturing its structured details for better traceability.
func (l *Logger) ServiceError(serviceError *errors.ServiceError) *Logger {
	attrs := []slog.Attr{
		slog.String("category", serviceError.Category().String()),
		slog.String("info", serviceError.Error()),
		slog.String("errorID", serviceError.ErrorID()),
	}

	// Since our ServiceErrors do not always wrap an error, we can omit it
	if innerErr := serviceError.InnerError(); innerErr != nil {
		attrs = append(attrs, slog.Any("innerError", innerErr))
	}

	return &Logger{
		logger: l.logger,
		level:  slog.LevelError,
		msg:    serviceError.ExposedMessage(),
		attrs:  attrs,
	}
}

func (l *Logger) Debug(msg string) *Logger {
	return &Logger{
		logger: l.logger,
		level:  slog.LevelDebug,
		msg:    msg,
	}
}

// WithAttributes adds structured attributes and logs the message
func (l *Logger) WithAttributes(attrs map[string]any) {
	convertedAttrs := make([]slog.Attr, 0, len(attrs))
	for key, value := range attrs {
		convertedAttrs = append(convertedAttrs, slog.Any(key, value))
	}
	l.logger.LogAttrs(nil, l.level, l.msg, convertedAttrs...)
}
