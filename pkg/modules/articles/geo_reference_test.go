package articles

import (
	"context"
	"net/http"
	"testing"
)

// The location reference is data, not code, and a gap in it is invisible until
// someone cannot find their own district — which is how Петродворцовый was
// found missing from Saint Petersburg after nine of its eighteen districts had
// been entered. These are spot checks on the places most likely to be used, so
// a future edit that drops rows fails here instead of in front of a seller.
func TestGeoReferenceCoversCityDistricts(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()

	// city -> how many districts it must have, and one that must be among them.
	cases := []struct {
		city    string
		want    int
		example string
	}{
		{"Санкт-Петербург", 18, "Петродворцовый"},
		{"Москва", 12, "Новомосковский АО"},
		{"Новосибирск", 10, "Заельцовский"},
		{"Екатеринбург", 7, "Верх-Исетский"},
		{"Казань", 7, "Ново-Савиновский"},
		// The home market: these were already complete and must stay so.
		{"Алматы", 8, "Наурызбай"},
		{"Астана", 5, "Сарыарка"},
		{"Шымкент", 5, "Аль-Фараби"},
		{"Караганда", 2, "Казыбекбийский"},
	}
	for _, c := range cases {
		var n int
		if err := app.pool.QueryRow(ctx, `
			SELECT count(*) FROM geo_nodes ch
			JOIN geo_nodes p ON p.id = ch.parent_id
			WHERE p.name_ru = $1`, c.city).Scan(&n); err != nil {
			t.Fatalf("%s: %v", c.city, err)
		}
		if n != c.want {
			t.Errorf("%s has %d districts, want %d", c.city, n, c.want)
		}
		var found bool
		if err := app.pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM geo_nodes ch
				JOIN geo_nodes p ON p.id = ch.parent_id
				WHERE p.name_ru = $1 AND ch.name_ru = $2)`, c.city, c.example).Scan(&found); err != nil {
			t.Fatalf("%s/%s: %v", c.city, c.example, err)
		}
		if !found {
			t.Errorf("%s is missing %q", c.city, c.example)
		}
	}
}

// Address fields are filled from a node's kind, never its depth, because the
// tree is not the same shape everywhere. A federal city sits exactly where an
// oblast sits, and filling by depth published "Область: Санкт-Петербург, Город:
// Петродворцовый" — an oblast that does not exist and a district called a city.
// The two shapes are pinned here side by side.
func TestListingAddressFollowsNodeKind(t *testing.T) {
	app := newTestApp(t)
	app.createUser("addr@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'addr@t.test'`)
	cookie := app.login("addr@t.test", "Parol12345")

	cases := []struct {
		name                            string
		node                            string // name_ru of the leaf the author picks
		country, region, city, district string
	}{
		// A federal city has no oblast above it — the region row must stay empty
		// rather than borrow the city's name.
		{"федеральный город", "Петродворцовый", "Россия", "", "Санкт-Петербург", "Петродворцовый"},
		// An ordinary region nests one level deeper and fills all four.
		{"обычная область", "Заельцовский", "Россия", "Новосибирская область", "Новосибирск", "Заельцовский"},
		// The same two shapes exist at home and must behave identically.
		{"Алматы", "Медеу", "Казахстан", "", "Алматы", "Медеу"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var node string
			if err := app.pool.QueryRow(context.Background(),
				`SELECT id FROM geo_nodes WHERE name_ru = $1 AND kind = 'district' LIMIT 1`, c.node).Scan(&node); err != nil {
				t.Fatalf("%s not in the reference: %v", c.node, err)
			}
			title := "Квартира " + c.node
			form := listingForm(node)
			form.Set("title_ru", title)
			if w := app.do(http.MethodPost, "/listings/new", form, withCookie(cookie)); w.Code != http.StatusSeeOther {
				t.Fatalf("post = %d (%s)", w.Code, w.Body.String())
			}
			var country, region, city, district string
			if err := app.pool.QueryRow(context.Background(),
				`SELECT country, region, city, district FROM listings WHERE title_ru = $1`, title).
				Scan(&country, &region, &city, &district); err != nil {
				t.Fatal(err)
			}
			got := []string{country, region, city, district}
			want := []string{c.country, c.region, c.city, c.district}
			for i, label := range []string{"страна", "область", "город", "район"} {
				if got[i] != want[i] {
					t.Errorf("%s = %q, want %q", label, got[i], want[i])
				}
			}
		})
	}
}

// Every node must hang off its parent. A district whose parent_code names a row
// that does not exist would be invisible in the cascade — present in the table,
// unreachable in the picker, and impossible to debug from the UI.
func TestGeoReferenceHasNoOrphans(t *testing.T) {
	app := newTestApp(t)
	var orphans int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM geo_nodes WHERE parent_code IS NOT NULL AND parent_id IS NULL`).
		Scan(&orphans); err != nil {
		t.Fatal(err)
	}
	if orphans != 0 {
		t.Errorf("%d location(s) have a parent_code that resolves to nothing", orphans)
	}
}
