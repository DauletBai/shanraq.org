package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// ruNodeID returns the id of the Russia country node from the geo reference.
func ruNodeID(t *testing.T, a *testApp) string {
	t.Helper()
	var id uuid.UUID
	if err := a.pool.QueryRow(context.Background(),
		`SELECT id FROM geo_nodes WHERE parent_id IS NULL AND country = 'RU'`).Scan(&id); err != nil {
		t.Fatalf("Russia must exist in the geo reference: %v", err)
	}
	return id.String()
}

// listingForm is a complete, valid submission. Tests copy it and change the one
// field under test, so a failure points at that field and not at a typo.
func listingForm(geoNode string) url.Values {
	return url.Values{
		"deal_type":      {"rent"},
		"property_type":  {"apartment"},
		"geo_node_id":    {geoNode},
		"price":          {"90000"},
		"area":           {"56"},
		"rooms":          {"2"},
		"title_kz":       {"Жалға пәтер"},
		"title_ru":       {"Квартира в аренду"},
		"title_en":       {"Apartment for rent"},
		"description_ru": {"Светлая квартира"},
		"contact":        {"+7 777 000 00 00"},
		"no_filters":     {"on"},
	}
}

// A price entered against a Russian address is a price in rubles. Storing it as
// tenge is not a cosmetic slip: the same number means a 5x different amount of
// money, so the listing lies to every reader. The currency must follow the
// country the author actually picked.
func TestListingCurrencyFollowsCountry(t *testing.T) {
	app := newTestApp(t)
	app.createUser("ru-listing@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'ru-listing@t.test'`)
	cookie := app.login("ru-listing@t.test", "Parol12345")

	form := listingForm(ruNodeID(t, app))
	w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("a Russian listing must save, got %d (%s)", w.Code, w.Body.String())
	}

	var currency, country string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT currency, country FROM listings WHERE title_ru = $1`, "Квартира в аренду").
		Scan(&currency, &country); err != nil {
		t.Fatalf("listing not stored: %v", err)
	}
	if currency != "RUB" {
		t.Errorf("currency = %q, want RUB — the author picked %q and typed a ruble price", currency, country)
	}
}

// A listing with no location has no country, and a price with no country
// defaults to tenge. Requiring the address is what makes the currency knowable
// at all, so the form must refuse rather than guess.
func TestListingRequiresLocation(t *testing.T) {
	app := newTestApp(t)
	app.createUser("no-geo@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'no-geo@t.test'`)
	cookie := app.login("no-geo@t.test", "Parol12345")

	form := listingForm("")
	form.Del("geo_node_id")
	w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("a listing with no location must be refused, got %d", w.Code)
	}
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM listings WHERE title_ru = $1`, "Квартира в аренду").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Error("a listing with no location was stored anyway")
	}
}

// Editing exists so a wrong price can be fixed. The owner hit exactly this and
// found no way out of the mistake but to wait for the listing to expire.
func TestListingEditFixesPriceAndCurrency(t *testing.T) {
	app := newTestApp(t)
	app.createUser("edit-mine@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'edit-mine@t.test'`)
	cookie := app.login("edit-mine@t.test", "Parol12345")

	ru := ruNodeID(t, app)
	if w := app.do(http.MethodPost, "/listings/new", listingForm(ru), withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("setup listing: %d (%s)", w.Code, w.Body.String())
	}
	var id string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM listings WHERE title_ru = $1`, "Квартира в аренду").Scan(&id); err != nil {
		t.Fatal(err)
	}

	// The edit form must open filled in, or "editing" means retyping everything.
	w := app.do(http.MethodGet, "/listings/"+id+"/edit", nil, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("the edit form must open for its owner, got %d", w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, "Квартира в аренду") || !strings.Contains(body, ru) {
		t.Error("the edit form opened blank instead of showing the listing")
	}

	form := listingForm(ru)
	form.Set("price", "120000")
	form.Set("currency", "KZT") // the author corrects the currency by hand
	if w := app.do(http.MethodPost, "/listings/"+id+"/edit", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("saving an edit = %d (%s)", w.Code, w.Body.String())
	}

	var price int64
	var currency string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT price, currency FROM listings WHERE id = $1`, id).Scan(&price, &currency); err != nil {
		t.Fatal(err)
	}
	if price != 120000 || currency != "KZT" {
		t.Errorf("after the edit: price=%d currency=%q, want 120000/KZT", price, currency)
	}
}

