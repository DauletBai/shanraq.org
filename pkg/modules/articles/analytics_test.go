package articles

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

func TestBotLabel(t *testing.T) {
	bots := map[string]string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)": "google",
		"Mozilla/5.0 (compatible; YandexBot/3.0)":                                  "yandex",
		"Mozilla/5.0 (compatible; bingbot/2.0)":                                    "bing",
		"facebookexternalhit/1.1":                                                  "facebook",
		"TelegramBot (like TwitterBot)":                                            "telegram",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0)":                                  "seo",
		"Mozilla/5.0 (compatible; GPTBot/1.0)":                                     "ai",
		"ClaudeBot/1.0":                                                            "ai",
		"curl/8.4.0":                                                               "other",
		"python-requests/2.31":                                                     "other",
		"":                                                                         "other",
	}
	for ua, want := range bots {
		if got := botLabel(ua); got != want {
			t.Errorf("botLabel(%q) = %q, want %q", ua, got, want)
		}
	}
	// Real browsers must never be flagged as bots.
	humans := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Version/17.0 Mobile/15E148 Safari/604.1",
	}
	for _, ua := range humans {
		if got := botLabel(ua); got != "" {
			t.Errorf("botLabel(%q) = %q, want human (empty)", ua, got)
		}
	}
}

func TestTrafficSource(t *testing.T) {
	const host = "shanraq.org"
	cases := map[string]string{
		"":                                  "direct",
		"https://shanraq.org/read/x":        "direct", // internal navigation
		"https://www.shanraq.org/":          "direct",
		"https://www.google.com/search?q=x": "google",
		"https://yandex.kz/":                "yandex",
		"https://t.me/shanraq_org":          "telegram",
		"https://l.facebook.com/l.php?u=x":  "facebook",
		"https://lnkd.in/abc":               "linkedin",
		"https://x.com/user":                "twitter",
		"https://example.com/blog":          "other",
	}
	for ref, want := range cases {
		if got := trafficSource(ref, host); got != want {
			t.Errorf("trafficSource(%q) = %q, want %q", ref, got, want)
		}
	}
}

