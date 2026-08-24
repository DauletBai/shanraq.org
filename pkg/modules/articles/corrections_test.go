package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRepairHomoglyphs(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		// A Latin A inside a Cyrillic word: a keyboard that changed mid-word.
		{"Aлматы", "Алматы", true},
		{"экoномика", "экономика", true},
		{"тeнгe", "тенге", true},
		// Pure Cyrillic and pure Latin are not this kind of mistake.
		{"Алматы", "", false},
		{"economy", "", false},
		// Latin in the majority means it is a Latin word, not a damaged Cyrillic
		// one: "Kazakh" must not be turned into Cyrillic because of one letter.
		{"Kazaхstan", "", false},
		// A Latin letter with no Cyrillic twin: guessing would corrupt the word.
		{"инфляgия", "", false},
	}
	for _, c := range cases {
		got, ok := repairHomoglyphs(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("repairHomoglyphs(%q) = %q, %v; want %q, %v", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestSafeFix(t *testing.T) {
	ok := []struct{ word, fixed string }{
		{"агенство", "агентство"},
		{"коффициент", "коэффициент"},
		{"вобщем", "в общем"},
		{"Тенге", "Теңге"},
	}
	for _, c := range ok {
		if err := safeFix(c.word, c.fixed); err != nil {
			t.Errorf("safeFix(%q, %q) refused a real fix: %v", c.word, c.fixed, err)
		}
	}

	bad := []struct {
		word, fixed, why string
	}{
		{"тенге", "тенге", "unchanged"},
		{"тенге", "", "empty"},
		{"тенге", "доллар США по официальному курсу", "a phrase, not a word"},
		{"тенге", "[тенге](http://evil)", "markup"},
		{"тенге", "тенге\n\nновый абзац", "newlines"},
		{"тенге", "манат", "a different word, not a spelling"},
		{"курс", "**курс**", "emphasis"},
		{"один", "один два три", "three words"},
	}
	for _, c := range bad {
		if err := safeFix(c.word, c.fixed); err == nil {
			t.Errorf("safeFix(%q, %q) allowed %s", c.word, c.fixed, c.why)
		}
	}
}

func TestMatchCase(t *testing.T) {
	for _, c := range []struct{ word, fixed, want string }{
		{"Агенство", "агентство", "Агентство"},
		{"агенство", "Агентство", "агентство"},
		{"НБРК", "нбк", "НБК"},
	} {
		if got := matchCase(c.word, c.fixed); got != c.want {
			t.Errorf("matchCase(%q, %q) = %q, want %q", c.word, c.fixed, got, c.want)
		}
	}
}

const corrBody = `## Курс и ставка

Национальный банк объявил базовую ставку. Ставка выше инфляции, и это агенство
подтверждает.

Ссылка на [агенство](https://example.com/agenstvo) в тексте, и ` + "`agenstvo`" + ` в коде.

## Второй раздел

Здесь слово агенство встречается ещё раз, в другом предложении.

` + "```" + `
var agenstvo = 1
` + "```" + `
`

func TestFindCorrectionSitePicksTheSentence(t *testing.T) {
	// The same word sits in four places: prose in the first section, a link
	// target, a code span, prose in the second section, and a fenced block. The
	// sentence is what decides which one the reader meant.
	site, ok := findCorrectionSite(corrBody,
		"Здесь слово агенство встречается ещё раз, в другом предложении.", "агенство")
	if !ok {
		t.Fatal("second-section sentence not located")
	}
	if got := corrBody[site.Start:site.End]; got != "агенство" {
		t.Fatalf("located %q", got)
	}
	if !strings.Contains(corrBody[:site.Start], "Второй раздел") {
		t.Error("landed in the first section, not the second")
	}

	first, ok := findCorrectionSite(corrBody,
		"Ставка выше инфляции, и это агенство подтверждает.", "агенство")
	if !ok {
		t.Fatal("first-section sentence not located")
	}
	if strings.Contains(corrBody[:first.Start], "Второй раздел") {
		t.Error("landed in the second section, not the first")
	}
	if first.Para == "" || !strings.Contains(first.Para, "Национальный банк") {
		t.Errorf("paragraph context is wrong: %q", first.Para)
	}
}

func TestFindCorrectionSiteRefusesUntouchableSpots(t *testing.T) {
	// A sentence that only matches inside a link target, a code span or a fence
	// must find nothing: correcting any of the three breaks something silently.
	for _, s := range []string{
		"Ссылка на агенство в тексте, и agenstvo в коде.",
		"var agenstvo = 1",
	} {
		if site, ok := findCorrectionSite(corrBody, s, "agenstvo"); ok {
			t.Errorf("located %q at %d, expected a refusal", s, site.Start)
		}
	}
}

func TestFindCorrectionSiteRefusesAForeignSentence(t *testing.T) {
	if _, ok := findCorrectionSite(corrBody,
		"Совершенно другое предложение про погоду в Караганде зимой.", "агенство"); ok {
		t.Error("a sentence that is not in the article was located anyway")
	}
}

func TestFindCorrectionSiteRefusesHeadings(t *testing.T) {
	// Headings are anchors: the table of contents and every deep link into the
	// article are built from them.
	if _, ok := findCorrectionSite("## Курс и ставка\n\nтекст\n", "Курс и ставка", "Курс"); ok {
		t.Error("a heading was offered for correction")
	}
}

func TestApplyAtTouchesOneOccurrence(t *testing.T) {
	body := "агенство и ещё раз агенство"
	site, ok := findCorrectionSite(body, "агенство и ещё раз агенство", "агенство")
	if !ok {
		t.Fatal("not located")
	}
	got := applyAt(body, site, "агентство")
	if strings.Count(got, "агентство") != 1 || strings.Count(got, "агенство") != 1 {
		t.Errorf("replaced the wrong number of occurrences: %q", got)
	}
}

func TestLevenshtein(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"", "abc", 3},
		{"агенство", "агентство", 1},
		{"тенге", "тенге", 0},
		{"коффициент", "коэффициент", 1},
	} {
		if got := levenshtein(c.a, c.b); got != c.want {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

// ---- over HTTP ----

func TestCorrectionFormNeedsNoLoginToRead(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("corr-author@example.com", "Passw0rd!x")
	_, slug := app.seedArticle(author, "published")

	rec := app.do(http.MethodGet, "/read/"+slug+"/typo", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET form = %d", rec.Code)
	}
	body := rec.Body.String()
	// A guest sees what the page is for and how to become able to use it.
	if !strings.Contains(body, T(LangRU, "corr.login_note")) {
		t.Error("the guest was not told an account is needed")
	}
	// All three fields carry a "?" hint, which is the whole point of the form.
	for _, k := range []string{"corr.h_chapter", "corr.h_sentence", "corr.h_word"} {
		if !strings.Contains(body, T(LangRU, k)) {
			t.Errorf("hint %s missing from the form", k)
		}
	}
}

func TestCorrectionSubmitRequiresLogin(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("corr-author2@example.com", "Passw0rd!x")
	_, slug := app.seedArticle(author, "published")

	rec := app.do(http.MethodPost, "/read/"+slug+"/typo", url.Values{
		"sentence": {"текст статьи"}, "word": {"текст"},
	})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("anonymous POST = %d, want a redirect to the login page", rec.Code)
	}
}

func TestCorrectionRejectsMalformedClaims(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("corr-author3@example.com", "Passw0rd!x")
	_, slug := app.seedArticle(author, "published")
	reader := "corr-reader3@example.com"
	app.createUser(reader, "Passw0rd!x")
	c := app.login(reader, "Passw0rd!x")

	cases := []struct {
		name  string
		form  url.Values
		wants string
	}{
		{"no sentence", url.Values{"word": {"текст"}}, "corr.e_sentence"},
		{"no word", url.Values{"sentence": {"текст статьи"}}, "corr.e_word"},
		{"a phrase instead of a word",
			url.Values{"sentence": {"текст статьи"}, "word": {"текст статьи"}}, "corr.e_one_word"},
		{"word absent from the sentence",
			url.Values{"sentence": {"текст статьи"}, "word": {"инфляция"}}, "corr.e_not_in_sentence"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.do(http.MethodPost, "/read/"+slug+"/typo", tc.form, withCookie(c))
			if rec.Code != http.StatusOK {
				t.Fatalf("POST = %d", rec.Code)
			}
			if want := T(LangRU, tc.wants); !strings.Contains(rec.Body.String(), want) {
				t.Errorf("expected %q in the response", want)
			}
		})
	}
}

// TestCorrectionAppliesHomoglyphFix walks the whole route with no model in play:
// a Latin letter wearing a Cyrillic face is settled by the server itself, so this
// exercises submission, location, application, and the outcome shown back.
func TestCorrectionAppliesHomoglyphFix(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("corr-author4@example.com", "Passw0rd!x")
	id, slug := app.seedArticle(author, "published")
	// A Latin "o" in the middle of a Cyrillic word.
	app.exec(`UPDATE article_translations SET body_md = $2 WHERE article_id = $1 AND lang = 'ru'`,
		id, "## Тело\n\nВ стране растёт экoномика и это заметно.")

	reader := "corr-reader4@example.com"
	app.createUser(reader, "Passw0rd!x")
	c := app.login(reader, "Passw0rd!x")

	rec := app.do(http.MethodPost, "/read/"+slug+"/typo", url.Values{
		"chapter":  {"Тело"},
		"sentence": {"В стране растёт экoномика и это заметно."},
		"word":     {"экoномика"},
	}, withCookie(c))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST = %d", rec.Code)
	}
	if want := T(LangRU, "corr.h_applied"); !strings.Contains(rec.Body.String(), want) {
		t.Fatalf("expected %q in the response, got:\n%s", want, clipForLog(rec.Body.String()))
	}

	art, err := NewStore(app.pool).GetPublishedBySlug(context.Background(), slug)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	body := art.Translations[LangRU].BodyMD
	if strings.Contains(body, "экoномика") {
		t.Error("the Latin letter is still in the article")
	}
	if !strings.Contains(body, "экономика") {
		t.Errorf("the word was not repaired: %q", body)
	}

	// And the claim is on the record, with what it became.
	rows, err := app.arts.corrections.ForArticle(context.Background(), id, 10)
	if err != nil {
		t.Fatalf("ForArticle: %v", err)
	}
	if len(rows) != 1 || rows[0].Status != CorrApplied || rows[0].Fixed != "экономика" {
		t.Errorf("audit row = %+v", rows)
	}
}

// TestCorrectionHeldWithoutChecker covers the honest-failure path: with no model
// configured a claim that needs judgement is kept, and the reader is told it is
// waiting rather than told it was rejected.
func TestCorrectionHeldWithoutChecker(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("corr-author5@example.com", "Passw0rd!x")
	id, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE article_translations SET body_md = $2 WHERE article_id = $1 AND lang = 'ru'`,
		id, "## Тело\n\nВ прошлом году агенство отчиталось о росте.")

	reader := "corr-reader5@example.com"
	app.createUser(reader, "Passw0rd!x")
	c := app.login(reader, "Passw0rd!x")

	rec := app.do(http.MethodPost, "/read/"+slug+"/typo", url.Values{
		"sentence": {"В прошлом году агенство отчиталось о росте."},
		"word":     {"агенство"},
	}, withCookie(c))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST = %d", rec.Code)
	}
	if want := T(LangRU, "corr.h_pending"); !strings.Contains(rec.Body.String(), want) {
		t.Errorf("expected %q, got:\n%s", want, clipForLog(rec.Body.String()))
	}
	art, _ := NewStore(app.pool).GetPublishedBySlug(context.Background(), slug)
	if !strings.Contains(art.Translations[LangRU].BodyMD, "агенство") {
		t.Error("the text was changed with no checker to authorise it")
	}
}

// clipForLog keeps a failure message readable.
func clipForLog(s string) string {
	if len(s) > 1200 {
		return s[:1200] + "…"
	}
	return s
}