// Delete and edit are owner-only. A listing id is in every public URL, so the
// only thing standing between a stranger and someone else's listing is this
// check — and a missing one would let anyone erase the classifieds section.
func TestListingEditAndDeleteAreOwnerOnly(t *testing.T) {
	app := newTestApp(t)
	app.createUser("owner@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'owner@t.test'`)
	owner := app.login("owner@t.test", "Parol12345")

	ru := ruNodeID(t, app)
	if w := app.do(http.MethodPost, "/listings/new", listingForm(ru), withCookie(owner)); w.Code != http.StatusSeeOther {
		t.Fatalf("setup listing: %d", w.Code)
	}
	var id string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM listings WHERE title_ru = $1`, "Квартира в аренду").Scan(&id); err != nil {
		t.Fatal(err)
	}

	app.createUser("stranger@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'stranger@t.test'`)
	stranger := app.login("stranger@t.test", "Parol12345")

	for _, c := range []struct {
		method, path string
		form         url.Values
	}{
		{http.MethodGet, "/listings/" + id + "/edit", nil},
		{http.MethodPost, "/listings/" + id + "/edit", listingForm(ru)},
		{http.MethodPost, "/listings/" + id + "/delete", url.Values{}},
	} {
		if w := app.do(c.method, c.path, c.form, withCookie(stranger)); w.Code != http.StatusNotFound {
			t.Errorf("%s %s by a stranger = %d, want 404", c.method, c.path, w.Code)
		}
	}
	var alive int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM listings WHERE id = $1`, id).Scan(&alive); err != nil {
		t.Fatal(err)
	}
	if alive != 1 {
		t.Fatal("a stranger destroyed or altered someone else's listing")
	}

	// The owner can delete, and the bookmark rows go with it: favorites carry no
	// foreign key, so nothing but this sweep stops them pointing at nothing.
	app.exec(`INSERT INTO favorites (user_id, item_type, item_id)
	          SELECT id, 'listing', $1 FROM auth_users WHERE lower(email) = 'stranger@t.test'`, id)
	if w := app.do(http.MethodPost, "/listings/"+id+"/delete", url.Values{}, withCookie(owner)); w.Code != http.StatusSeeOther {
		t.Fatalf("the owner must be able to delete, got %d", w.Code)
	}
	var gone, orphans int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT (SELECT count(*) FROM listings WHERE id = $1),
		        (SELECT count(*) FROM favorites WHERE item_type = 'listing' AND item_id = $1)`, id).
		Scan(&gone, &orphans); err != nil {
		t.Fatal(err)
	}
	if gone != 0 {
		t.Error("the listing survived its owner's delete")
	}
	if orphans != 0 {
		t.Error("bookmarks were left pointing at a deleted listing")
	}
}

// The bug the owner hit: a submission that bounces off validation comes back
// with the location silently cleared, because the hidden geo field is rendered
// without its value. The author fixes the flagged field, submits again, and the
// listing saves against no country at all — which means tenge, whatever they
// picked the first time. The re-rendered form must carry the choice back.
func TestListingFormKeepsLocationAfterError(t *testing.T) {
	app := newTestApp(t)
	app.createUser("keep-geo@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'keep-geo@t.test'`)
	cookie := app.login("keep-geo@t.test", "Parol12345")

	ru := ruNodeID(t, app)
	form := listingForm(ru)
	form.Del("contact") // one missing required field is enough to bounce it

	w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("an incomplete form must re-render, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, ru) {
		t.Errorf("the re-rendered form dropped the chosen location (%s): the author's "+
			"country is gone and the next submit will be priced in tenge", ru)
	}
	if !strings.Contains(body, "90000") {
		t.Error("the re-rendered form dropped the price the author typed")
	}
}
