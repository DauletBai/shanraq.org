package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type poolWrapper struct {
	*pgxpool.Pool
}

func (p poolWrapper) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, arguments...)
}

func (p poolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

func (p poolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p poolWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.Pool.BeginTx(ctx, txOptions)
}

// Store manages DB operations for the job queue.
type Store struct {
	db pgxPool
}

func NewStore(db *pgxpool.Pool) *Store {
	return newStoreWithPool(poolWrapper{db})
}

func newStoreWithPool(db pgxPool) *Store {
	return &Store{db: db}
}

// MetricsSnapshot captures aggregate queue statistics.
type MetricsSnapshot struct {
	Total           int
	Pending         int
	Running         int
	Retry           int
	Failed          int
	Done            int
	DoneLastHour    int
	FailedLastHour  int
	NextScheduled   time.Time
	NextScheduledOk bool
}

func (s *Store) Enqueue(ctx context.Context, job Job) error {
	var userID any
	if job.UserID != uuid.Nil {
		userID = job.UserID
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO job_queue (id, user_id, name, payload, run_at, max_attempts)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, job.ID, userID, job.Name, job.Payload, job.RunAt, job.MaxAttempts)
	if err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}
	return nil
}

// Metrics returns aggregate queue statistics for dashboards.
func (s *Store) Metrics(ctx context.Context, userID *uuid.UUID) (MetricsSnapshot, error) {
	var snap MetricsSnapshot
	var next sql.NullTime
	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'pending') AS pending,
			COUNT(*) FILTER (WHERE status = 'running') AS running,
			COUNT(*) FILTER (WHERE status = 'retry') AS retry,
			COUNT(*) FILTER (WHERE status = 'failed') AS failed,
			COUNT(*) FILTER (WHERE status = 'done') AS done,
			COUNT(*) FILTER (WHERE status = 'done' AND updated_at >= NOW() - INTERVAL '1 hour') AS done_last_hour,
			COUNT(*) FILTER (WHERE status = 'failed' AND updated_at >= NOW() - INTERVAL '1 hour') AS failed_last_hour,
			MIN(run_at) FILTER (WHERE status IN ('pending', 'retry')) AS next_scheduled
		FROM job_queue
	`
	args := []any{}
	if userID != nil && *userID != uuid.Nil {
		query += " WHERE user_id = $1"
		args = append(args, *userID)
	}

	err := s.db.QueryRow(ctx, query, args...).Scan(
		&snap.Total,
		&snap.Pending,
		&snap.Running,
		&snap.Retry,
		&snap.Failed,
		&snap.Done,
		&snap.DoneLastHour,
		&snap.FailedLastHour,
		&next,
	)
	if err != nil {
		return MetricsSnapshot{}, fmt.Errorf("queue metrics: %w", err)
	}
	if next.Valid {
		snap.NextScheduled = next.Time
		snap.NextScheduledOk = true
	}
	return snap, nil
}

func (s *Store) ClaimNextJob(ctx context.Context) (Job, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Job{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	job, err := scanJob(tx.QueryRow(ctx, `
		SELECT id, name, payload, run_at, attempts, max_attempts, status, last_error, created_at, updated_at
		FROM job_queue
		WHERE status IN ('pending', 'retry')
		  AND run_at <= NOW()
		ORDER BY run_at
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`))
	if err != nil {
		if err == pgx.ErrNoRows {
			return Job{}, ErrNoJobs
		}
		return Job{}, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE job_queue
		SET status = 'running',
		    attempts = attempts + 1,
		    updated_at = NOW()
		WHERE id = $1
	`, job.ID)
	if err != nil {
		return Job{}, fmt.Errorf("mark running: %w", err)
	}
	job.Attempts++
	job.Status = "running"

	if err := tx.Commit(ctx); err != nil {
		return Job{}, fmt.Errorf("commit: %w", err)
	}
	return job, nil
}

func (s *Store) MarkDone(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE job_queue
		SET status = 'done',
		    last_error = NULL,
		    updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark done: %w", err)
	}
	return nil
}

