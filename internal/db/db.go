package db

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"shanraq.org/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Connect builds a pgx pool with sane defaults and observability hooks.
func Connect(ctx context.Context, cfg config.DatabaseConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse db url: %w", err)
	}

	pcfg.MaxConns = cfg.MaxConns
	pcfg.MinConns = cfg.MinConns
	pcfg.MaxConnLifetime = cfg.MaxConnLifetime
	pcfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	pcfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pcfg.ConnConfig.Tracer = NewTracer(logger)

	deadlineCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var (
		pool    *pgxpool.Pool
		lastErr error
		backoff = 500 * time.Millisecond
	)

	for attempt := 1; attempt <= 3; attempt++ {
		attemptCtx, attemptCancel := context.WithTimeout(deadlineCtx, 10*time.Second)
		pool, lastErr = pgxpool.NewWithConfig(attemptCtx, pcfg)
		attemptCancel()

		if lastErr == nil {
			return pool, nil
		}

		if logger != nil {
			logger.Warn("postgres connect failed", zap.Error(lastErr), zap.Int("attempt", attempt))
		}

		if attempt == 3 || !retryable(lastErr) {
			break
		}

		select {
		case <-deadlineCtx.Done():
			return nil, fmt.Errorf("connect postgres: %w", deadlineCtx.Err())
		case <-time.After(backoff):
		}
		backoff *= 2
	}

	return nil, fmt.Errorf("connect postgres: %w", lastErr)
}

func retryable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "operation was canceled") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "temporary")
}
