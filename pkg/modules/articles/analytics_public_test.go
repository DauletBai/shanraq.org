package articles

import (
	"html/template"
	"net/http"
	"strings"
	"testing"
	"time"
)

// seedAudience lays down the minimum at which the page shows charts rather than
// "no data yet". Without it the test would only exercise the empty branch — and
// would pass on other tests' leftover rows in the database.
func seedAudience(app *testApp) {
	app.exec(`INSERT INTO analytics_daily (day, kind, label, is_guest, n) VALUES
		(CURRENT_DATE, 'page', 'article', true, 12),
		(CURRENT_DATE, 'page', 'home', true, 7),
		(CURRENT_DATE, 'page', 'home', false, 2),
		(CURRENT_DATE, 'country', 'KZ', true, 9),
		(CURRENT_DATE, 'lang', 'ru', true, 15),
		(CURRENT_DATE, 'browser', 'chrome', true, 11),
		(CURRENT_DATE, 'os', 'windows', true, 11),
		(CURRENT_DATE, 'device', 'desktop', true, 11),
		(CURRENT_DATE, 'source', 'direct', true, 14),
		(CURRENT_DATE, 'bot', 'google', true, 40)
		ON CONFLICT DO NOTHING`)
	// The cache lives an hour and is shared across the process, so it has to be
	// dropped between tests: otherwise the page shows the previous database's numbers.
	publicStats.mu.Lock()
	publicStats.byLg = map[string]PublicStats{}
	publicStats.at = time.Time{}
	publicStats.mu.Unlock()
}

// The analytics page is public: it exists precisely so that a stranger can see it.
// If it demands a login, there is no point to it.
func TestTheAudiencePageIsOpenToEveryone(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedAudience(app)

	w := app.do(http.MethodGet, "/analytics", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница аналитики отдала %d без входа", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, T(LangRU, "stats.title")) {
		t.Error("на странице нет её собственного заголовка")
	}
	// The method is stated on the page itself: a number without a definition is the
	// very thing our own articles argue with.
	if !strings.Contains(body, T(LangRU, "stats.method")) {
		t.Error("на странице не сказано, как мы считаем")
	}
	// Crawlers are shown separately rather than blended in with people: there are more
	// of them, and mixing the two would double the reported audience.
	if !strings.Contains(body, T(LangRU, "stats.bots_note")) {
		t.Error("нет оговорки о том, что краулеры не входят в число людей")
	}
}

// Not one line of the page may point at a person. We do not write specific
// addresses — only the kind of page — and the page must not cross the country with
// anything.
func TestTheAudiencePageNamesNoReaderAndNoArticle(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	seedAudience(app)
	author := app.createUser("statspriv@example.com", "Parol123!")
	_, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() WHERE slug = $1`, slug)

	body := app.do(http.MethodGet, "/analytics", nil).Body.String()
	if strings.Contains(body, "/read/"+slug) {
		t.Error("страница аналитики выдаёт адреса отдельных статей")
	}
	if strings.Contains(body, "statspriv@example.com") {
		t.Error("на странице оказался адрес электронной почты")
	}
}

// The page has to open in all three languages: it is part of the publication, not
// an internal export.
func TestTheAudiencePageSpeaksAllThreeLanguages(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	seedAudience(app)
	for _, lang := range []string{LangKZ, LangRU, LangEN} {
		w := app.do(http.MethodGet, "/analytics?lang="+lang, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("язык %s: %d", lang, w.Code)
		}
		// The comparison is against the escaped form: an apostrophe, an ampersand and a
		// plus look different in markup than in the dictionary, and a literal search would
		// find not a missing translation but the template engine doing its job.
		if want := template.HTMLEscapeString(T(lang, "stats.lead")); !strings.Contains(w.Body.String(), want) {
			t.Errorf("язык %s: нет вводной строки на этом языке", lang)
		}
	}
}

// Countries is the panel a publisher acts on: it decides what gets written and
// who the advertiser is sold. Twenty rows cut it off while the tail still held
// countries worth naming, so this one panel keeps thirty while the rest stay at
// twenty.
func TestTheCountryPanelKeepsMoreRowsThanTheOthers(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	// Twenty-six countries, each with a different count so the order is settled,
	// plus one more so the tail is folded rather than merely truncated.
	for i := 0; i < 26; i++ {
		code := string(rune('A'+i)) + "Z"
		app.exec(`INSERT INTO analytics_daily (day, kind, label, is_guest, n)
			VALUES (CURRENT_DATE, 'country', $1, true, $2)
			ON CONFLICT (day, kind, label, is_guest) DO UPDATE SET n = EXCLUDED.n`, code, 100-i)
		app.exec(`INSERT INTO analytics_daily (day, kind, label, is_guest, n)
			VALUES (CURRENT_DATE, 'browser', $1, true, $2)
			ON CONFLICT (day, kind, label, is_guest) DO UPDATE SET n = EXCLUDED.n`, code, 100-i)
	}

	rows := app.module().simpleRowsN(t.Context(), metricCountry, "ag.country.", "ru", countryRowsMax)
	if len(rows) < 21 {
		t.Fatalf("countries returned %d rows; the cap did not rise above twenty", len(rows))
	}
	if len(rows) > countryRowsMax+1 { // +1 for the folded remainder
		t.Errorf("countries returned %d rows, more than %d plus a remainder", len(rows), countryRowsMax)
	}

	// The other panels are unchanged, on the same shape of data.
	other := app.module().simpleRows(t.Context(), metricBrowser, "ag.browser.", "ru")
	if len(other) > simpleRowsMax+1 {
		t.Errorf("browsers returned %d rows, more than %d plus a remainder", len(other), simpleRowsMax)
	}
	if len(other) >= len(rows) {
		t.Errorf("browsers kept %d rows and countries %d; countries should keep more", len(other), len(rows))
	}
}
