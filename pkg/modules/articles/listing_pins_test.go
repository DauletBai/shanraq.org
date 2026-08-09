package articles

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// A listing must reach the map from its address alone. Two ways this failed at
// once: the query still named a renamed column, so it errored and the map went
// blank; and no district in the reference carries coordinates, so a flat in
// Медеу or Петродворцовый had nothing to plot even once the query ran. The pin
// now climbs to the nearest ancestor that has coordinates.
func TestListingPinsReachTheMapFromADistrict(t *testing.T) {
	app := newTestApp(t)
	app.createUser("pins@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'pins@t.test'`)
	cookie := app.login("pins@t.test", "Parol12345")

	// Медеу is a district of Almaty. The district has no coordinates of its own;
	// Almaty above it does.
	var node uuid.UUID
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM geo_nodes WHERE name_ru = 'Медеу' AND kind = 'district'`).Scan(&node); err != nil {
		t.Fatal(err)
	}
	var districtHasCoords bool
	if err := app.pool.QueryRow(context.Background(),
		`SELECT lat IS NOT NULL FROM geo_nodes WHERE id = $1`, node).Scan(&districtHasCoords); err != nil {
		t.Fatal(err)
	}
	if districtHasCoords {
		t.Skip("the district gained coordinates; this test guards the case where it has none")
	}

	form := listingForm(node.String())
	form.Set("title_ru", "Квартира у Медеу")
	form.Set("currency", "KZT")
	if w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("post = %d (%s)", w.Code, w.Body.String())
	}

	pins, err := NewListingStore(app.pool).ListingPins(context.Background(), 100)
	if err != nil {
		t.Fatalf("the pins query must run — a broken one empties the whole map: %v", err)
	}
	var found *ListingPin
	for i := range pins {
		if pins[i].Title == "Квартира у Медеу" {
			found = &pins[i]
			break
		}
	}
	if found == nil {
		t.Fatal("a listing with a full address produced no marker")
	}
	if found.Lat == 0 || found.Lng == 0 {
		t.Errorf("marker has no position: %+v", found)
	}
	// Almaty is at roughly 43.24 N, 76.89 E — the pin must land on the city, not
	// at the origin or on some other continent.
	if found.Lat < 42 || found.Lat > 44 || found.Lng < 76 || found.Lng > 78 {
		t.Errorf("marker at %.2f,%.2f is not near Almaty", found.Lat, found.Lng)
	}
	if found.Cur != "KZT" {
		t.Errorf("pin currency = %q, want KZT — the popup prints the symbol from it", found.Cur)
	}
}
