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

	// Exactly one. A second call sat at the foot of the sidebar and had rendered
	// nothing for as long as the slot returned nil on an unsold surface; the
	// moment it stopped doing that, the page carried two.
	if n := strings.Count(body, `class="adcarousel adslot"`); n != 1 {
		t.Fatalf("the page has %d advertising slots, want exactly 1", n)
	}
	if strings.Contains(body, "adslide--news") {
		t.Error("the duplicate news carousel is still in the sidebar")
	}
	// One carousel of headlines, not two.
	if n := strings.Count(body, "newshero__slide"); n == 0 {
		t.Error("the hero carousel disappeared along with the duplicate")
	}
	// Marked as advertising. A house slide advertises the slot rather than a
	// customer, but it is still advertising and still has to say so.
	if !strings.Contains(body, `adslot__ribbon adslot__ribbon--house">Реклама`) {
		t.Error("the slot carries no advertising mark")
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
