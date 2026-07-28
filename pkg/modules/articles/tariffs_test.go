package articles

import "testing"

// With no operator overrides loaded, every pricing read must return the exact
// figure the code shipped with — the move to an editable table must not change
// a single default price or period.
func TestTariffDefaultsUnchanged(t *testing.T) {
	tariffCache.Store(nil)
	cases := map[string]int64{
		"ad.horizontal.30": 90000, "ad.vertical.3": 14000, "ad.square.14": 30000, "ad.rectangle.7": 16000,
		"weight.high": 20, "weight.mid": 13, "weight.base": 10,
		"banner.1": 990, "banner.3": 2690, "banner.7": 4990,
		"promote.price": 299, "feature.price": 499,
		"listing.free_days": 21, "promote.days": 7, "feature.days": 7,
	}
	for key, want := range cases {
		if got := tariffVal(key); got != want {
			t.Errorf("tariffVal(%q) = %d, want %d", key, got, want)
		}
	}
}

func TestPricingAccessorsMatchDefaults(t *testing.T) {
	tariffCache.Store(nil)
	if got := adFormatPriceVal("rectangle", 30); got != 40000 {
		t.Errorf("adFormatPriceVal(rectangle,30) = %d, want 40000", got)
	}
	if got := listingBannerPrice(5); got != 3990 {
		t.Errorf("listingBannerPrice(5) = %d, want 3990", got)
	}
	if got := adSurfaceWeight10(surfaceHome); got != 20 {
		t.Errorf("adSurfaceWeight10(home) = %d, want 20", got)
	}
	if got := PromotePrice(); got != 299 {
		t.Errorf("PromotePrice() = %d, want 299", got)
	}
	if got := FeaturePrice(); got != 499 {
		t.Errorf("FeaturePrice() = %d, want 499", got)
	}
	if got := freeDaysVal(); got != 21 {
		t.Errorf("freeDaysVal() = %d, want 21", got)
	}
}

// The pricing formula (base × summed surface weights / 10) must be unchanged.
func TestAdOrderTotalUnchanged(t *testing.T) {
	tariffCache.Store(nil)
	// horizontal 30d = 90000 base; home(20) + realestate(20) = weight 40; /10.
	p := AdOrderTotal("horizontal", []string{surfaceHome, surfaceRealestate}, 30)
	if want := int64(90000 * 40 / 10); p.Total != want {
		t.Errorf("AdOrderTotal.Total = %d, want %d", p.Total, want)
	}
}

// An override in the cache must win over the built-in default.
func TestTariffOverrideWins(t *testing.T) {
	m := map[string]int64{"promote.price": 777}
	tariffCache.Store(&m)
	defer tariffCache.Store(nil)
	if got := PromotePrice(); got != 777 {
		t.Errorf("PromotePrice() with override = %d, want 777", got)
	}
	// A key absent from the override still falls back to the built-in default.
	if got := tariffVal("feature.price"); got != 499 {
		t.Errorf("tariffVal(feature.price) fallback = %d, want 499", got)
	}
}
