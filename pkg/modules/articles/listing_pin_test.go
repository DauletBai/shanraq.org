package articles

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestListingPinByIDIntegration covers the three answers the map on a listing's
// own page can get, because each one draws something different: the building,
// the settlement, or nothing at all.
//
// The middle case is the one worth a test. A listing with no coordinates of its
// own still has a location — the settlement it was filed under — and the page
// must show it as a settlement, not as an address. Getting Exact wrong there
// puts a sharp marker on a district centroid and tells a buyer a flat is on a
// street nobody typed.
func TestListingPinByIDIntegration(t *testing.T) {
	dsn := requireTestDB(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	author := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO auth_users (id,email,password_hash,role) VALUES ($1,$2,'x','user')`,
		author, "pin-"+author.String()+"@t.test"); err != nil {
		t.Fatalf("user: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, author) })

	// A settlement that knows where it is, and a district that does not — the
	// district is what makes the query climb to its parent.
	town, district := uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, `
		INSERT INTO geo_nodes (id,code,country,level,kind,name_ru,lat,lng)
		VALUES ($1,$2,'KZ',3,'city','Тестгород',53.0000,63.0000)`, town, "pin-town-"+town.String()[:8]); err != nil {
		t.Fatalf("town: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO geo_nodes (id,code,country,level,kind,name_ru,parent_id)
		VALUES ($1,$2,'KZ',4,'district','Тестрайон',$3)`, district, "pin-dist-"+district.String()[:8], town); err != nil {
		t.Fatalf("district: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM geo_nodes WHERE id IN ($1,$2)`, district, town) })

	store := &ListingStore{db: pool}

	mk := func(t *testing.T, geo *uuid.UUID, lat, lng *float64) uuid.UUID {
		t.Helper()
		id := uuid.New()
		if _, err := pool.Exec(ctx, `
			INSERT INTO listings (id,author_id,geo_node_id,lat,lng,title,status,expires_at)
			VALUES ($1,$2,$3,$4,$5,'Тест',
			        'published', NOW() + INTERVAL '7 days')`, id, author, geo, lat, lng); err != nil {
			t.Fatalf("listing: %v", err)
		}
		t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id=$1`, id) })
		return id
	}
	f := func(v float64) *float64 { return &v }

	t.Run("own coordinates win and are exact", func(t *testing.T) {
		id := mk(t, &town, f(53.1234), f(63.4321))
		pin, ok := store.ListingPinByID(ctx, id)
		if !ok {
			t.Fatal("no pin for a listing that carries its own coordinates")
		}
		if !pin.Exact {
			t.Error("Exact = false; a listing with its own lat/lng is exactly placed")
		}
		if pin.Lat != 53.1234 || pin.Lng != 63.4321 {
			t.Errorf("pin at %v,%v — the listing's own coordinates were overridden", pin.Lat, pin.Lng)
		}
	})

	t.Run("district climbs to its parent and is not exact", func(t *testing.T) {
		id := mk(t, &district, nil, nil)
		pin, ok := store.ListingPinByID(ctx, id)
		if !ok {
			t.Fatal("no pin; the district's parent carries coordinates and should have been found")
		}
		if pin.Exact {
			t.Error("Exact = true for a listing that never gave an address — the page would draw a building marker on a town centre")
		}
		if pin.Lat != 53.0 || pin.Lng != 63.0 {
			t.Errorf("pin at %v,%v, want the parent town's 53,63", pin.Lat, pin.Lng)
		}
	})

	t.Run("nothing anywhere means no map", func(t *testing.T) {
		id := mk(t, nil, nil, nil)
		if _, ok := store.ListingPinByID(ctx, id); ok {
			t.Error("got a pin for a listing with no coordinates and no geo node")
		}
	})

	t.Run("expired listings still resolve", func(t *testing.T) {
		id := mk(t, &town, f(53.5), f(63.5))
		if _, err := pool.Exec(ctx, `UPDATE listings SET expires_at = NOW() - INTERVAL '1 day' WHERE id=$1`, id); err != nil {
			t.Fatalf("expire: %v", err)
		}
		if _, ok := store.ListingPinByID(ctx, id); !ok {
			t.Error("no pin for an expired listing; its owner still opens the page and deserves the map")
		}
	})
}
