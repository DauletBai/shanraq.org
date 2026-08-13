package articles

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// seedCitable publishes an article that llms.txt is willing to list: published,
// indexable, and dated. seedArticle leaves published_at null, which is fine for
// the reader but is exactly what the index filters on.
func (a *testApp) seedCitable(author uuid.UUID) (uuid.UUID, string) {
	a.t.Helper()
	id, slug := a.seedArticle(author, "published")
	a.exec(`UPDATE articles SET published_at = NOW() WHERE id = $1`, id)
	return id, slug
}

// The whole point of the AI group is that a machine-written column, already out
// of the search index, is also out of reach of the assistants. If robots.txt
// stops naming it, the de-indexing is only half done and nobody would notice.
func TestRobotsKeepsAICrawlersOffTheAIColumns(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("robots-ai@example.com", "Sup3r-Secret-Pass!")
	_, open := app.seedCitable(author)
	hiddenID, hidden := app.seedCitable(author)
	app.exec(`UPDATE articles SET indexable = FALSE WHERE id = $1`, hiddenID)

	body := app.do(http.MethodGet, "/robots.txt", nil).Body.String()
	group, ok := aiGroupOf(body)
	if !ok {
		t.Fatal("robots.txt has no group naming the AI crawlers")
	}
	if !strings.Contains(group, "Disallow: /read/"+hidden) {
		t.Errorf("the non-indexable article is not disallowed for AI crawlers:\n%s", group)
	}
	if strings.Contains(group, "Disallow: /read/"+open) {
		t.Error("an indexable article was disallowed for AI crawlers; the channel is the point")
	}
}

// A crawler obeys exactly one group: the most specific one that names it. Once
// GPTBot finds itself by name it stops reading User-agent: *, so every private
// path has to be repeated here or it is served to the whole list.
func TestAIGroupRepeatsThePrivatePaths(t *testing.T) {
	app := newTestApp(t)
	group, ok := aiGroupOf(app.do(http.MethodGet, "/robots.txt", nil).Body.String())
	if !ok {
		t.Fatal("robots.txt has no AI crawler group")
	}
	for _, p := range []string{"/admin", "/studio", "/api/", "/jobs"} {
		if !strings.Contains(group, "Disallow: "+p) {
			t.Errorf("the AI group does not repeat %q, so it is open to every agent it names", p)
		}
	}
}

// aiGroupOf returns the robots.txt group that names GPTBot: everything from its
// User-agent line to the blank line that ends the group.
func aiGroupOf(robots string) (string, bool) {
	i := strings.Index(robots, "User-agent: GPTBot")
	if i < 0 {
		return "", false
	}
	rest := robots[i:]
	if end := strings.Index(rest, "\n\n"); end >= 0 {
		return rest[:end], true
	}
	return rest, true
}

// A page kept out of the index has to say so in a header too: AI crawlers do
// not parse a robots meta tag, and a HEAD request has no body to parse anyway.
func TestNonIndexableArticleSendsTheHeader(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("robots-hdr@example.com", "Sup3r-Secret-Pass!")
	id, slug := app.seedCitable(author)
	app.exec(`UPDATE articles SET indexable = FALSE WHERE id = $1`, id)

	got := app.do(http.MethodGet, "/read/"+slug, nil).Header().Get("X-Robots-Tag")
	for _, want := range []string{"noindex", "noai"} {
		if !strings.Contains(got, want) {
			t.Errorf("X-Robots-Tag = %q, missing %q", got, want)
		}
	}
	_, openSlug := app.seedCitable(author)
	if h := app.do(http.MethodGet, "/read/"+openSlug, nil).Header().Get("X-Robots-Tag"); h != "" {
		t.Errorf("an indexable article should carry no X-Robots-Tag, got %q", h)
	}
}

// llms.txt is an invitation to cite us. Listing a machine-written column in it
// would invite exactly what the rest of this change prevents.
func TestLLMSListsOnlyTheHumanWork(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("llms@example.com", "Sup3r-Secret-Pass!")
	_, open := app.seedCitable(author)
	hiddenID, hidden := app.seedCitable(author)
	app.exec(`UPDATE articles SET indexable = FALSE WHERE id = $1`, hiddenID)

	w := app.do(http.MethodGet, "/llms.txt", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("llms.txt = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("llms.txt Content-Type = %q, want text/plain", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "/read/"+open) {
		t.Error("llms.txt does not list the human-written article")
	}
	if strings.Contains(body, "/read/"+hidden) {
		t.Error("llms.txt lists a machine-written column")
	}
}

// The reference is the deliverable: an author, a title, a date in words and an
// absolute URL. Losing any one of them turns a citation into a shrug.
func TestCiteLine(t *testing.T) {
	when := time.Date(2026, 8, 13, 7, 45, 0, 0, time.UTC)
	got := citeLine(LangRU, "Даулет Баймурза", "Украина отказалась быть стадом",
		"https://shanraq.org", "/read/ukraina", &when)
	for _, want := range []string{
		"Даулет Баймурза",
		"«Украина отказалась быть стадом»",
		"Shanraq.org",
		"13 августа 2026",
		"https://shanraq.org/read/ukraina",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("citation line is missing %q:\n%s", want, got)
		}
	}
	if en := citeLine(LangEN, "D. B.", "Title", "https://shanraq.org", "/read/x", &when); !strings.Contains(en, "August 13, 2026") {
		t.Errorf("English date should read as English: %s", en)
	}
	if kz := citeLine(LangKZ, "D. B.", "Title", "https://shanraq.org", "/read/x", &when); !strings.Contains(kz, "13 тамыз 2026") {
		t.Errorf("Kazakh date should read as Kazakh: %s", kz)
	}
	// An article with no date still yields a usable reference rather than the
	// word "January" invented out of a zero time.
	if none := citeLine(LangRU, "A", "T", "https://shanraq.org", "/read/x", nil); strings.Contains(none, "января") {
		t.Errorf("a dateless article invented a date: %s", none)
	}
}

// The block is only offered where we stand behind the text.
func TestCitationBlockOnlyOnCitableArticles(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("cite@example.com", "Sup3r-Secret-Pass!")
	_, open := app.seedCitable(author)
	hiddenID, hidden := app.seedCitable(author)
	app.exec(`UPDATE articles SET indexable = FALSE WHERE id = $1`, hiddenID)

	if body := app.do(http.MethodGet, "/read/"+open, nil).Body.String(); !strings.Contains(body, `class="cite"`) {
		t.Error("a citable article shows no citation block")
	}
	if body := app.do(http.MethodGet, "/read/"+hidden, nil).Body.String(); strings.Contains(body, `class="cite"`) {
		t.Error("a machine-written column offers a citation block")
	}
}
