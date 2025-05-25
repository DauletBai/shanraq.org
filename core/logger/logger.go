// github.com/DauletBai/shanraq.org/core/logger/logger.go
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// LogLevel type for defining log levels
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// Interface defines a common logging interface for the framework.
// It's based on slog but abstracts it slightly to allow different implementations.
type Interface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	ErrorContext(ctx context.Context, msg string, args ...interface{}) // For context-aware logging
	With(args ...interface{}) Interface                               // For adding contextual attributes
	Handler() slog.Handler                                            // Expose the underlying handler if needed
}

// SlogLogger is an implementation of Logger that uses Go's standard slog package.
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new SlogLogger.
// 'level' can be "DEBUG", "INFO", "WARN", "ERROR".
// 'output' is where the logs will be written (e.g., os.Stdout, a file).
// 'isJSON' determines if the output should be JSON formatted.
func NewSlogLogger(level LogLevel, output io.Writer, isJSON bool, addSource bool) *SlogLogger {
	if output == nil {
		output = os.Stdout
	}

	var slogLevel slog.Level
	switch strings.ToUpper(string(level)) {
	case string(LevelDebug):
		slogLevel = slog.LevelDebug
	case string(LevelInfo):
		slogLevel = slog.LevelInfo
	case string(LevelWarn):
		slogLevel = slog.LevelWarn
	case string(LevelError):
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo // Default to INFO
	}

	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: addSource, // Adds source file and line number
	}

	var handler slog.Handler
	if isJSON {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	return &SlogLogger{logger: slog.New(handler)}
}

// Default returns a basic SlogLogger with INFO level, text output to os.Stdout.
func Default() *SlogLogger {
	return NewSlogLogger(LevelInfo, os.Stdout, false, false)
}

func (l *SlogLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	l.logger.ErrorContext(ctx, msg, args...)
}


// With returns a new Logger with the given arguments added as contextual attributes.
func (l *SlogLogger) With(args ...interface{}) Interface {
	return &SlogLogger{logger: l.logger.With(args...)}
}

// Handler returns the underlying slog.Handler.
func (l *SlogLogger) Handler() slog.Handler {
	return l.logger.Handler()
}


// Ensure SlogLogger implements Interface
var _ Interface = (*SlogLogger)(nil)

// Global instance (optional, dependency injection is generally preferred)
// We will manage this through the Kernel later.
// var globalLogger Interface = Default()

// SetGlobalLogger sets the global logger instance.
// func SetGlobalLogger(logger Interface) {
// 	globalLogger = logger
// }

// GetGlobalLogger returns the global logger instance.
// func GetGlobalLogger() Interface {
// 	return globalLogger
// }

// Convenience functions to use the global logger (if you choose to have one and set it)
// func Debug(msg string, args ...interface{}) { globalLogger.Debug(msg, args...) }
// func Info(msg string, args ...interface{})  { globalLogger.Info(msg, args...) }
// func Warn(msg string, args ...interface{})  { globalLogger.Warn(msg, args...) }
// func Error(msg string, args ...interface{}) { globalLogger.Error(msg, args...) }
// func ErrorContext(ctx context.Context, msg string, args ...interface{}) { globalLogger.ErrorContext(ctx, msg, args...) }
// func With(args ...interface{}) Interface { return globalLogger.With(args...) }