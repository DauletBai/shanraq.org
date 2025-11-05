package db

import (
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

// NewTracer wires pgx tracing events into zap logging and OpenTelemetry spans.
func NewTracer(logger *zap.Logger) pgx.QueryTracer {
	var logTracer *tracelog.TraceLog
	if logger != nil {
		logTracer = &tracelog.TraceLog{
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

	otelTracer := otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())

	return &combinedTracer{otel: otelTracer, log: logTracer}
}

type combinedTracer struct {
	otel *otelpgx.Tracer
	log  *tracelog.TraceLog
}

func (t *combinedTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if t.otel != nil {
		ctx = t.otel.TraceQueryStart(ctx, conn, data)
	}
	if t.log != nil {
		ctx = t.log.TraceQueryStart(ctx, conn, data)
	}
	return ctx
}

func (t *combinedTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if t.log != nil {
		t.log.TraceQueryEnd(ctx, conn, data)
	}
	if t.otel != nil {
		t.otel.TraceQueryEnd(ctx, conn, data)
	}
}

func (t *combinedTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	if t.otel != nil {
		ctx = t.otel.TraceBatchStart(ctx, conn, data)
	}
	if t.log != nil {
		ctx = t.log.TraceBatchStart(ctx, conn, data)
	}
	return ctx
}

func (t *combinedTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	if t.log != nil {
		t.log.TraceBatchQuery(ctx, conn, data)
	}
	if t.otel != nil {
		t.otel.TraceBatchQuery(ctx, conn, data)
	}
}

func (t *combinedTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	if t.log != nil {
		t.log.TraceBatchEnd(ctx, conn, data)
	}
	if t.otel != nil {
		t.otel.TraceBatchEnd(ctx, conn, data)
	}
}

func (t *combinedTracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	if t.otel != nil {
		ctx = t.otel.TraceCopyFromStart(ctx, conn, data)
	}
	if t.log != nil {
		ctx = t.log.TraceCopyFromStart(ctx, conn, data)
	}
	return ctx
}

func (t *combinedTracer) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
	if t.log != nil {
		t.log.TraceCopyFromEnd(ctx, conn, data)
	}
	if t.otel != nil {
		t.otel.TraceCopyFromEnd(ctx, conn, data)
	}
}

func (t *combinedTracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	if t.otel != nil {
		ctx = t.otel.TracePrepareStart(ctx, conn, data)
	}
	if t.log != nil {
		ctx = t.log.TracePrepareStart(ctx, conn, data)
	}
	return ctx
}

func (t *combinedTracer) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	if t.log != nil {
		t.log.TracePrepareEnd(ctx, conn, data)
	}
	if t.otel != nil {
		t.otel.TracePrepareEnd(ctx, conn, data)
	}
}

func (t *combinedTracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	if t.otel != nil {
		ctx = t.otel.TraceConnectStart(ctx, data)
	}
	if t.log != nil {
		ctx = t.log.TraceConnectStart(ctx, data)
	}
	return ctx
}

func (t *combinedTracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	if t.log != nil {
		t.log.TraceConnectEnd(ctx, data)
	}
	if t.otel != nil {
		t.otel.TraceConnectEnd(ctx, data)
	}
}

func (t *combinedTracer) TraceAcquireStart(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireStartData) context.Context {
	if t.log != nil {
		ctx = t.log.TraceAcquireStart(ctx, pool, data)
	}
	return ctx
}

func (t *combinedTracer) TraceAcquireEnd(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireEndData) {
	if t.log != nil {
		t.log.TraceAcquireEnd(ctx, pool, data)
	}
}

func (t *combinedTracer) TraceRelease(pool *pgxpool.Pool, data pgxpool.TraceReleaseData) {
	if t.log != nil {
		t.log.TraceRelease(pool, data)
	}
}
