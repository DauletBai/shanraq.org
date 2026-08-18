package articles

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
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

// The cap and the label the seller reads come from one constant. They used to
// be two numbers — the const said ten and a translated string said ten — and
// nothing would have caught them drifting apart.
func TestPhotoCapMatchesWhatTheFormPromises(t *testing.T) {
	app := newTestApp(t)
	app.createUser("photos@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'photos@t.test'`)
	cookie := app.login("photos@t.test", "Parol12345")

	form := app.do(http.MethodGet, "/listings/new", nil, withCookie(cookie)).Body.String()
	promise := fmt.Sprintf(T(LangRU, "re.photos_max"), maxListingPhotos)
	if !strings.Contains(form, promise) {
		t.Errorf("the form does not promise %q", promise)
	}
	if strings.Contains(form, "%d") {
		t.Error("the label still carries an unrendered placeholder")
	}

	// Five more than allowed: the extra ones are dropped, the listing still saves.
	f := listingForm(ruNodeID(t, app))
	for i := 0; i < maxListingPhotos+5; i++ {
		f.Add("image", fmt.Sprintf("/media/p%02d.jpg", i))
	}
	if w := app.do(http.MethodPost, "/listings/new", f, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("listing must save, got %d (%s)", w.Code, w.Body.String())
	}
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT cardinality(images) FROM listings WHERE title_ru = $1`, "Квартира в аренду").Scan(&n); err != nil {
		t.Fatalf("listing not stored: %v", err)
	}
	if n != maxListingPhotos {
		t.Errorf("stored %d photos, want the cap of %d", n, maxListingPhotos)
	}
}

// The cap lives in one constant, and the browser has to hear about it too. The
// server was raised to twenty while the uploader kept data-max="10", so a
// seller was stopped at ten by a widget while the form promised twenty and the
// handler would have accepted them. The earlier test checked the handler and
// the label and missed the only thing the seller actually touches.
func TestUploaderHonoursTheSameCap(t *testing.T) {
	app := newTestApp(t)
	app.createUser("upcap@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'upcap@t.test'`)
	cookie := app.login("upcap@t.test", "Parol12345")

	body := app.do(http.MethodGet, "/listings/new", nil, withCookie(cookie)).Body.String()

	// Only the file uploaders, in page order: photos, documents, contract. The
	// room-breakdown widget carries a data-max of its own and must not be
	// mistaken for one of them just because the numbers happen to agree.
	found := regexp.MustCompile(`data-gallery data-max="(\d+)"`).FindAllStringSubmatch(body, -1)
	want := []int{maxListingPhotos, maxListingDocs, 1}
	if len(found) != len(want) {
		t.Fatalf("found %d file uploaders, want %d", len(found), len(want))
	}
	for i, m := range found {
		got, _ := strconv.Atoi(m[1])
		if got != want[i] {
			t.Errorf("uploader %d allows %d files, want %d", i+1, got, want[i])
		}
	}
}

// A listing in Russia showed no flag: the table held Kazakhstan and nothing
// else, so every Russian address rendered a blank where the flag belongs. The
// cascade also stores the country in whichever language the author was reading,
// so all three spellings have to resolve.
func TestEveryCountryWeServeHasAFlag(t *testing.T) {
	for _, name := range []string{"Казахстан", "Қазақстан", "Kazakhstan", "Россия", "Ресей", "Russia"} {
		if got := string(countryFlag(name)); got == "" {
			t.Errorf("countryFlag(%q) is empty — a listing there shows no flag", name)
		} else if !strings.Contains(got, "<svg") {
			t.Errorf("countryFlag(%q) = %q, want an svg", name, got)
		}
	}
	if countryFlag("Атлантида") != "" {
		t.Error("an unknown country must render nothing rather than a wrong flag")
	}
}
