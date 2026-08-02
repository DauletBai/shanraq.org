package articles

import (
	"context"
	"net/http"
	"net/url"
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
	metricPage   = "page"
	metricClick  = "click"
	metricBot    = "bot"    // crawler/preview-bot page hit, by family
	metricSource = "source" // where a human visit came from, by referrer
)

// botLabel classifies a User-Agent as a crawler/bot and returns its family
// ("google", "yandex", "facebook", …), or "" for an ordinary human browser.
// The UA string is read and immediately discarded — nothing is stored, so this
// keeps the aggregate-only, no-profiling promise. Kept pure for testing.
func botLabel(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.TrimSpace(u) == "":
		return "other" // no UA at all → a script/bot, never a real browser
	case strings.Contains(u, "googlebot"), strings.Contains(u, "google-inspectiontool"), strings.Contains(u, "storebot-google"), strings.Contains(u, "apis-google"):
		return "google"
	case strings.Contains(u, "yandex"):
		return "yandex"
	case strings.Contains(u, "bingbot"), strings.Contains(u, "bingpreview"), strings.Contains(u, "msnbot"):
		return "bing"
	case strings.Contains(u, "facebookexternalhit"), strings.Contains(u, "facebot"), strings.Contains(u, "meta-external"):
		return "facebook"
	case strings.Contains(u, "telegrambot"), strings.Contains(u, "telegram"):
		return "telegram"
	case strings.Contains(u, "whatsapp"):
		return "whatsapp"
	case strings.Contains(u, "twitterbot"):
		return "twitter"
	case strings.Contains(u, "applebot"):
		return "apple"
	case strings.Contains(u, "ahrefsbot"), strings.Contains(u, "semrushbot"), strings.Contains(u, "mj12bot"), strings.Contains(u, "dotbot"), strings.Contains(u, "petalbot"), strings.Contains(u, "dataforseo"), strings.Contains(u, "blexbot"):
		return "seo"
	case strings.Contains(u, "gptbot"), strings.Contains(u, "oai-searchbot"), strings.Contains(u, "chatgpt"), strings.Contains(u, "claudebot"), strings.Contains(u, "anthropic"), strings.Contains(u, "ccbot"), strings.Contains(u, "perplexity"), strings.Contains(u, "bytespider"), strings.Contains(u, "amazonbot"):
		return "ai"
	case strings.Contains(u, "bot"), strings.Contains(u, "crawler"), strings.Contains(u, "spider"), strings.Contains(u, "slurp"),
		strings.Contains(u, "curl"), strings.Contains(u, "wget"), strings.Contains(u, "python"), strings.Contains(u, "go-http"),
		strings.Contains(u, "java/"), strings.Contains(u, "okhttp"), strings.Contains(u, "headless"), strings.Contains(u, "phantom"):
		return "other"
	}
	return ""
}

// trafficSource classifies where a human visit came from, using only the
// Referer host (no query, no path, nothing stored). Empty and same-site
// referrers are "direct". Kept pure for testing.
func trafficSource(referer, host string) string {
	if strings.TrimSpace(referer) == "" {
		return "direct"
	}
	u, err := url.Parse(referer)
	if err != nil || u.Host == "" {
		return "direct"
	}
	h := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	host = strings.TrimPrefix(strings.ToLower(host), "www.")
	if h == host || strings.HasSuffix(h, "shanraq.org") {
		return "direct" // internal navigation, not an external source
	}
	switch {
	case strings.Contains(h, "google"):
		return "google"
	case strings.Contains(h, "yandex"), h == "ya.ru":
		return "yandex"
	case h == "t.me", strings.Contains(h, "telegram"):
		return "telegram"
	case strings.Contains(h, "facebook"), h == "fb.com", strings.Contains(h, "fb.me"), strings.Contains(h, "lm.facebook"):
		return "facebook"
	case strings.Contains(h, "instagram"):
		return "instagram"
	case strings.Contains(h, "linkedin"), h == "lnkd.in":
		return "linkedin"
	case strings.Contains(h, "twitter"), h == "x.com", h == "t.co":
		return "twitter"
	case strings.Contains(h, "youtube"), h == "youtu.be":
		return "youtube"
	case strings.Contains(h, "whatsapp"):
		return "whatsapp"
	case strings.Contains(h, "bing"):
		return "bing"
	case strings.Contains(h, "duckduckgo"):
		return "duckduckgo"
	}
	return "other"
}

