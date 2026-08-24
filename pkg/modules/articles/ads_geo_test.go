package articles

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// A share is the place's population over its country's, in hundred-thousandths.
// The clamps matter more than the arithmetic: the two figures come from counts
// taken in different years, so a region can legitimately read as larger than the
// country it sits in, and a share above one would price it above nationwide.
func TestAdGeoShareIsClamped(t *testing.T) {
	cases := []struct {
		name         string
		pop, country int64
		want         int64
	}{
		{"вся страна", 20_843_754, 20_843_754, adGeoShareUnit},
		{"область", 833_643, 20_843_754, 3999},
		{"посёлок", 9_625, 20_843_754, 46},
		{"перепись разных лет", 21_000_000, 20_843_754, adGeoShareUnit},
		{"население неизвестно", 0, 20_843_754, adGeoShareUnit},
		{"страна неизвестна", 9_625, 0, adGeoShareUnit},
		{"меньше единицы округления", 1, 20_843_754, 1},
	}
	for _, c := range cases {
		if got := adGeoShare(c.pop, c.country); got != c.want {
			t.Errorf("%s: adGeoShare(%d, %d) = %d, want %d", c.name, c.pop, c.country, got, c.want)
		}
	}
}

// The ladder runs on the square root of the share, so a place is never as cheap
// as its bare population would make it. Straight proportion put Kachar at forty
// tenge a month.
func TestAdGeoMultTakesTheSquareRoot(t *testing.T) {
	tariffCache.Store(nil)
	cases := []struct{ share, want int64 }{
		{adGeoShareUnit, 10000}, // вся страна — без скидки
		{10692, 3270},           // Алматы
		{3999, 2000},            // Костанайская область
		{46, 214},               // Качар
		// A share of zero cannot happen — adGeoShare floors it at one — but the
		// multiplier guards against it anyway and treats it as the smallest
		// share there is. The order floor then lifts the price regardless.
		{0, 32},
	}
	for _, c := range cases {
		if got := adGeoMult(c.share); got != c.want {
			t.Errorf("adGeoMult(%d) = %d, want %d", c.share, got, c.want)
		}
	}
	// Narrowing must never raise a price: the card rate is the nationwide rate
	// and geography is a discount, not a tier above it.
	prev := int64(0)
	for _, share := range []int64{1, 46, 3999, 10692, adGeoShareUnit} {
		m := adGeoMult(share)
		if m < prev {
			t.Errorf("множитель падает с ростом доли: доля %d дала %d после %d", share, m, prev)
		}
		if m > adGeoMultUnit {
			t.Errorf("доля %d дала множитель %d — дороже, чем вся страна", share, m)
		}
		prev = m
	}
}

// The exponent is a tariff so the ladder can be retuned from the admin panel.
func TestAdGeoExponentIsATariff(t *testing.T) {
	m := map[string]int64{"geo.exponent": 100} // straight proportion
	tariffCache.Store(&m)
	defer tariffCache.Store(nil)
	if got := adGeoMult(3999); got != 400 {
		t.Errorf("при показателе 100 доля 3999 должна дать 400, вышло %d", got)
	}
}

// The floor is the price of handling an order, so it lifts a tiny booking — but
// it must not appear on an empty form, where nothing has been chosen yet.
func TestGeoFloorAppliesOnlyToRealOrders(t *testing.T) {
	tariffCache.Store(nil)
	kachar := adGeoMult(46)

	empty := AdOrderTotalGeo("horizontal", nil, 30, kachar)
	if empty.Total != 0 {
		t.Errorf("пустой заказ стоит %d, ожидался 0", empty.Total)
	}

	tiny := AdOrderTotalGeo("rectangle", []string{surfaceHome}, 3, kachar)
	if !tiny.AtFloor || tiny.Total != geoMinPriceVal() {
		t.Errorf("мелкий заказ должен упереться в пол %d: %+v", geoMinPriceVal(), tiny)
	}

	// Nationwide is reported beside the discounted total, so an advertiser sees
	// what the narrowing saved rather than only a smaller number.
	big := AdOrderTotalGeo("horizontal", []string{surfaceHome}, 30, kachar)
	if big.Nationwide != 135000 {
		t.Errorf("цена за всю страну = %d, ожидалось 135000", big.Nationwide)
	}
	if big.Total >= big.Nationwide {
		t.Errorf("скидка не применилась: %+v", big)
	}
}

// Buying a region buys everything inside it. This is the whole of what geo
// targeting sells, and it is asserted against the real tree.
func TestAdGeoCoversTheSubtree(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()
	g := app.module().geo

	id := func(slug string) uuid.UUID {
		n, err := g.BySlug(ctx, slug, LangRU)
		if err != nil || n == nil {
			t.Skipf("в справочнике нет места %q", slug)
		}
		return uuid.MustParse(n.ID)
	}
	oblast, city, kachar := id("kostanaiskaya-oblast"), id("kostanai"), id("kachar")
	almaty := id("almaty")

	cases := []struct {
		name         string
		bought, page uuid.UUID
		want         bool
	}{
		{"область покрывает свой город", oblast, city, true},
		{"область покрывает посёлок", oblast, kachar, true},
		{"область покрывает саму себя", oblast, oblast, true},
		{"посёлок не покрывает область", kachar, oblast, false},
		{"область не покрывает чужой город", oblast, almaty, false},
		{"вся страна покрывает всё", uuid.Nil, kachar, true},
		{"страница без места — только для заказов без места", oblast, uuid.Nil, false},
	}
	for _, c := range cases {
		got, err := g.AdGeoCovers(ctx, c.bought, c.page)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%s: получилось %v, ожидалось %v", c.name, got, c.want)
		}
	}
}

// A place we have no population for cannot be priced, and is refused rather
// than guessed: an invented number here becomes an invoice.
func TestAdGeoPlaceRefusesUnknownPopulation(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()
	g := app.module().geo

	var id uuid.UUID
	err := app.pool.QueryRow(ctx,
		`SELECT id FROM geo_nodes WHERE population IS NULL AND country = 'KZ' LIMIT 1`).Scan(&id)
	if err != nil {
		t.Skip("в справочнике не осталось мест без населения")
	}
	got, err := g.AdGeoPlace(ctx, id, LangRU)
	if err != nil {
		t.Fatalf("AdGeoPlace: %v", err)
	}
	if got != nil {
		t.Errorf("место без населения получило цену: %+v", got)
	}
}
