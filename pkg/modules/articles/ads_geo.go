package articles

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// Geography as the third dimension of the ad rate card.
//
// Format says how big the banner is, surface says which section it runs in, and
// geography says how much of the country it reaches. An advertiser who wants
// the whole of Kazakhstan pays the card rate; one who wants only Kachar pays a
// fraction of it, because only a fraction of the audience is being bought.
//
// The discount runs on the SQUARE ROOT of the population share rather than the
// share itself. Straight proportion prices Kachar at forty tenge a month, which
// is less than it costs to look at the order — and it would also be wrong on
// its own terms, because a banner on a place's pages is seen by more than that
// place's residents, and because reviewing, scheduling and serving a booking
// costs the same whatever its reach. Price therefore falls more slowly than
// population: halve the audience and you pay about seven tenths, not a half.
//
// What is sold is stated plainly and is exactly what is delivered: the banner
// runs on the pages of that place and of every place inside it, for the booked
// period. Nobody is promised a number of impressions, because that number
// depends on what readers do.
//
// The geography also does not exist to raise anyone's price. The card rate is
// the nationwide rate; every choice below it is a discount, and there is no
// tier above it. Kazakhstan's own market runs the other way — informburo.kz
// charges a twenty-five percent premium for targeting — which is defensible for
// a publication whose audience is already national and who therefore gives
// something up by narrowing it. We are not that publication, and a local shop
// in Kachar should not be paying a surcharge to reach its own street.

// adGeoShareUnit is the denominator of a stored population share: shares are
// held in hundred-thousandths so that a village of ten thousand in a country of
// twenty million still lands on a non-zero integer.
const adGeoShareUnit = 100_000

// adGeoMultUnit is the denominator of a price multiplier (×10 000).
const adGeoMultUnit = 10_000

// AdGeo is a place an ad can be bought for, priced.
type AdGeo struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Level int    `json:"level"`
	// Pop is the place's population, and Year the year it was counted. Both are
	// shown to the advertiser: a price derived from a figure has to be able to
	// name the figure and its date.
	Pop  int64 `json:"pop"`
	Year int   `json:"year"`
	// Share of the country's population, in hundred-thousandths.
	Share int64 `json:"share"`
	// Mult is the price multiplier ×10 000.
	Mult int64 `json:"mult"`
}

// adGeoShare converts a population and its country total into a stored share.
// A place larger than the total it is measured against (which happens when the
// two figures come from different years) is clamped: a share above one would
// price a region higher than the country it sits in.
func adGeoShare(pop, countryPop int64) int64 {
	if pop <= 0 || countryPop <= 0 {
		return adGeoShareUnit
	}
	share := int64(math.Round(float64(pop) / float64(countryPop) * adGeoShareUnit))
	if share > adGeoShareUnit {
		return adGeoShareUnit
	}
	if share < 1 {
		return 1
	}
	return share
}

// adGeoMult turns a population share into a price multiplier ×10 000.
//
// The exponent is an editable tariff so the ladder can be retuned without a
// deploy: 100 is straight proportion, 50 is the square root the card is built
// on, 33 a gentler slope still.
func adGeoMult(share int64) int64 {
	if share >= adGeoShareUnit {
		return adGeoMultUnit
	}
	if share <= 0 {
		share = 1
	}
	exp := float64(geoExponentVal()) / 100
	if exp <= 0 || exp > 1 {
		exp = 0.5
	}
	m := math.Pow(float64(share)/adGeoShareUnit, exp) * adGeoMultUnit
	out := int64(math.Round(m))
	if out < 1 {
		return 1
	}
	return out
}

// GeoPlace is the raw population fact behind a place's price.
type GeoPlace struct {
	ID         uuid.UUID
	Slug       string
	Name       string
	Kind       string
	Level      int
	Country    string
	Population int64
	Year       int
}

