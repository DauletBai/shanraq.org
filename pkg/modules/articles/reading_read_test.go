package articles

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A fifteen-minute article: 2700-odd words at the 180 a minute the page prints
// under the title.
const fifteenMin = 15 * 60

func TestReadCounts(t *testing.T) {
	cases := []struct {
		name   string
		depth  int
		secs   int
		expect int
		want   bool
	}{
		{"flicked to the bottom in four seconds", 100, 4, fifteenMin, false},
		{"read to the end in nine minutes", 100, 540, fifteenMin, true},
		{"nine minutes but stopped halfway", 50, 540, fifteenMin, false},
		{"exactly half the estimate, at the end", 100, fifteenMin / 2, fifteenMin, true},
		{"a second short of half", 100, fifteenMin/2 - 1, fifteenMin, false},
		{"no estimate to measure against", 100, 600, 0, false},
		{"a short article read quickly", 100, 40, 60, true},
	}
	for _, c := range cases {
		if got := readCounts(c.depth, c.secs, c.expect); got != c.want {
			t.Errorf("%s: readCounts(%d, %ds, %ds) = %v, want %v",
				c.name, c.depth, c.secs, c.expect, got, c.want)
		}
	}
}

// The estimate the threshold is measured against is the one the article shows
// its reader, so the two can never drift apart.
func TestReadingMinutesMatchesWhatThePageClaims(t *testing.T) {
	body := ""
	for i := 0; i < 2700; i++ {
		body += "слово "
	}
	if got, want := readingMinutes(body), 15; got != want {
		t.Errorf("readingMinutes(2700 words) = %d, want %d", got, want)
	}
}

// PostgreSQL returns SUM(bigint) as numeric. Dividing that value by 60 keeps a
// fractional scale, which pgx cannot scan into the integer shown by the admin
// dashboard. Exercise the real result type so this does not silently return to
// a warning and an empty reading total.
func TestReadTotalsConvertsSecondsToWholeMinutes(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, requireTestDB(t))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `
		CREATE TEMP TABLE article_reads (
			article_id uuid PRIMARY KEY,
			finished bigint NOT NULL,
			samples bigint NOT NULL,
			seconds bigint NOT NULL
		);
		INSERT INTO article_reads(article_id, finished, samples, seconds)
		VALUES ($1, 1, 2, 119)`, uuid.New()); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	defer func() { _, _ = conn.Exec(ctx, `DROP TABLE article_reads`) }()

	got, err := (&Store{db: conn}).ReadTotals(ctx)
	if err != nil {
		t.Fatalf("ReadTotals: %v", err)
	}
	want := (ReadTotals{Samples: 2, Finished: 1, Minutes: 1})
	if got != want {
		t.Fatalf("ReadTotals = %+v, want %+v", got, want)
	}
}
