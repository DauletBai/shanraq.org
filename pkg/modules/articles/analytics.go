package articles

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// Guest analytics is deliberately AGGREGATE ONLY: every hit is folded into a
// per-day counter keyed by a coarse page kind (or a named click event) and a
// single guest/registered flag. No visitor identifier, IP, or session is ever
// stored, so this honours the Privacy Policy's "minimal analytics, no
// behavioral profiling" promise while still telling us how many people read
// what, and what they click.

// pageKind maps a request path to the coarse bucket shown in the dashboard, or
// "" for paths that should not be counted (forms handled elsewhere, assets,
// API, staff areas). Kept pure for easy testing.
func pageKind(path string) string {
	switch path {
	case "/":
		return "home"
	case "/read":
		return "article"
	case "/listings", "/listings/new", "/listings/my":
		return "listings"
	case "/favorites":
		return "favorites"
	case "/about", "/guide", "/formatting", "/pricing", "/support", "/privacy", "/terms":
		return "static"
	}
	switch {
	case strings.HasPrefix(path, "/read/"):
		return "article"
	case strings.HasPrefix(path, "/author/"):
		return "author"
	case strings.HasPrefix(path, "/agent/"):
		return "agent"
	case strings.HasPrefix(path, "/listings/"):
		return "listing"
	}
	return ""
}

// trackedEvents is the closed set of click events the beacon may record. A
// closed set keeps arbitrary strings out of the counter table.
var trackedEvents = map[string]bool{
	"show_contact":  true, // revealed a listing's phone/contact
	"register_cta":  true, // clicked a sign-up call to action
	"login_cta":     true, // clicked a sign-in call to action
	"post_listing":  true, // clicked "post a listing"
	"write_article": true, // clicked "write"
}

const (
	metricPage  = "page"
	metricClick = "click"
)

// metricKey identifies one counter bucket.
type metricKey struct {
	kind  string
	label string
	guest bool
}

// Metrics buffers analytics increments in memory and flushes them to
// analytics_daily on a ticker, so a page view costs a map write under a mutex
// rather than a DB round-trip on the request path.
type Metrics struct {
	db  *pgxpool.Pool
	log *zap.Logger
	mu  sync.Mutex
	buf map[metricKey]int64
}

// NewMetrics returns a collector bound to the pool.
func NewMetrics(db *pgxpool.Pool, log *zap.Logger) *Metrics {
	return &Metrics{db: db, log: log, buf: map[metricKey]int64{}}
}

// inc buffers one hit. A nil collector or empty label is a no-op, so callers
// never have to guard.
func (mt *Metrics) inc(kind, label string, guest bool) {
	if mt == nil || label == "" {
		return
	}
	mt.mu.Lock()
	mt.buf[metricKey{kind, label, guest}]++
	mt.mu.Unlock()
}

// Flush writes the buffered deltas to analytics_daily and clears the buffer.
// It is best-effort: on a DB error the batch is dropped (aggregate counters,
// not billing) rather than retried forever.
func (mt *Metrics) Flush(ctx context.Context) {
	if mt == nil {
		return
	}
	mt.mu.Lock()
	if len(mt.buf) == 0 {
		mt.mu.Unlock()
		return
	}
	batch := mt.buf
	mt.buf = map[metricKey]int64{}
	mt.mu.Unlock()

	b := &pgx.Batch{}
	for k, n := range batch {
		b.Queue(`
			INSERT INTO analytics_daily (day, kind, label, is_guest, n)
			VALUES (CURRENT_DATE, $1, $2, $3, $4)
			ON CONFLICT (day, kind, label, is_guest)
			DO UPDATE SET n = analytics_daily.n + EXCLUDED.n`,
			k.kind, k.label, k.guest, n)
	}
	res := mt.db.SendBatch(ctx, b)
	defer res.Close()
	for range batch {
		if _, err := res.Exec(); err != nil {
			mt.log.Warn("flush analytics", zap.Error(err))
			return
		}
	}
}

// trackTraffic counts one page view per GET of a countable page. It reads the
// (soft-loaded) session only to tell guests from signed-in users — nothing
// about who they are is recorded.
func (m *Module) trackTraffic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if kind := pageKind(r.URL.Path); kind != "" {
				_, ok := auth.ClaimsFromContext(r.Context())
				m.metrics.inc(metricPage, kind, !ok)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// handleTrack records a named click event sent via navigator.sendBeacon. It is
// fire-and-forget and always answers 204, even for unknown events.
func (m *Module) handleTrack(w http.ResponseWriter, r *http.Request) {
	e := strings.TrimSpace(r.URL.Query().Get("e"))
	if trackedEvents[e] {
		_, ok := auth.ClaimsFromContext(r.Context())
		m.metrics.inc(metricClick, e, !ok)
	}
	w.WriteHeader(http.StatusNoContent)
}

// metricsFlushLoop persists buffered counters periodically and once more on
// shutdown, so the last partial window is not lost.
func (m *Module) metricsFlushLoop(ctx context.Context) {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			// Detach from the cancelled request context for the final write.
			m.metrics.Flush(context.Background())
			return
		case <-t.C:
			m.metrics.Flush(ctx)
		}
	}
}

// ---- dashboard read side ----

// Audience splits a count between anonymous guests and signed-in users.
type Audience struct {
	Guest      int64
	Registered int64
}

// Total is the combined count.
func (a Audience) Total() int64 { return a.Guest + a.Registered }

// GuestPageRow is one page-kind's traffic over the reporting window.
type GuestPageRow struct {
	Kind  string
	Title string
	A     Audience
	Pct   int // guest share of the busiest kind, for the bar width
}

