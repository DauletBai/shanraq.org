package articles

import (
	"net/url"
	"strings"
	"testing"
)

// The rate card's worked example is recomputed from the same populations the
// picker prices with, so a figure an advertiser reads on the card is the figure
// the form will quote. This asserts the two have not drifted apart — the card
// going stale while the form quotes something else is the failure that costs
// trust rather than money.

func TestAdvertisePageShowsTheGeoLadder(t *testing.T) {
	app := newTestApp(t)
	app.createUser("adv-geo@example.com", "sup3rSecret!42")
	c := app.login("adv-geo@example.com", "sup3rSecret!42")
	app.do("POST", "/advertise/company", url.Values{
		"company_name":  {"ТОО Качар-Сервис"},
		"bin":           {"123456789012"},
		"legal_form":    {"too"},
		"address":       {"Качар"},
		"contact_name":  {"Иван"},
		"contact_phone": {"+77010000000"},
		"contact_email": {"adv-geo@example.com"},
	}, withCookie(c))
	rec := app.do("GET", "/advertise", nil, withCookie(c))
	t.Logf("HTTP %d, %d байт", rec.Code, rec.Body.Len())
	b := rec.Body.String()
	for _, probe := range []string{
		"Как география меняет цену", "Вся страна", "data-ad-geo",
		"Что именно продаётся", "Костанайская область", "Качар",
	} {
		if !strings.Contains(b, probe) {
			t.Errorf("нет на странице: %q", probe)
		}
	}
	// The ladder must show the same figures the pricing code produces.
	for _, want := range []string{"135 000", "27 000", "2 889"} {
		if !strings.Contains(b, want) {
			t.Errorf("нет цены %q", want)
		}
	}
}
