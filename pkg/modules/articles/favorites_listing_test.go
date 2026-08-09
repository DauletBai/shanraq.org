package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Saving a listing and then finding the favourites page empty is the whole
// feature failing silently. The list query selected the agent-badge columns
// without joining the agent table, so it errored on every call; the handler
// logged it and rendered an empty page, which reads to the user as "the button
// did nothing" even though the bookmark was stored correctly.
func TestFavoriteListingAppearsOnTheFavoritesPage(t *testing.T) {
	app := newTestApp(t)
	owner := app.createUser("favowner@t.test", "Parol12345")
	id := app.seedListing(owner)

	app.createUser("favreader@t.test", "Parol12345")
	cookie := app.login("favreader@t.test", "Parol12345")

	if w := app.do(http.MethodPost, "/favorites/listing/"+id.String(), url.Values{}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("saving a listing = %d, want 303", w.Code)
	}
	// The bookmark is stored — that half always worked.
	var saved int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM favorites WHERE item_type = 'listing' AND item_id = $1`, id).Scan(&saved); err != nil {
		t.Fatal(err)
	}
	if saved != 1 {
		t.Fatalf("the bookmark was not stored (%d rows)", saved)
	}

	// The half that did not: reading them back.
	ls, err := NewListingStore(app.pool).ListFavorited(context.Background(),
		mustUserID(t, app, "favreader@t.test"))
	if err != nil {
		t.Fatalf("listing favourites must load, got: %v", err)
	}
	if len(ls) != 1 {
		t.Fatalf("saved 1 listing, the favourites query returned %d", len(ls))
	}

	w := app.do(http.MethodGet, "/favorites", nil, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("favourites page = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "/listings/"+id.String()) {
		t.Error("the saved listing is missing from the favourites page")
	}
}

func mustUserID(t *testing.T, app *testApp, email string) (id uuid.UUID) {
	t.Helper()
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM auth_users WHERE lower(email) = lower($1)`, email).Scan(&id); err != nil {
		t.Fatalf("user %s: %v", email, err)
	}
	return id
}
