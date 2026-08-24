package articles

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Tariffs — the ad rate card, listing service prices and durations — move out of
// Go into an editable key→value table so an operator can retune them from the
// admin panel without a redeploy. Reads go through tariffVal, which returns the
// operator's value from a cached snapshot or falls back to the built-in default,
// so a missing row (or an unloaded cache) never breaks pricing.

// maxTariffValue caps any editable tariff. Real prices are ≤ ~90 000 and
// weights ≤ ~20, so 10 000 000 is generous headroom yet keeps rate×weight far
// inside int64.
const maxTariffValue = 10_000_000

type tariffDef struct {
	key    string
	def    int64
	minOne bool // durations and weights must be >= 1; prices may be 0 (free)
}

// tariffDefs is the single source of built-in defaults, seed rows and form
// order. Keys are stable identifiers; the built-in values mirror the figures the
// code shipped with.
var tariffDefs = []tariffDef{
	// Ad rate card: ad.<format>.<days> (nationwide price on a ×1.0 surface,
	// tenge). These are the rates for the whole of Kazakhstan; geography only
	// discounts them.
	//
	// The card was cut to roughly a seventh of its first version on 24 August
	// 2026, after the market was actually looked at rather than guessed. What
	// the market charges today:
	//
	//   zakon.kz          300 000 ₸ a day for a text-graphic block at full traffic
	//   informburo.kz     1 650 ₸ per thousand impressions, plus 25% for targeting
	//   cifrum.kz         60 000 ₸ for seven days of a banner in one rubric
	//   ng.kz (Kostanay)  54 000 ₸ a month at the top of the front page,
	//                     39 600 in the middle, 12 600 at the foot — against
	//                     330 000 unique visitors a month
	//   ibirzha.kz        4 000 ₸ a week, about 17 000 a month
	//
	// The old card asked 180 000 ₸ for a month of the top banner on our own
	// front page: three and a third times what the strongest paper in Kostanay
	// charges for the same month, on a fraction of its readership. Nobody who
	// checked would have paid it.
	//
	// The top format nationwide is now 25 000 ₸ a month on the front page
	// (12 500 × the ×2.0 surface weight) — under half of ng.kz's regional rate,
	// for the whole country. Bought for Kostanay region it comes to about
	// 5 000 ₸, eleven times under their price in their own city.
	//
	// The card prices a slot, not an impression, which is how a local paper
	// sells and what an advertiser here is actually buying: a place in the
	// publication for a period. These are tariffs, and they move from the admin
	// panel as the audience grows.
	{"ad.horizontal.3", 2500, false}, {"ad.horizontal.7", 5000, false}, {"ad.horizontal.14", 8300, false}, {"ad.horizontal.30", 12500, false},
	{"ad.vertical.3", 1900, false}, {"ad.vertical.7", 3900, false}, {"ad.vertical.14", 6500, false}, {"ad.vertical.30", 9800, false},
	{"ad.square.3", 1300, false}, {"ad.square.7", 2500, false}, {"ad.square.14", 4200, false}, {"ad.square.30", 6300, false},
	{"ad.rectangle.3", 1100, false}, {"ad.rectangle.7", 2200, false}, {"ad.rectangle.14", 3700, false}, {"ad.rectangle.30", 5500, false},
	// Geography. exponent is a percentage: 100 prices strictly in proportion to
	// population, 50 takes the square root of the share, 33 the cube root. The
	// card is built on 50 — see ads_geo.go for why proportion is the wrong
	// answer. min_price is the floor under any order, because below it the
	// booking costs more to handle than it brings.
	{"geo.exponent", 50, true}, {"geo.min_price", 1000, false},
	// Surface weight tiers (×10).
	{"weight.high", 20, true}, {"weight.mid", 13, true}, {"weight.base", 10, true},
	// Listing sidebar banner: banner.<days> (tenge).
	{"banner.1", 990, false}, {"banner.2", 1890, false}, {"banner.3", 2690, false}, {"banner.4", 3390, false}, {"banner.5", 3990, false}, {"banner.6", 4490, false}, {"banner.7", 4990, false},
	// Listing services (tenge).
	{"promote.price", 299, false}, {"feature.price", 499, false},
	// Durations (days).
	{"listing.free_days", 21, true}, {"promote.days", 7, true}, {"feature.days", 7, true},
}

var builtinTariffs = func() map[string]int64 {
	m := make(map[string]int64, len(tariffDefs))
	for _, d := range tariffDefs {
		m[d.key] = d.def
	}
	return m
}()

func isTariffKey(key string) bool { _, ok := builtinTariffs[key]; return ok }

func tariffMinOne(key string) bool {
	for _, d := range tariffDefs {
		if d.key == key {
			return d.minOne
		}
	}
	return false
}

