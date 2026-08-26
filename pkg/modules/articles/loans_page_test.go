package articles

import (
	"strings"
	"testing"
)

// The page has to render in all three languages and carry the programme table
// as plain HTML: that table is what a reader without script sees and what a
// crawler indexes.
func TestCalculatorPageRendersInEveryLanguage(t *testing.T) {
	app := newTestApp(t)
	for _, lang := range []string{"kz", "ru", "en"} {
		rec := app.do("GET", "/calculator?lang="+lang, nil)
		if rec.Code != 200 {
			t.Fatalf("%s: код %d", lang, rec.Code)
		}
		body := rec.Body.String()
		for _, want := range []string{T(lang, "calc.title"), T(lang, "calc.who"), "data-programs", "7-20-25"} {
			if !strings.Contains(body, want) {
				t.Errorf("%s: на странице нет %q", lang, want)
			}
		}
		// Unverified rows must be visibly marked, not quietly shown as fact.
		if !strings.Contains(body, T(lang, "calc.t_unchecked")) {
			t.Errorf("%s: непроверенные ставки не помечены", lang)
		}
	}
}

// The sidebar of the classifieds points at the calculator.
func TestListingsSidebarLinksToTheCalculator(t *testing.T) {
	app := newTestApp(t)
	rec := app.do("GET", "/listings?lang=ru", nil)
	if rec.Code != 200 {
		t.Fatalf("код %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "/calculator?lang=ru") {
		t.Error("в сайдбаре объявлений нет ссылки на калькулятор")
	}
}
