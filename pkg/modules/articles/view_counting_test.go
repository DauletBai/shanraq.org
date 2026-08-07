package articles

import (
	"context"
	"net/http"
	"testing"
)

// TestCrawlersDoNotCountAsViews pins the fix for the defect that made the whole
// dashboard misleading: every crawler hit used to bump views_count. Two thirds
// of the "views" on the live site were Googlebot, the Facebook link scraper and
// AI crawlers — and because that count is the denominator of the reading-depth
// funnel (whose numerator needs JavaScript, so only humans reach it), a real 23%
// read-through was displayed as 2%.
//
// The test drives real HTTP through the router, so it covers the handler path
// rather than the helper.
func TestCrawlersDoNotCountAsViews(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("viewcount@example.com", "Sup3r-Secret-Pass!")
	id, slug := app.seedArticle(author, "published")

	views := func() int64 {
		app.t.Helper()
		var n int64
		if err := app.pool.QueryRow(context.Background(),
			`SELECT views_count FROM articles WHERE id = $1`, id).Scan(&n); err != nil {
			t.Fatalf("read views: %v", err)
		}
		return n
	}

	if got := views(); got != 0 {
		t.Fatalf("fresh article starts at %d views, want 0", got)
	}

	// Every one of these used to add a view. They are the exact families the
	// live analytics panel reports as the heaviest crawlers of the site.
	crawlers := map[string]string{
		"googlebot": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"facebook":  "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"ai":        "Mozilla/5.0 (compatible; GPTBot/1.1; +https://openai.com/gptbot)",
		"claude":    "Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)",
		"applebot":  "Mozilla/5.0 (compatible; Applebot/0.1)",
		"headless":  "Mozilla/5.0 (X11; Linux x86_64) HeadlessChrome/120.0.0.0 Safari/537.36",
		"curl":      "curl/8.4.0",
		"empty-ua":  "",
	}
	for name, ua := range crawlers {
		w := app.do(http.MethodGet, "/read/"+slug, nil, withHeader("User-Agent", ua))
		if w.Code != http.StatusOK {
			t.Fatalf("%s: GET article = %d, want 200", name, w.Code)
		}
		if got := views(); got != 0 {
			t.Fatalf("%s bumped views to %d — crawlers must not be counted", name, got)
		}
	}

	// A real browser still counts, or the fix would have replaced one wrong
	// number with another.
	const chrome = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	if w := app.do(http.MethodGet, "/read/"+slug, nil, withHeader("User-Agent", chrome)); w.Code != http.StatusOK {
		t.Fatalf("GET article as browser = %d, want 200", w.Code)
	}
	if got := views(); got != 1 {
		t.Fatalf("browser visit left views at %d, want 1", got)
	}
	if w := app.do(http.MethodGet, "/read/"+slug, nil, withHeader("User-Agent", chrome)); w.Code != http.StatusOK {
		t.Fatalf("second browser GET = %d, want 200", w.Code)
	}
	if got := views(); got != 2 {
		t.Fatalf("second browser visit left views at %d, want 2", got)
	}
}

// TestFinishRateIgnoresViewCount checks the figure the studio now leads with.
// It divides finishers by starters, so it cannot be pushed around by whatever
// the view counter happens to include — which is the whole reason it exists.
func TestFinishRateIgnoresViewCount(t *testing.T) {
	cases := []struct {
		name             string
		d25, d100, views int64
		want             int
	}{
		{"three of four finish", 4, 3, 1000, 75},
		{"all finish", 10, 10, 99999, 100},
		{"nobody starts", 0, 0, 500, 0},
		{"inflated views do not matter", 8, 6, 10_000_000, 75},
	}
	for _, c := range cases {
		if got := pctOf(c.d100, c.d25); got != c.want {
			t.Errorf("%s: finish rate = %d%%, want %d%%", c.name, got, c.want)
		}
	}
}