// utmSource maps an explicit ?utm_source= tag to one of our known source
// labels, so a link shared with ?utm_source=telegram is attributed even when the
// browser strips the Referer (messengers, in-app browsers). Unknown values
// return "" so the caller falls back to referrer-based classification. The
// closed mapping keeps arbitrary strings out of the counter table.
func utmSource(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "instagram", "ig":
		return "instagram"
	case "telegram", "tg":
		return "telegram"
	case "facebook", "fb", "meta":
		return "facebook"
	case "youtube", "yt":
		return "youtube"
	case "google":
		return "google"
	case "yandex":
		return "yandex"
	case "twitter", "x":
		return "twitter"
	case "linkedin":
		return "linkedin"
	case "whatsapp", "wa":
		return "whatsapp"
	}
	return ""
}

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
				bot := botLabel(r.Header.Get("User-Agent"))
				switch {
				case bot == "seo":
					// Commercial SEO scanners are turned away in robots.txt and
					// excluded from analytics entirely, so they neither count as
					// guests nor clutter the bot panel.
				case bot != "":
					// Other crawlers are counted apart so they never inflate the
					// real human audience — that was the whole point.
					m.metrics.inc(metricBot, bot, true)
				default:
					_, ok := auth.ClaimsFromContext(r.Context())
					m.metrics.inc(metricPage, kind, !ok)
					// Prefer an explicit utm_source (survives the referrer being
					// stripped by messengers/apps); fall back to the Referer host.
					src := utmSource(r.URL.Query().Get("utm_source"))
					if src == "" {
						src = trafficSource(r.Header.Get("Referer"), r.Host)
					}
					m.metrics.inc(metricSource, src, !ok)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// handleTrack records a named click event sent via navigator.sendBeacon. It is
// fire-and-forget and always answers 204, even for unknown events.
func (m *Module) handleTrack(w http.ResponseWriter, r *http.Request) {
	e := strings.TrimSpace(r.URL.Query().Get("e"))
	if trackedEvents[e] && botLabel(r.Header.Get("User-Agent")) == "" {
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

// GuestSimpleRow is one labeled aggregate over the window — a traffic source or
// a bot family — carrying a single count (no guest/registered split).
type GuestSimpleRow struct {
	Name  string
	Title string
	N     int64
	Pct   int
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
	Week      Audience // views in the last 7 days
	Month     Audience // views in the last 30 days
	Year      Audience // views in the last 365 days
	Pages     []GuestPageRow
	Clicks    []GuestClickRow
	Sources   []GuestSimpleRow // human visits by referrer (30 days)
	Bots      []GuestSimpleRow // crawler hits by family (30 days)
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

	// Rolling windows (last N days), not calendar weeks/months, so the tiles
	// always nest: today ≤ 7 days ≤ 30 days ≤ 365 days. Calendar windows made the
	// current week exceed the month-to-date at the start of a month, which reads
	// as a bug even though it was arithmetically correct.
	_ = db.QueryRow(ctx, `
		SELECT
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day = CURRENT_DATE AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '6 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '6 days' AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '29 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '29 days' AND NOT is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '364 days' AND is_guest), 0),
		  COALESCE(SUM(n) FILTER (WHERE day >= CURRENT_DATE - INTERVAL '364 days' AND NOT is_guest), 0)
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

	// Traffic sources and bots, last 30 days.
	g.Sources = m.simpleRows(ctx, metricSource, "ag.source.", lang)
	g.Bots = m.simpleRows(ctx, metricBot, "ag.bot.", lang)

	// 14-day guest-views sparkline, gaps filled with zero.
	g.Trend = m.guestTrend(ctx)
	if len(g.Trend) > 0 {
		g.TrendFrom = g.Trend[0].Label
		g.TrendTo = g.Trend[len(g.Trend)-1].Label
	}

	g.HasData = g.Year.Total() > 0 || len(g.Pages) > 0 || len(g.Bots) > 0
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

// simpleRows loads a 30-day aggregate for a single-count metric kind (bots or
// sources), grouped by label, busiest first, with bar percentages and localized
// titles. Unknown labels fall back to the raw name.
func (m *Module) simpleRows(ctx context.Context, kind, i18nPrefix, lang string) []GuestSimpleRow {
	rows, err := m.rt.DB.Query(ctx, `
		SELECT label, COALESCE(SUM(n), 0)
		FROM analytics_daily
		WHERE kind = $1 AND day >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY label`, kind)
	if err != nil {
		m.rt.Logger.Warn("guest analytics "+kind, zap.Error(err))
		return nil
	}
	defer rows.Close()

	var out []GuestSimpleRow
	var max int64
	for rows.Next() {
		var r GuestSimpleRow
		if err := rows.Scan(&r.Name, &r.N); err != nil {
			continue
		}
		r.Title = T(lang, i18nPrefix+r.Name)
		if strings.HasPrefix(r.Title, i18nPrefix) {
			r.Title = r.Name
		}
		out = append(out, r)
		if r.N > max {
			max = r.N
		}
	}
	for i := 1; i < len(out); i++ { // busiest first (tiny list, insertion sort)
		for j := i; j > 0 && out[j].N > out[j-1].N; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
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
