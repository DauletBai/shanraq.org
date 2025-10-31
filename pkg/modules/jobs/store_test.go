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
	counts, err := store.CountByStatus(context.Background())
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
		"id", "name", "payload", "run_at", "attempts", "max_attempts", "status", "last_error", "created_at", "updated_at",
	}).
		AddRow(uuid.New(), "demo", []byte(`{"foo":"bar"}`), now, 1, 3, "done", nil, now, now)

	mock.ExpectQuery("SELECT id, name, payload, run_at, attempts, max_attempts, status, last_error, created_at, updated_at FROM job_queue").
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