// GuestClickRow is one named click event's counts over the window.
type GuestClickRow struct {
	Name  string
	Title string
	A     Audience
}

// GuestTrendDay is one day in the 14-day guest-views sparkline.
type GuestTrendDay struct {
	Label string
	N     int64
	Pct   int
}

// GuestAnalytics is the aggregate audience block shown at the bottom of the
// admin dashboard.
type GuestAnalytics struct {
	Day       Audience // views today
	Week      Audience // views this calendar week (Mon–today)
	Month     Audience // views this calendar month
	Year      Audience // views this calendar year
	Pages     []GuestPageRow
	Clicks    []GuestClickRow
	Trend     []GuestTrendDay
	TrendFrom string // first day label in the trend (oldest)
	TrendTo   string // last day label in the trend (today)
	HasData   bool
}

// guestAnalytics assembles the audience dashboard. Titles are localized here so
// the template stays declarative. Errors on any part degrade to empty rather
// than failing the whole admin page.
func (m *Module) guestAnalytics(ctx context.Context, lang string) GuestAnalytics {
	var g GuestAnalytics
	db := m.rt.DB

	_ = db.QueryRow(ctx, `
		SELECT
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('week', CURRENT_DATE)::date AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('week', CURRENT_DATE)::date AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('month', CURRENT_DATE)::date AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('month', CURRENT_DATE)::date AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('year', CURRENT_DATE)::date AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= date_trunc('year', CURRENT_DATE)::date AND NOT is_guest), 0)
		FROM analytics_daily WHERE kind = 'page'`).
		Scan(&g.Day.Guest, &g.Day.Registered, &g.Week.Guest, &g.Week.Registered, &g.Month.Guest, &g.Month.Registered, &g.Year.Guest, &g.Year.Registered)

	// Pages, last 30 days.
	if rows, err := db.Query(ctx, `
		SELECT label,
		       COALESCE(SUM(n) FILTER (WHERE is_guest), 0),
		       COALESCE(SUM(n) FILTER (WHERE NOT is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'page' AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`); err == nil {
		var maxTotal int64
		for rows.Next() {
			var r GuestPageRow
			if err := rows.Scan(&r.Kind, &r.A.Guest, &r.A.Registered); err == nil {
				r.Title = T(lang, "ag.page."+r.Kind)
				if strings.HasPrefix(r.Title, "ag.page.") {
					r.Title = r.Kind
				}
				g.Pages = append(g.Pages, r)
				if r.A.Total() > maxTotal {
					maxTotal = r.A.Total()
				}
			}
		}
		rows.Close()
		sortPageRows(g.Pages)
		for i := range g.Pages {
			g.Pages[i].Pct = barPct(g.Pages[i].A.Total(), maxTotal)
		}
	} else {
		m.rt.Logger.Warn("guest analytics pages", zap.Error(err))
	}

	// Clicks, last 30 days.
	if rows, err := db.Query(ctx, `
		SELECT label,
		       COALESCE(SUM(n) FILTER (WHERE is_guest), 0),
		       COALESCE(SUM(n) FILTER (WHERE NOT is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'click' AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`); err == nil {
		for rows.Next() {
			var r GuestClickRow
			if err := rows.Scan(&r.Name, &r.A.Guest, &r.A.Registered); err == nil {
				r.Title = T(lang, "ag.click."+r.Name)
				if strings.HasPrefix(r.Title, "ag.click.") {
					r.Title = r.Name
				}
				g.Clicks = append(g.Clicks, r)
			}
		}
		rows.Close()
		sortClickRows(g.Clicks)
	} else {
		m.rt.Logger.Warn("guest analytics clicks", zap.Error(err))
	}

	// 14-day guest-views sparkline, gaps filled with zero.
	g.Trend = m.guestTrend(ctx)
	if len(g.Trend) > 0 {
		g.TrendFrom = g.Trend[0].Label
		g.TrendTo = g.Trend[len(g.Trend)-1].Label
	}

	g.HasData = g.Year.Total() > 0 || len(g.Pages) > 0
	return g
}

// guestTrend returns the last 14 calendar days of guest page views, oldest
// first, with missing days present as zero so the sparkline is continuous.
func (m *Module) guestTrend(ctx context.Context) []GuestTrendDay {
	counts := map[string]int64{}
	if rows, err := m.rt.DB.Query(ctx, `
		SELECT day, COALESCE(SUM(n) FILTER (WHERE is_guest), 0)
		FROM analytics_daily
		WHERE kind = 'page' AND day >= CURRENT_DATE - INTERVAL '13 days'
		GROUP BY day`); err == nil {
		for rows.Next() {
			var d time.Time
			var n int64
			if err := rows.Scan(&d, &n); err == nil {
				counts[d.Format("2006-01-02")] = n
			}
		}
		rows.Close()
	} else {
		m.rt.Logger.Warn("guest analytics trend", zap.Error(err))
	}

	today := time.Now()
	out := make([]GuestTrendDay, 0, 14)
	var max int64
	for i := 13; i >= 0; i-- {
		d := today.AddDate(0, 0, -i)
		n := counts[d.Format("2006-01-02")]
		if n > max {
			max = n
		}
		out = append(out, GuestTrendDay{Label: d.Format("02.01"), N: n})
	}
	for i := range out {
		out[i].Pct = barPct(out[i].N, max)
	}
	return out
}

// sortPageRows orders page kinds by total traffic, busiest first.
func sortPageRows(rs []GuestPageRow) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].A.Total() > rs[j-1].A.Total(); j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}

// sortClickRows orders click events by total count, most first.
func sortClickRows(rs []GuestClickRow) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].A.Total() > rs[j-1].A.Total(); j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}
