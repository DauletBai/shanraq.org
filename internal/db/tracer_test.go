package db

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// logSink returns a logger that keeps what was written, and the record of it.
func logSink(t *testing.T) (*zap.Logger, *observer.ObservedLogs) {
	t.Helper()
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(core), logs
}

// The rule the tracer exists to keep: query arguments never reach the log.
// They carry emails, contacts, article bodies and password hashes, and a log
// line is the easiest place in a system for those to end up somewhere they
// should not be.
func TestQueryArgumentsAreNeverLogged(t *testing.T) {
	logger, logs := logSink(t)
	tr := NewTracer(logger).(*combinedTracer)

	tr.log.Logger.Log(context.Background(), tracelog.LogLevelError, "query failed", map[string]any{
		"sql":  "SELECT * FROM auth_users WHERE email = $1",
		"args": []any{"someone@example.com", "$2b$12$hashthatmustnotappear"},
		"err":  "boom",
	})

	if logs.Len() != 1 {
		t.Fatalf("expected one line, got %d", logs.Len())
	}
	entry := logs.All()[0]
	for _, f := range entry.Context {
		if f.Key == "args" {
			t.Fatal("query arguments were written to the log")
		}
	}
	// The rest of the line still has to be useful, or the query cannot be found.
	whole := entry.Message + entry.ContextMap()["sql"].(string)
	if !strings.Contains(whole, "auth_users") {
		t.Error("the statement itself was dropped along with its arguments")
	}
	if _, ok := entry.ContextMap()["err"]; !ok {
		t.Error("the error was dropped, leaving a line that says nothing")
	}
	if got := entry.ContextMap()["component"]; got != "pgx" {
		t.Errorf("component = %v, want pgx", got)
	}
}

// A successful query at Info would be several lines per request. It is written
// at Debug on purpose, so production logs carry failures and not a transcript.
func TestOnlyFailuresRiseAboveDebug(t *testing.T) {
	logger, logs := logSink(t)
	tr := NewTracer(logger).(*combinedTracer)

	for _, c := range []struct {
		in   tracelog.LogLevel
		want zapcore.Level
	}{
		{tracelog.LogLevelTrace, zapcore.DebugLevel},
		{tracelog.LogLevelDebug, zapcore.DebugLevel},
		{tracelog.LogLevelInfo, zapcore.DebugLevel},
		{tracelog.LogLevelWarn, zapcore.WarnLevel},
		{tracelog.LogLevelError, zapcore.ErrorLevel},
	} {
		logs.TakeAll()
		tr.log.Logger.Log(context.Background(), c.in, "m", nil)
		all := logs.All()
		if len(all) != 1 {
			t.Fatalf("%v produced %d lines", c.in, len(all))
		}
		if all[0].Level != c.want {
			t.Errorf("pgx %v was logged at %v, want %v", c.in, all[0].Level, c.want)
		}
	}
	// And the tracer only asks pgx for the failures in the first place.
	if tr.log.LogLevel != tracelog.LogLevelError {
		t.Errorf("tracer subscribes at %v; successful queries should not be reported at all", tr.log.LogLevel)
	}
}

// Without a logger the tracer still has to be a tracer: telemetry alone is a
// valid configuration, and every hook is called on the hot path.
func TestATracerWithNoLoggerStillWorks(t *testing.T) {
	tr, ok := NewTracer(nil).(*combinedTracer)
	if !ok {
		t.Fatal("NewTracer(nil) did not return a tracer")
	}
	if tr.log != nil {
		t.Error("a nil logger should leave no log tracer to call")
	}
	if tr.otel == nil {
		t.Error("telemetry should be wired even with no logger")
	}
	// Each hook, in the order a query goes through them. A panic here is a
	// panic inside every database call the process makes.
	ctx := context.Background()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("a trace hook panicked with no logger configured: %v", r)
		}
	}()
	tr.TraceConnectEnd(ctx, pgx.TraceConnectEndData{})
	tr.TraceBatchEnd(ctx, nil, pgx.TraceBatchEndData{})
	tr.TraceCopyFromEnd(ctx, nil, pgx.TraceCopyFromEndData{})
	tr.TracePrepareEnd(ctx, nil, pgx.TracePrepareEndData{})
}
