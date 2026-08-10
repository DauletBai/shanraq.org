package articles

import (
	"net/http"
	"strings"
	"testing"
)

// The home page ran two carousels: the hero across the top and, in the sidebar,
// a second one cycling the same six headlines. Half the sidebar was spent
// saying what the page already said. The sidebar now carries the advertising
// slot instead — and because nothing is sold yet, it sells itself.
func TestHomeSidebarCarriesTheAdSlotNotACopyOfTheHero(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("adslot@t.test", "Parol12345")
	app.seedArticle(author, "published")

	body := app.do(http.MethodGet, "/", nil).Body.String()

	if !strings.Contains(body, `class="adcarousel adslot"`) {
		t.Fatal("the sidebar has no advertising slot")
	}
	if strings.Contains(body, "adslide--news") {
		t.Error("the duplicate news carousel is still in the sidebar")
	}
	// One carousel of headlines, not two.
	if n := strings.Count(body, "newshero__slide"); n == 0 {
		t.Error("the hero carousel disappeared along with the duplicate")
	}
	// Unsold slots must not pose as sold ones.
	if !strings.Contains(body, "Место свободно") {
		t.Error("an unsold slot does not say so")
	}
	if strings.Contains(body, `adslot__ribbon">Реклама`) {
		t.Error("an unsold slot is labelled as an advertisement")
	}
	// House slides point inward, so they must not be marked as paid outbound
	// links: rel="sponsored" on our own /advertise page would be a lie to
	// crawlers about a relationship that does not exist.
	if strings.Contains(body, `href="/advertise" rel="nofollow sponsored`) {
		t.Error("a house slide is marked as a sponsored outbound link")
	}
}

// Every house slide has to carry all three of its texts in every language —
// a missing key renders as the key itself, which would ship "house.free_desc"
// into the sidebar of the front page.
func TestHouseAdsAreFullyTranslated(t *testing.T) {
	for _, lang := range []string{LangRU, LangKZ, LangEN} {
		ads := houseAds(lang)
		if len(ads) < 3 {
			t.Fatalf("%s: only %d house slides — the carousel needs something to cycle", lang, len(ads))
		}
		for i, a := range ads {
			for name, v := range map[string]string{"title": a.Title, "desc": a.Desc, "cta": a.Price} {
				if v == "" || strings.HasPrefix(v, "house.") {
					t.Errorf("%s slide %d: %s is untranslated (%q)", lang, i, name, v)
				}
			}
			if !a.House {
				t.Errorf("%s slide %d is not marked as a house slide, so it would be labelled an advertisement", lang, i)
			}
			if a.URL == "" {
				t.Errorf("%s slide %d has nowhere to click", lang, i)
			}
		}
	}
}
