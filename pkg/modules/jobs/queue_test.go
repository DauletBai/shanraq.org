package jobs

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The queue's promises are all about what two workers see at once and what a
// failure leaves behind, so these run against a real database. SKIP LOCKED
// cannot be demonstrated to a mock.

type qFixture struct {
	pool  *pgxpool.Pool
	store *Store
	ctx   context.Context
}

func newQueue(t *testing.T) *qFixture {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the queue tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	// A run starts from an empty queue, or a leftover job from elsewhere would
	// be claimed instead of the one under test.
	if _, err := pool.Exec(ctx, `DELETE FROM job_queue`); err != nil {
		t.Fatalf("clear queue: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM job_queue`) })
	return &qFixture{pool: pool, store: NewStore(pool), ctx: ctx}
}

func (f *qFixture) put(t *testing.T, name string, runAt time.Time) uuid.UUID {
	t.Helper()
	j := Job{ID: uuid.New(), Name: name, Payload: json.RawMessage(`{"n":1}`),
		RunAt: runAt, MaxAttempts: 3, Status: "pending"}
	if err := f.store.Enqueue(f.ctx, j); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	return j.ID
}

func (f *qFixture) status(t *testing.T, id uuid.UUID) (string, int) {
	t.Helper()
	var s string
	var attempts int
	if err := f.pool.QueryRow(f.ctx,
		`SELECT status::text, attempts FROM job_queue WHERE id=$1`, id).Scan(&s, &attempts); err != nil {
		t.Fatalf("read job: %v", err)
	}
	return s, attempts
}

func TestAClaimedJobIsRunningAndCountsTheAttempt(t *testing.T) {
	f := newQueue(t)
	id := f.put(t, "send", time.Now().Add(-time.Minute))

	job, err := f.store.ClaimNextJob(f.ctx)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if job.ID != id {
		t.Fatalf("claimed %s, want %s", job.ID, id)
	}
	if job.Status != "running" || job.Attempts != 1 {
		t.Errorf("claimed job is %s on attempt %d; want running on 1", job.Status, job.Attempts)
	}
	if s, a := f.status(t, id); s != "running" || a != 1 {
		t.Errorf("the row says %s on attempt %d; want running on 1", s, a)
	}
}

// The guarantee the whole worker pool rests on. Ten workers claiming at once
// must between them take each job exactly once; a job handed out twice is an
// email sent twice, or a payment taken twice.
func TestNoJobIsClaimedTwice(t *testing.T) {
	f := newQueue(t)
	const jobs, workers = 12, 10
	for i := 0; i < jobs; i++ {
		f.put(t, "send", time.Now().Add(-time.Minute))
	}

	var mu sync.Mutex
	seen := map[uuid.UUID]int{}
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				job, err := f.store.ClaimNextJob(context.Background())
				if err != nil {
					return // ErrNoJobs, or the queue is drained
				}
				mu.Lock()
				seen[job.ID]++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(seen) != jobs {
		t.Errorf("%d of %d jobs were claimed", len(seen), jobs)
	}
	for id, n := range seen {
		if n != 1 {
			t.Errorf("job %s was handed to %d workers", id, n)
		}
	}
}

// A job whose time has not come is not work yet.
func TestAJobScheduledForLaterIsNotClaimed(t *testing.T) {
	f := newQueue(t)
	f.put(t, "later", time.Now().Add(time.Hour))
	if _, err := f.store.ClaimNextJob(f.ctx); err != ErrNoJobs {
		t.Errorf("a job due in an hour was claimed now (err %v)", err)
	}
}

func TestAnEmptyQueueSaysSoRatherThanFailing(t *testing.T) {
	f := newQueue(t)
	if _, err := f.store.ClaimNextJob(f.ctx); err != ErrNoJobs {
		t.Errorf("claiming from an empty queue gave %v, want ErrNoJobs", err)
	}
}

// Retry returns work to the queue; done and failed take it out. The difference
// is whether a second worker will ever see it again.
func TestWhatEachEndingLeavesBehind(t *testing.T) {
	f := newQueue(t)

	retried := f.put(t, "flaky", time.Now().Add(-time.Minute))
	if _, err := f.store.ClaimNextJob(f.ctx); err != nil {
		t.Fatal(err)
	}
	if err := f.store.MarkRetry(f.ctx, retried, "upstream timeout", nil); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if s, _ := f.status(t, retried); s != "retry" {
		t.Errorf("after MarkRetry the job is %s, want retry", s)
	}
	// And it is claimable again, which is what retry is for. run_at may have
	// been pushed out, so it is brought back to now first.
	if _, err := f.pool.Exec(f.ctx, `UPDATE job_queue SET run_at = NOW() - INTERVAL '1 minute' WHERE id=$1`, retried); err != nil {
		t.Fatal(err)
	}
	again, err := f.store.ClaimNextJob(f.ctx)
	if err != nil || again.ID != retried {
		t.Fatalf("a retried job was not offered again: %v", err)
	}
	if again.Attempts != 2 {
		t.Errorf("the second run is attempt %d, want 2", again.Attempts)
	}

	done := f.put(t, "ok", time.Now().Add(-time.Minute))
	if err := f.store.MarkDone(f.ctx, done); err != nil {
		t.Fatalf("done: %v", err)
	}
	if s, _ := f.status(t, done); s != "done" {
		t.Errorf("after MarkDone the job is %s, want done", s)
	}

	failed := f.put(t, "broken", time.Now().Add(-time.Minute))
	if err := f.store.MarkFailed(f.ctx, failed, "handler missing"); err != nil {
		t.Fatalf("failed: %v", err)
	}
	if s, _ := f.status(t, failed); s != "failed" {
		t.Errorf("after MarkFailed the job is %s, want failed", s)
	}

	// Neither of the finished ones comes back.
	for {
		job, err := f.store.ClaimNextJob(f.ctx)
		if err != nil {
			break
		}
		if job.ID == done || job.ID == failed {
			t.Errorf("a finished job (%s) was claimed again", job.Status)
		}
	}
}

func TestCancellingTakesAJobOutOfTheQueue(t *testing.T) {
	f := newQueue(t)
	id := f.put(t, "unwanted", time.Now().Add(-time.Minute))
	if err := f.store.Cancel(f.ctx, id, "changed my mind", nil); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if _, err := f.store.ClaimNextJob(f.ctx); err != ErrNoJobs {
		t.Error("a cancelled job was still offered to a worker")
	}
}

// The payload survives the round trip, which is the only reason to enqueue
// anything at all.
func TestThePayloadComesBackAsItWentIn(t *testing.T) {
	f := newQueue(t)
	want := map[string]any{"to": "someone@example.com", "tries": float64(3)}
	body, _ := json.Marshal(want)
	j := Job{ID: uuid.New(), Name: "mail", Payload: body,
		RunAt: time.Now().Add(-time.Minute), MaxAttempts: 3, Status: "pending"}
	if err := f.store.Enqueue(f.ctx, j); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	claimed, err := f.store.ClaimNextJob(f.ctx)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	var got map[string]any
	if err := claimed.Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["to"] != want["to"] || got["tries"] != want["tries"] {
		t.Errorf("payload came back as %v, want %v", got, want)
	}
}
