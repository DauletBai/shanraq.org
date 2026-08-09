package articles

import (
	"context"
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
