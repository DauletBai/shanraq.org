package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestStoreCountByStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock: %v", err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"status", "count"}).
		AddRow("pending", 2).
		AddRow("done", 3)

	mock.ExpectQuery("SELECT status, COUNT\\(\\*\\)\\s+FROM job_queue").
		WillReturnRows(rows)

	store := newStoreWithPool(mock)
	counts, err := store.CountByStatus(context.Background(), nil)
	if err != nil {
		t.Fatalf("count by status: %v", err)
	}

	if counts["pending"] != 2 || counts["done"] != 3 {
		t.Fatalf("unexpected counts: %#v", counts)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestStoreListRecent(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock: %v", err)
	}
	defer mock.Close()

	now := time.Now()
	rows := pgxmock.NewRows([]string{
		"id", "user_id", "name", "payload", "run_at", "attempts", "max_attempts", "status", "last_error", "created_at", "updated_at",
	}).
		AddRow(uuid.New(), uuid.New().String(), "demo", []byte(`{"foo":"bar"}`), now, 1, 3, "done", nil, now, now)

	mock.ExpectQuery("SELECT id, user_id, name, payload, run_at, attempts, max_attempts, status, last_error, created_at, updated_at FROM job_queue").
		WithArgs(8, 0).
		WillReturnRows(rows)

	store := newStoreWithPool(mock)
	jobs, err := store.List(context.Background(), ListOptions{Limit: 8})
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Name != "demo" {
		t.Fatalf("unexpected job name: %s", jobs[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestStoreMetrics(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock: %v", err)
	}
	defer mock.Close()

	next := time.Now().Add(5 * time.Minute)
	rows := pgxmock.NewRows([]string{
		"total", "pending", "running", "retry", "failed", "done",
		"done_last_hour", "failed_last_hour", "next_scheduled",
	}).AddRow(10, 3, 2, 1, 0, 4, 2, 0, next)

	mock.ExpectQuery("SELECT\\s+COUNT\\(\\*\\) AS total").
		WillReturnRows(rows)

	store := newStoreWithPool(mock)
	snap, err := store.Metrics(context.Background(), nil)
	if err != nil {
		t.Fatalf("metrics: %v", err)
	}
	if snap.Total != 10 || snap.Pending != 3 || snap.Done != 4 {
		t.Fatalf("unexpected snapshot: %#v", snap)
	}
	if !snap.NextScheduledOk || !snap.NextScheduled.Equal(next) {
		t.Fatalf("expected next scheduled %v, got %#v", next, snap)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
