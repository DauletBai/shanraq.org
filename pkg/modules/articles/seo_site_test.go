package articles

import (
	"encoding/json"
	"strings"
	"testing"
)

// unwrapLD pulls the JSON out of a <script type="application/ld+json"> block and
// insists that it parses. A block that does not parse is worse than no block:
// the page keeps rendering, search keeps reading nothing, and the fault shows up
// only in a validator nobody runs.
func unwrapLD(t *testing.T, s string) map[string]any {
	t.Helper()
	open := strings.Index(s, ">")
	close := strings.LastIndex(s, "</script>")
	if open < 0 || close < 0 || close < open {
		t.Fatalf("not a script block: %q", s)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(s[open+1:close]), &out); err != nil {
		t.Fatalf("structured data is not JSON: %v\n%s", err, s[open+1:close])
	}
	return out
}

// The site's card used to be four fields typed into the template, and it said
// the name of the site and nothing else. What a search engine needs is what the
// site publishes, and that sentence already exists as the meta description --
// so the two are now the same string, and this is the check that they stay so.
func TestSiteCardCarriesWhatTheSitePublishes(t *testing.T) {
	for _, lang := range []string{LangKZ, LangRU, LangEN} {
		ld := unwrapLD(t, string(siteLD("n0nce", "https://shanraq.org", lang)))
		if got := ld["description"]; got != T(lang, "seo.site_desc") {
			t.Errorf("%s: the card describes the site as %q, the meta tag as %q",
				lang, got, T(lang, "seo.site_desc"))
		}
		pub, ok := ld["publisher"].(map[string]any)
		if !ok {
			t.Fatalf("%s: no publisher in the site card", lang)
		}
		if pub["name"] != "Shanraq.org" {
			t.Errorf("%s: publisher is %v", lang, pub["name"])
		}
	}
}

// A course is free, and that is the whole offer. Stated as a flag alone it is a
// word inside a sentence; stated as an offer of zero it is a fact search can
// read, which is why both are emitted and why both are checked.
func TestCourseCardSaysTheCourseIsFree(t *testing.T) {
	page := &CoursePage{
		Series: &Series{
			Slug:    "go",
			Title:   map[string]string{LangRU: "Go: с нуля до своего блога"},
			Summary: map[string]string{LangRU: "Пятьдесят уроков."},
		},
		Minutes: 620,
	}
	page.Lang = LangRU
	page.SiteURL = "https://shanraq.org"
	page.Nonce = "n0nce"

	m := &Module{}
	m.applyCourseSEO(page)
	ld := unwrapLD(t, string(page.JSONLD))

	if ld["@type"] != "Course" {
		t.Errorf("the course page is marked up as %v", ld["@type"])
	}
	if ld["name"] != "Go: с нуля до своего блога" {
		t.Errorf("course name is %v", ld["name"])
	}
	if ld["isAccessibleForFree"] != true {
		t.Error("the course is not marked as free")
	}
	offer, ok := ld["offers"].(map[string]any)
	if !ok || offer["price"] != "0" {
		t.Errorf("the offer is %v, and a free course has to cost zero", ld["offers"])
	}
	inst, ok := ld["hasCourseInstance"].(map[string]any)
	if !ok || inst["courseWorkload"] != "PT620M" {
		t.Errorf("the workload is %v, not the reading time the page shows", ld["hasCourseInstance"])
	}
	if ld["url"] != "https://shanraq.org/course/go?lang=ru" {
		t.Errorf("the course URL is %v", ld["url"])
	}
}

// An info page used to be handed to search with the site's own description, so
// /about and /pricing were offered under the same sentence and neither said
// what it held. The page's opening lines are what it says about itself.
func TestInfoPageDescribesItselfNotTheSite(t *testing.T) {
	body := staticContent("about", LangRU).Body
	lead := excerpt(stripMD(firstParagraphs(body, 2)), 200)
	if lead == "" {
		t.Fatal("the About page yields no description")
	}
	if lead == T(LangRU, "seo.site_desc") {
		t.Error("the About page still describes the site rather than itself")
	}
	if strings.HasPrefix(lead, "#") || strings.Contains(lead, "##") {
		t.Errorf("a heading leaked into the description: %q", lead)
	}
	if n := len([]rune(lead)); n > 201 {
		t.Errorf("the description is %d runes; search shows a fraction of that", n)
	}
}

// A currency variant of /rates used to be the whole page with a different
// chart: same heading, same two thousand words, different numbers. Sixteen of
// those are one document to a search engine, which then picks a canonical of
// its own and drops the rest -- which is exactly what Search Console reported.
// The heading and the opening sentence are what make the page its own.
func TestCurrencyPageIsAboutThatCurrency(t *testing.T) {
	if got := T(LangRU, "fx.title_cur"); !strings.Contains(got, "%s") {
		t.Fatalf("the per-currency title has no place for the code: %q", got)
	}
	lead := T(LangRU, "fx.lead_cur")
	if n := strings.Count(lead, "%s"); n != 4 {
		t.Errorf("the per-currency lead takes %d values, the page fills four", n)
	}
	for _, lang := range []string{LangKZ, LangRU, LangEN} {
		if T(lang, "fx.lead_cur") == "" {
			t.Errorf("%s: no per-currency lead", lang)
		}
		if T(lang, "fx.title_cur") == T(lang, "fx.title") {
			t.Errorf("%s: a currency page is headed like the index", lang)
		}
	}
}