func (s *Store) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE job_queue
		SET status = 'failed',
		    last_error = $2,
		    updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	if err != nil {
		return fmt.Errorf("mark failed: %w", err)
	}
	return nil
}

func (s *Store) MarkRetry(ctx context.Context, id uuid.UUID, reason string, userID *uuid.UUID) error {
	query := `
		UPDATE job_queue
		SET status = 'retry',
		    last_error = $2,
		    run_at = NOW() + INTERVAL '15 seconds',
		    updated_at = NOW()
		WHERE id = $1
	`
	args := []any{id, reason}
	query = addUserFilter(query, &args, userID)
	_, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("mark retry: %w", err)
	}
	return nil
}

func (s *Store) MarkPending(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	query := `
		UPDATE job_queue
		SET status = 'pending',
		    run_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`
	args := []any{id}
	query = addUserFilter(query, &args, userID)
	_, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("mark pending: %w", err)
	}
	return nil
}

func (s *Store) Cancel(ctx context.Context, id uuid.UUID, reason string, userID *uuid.UUID) error {
	query := `
		UPDATE job_queue
		SET status = 'failed',
		    last_error = $2,
		    updated_at = NOW()
		WHERE id = $1
	`
	args := []any{id, reason}
	query = addUserFilter(query, &args, userID)
	_, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cancel job: %w", err)
	}
	return nil
}

func (s *Store) ListRecent(ctx context.Context, limit int, userID *uuid.UUID) ([]Job, error) {
	return s.List(ctx, ListOptions{Limit: limit, UserID: userID})
}

func (s *Store) List(ctx context.Context, opts ListOptions) ([]Job, error) {
	limit := opts.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}

	var builder strings.Builder
	builder.WriteString("SELECT id, user_id, name, payload, run_at, attempts, max_attempts, status, last_error, created_at, updated_at FROM job_queue")

	args := make([]any, 0, 3)
	idx := 1
	conditions := make([]string, 0, 2)
	if opts.UserID != nil && *opts.UserID != uuid.Nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(idx))
		args = append(args, *opts.UserID)
		idx++
	}
	if opts.Status != "" {
		conditions = append(conditions, "status = $"+strconv.Itoa(idx))
		args = append(args, opts.Status)
		idx++
	}
	if len(conditions) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(conditions, " AND "))
	}

	builder.WriteString(" ORDER BY created_at DESC LIMIT $")
	builder.WriteString(strconv.Itoa(idx))
	args = append(args, limit)
	idx++

	builder.WriteString(" OFFSET $")
	builder.WriteString(strconv.Itoa(idx))
	args = append(args, offset)

	rows, err := s.db.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *Store) CountByStatus(ctx context.Context, userID *uuid.UUID) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*)
		FROM job_queue
	`
	args := []any{}
	if userID != nil && *userID != uuid.Nil {
		query += " WHERE user_id = $1"
		args = append(args, *userID)
	}
	query += " GROUP BY status"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("count jobs: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan count: %w", err)
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func addUserFilter(query string, args *[]any, userID *uuid.UUID) string {
	if userID != nil && *userID != uuid.Nil {
		idx := len(*args) + 1
		query += " AND user_id = $" + strconv.Itoa(idx)
		*args = append(*args, *userID)
	}
	return query
}

func scanJob(row pgx.Row) (Job, error) {
	var job Job
	var lastError *string
	var user pgtype.UUID
	if err := row.Scan(&job.ID, &user, &job.Name, &job.Payload, &job.RunAt, &job.Attempts, &job.MaxAttempts, &job.Status, &lastError, &job.CreatedAt, &job.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return Job{}, err
		}
		return Job{}, fmt.Errorf("scan job: %w", err)
	}
	if user.Valid {
		u, err := uuid.FromBytes(user.Bytes[:])
		if err == nil {
			job.UserID = u
		}
	}
	job.LastError = lastError
	return job, nil
}
