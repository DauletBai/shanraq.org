package articles

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The three things a buyer asks before agreeing to a viewing. They are optional,
// but when given they must survive the round trip and appear on the page: a
// field that quietly drops what the seller typed is worse than no field.
func TestBuildingDetailsSurviveAndShow(t *testing.T) {
	app := newTestApp(t)
	app.createUser("build@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'build@t.test'`)
	cookie := app.login("build@t.test", "Parol12345")

	form := listingForm(ruNodeID(t, app))
	form.Set("build_year", "1998")
	form.Set("wall_material", "panel")
	form.Set("ceiling_height", "2,7") // a comma is what a Russian keyboard produces
	if w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("listing must save, got %d (%s)", w.Code, w.Body.String())
	}

	var year int
	var wall string
	var ceiling float64
	var id string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id::text, build_year, wall_material, ceiling_height FROM listings WHERE title_ru = $1`,
		"Квартира в аренду").Scan(&id, &year, &wall, &ceiling); err != nil {
		t.Fatalf("listing not stored: %v", err)
	}
	if year != 1998 || wall != "panel" || ceiling != 2.7 {
		t.Errorf("stored %d / %q / %v, want 1998 / panel / 2.7", year, wall, ceiling)
	}

	// The form is what the seller touches, so it matters more than the page.
	form2 := app.do(http.MethodGet, "/listings/new", nil, withCookie(cookie)).Body.String()
	for _, want := range []string{`name="build_year"`, `name="wall_material"`, `name="ceiling_height"`,
		T(LangRU, "re.wall.frame_reed"), T(LangRU, "re.not_stated")} {
		if !strings.Contains(form2, want) {
			t.Errorf("the posting form is missing %q", want)
		}
	}

	body := app.do(http.MethodGet, "/listings/"+id, nil).Body.String()
	for _, want := range []string{T(LangRU, "re.build_year"), "1998",
		T(LangRU, "re.wall.panel"), T(LangRU, "re.ceiling_height")} {
		if !strings.Contains(body, want) {
			t.Errorf("the listing page does not show %q", want)
		}
	}
}

// Out of range is dropped, not clamped. Turning a typo into a plausible number
// would publish a measurement nobody took, and a reader cannot tell the two
// apart afterwards.
func TestImplausibleBuildingDetailsAreDropped(t *testing.T) {
	app := newTestApp(t)
	app.createUser("build2@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'build2@t.test'`)
	cookie := app.login("build2@t.test", "Parol12345")

	form := listingForm(ruNodeID(t, app))
	form.Set("build_year", "1200")          // older than any dwelling advertised here
	form.Set("wall_material", "adamantium") // not in the list
	form.Set("ceiling_height", "27")        // metres, not centimetres
	if w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("the listing itself is valid and must save, got %d", w.Code)
	}

	var year int
	var wall string
	var ceiling float64
	if err := app.pool.QueryRow(context.Background(),
		`SELECT build_year, wall_material, ceiling_height FROM listings WHERE title_ru = $1`,
		"Квартира в аренду").Scan(&year, &wall, &ceiling); err != nil {
		t.Fatalf("listing not stored: %v", err)
	}
	if year != 0 || wall != "" || ceiling != 0 {
		t.Errorf("stored %d / %q / %v, want everything unset", year, wall, ceiling)
	}
}

// A flat in an unfinished building is advertised with the year it will be
// finished, so the ceiling has to sit ahead of today — but not indefinitely.
func TestBuildYearWindow(t *testing.T) {
	now := time.Now().Year()
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"1850", 1850}, {"1849", 0},
		{"2015", 2015},
		{strconv.Itoa(now + 5), now + 5}, {strconv.Itoa(now + 6), 0},
		{"", 0}, {"не знаю", 0},
	} {
		if got := parseBuildYear(tc.in); got != tc.want {
			t.Errorf("parseBuildYear(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
