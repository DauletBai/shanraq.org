package jobs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store manages DB operations for the job queue.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Enqueue(ctx context.Context, job Job) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO job_queue (id, name, payload, run_at, max_attempts)
		VALUES ($1, $2, $3, $4, $5)
	`, job.ID, job.Name, job.Payload, job.RunAt, job.MaxAttempts)
	if err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}
	return nil
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

func (s *Store) MarkRetry(ctx context.Context, id uuid.UUID, reason string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE job_queue
		SET status = 'retry',
		    last_error = $2,
		    run_at = NOW() + INTERVAL '15 seconds',
		    updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	if err != nil {
		return fmt.Errorf("mark retry: %w", err)
	}
	return nil
}

func (s *Store) MarkPending(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE job_queue
		SET status = 'pending',
		    run_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark pending: %w", err)
	}
	return nil
}

func (s *Store) Cancel(ctx context.Context, id uuid.UUID, reason string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE job_queue
		SET status = 'failed',
		    last_error = $2,
		    updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	if err != nil {
		return fmt.Errorf("cancel job: %w", err)
	}
	return nil
}

func (s *Store) ListRecent(ctx context.Context, limit int) ([]Job, error) {
	return s.List(ctx, ListOptions{Limit: limit})
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
	builder.WriteString("SELECT id, name, payload, run_at, attempts, max_attempts, status, last_error, created_at, updated_at FROM job_queue")

	args := make([]any, 0, 3)
	idx := 1
	if opts.Status != "" {
		builder.WriteString(" WHERE status = $")
		builder.WriteString(strconv.Itoa(idx))
		args = append(args, opts.Status)
		idx++
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

// CountByStatus aggregates jobs by status for dashboards.
func (s *Store) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := s.db.Query(ctx, `
		SELECT status, COUNT(*)
		FROM job_queue
		GROUP BY status
	`)
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

func scanJob(row pgx.Row) (Job, error) {
	var job Job
	var lastError *string
	if err := row.Scan(&job.ID, &job.Name, &job.Payload, &job.RunAt, &job.Attempts, &job.MaxAttempts, &job.Status, &lastError, &job.CreatedAt, &job.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return Job{}, err
		}
		return Job{}, fmt.Errorf("scan job: %w", err)
	}
	job.LastError = lastError
	return job, nil
}
