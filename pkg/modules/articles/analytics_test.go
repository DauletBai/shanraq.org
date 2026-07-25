package articles

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
)

func TestPageKind(t *testing.T) {
	cases := map[string]string{
		"/":                   "home",
		"/read":               "article",
		"/read/some-slug":     "article",
		"/listings":           "listings",
		"/listings/new":       "listings",
		"/listings/my":        "listings",
		"/listings/abc-123":   "listing",
		"/author/42":          "author",
		"/agent/7":            "agent",
		"/favorites":          "favorites",
		"/about":              "static",
		"/privacy":            "static",
		"/terms":              "static",
		"/api/geo/roots":      "",        // not counted
		"/static/css/x.css":   "",        // not counted
		"/studio/new":         "",        // staff area, not a public page
		"/admin":              "",        // not counted
		"/read/slug/progress": "article", // still under /read/ prefix (POST filtered by method upstream)
	}
	for path, want := range cases {
		if got := pageKind(path); got != want {
			t.Errorf("pageKind(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestTrackedEventsClosedSet(t *testing.T) {
	if !trackedEvents["show_contact"] {
		t.Error("show_contact should be a tracked event")
	}
	if trackedEvents["arbitrary_injected"] {
		t.Error("unknown events must not be tracked")
	}
}

func TestAudienceTotal(t *testing.T) {
	a := Audience{Guest: 7, Registered: 3}
	if a.Total() != 10 {
		t.Errorf("Total() = %d, want 10", a.Total())
	}
	if (Audience{}).Total() != 0 {
		t.Error("empty Audience total should be 0")
	}
}

func TestSortPageRowsByTotalDesc(t *testing.T) {
	rows := []GuestPageRow{
		{Kind: "a", A: Audience{Guest: 1, Registered: 1}}, // 2
		{Kind: "b", A: Audience{Guest: 5, Registered: 5}}, // 10
		{Kind: "c", A: Audience{Guest: 4, Registered: 0}}, // 4
	}
	sortPageRows(rows)
	if rows[0].Kind != "b" || rows[1].Kind != "c" || rows[2].Kind != "a" {
		t.Errorf("page rows not sorted by total desc: %v", []string{rows[0].Kind, rows[1].Kind, rows[2].Kind})
	}
}

func TestSortClickRowsByTotalDesc(t *testing.T) {
	rows := []GuestClickRow{
		{Name: "x", A: Audience{Guest: 2}},
		{Name: "y", A: Audience{Guest: 9}},
		{Name: "z", A: Audience{Guest: 5}},
	}
	sortClickRows(rows)
	if rows[0].Name != "y" || rows[1].Name != "z" || rows[2].Name != "x" {
		t.Errorf("click rows not sorted by total desc: %v", []string{rows[0].Name, rows[1].Name, rows[2].Name})
	}
}

func TestMetricsIncBuffering(t *testing.T) {
	// nil collector and empty label must be safe no-ops.
	var nilM *Metrics
	nilM.inc(metricPage, "home", true)

	m := NewMetrics(nil, zap.NewNop())
	m.inc(metricPage, "", true) // empty label ignored
	m.inc(metricPage, "home", true)
	m.inc(metricPage, "home", true)
	m.inc(metricPage, "home", false) // different audience → different bucket
	m.inc(metricClick, "login_cta", true)

	if got := m.buf[metricKey{metricPage, "home", true}]; got != 2 {
		t.Errorf("guest home = %d, want 2", got)
	}
	if got := m.buf[metricKey{metricPage, "home", false}]; got != 1 {
		t.Errorf("registered home = %d, want 1", got)
	}
	if got := m.buf[metricKey{metricClick, "login_cta", true}]; got != 1 {
		t.Errorf("click login_cta = %d, want 1", got)
	}
	if _, ok := m.buf[metricKey{metricPage, "", true}]; ok {
		t.Error("empty label must not be buffered")
	}
}

// TestMetricsFlushAndGuestAnalytics exercises the full write→read path against
// the test DB: buffered counters flush into analytics_daily, then guestAnalytics
// rolls them up (including the new "week" bucket) for the dashboard.
func TestMetricsFlushAndGuestAnalytics(t *testing.T) {
	dsn := requireTestDB(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	// Dedicated table; safe to reset in the test DB for a deterministic assert.
	if _, err := pool.Exec(ctx, `TRUNCATE analytics_daily`); err != nil {
		t.Fatalf("truncate: %v", err)
	}

	m := NewMetrics(pool, zap.NewNop())
	// 3 guest home views, 1 registered home view, 2 guest article views, 1 click.
	for i := 0; i < 3; i++ {
		m.inc(metricPage, "home", true)
	}
	m.inc(metricPage, "home", false)
	m.inc(metricPage, "article", true)
	m.inc(metricPage, "article", true)
	m.inc(metricClick, "show_contact", true)
	m.Flush(ctx)

	// A second flush of nothing must be a harmless no-op.
	m.Flush(ctx)

	// Idempotent accumulation: another guest home view adds to the same row.
	m.inc(metricPage, "home", true)
	m.Flush(ctx)

	var homeGuest int64
	if err := pool.QueryRow(ctx,
		`SELECT n FROM analytics_daily WHERE kind='page' AND label='home' AND is_guest=true AND day=CURRENT_DATE`).
		Scan(&homeGuest); err != nil {
		t.Fatalf("read home guest: %v", err)
	}
	if homeGuest != 4 {
		t.Errorf("home guest views = %d, want 4 (3 + 1 across two flushes)", homeGuest)
	}

	mod := &Module{rt: &shanraq.Runtime{DB: pool, Logger: zap.NewNop()}}
	g := mod.guestAnalytics(ctx, LangRU)

	if !g.HasData {
		t.Fatal("HasData should be true after recording views")
	}
	// Today's totals: pages = 4 home + 2 article + 1 registered home = 7.
	if g.Day.Guest != 6 || g.Day.Registered != 1 {
		t.Errorf("Day = %+v, want guest 6 / registered 1", g.Day)
	}
	// Week/month/year include today, so they are at least today's totals. Since
	// we truncated, they equal today exactly here.
	if g.Week.Total() != 7 || g.Month.Total() != 7 || g.Year.Total() != 7 {
		t.Errorf("roll-ups week=%d month=%d year=%d, want 7 each", g.Week.Total(), g.Month.Total(), g.Year.Total())
	}
	// Pages sorted by total: home (4+1=5) before article (2).
	if len(g.Pages) != 2 || g.Pages[0].Kind != "home" || g.Pages[1].Kind != "article" {
		t.Errorf("pages = %+v, want [home, article]", g.Pages)
	}
	// Page titles are localized (not the raw kind).
	if g.Pages[0].Title != T(LangRU, "ag.page.home") {
		t.Errorf("home title = %q, want localized", g.Pages[0].Title)
	}
	// One click event recorded.
	if len(g.Clicks) != 1 || g.Clicks[0].Name != "show_contact" || g.Clicks[0].A.Guest != 1 {
		t.Errorf("clicks = %+v, want one show_contact guest=1", g.Clicks)
	}
	// Trend is always a continuous 14-day window, today last.
	if len(g.Trend) != 14 {
		t.Errorf("trend length = %d, want 14", len(g.Trend))
	}
	if g.TrendTo == "" || g.TrendFrom == "" {
		t.Error("trend endpoints should be set")
	}
	if g.Trend[13].N != 6 { // today's guest page views
		t.Errorf("trend today guest = %d, want 6", g.Trend[13].N)
	}
}