// tariffCache holds the latest snapshot of operator values, swapped atomically
// so the pure pricing functions (called from templates) stay lock-free.
var tariffCache atomic.Pointer[map[string]int64]

// tariffVal is the effective value for a key: the operator override, else the
// built-in default.
func tariffVal(key string) int64 {
	if m := tariffCache.Load(); m != nil {
		if v, ok := (*m)[key]; ok {
			return v
		}
	}
	return builtinTariffs[key]
}

// ---- typed accessors used by the pricing code ----

func adFormatPriceVal(format string, days int) int64 {
	return tariffVal("ad." + format + "." + strconv.Itoa(days))
}
func surfaceWeightVal(tier string) int64 { return tariffVal("weight." + tier) }

// geoExponentVal is the ladder's steepness, as a percentage.
func geoExponentVal() int64 { return tariffVal("geo.exponent") }

// geoMinPriceVal is the floor under a geo-discounted order, in tenge.
func geoMinPriceVal() int64 { return tariffVal("geo.min_price") }
func bannerPriceVal(days int) int64 {
	if days < 1 || days > 7 {
		days = 1
	}
	return tariffVal("banner." + strconv.Itoa(days))
}
func promotePriceVal() int64 { return tariffVal("promote.price") }
func featurePriceVal() int64 { return tariffVal("feature.price") }
func freeDaysVal() int64     { return clampDays(tariffVal("listing.free_days")) }
func promoteDaysVal() int64  { return clampDays(tariffVal("promote.days")) }
func featureDaysVal() int64  { return clampDays(tariffVal("feature.days")) }

// clampDays keeps a duration a sane positive integer even if a bad value slipped
// into the table — it is inlined into SQL, so it must never be negative.
func clampDays(d int64) int64 {
	if d < 1 {
		return 1
	}
	if d > 3650 {
		return 3650
	}
	return d
}

// PromotePrice/FeaturePrice expose the listing service prices to templates.
func PromotePrice() int64 { return promotePriceVal() }
func FeaturePrice() int64 { return featurePriceVal() }

// ---- store ----

// TariffStore seeds, caches and persists the tariffs.
type TariffStore struct{ db *pgxpool.Pool }

// NewTariffStore builds a TariffStore over the shared pgx pool.
func NewTariffStore(db *pgxpool.Pool) *TariffStore { return &TariffStore{db: db} }

// Load seeds any missing defaults (idempotent) and loads all rows into the
// cache. Best-effort at the call site: on failure the cache stays nil and
// pricing serves the built-in defaults.
func (s *TariffStore) Load(ctx context.Context) error {
	for _, d := range tariffDefs {
		if _, err := s.db.Exec(ctx,
			`INSERT INTO tariffs (key, value) VALUES ($1,$2) ON CONFLICT (key) DO NOTHING`,
			d.key, d.def); err != nil {
			return fmt.Errorf("seed tariff %s: %w", d.key, err)
		}
	}
	return s.refresh(ctx)
}

func (s *TariffStore) refresh(ctx context.Context) error {
	rows, err := s.db.Query(ctx, `SELECT key, value FROM tariffs`)
	if err != nil {
		return fmt.Errorf("load tariffs: %w", err)
	}
	defer rows.Close()
	m := make(map[string]int64)
	for rows.Next() {
		var k string
		var v int64
		if err := rows.Scan(&k, &v); err != nil {
			return err
		}
		m[k] = v
	}
	if err := rows.Err(); err != nil {
		return err
	}
	tariffCache.Store(&m)
	return nil
}

// SaveMany validates and upserts a batch of tariffs in a single transaction, so
// a mid-batch failure leaves NO partial price change, then refreshes the cache
// only after the commit. Unknown keys are ignored; values are clamped (>= 0,
// >= 1 for durations/weights, and <= maxTariffValue).
func (s *TariffStore) SaveMany(ctx context.Context, values map[string]int64, by *uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tariffs tx: %w", err)
	}
	defer tx.Rollback(ctx)
	for key, value := range values {
		if !isTariffKey(key) {
			continue
		}
		if value < 0 {
			value = 0
		}
		if tariffMinOne(key) && value < 1 {
			value = 1
		}
		if value > maxTariffValue {
			value = maxTariffValue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO tariffs (key, value, updated_at, updated_by) VALUES ($1,$2,NOW(),$3)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW(), updated_by = EXCLUDED.updated_by`,
			key, value, by); err != nil {
			return fmt.Errorf("save tariff %s: %w", key, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tariffs: %w", err)
	}
	// Cache is refreshed only after a successful commit, so a rollback can never
	// leave the running prices out of sync with the table.
	return s.refresh(ctx)
}