// AdGeoPlace loads one place and prices it. A place whose population is unknown
// is refused rather than guessed: an invented number here becomes an invoice.
func (s *GeoStore) AdGeoPlace(ctx context.Context, id uuid.UUID, lang string) (*AdGeo, error) {
	name := fmt.Sprintf("COALESCE(NULLIF(c.%s,''), c.name_ru)", geoNameCol(lang))
	var (
		g       AdGeo
		pop     *int64
		year    *int
		country string
	)
	err := s.db.QueryRow(ctx, fmt.Sprintf(`
		SELECT c.id::text, COALESCE(c.slug,''), %s, c.kind, c.level, c.country, c.population, c.population_year
		FROM geo_nodes c WHERE c.id = $1`, name), id).
		Scan(&g.ID, &g.Slug, &g.Name, &g.Kind, &g.Level, &country, &pop, &year)
	if err != nil {
		return nil, fmt.Errorf("ad geo place: %w", err)
	}
	if pop == nil || *pop <= 0 {
		return nil, nil
	}
	g.Pop = *pop
	if year != nil {
		g.Year = *year
	}

	countryPop, err := s.countryPopulation(ctx, country)
	if err != nil {
		return nil, err
	}
	g.Share = adGeoShare(g.Pop, countryPop)
	g.Mult = adGeoMult(g.Share)
	return &g, nil
}

// countryPopulation is the total a place's share is measured against.
func (s *GeoStore) countryPopulation(ctx context.Context, country string) (int64, error) {
	var pop *int64
	err := s.db.QueryRow(ctx,
		`SELECT population FROM geo_nodes WHERE level = 0 AND country = $1`, country).Scan(&pop)
	if err != nil {
		return 0, fmt.Errorf("country population %s: %w", country, err)
	}
	if pop == nil || *pop <= 0 {
		return 0, fmt.Errorf("country population %s is unknown", country)
	}
	return *pop, nil
}

// AdGeoOptions lists the places worth offering an advertiser: everything with a
// known population, largest first inside each level. Districts of a city and
// places we have no count for are left out — they cannot be priced, and a
// picker that offers a choice it cannot quote is worse than one that does not.
func (s *GeoStore) AdGeoOptions(ctx context.Context, country, lang string) ([]AdGeo, error) {
	name := fmt.Sprintf("COALESCE(NULLIF(c.%s,''), c.name_ru)", geoNameCol(lang))
	rows, err := s.db.Query(ctx, fmt.Sprintf(`
		SELECT c.id::text, COALESCE(c.slug,''), %s, c.kind, c.level, c.population,
		       COALESCE(c.population_year, 0)
		FROM geo_nodes c
		WHERE c.country = $1 AND c.population > 0
		ORDER BY c.level, c.population DESC`, name), country)
	if err != nil {
		return nil, fmt.Errorf("ad geo options: %w", err)
	}
	defer rows.Close()

	countryPop, err := s.countryPopulation(ctx, country)
	if err != nil {
		return nil, err
	}
	out := []AdGeo{}
	for rows.Next() {
		var g AdGeo
		if err := rows.Scan(&g.ID, &g.Slug, &g.Name, &g.Kind, &g.Level, &g.Pop, &g.Year); err != nil {
			return nil, err
		}
		g.Share = adGeoShare(g.Pop, countryPop)
		g.Mult = adGeoMult(g.Share)
		out = append(out, g)
	}
	return out, rows.Err()
}

// AdGeoCovers reports whether a booking bought for `bought` should run on a page
// about `page`. Buying a region buys everything inside it, so the test is
// ancestry: the page's place must be the bought place or sit under it.
func (s *GeoStore) AdGeoCovers(ctx context.Context, bought, page uuid.UUID) (bool, error) {
	if bought == uuid.Nil {
		return true, nil // nationwide
	}
	if page == uuid.Nil {
		return false, nil // a page with no place is not inside any region
	}
	var ok bool
	err := s.db.QueryRow(ctx, `
		WITH RECURSIVE up AS (
			SELECT id, parent_id FROM geo_nodes WHERE id = $1
			UNION ALL
			SELECT g.id, g.parent_id FROM geo_nodes g JOIN up ON g.id = up.parent_id
		)
		SELECT EXISTS(SELECT 1 FROM up WHERE id = $2)`, page, bought).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("ad geo covers: %w", err)
	}
	return ok, nil
}