func TestPageKind(t *testing.T) {
	cases := map[string]string{
		"/":                 "home",
		"/read":             "article",
		"/read/some-slug":   "article",
		"/listings":         "listings",
		"/listings/new":     "listings",
		"/listings/my":      "listings",
		"/listings/abc-123": "listing",
		"/author/42":        "author",
		"/agent/7":          "agent",
		"/favorites":        "favorites",
		// Each standing page counts under its own name: a shared "static"
		// bucket said how many people read something and never which.
		"/about":     "about",
		"/privacy":   "privacy",
		"/terms":     "terms",
		"/adam":      "adam",
		"/framework": "framework",
		// The forecast pages are most of the site and used to count for
		// nothing at all, taking the visitor behind them with it.
		"/weather":            "weather",
		"/weather/karaotkel":  "weather",
		"/rates":              "rates",
		"/archive/2026-08-29": "archive",
		"/place/kachar":       "place",
		"/predictions":        "predictions",
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
	// Trend is a continuous window ending today, with no gaps: a day nobody
	// visited has to appear as a zero-height bar, not disappear and shift the
	// dates under every other bar.
	if len(g.Trend) != guestTrendDays {
		t.Errorf("trend length = %d, want %d", len(g.Trend), guestTrendDays)
	}
	if g.TrendTo == "" || g.TrendFrom == "" {
		t.Error("trend endpoints should be set")
	}
	if last := g.Trend[len(g.Trend)-1]; last.N != 6 { // today's guest page views
		t.Errorf("trend today guest = %d, want 6", last.N)
	}
	// The last day is always labelled, and the axis carries a handful of dates
	// rather than one under every bar.
	if !g.Trend[len(g.Trend)-1].Tick {
		t.Error("the most recent day is not labelled on the axis")
	}
	ticks := 0
	for _, d := range g.Trend {
		if d.Tick {
			ticks++
		}
	}
	if ticks < 4 || ticks > 7 {
		t.Errorf("dates on the axis = %d, want about %d", ticks, guestTrendLabels)
	}
}

// TestCountryFlagEmoji covers the derivation, not a table: every country the
// geoip can report has to get a flag without anyone adding it by hand, which is
// the entire reason the code derives one from the ISO letters.
func TestCountryFlagEmoji(t *testing.T) {
	cases := map[string]string{
		"KZ":         "\U0001F1F0\U0001F1FF",
		"NL":         "\U0001F1F3\U0001F1F1",
		"US":         "\U0001F1FA\U0001F1F8",
		"gb":         "\U0001F1EC\U0001F1E7", // case-insensitive
		"datacenter": "☁️",                   // hosting is not a place
		"":           "",
		"K":          "",
		"KAZ":        "",
		"K1":         "", // digits are not regional indicators
	}
	for in, want := range cases {
		if got := countryFlagEmoji(in); got != want {
			t.Errorf("countryFlagEmoji(%q) = %q, want %q", in, got, want)
		}
	}

	// Every label the live site has actually recorded must render something, or
	// the panel shows a ragged column of half-flagged rows.
	for _, code := range []string{"AU", "CN", "ES", "GB", "HU", "ID", "KZ", "NL", "RU", "SE", "US", "datacenter"} {
		if countryFlagEmoji(code) == "" {
			t.Errorf("no flag for %q, which the live analytics does report", code)
		}
		if name := T(LangRU, "ag.country."+code); strings.HasPrefix(name, "ag.country.") {
			t.Errorf("no Russian name for %q — it would render as a bare ISO code", code)
		}
	}
}

// TestGoogleIsNotOneThing pins the split that mattered: the panel used to file
// every host containing "google" as organic search, so a link opened from Gmail
// counted as SEO working. It reported 20 Google visits in a month while Search
// Console counted one click in three.
func TestGoogleIsNotOneThing(t *testing.T) {
	cases := map[string]string{
		"https://www.google.com/":             "google",
		"https://google.kz/search?q=x":        "google",
		"https://google.com.tr/":              "google",
		"https://news.google.com/":            "google",
		"https://mail.google.com/mail/u/0/":   "email",
		"https://translate.google.com/":       "translate",
		"https://shanraq-org.translate.goog/": "translate",
		"https://docs.google.com/document/d":  "other",
		"https://drive.google.com/file/d/1":   "other",
		"https://lh3.googleusercontent.com/":  "other",
	}
	for ref, want := range cases {
		if got := trafficSource(ref, "shanraq.org"); got != want {
			t.Errorf("trafficSource(%q) = %q, want %q", ref, got, want)
		}
	}
	// Every bucket the classifier can emit needs a label, or the panel prints
	// the raw key at the reader.
	for _, src := range []string{"google", "email", "translate", "share", "other"} {
		if name := T(LangRU, "ag.source."+src); strings.HasPrefix(name, "ag.source.") {
			t.Errorf("no Russian label for source %q", src)
		}
	}
}

// TestSitemapListsEveryLanguage pins Google's sitemap requirement for
// multilingual sites: a separate <url> element per language version, each
// listing every alternate including itself. We used to emit only the Russian
// <loc>, so the Kazakh and English versions of every article were never
// actually submitted for crawling — two thirds of a trilingual catalogue.
func TestSitemapListsEveryLanguage(t *testing.T) {
	mod := &Module{rt: &shanraq.Runtime{Config: config.Config{PublicBaseURL: "https://shanraq.org"}}}
	doc := mod.sitemapDoc(func(emit func(string, time.Time)) {
		emit("/read/example", time.Time{})
	})
	x := string(doc)

	for _, lang := range Langs {
		loc := "<loc>https://shanraq.org/read/example?lang=" + lang + "</loc>"
		if !strings.Contains(x, loc) {
			t.Errorf("sitemap has no <loc> for %s — that version is never submitted", lang)
		}
	}
	if n := strings.Count(x, "<url>"); n != len(Langs) {
		t.Errorf("sitemap emitted %d <url> elements for one page, want %d (one per language)", n, len(Langs))
	}
	// Reciprocity: each entry must name all three languages plus x-default, or
	// Google is entitled to ignore the annotations entirely.
	if n := strings.Count(x, `hreflang="kk"`); n != len(Langs) {
		t.Errorf("kk alternate appears %d times, want %d (once in each entry)", n, len(Langs))
	}
	if n := strings.Count(x, `hreflang="x-default"`); n != len(Langs) {
		t.Errorf("x-default appears %d times, want %d", n, len(Langs))
	}
}

// The Y axis is built from the data: a step people count in, and a top that is
// the first multiple of it at or above the busiest day.
func TestAxisTicks(t *testing.T) {
	cases := []struct {
		peak int64
		want []int64
	}{
		{414, []int64{500, 400, 300, 200, 100, 0}},
		{291, []int64{300, 200, 100, 0}},
		{41, []int64{50, 40, 30, 20, 10, 0}},
		{7, []int64{8, 6, 4, 2, 0}},
		{1040, []int64{1250, 1000, 750, 500, 250, 0}},
	}
	for _, c := range cases {
		got := axisTicks(c.peak)
		var vals []int64
		for _, tk := range got {
			vals = append(vals, tk.N)
		}
		if len(vals) != len(c.want) {
			t.Errorf("axisTicks(%d) = %v, want %v", c.peak, vals, c.want)
			continue
		}
		for i := range c.want {
			if vals[i] != c.want[i] {
				t.Errorf("axisTicks(%d) = %v, want %v", c.peak, vals, c.want)
				break
			}
		}
	}
}

// Whatever the data, the axis must contain it and stay readable.
func TestAxisTicksAlwaysCoverThePeak(t *testing.T) {
	for peak := int64(0); peak < 5000; peak += 13 {
		ticks := axisTicks(peak)
		if len(ticks) < 2 {
			t.Fatalf("peak %d produced %d ticks", peak, len(ticks))
		}
		if ticks[0].N < peak {
			t.Fatalf("peak %d overflows the axis top %d", peak, ticks[0].N)
		}
		if ticks[len(ticks)-1].N != 0 {
			t.Fatalf("peak %d: axis does not reach zero (%v)", peak, ticks[len(ticks)-1].N)
		}
		if ticks[0].Pct != 100 || ticks[len(ticks)-1].Pct != 0 {
			t.Fatalf("peak %d: ends at %d%%..%d%%", peak, ticks[len(ticks)-1].Pct, ticks[0].Pct)
		}
		if len(ticks) > 9 {
			t.Fatalf("peak %d: %d gridlines is too crowded", peak, len(ticks))
		}
	}
}

// Headroom above the peak must stay modest, or the chart reads as half empty.
func TestAxisTopHugsTheData(t *testing.T) {
	for peak := int64(20); peak < 5000; peak += 7 {
		top := axisTicks(peak)[0].N
		if float64(top) > float64(peak)*1.5 {
			t.Fatalf("peak %d got an axis top of %d — too much empty space", peak, top)
		}
	}
}

// The dropped counter has one job: to be right about how much was lost while
// the database was unreachable. It counted rows in the batch — including ones
// already handed back — instead of the tallies still ahead of it.
func TestGiveBackCountsLostTallies(t *testing.T) {
	mt := &Metrics{log: zap.NewNop(), buf: map[metricKey]int64{}, slots: map[slotKey]int64{}}
	for i := 0; i < metricsBufMax; i++ {
		mt.buf[metricKey{metricPage, string(rune('a'+i%26)) + string(rune('a'+i/26)), true}] = 1
	}
	mt.giveBack([]metricDelta{
		{k: metricKey{metricPage, "home", true}, n: 7},
		{k: metricKey{metricPage, "article", true}, n: 11},
	})
	if got := mt.Dropped(); got != 18 {
		t.Errorf("Dropped() = %d, хотели 18 (7+11), а не число строк", got)
	}
}
