package db

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

// NewTracer wires pgx tracing events into zap for easier observability.
func NewTracer(logger *zap.Logger) *tracelog.TraceLog {
	if logger == nil {
		return nil
	}

	return &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
			fields := make([]zap.Field, 0, len(data)+2)
			fields = append(fields, zap.String("component", "pgx"))
			for k, v := range data {
				fields = append(fields, zap.Any(k, v))
			}

			switch level {
			case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
				logger.Debug(msg, fields...)
			case tracelog.LogLevelInfo:
				logger.Info(msg, fields...)
			case tracelog.LogLevelWarn:
				logger.Warn(msg, fields...)
			case tracelog.LogLevelError:
				logger.Error(msg, fields...)
			default:
				logger.Info(msg, fields...)
			}
		}),
		LogLevel: tracelog.LogLevelInfo,
	}
}
