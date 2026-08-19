package articles

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// Only the feed's own cards: the sidebar's "latest" list links to the newest
// articles on every page, and counting those would make any two pages look
// like they overlapped.
var slugRe = regexp.MustCompile(`class="post__title"><a href="/read/([a-z0-9-]+)\?`)

func feedSlugs(body string) map[string]bool {
	out := map[string]bool{}
	for _, m := range slugRe.FindAllStringSubmatch(body, -1) {
		out[m[1]] = true
	}
	return out
}

// The feed shows 21 articles and the site holds several times that. Until this
// existed, everything past the twenty-first was reachable only by search or a
// direct link: no path from the site, and no path for a crawler either — which
// is one answer to why the archive draws no search traffic.
func TestHomeFeedPaginates(t *testing.T) {
	app := newTestApp(t)

	var published int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM articles WHERE status = 'published'`).Scan(&published); err != nil {
		t.Fatalf("count published: %v", err)
	}
	if published <= homePageSize {
		t.Skipf("only %d published articles; pagination needs more than %d", published, homePageSize)
	}

	first := app.do(http.MethodGet, "/", nil)
	if first.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", first.Code)
	}
	body1 := first.Body.String()
	if !strings.Contains(body1, `rel="next"`) {
		t.Error("page 1 offers no way forward")
	}
	if strings.Contains(body1, `rel="prev"`) {
		t.Error("page 1 offers a way back from the beginning")
	}

	second := app.do(http.MethodGet, "/?page=2", nil)
	if second.Code != http.StatusOK {
		t.Fatalf("GET /?page=2 = %d, want 200", second.Code)
	}
	body2 := second.Body.String()
	if !strings.Contains(body2, `rel="prev"`) {
		t.Error("page 2 offers no way back")
	}

	// Different articles, or the second page is the first one again.
	s1, s2 := feedSlugs(body1), feedSlugs(body2)
	if len(s1) == 0 || len(s2) == 0 {
		t.Fatalf("no articles parsed: page 1 had %d, page 2 had %d", len(s1), len(s2))
	}
	for slug := range s2 {
		if s1[slug] {
			t.Errorf("article %q appears on both pages", slug)
		}
	}

	// Canonical and hreflang must name the page. Pointing page 2 at "/" tells a
	// crawler the two are the same page, which is how an archive disappears.
	if !strings.Contains(body2, "page=2") {
		t.Error("page 2 does not name itself in its canonical or its links")
	}
	if strings.Contains(body1, `rel="canonical" href="`+testOrigin+`/?page=`) {
		t.Error("page 1 canonicalised itself with a page number")
	}
}

// A page number is user input and reaches the database as an offset.
func TestHomeFeedSurvivesAJunkPageNumber(t *testing.T) {
	app := newTestApp(t)
	for _, q := range []string{"?page=abc", "?page=-5", "?page=0", "?page=99999999999999999999", "?page=1e9", "?page="} {
		w := app.do(http.MethodGet, "/"+q, nil)
		if w.Code != http.StatusOK {
			t.Errorf("GET /%s = %d, want 200", q, w.Code)
		}
		if strings.Contains(w.Body.String(), `rel="prev"`) {
			t.Errorf("GET /%s was treated as a later page", q)
		}
	}
	// Far past the end is empty, not an error.
	w := app.do(http.MethodGet, "/?page=999", nil)
	if w.Code != http.StatusOK {
		t.Errorf("GET /?page=999 = %d, want 200", w.Code)
	}
	if strings.Contains(w.Body.String(), `rel="next"`) {
		t.Error("a page past the end still offers a next page")
	}
}

// The pager has to carry the reader's filter, or paging out of a category
// silently drops them back into everything.
func TestPagerKeepsTheCategory(t *testing.T) {
	app := newTestApp(t)
	w := app.do(http.MethodGet, "/?cat=economy&page=2", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /?cat=economy&page=2 = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `rel="prev"`) {
		t.Fatal("no way back from page 2 of a category")
	}
	prev := regexp.MustCompile(`href="([^"]+)" rel="prev"`).FindStringSubmatch(body)
	if prev == nil {
		t.Fatal("could not read the previous-page link")
	}
	if !strings.Contains(prev[1], "cat=economy") {
		t.Errorf("previous-page link %q dropped the category", prev[1])
	}
}
