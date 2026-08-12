package articles

import (
	"context"
	"testing"
)

// TestSourceResetSparesEveryOtherMetric pins the shape of the clean-up: the
// traffic-source counters were the only ones ever computed wrongly, so they are
// the only ones removed.
//
// The temptation was to wipe the whole table and start everything from one
// date. That would have destroyed the audience breakdown, the bot statistics
// and the only record of how the site launched — nine metrics that were always
// correct — in exchange for a matching date. Crawlers were never counted as
// guests, and that separation is what produced the honest 810 real article
// views against a counter reading 2 284.
func TestSourceResetSparesEveryOtherMetric(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()

	// Every kind the tracker can write, one row each.
	kinds := []string{
		metricPage, metricClick, metricBot, metricSource, metricDevice,
		metricOS, metricBrowser, metricCountry, metricLang, metricGeoLang,
	}
	// The legacy table has no unique key, so the archive INSERT below appends
	// rather than upserting. Without this the probe rows survive the run and the
	// next one against the same database reads 14 where it expects 7 — the test
	// passes once on a fresh database and fails ever after.
	t.Cleanup(func() {
		app.exec(`DELETE FROM analytics_daily WHERE label = 'probe'`)
		app.exec(`DELETE FROM analytics_daily_legacy WHERE label = 'probe'`)
	})
	for _, k := range kinds {
		app.exec(`INSERT INTO analytics_daily (day, kind, label, is_guest, n)
		          VALUES (CURRENT_DATE - 5, $1, 'probe', TRUE, 7)
		          ON CONFLICT (day, kind, label, is_guest) DO UPDATE SET n = 7`, k)
	}

	// Re-run the migration's effect. The migration itself already ran against
	// this database; this repeats it on the rows the test just inserted.
	app.exec(`INSERT INTO analytics_daily_legacy
	          SELECT * FROM analytics_daily WHERE kind = 'source' ON CONFLICT DO NOTHING`)
	app.exec(`DELETE FROM analytics_daily WHERE kind = 'source'`)

	count := func(kind string) int64 {
		t.Helper()
		var n int64
		if err := app.pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(n),0) FROM analytics_daily WHERE kind = $1 AND label = 'probe'`,
			kind).Scan(&n); err != nil {
			t.Fatalf("count %s: %v", kind, err)
		}
		return n
	}

	if got := count(metricSource); got != 0 {
		t.Errorf("source rows survived the reset (%d) — they are the contaminated ones", got)
	}
	for _, k := range kinds {
		if k == metricSource {
			continue
		}
		if got := count(k); got != 7 {
			t.Errorf("metric %q lost data (%d, want 7) — it was never wrong and must be kept", k, got)
		}
	}

	// Nothing is destroyed: the removed rows are still readable.
	var archived int64
	if err := app.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(n),0) FROM analytics_daily_legacy WHERE kind = 'source' AND label = 'probe'`,
	).Scan(&archived); err != nil {
		t.Fatalf("read legacy: %v", err)
	}
	if archived != 7 {
		t.Errorf("archived source rows = %d, want 7 — the reset must preserve, not delete", archived)
	}
}
