package articles

import (
	"strings"
	"time"
)

// Real-estate taxonomy (property types are the "real estate categories").
var (
	DealTypes     = []string{"sale", "rent"}
	PropertyTypes = []string{"apartment", "house", "land", "commercial", "dacha"}
)

func isDealType(s string) bool {
	for _, d := range DealTypes {
		if d == s {
			return true
		}
	}
	return false
}

func isPropertyType(s string) bool {
	for _, p := range PropertyTypes {
		if p == s {
			return true
		}
	}
	return false
}

// Listing is one real-estate advert.
type Listing struct {
	ID            string
	AuthorID      string
	AuthorEmail   string
	DealType      string
	PropertyType  string
	Country       string
	Region        string
	City          string
	District      string
	Microdistrict string
	Street        string
	House         string
	Lat           *float64
	Lng           *float64
	Price         int64
	Currency      string // "KZT" | "RUB" — display currency, derived from country
	Area          float64
	LandArea      float64
	Rooms         int
	RoomSpecs     []RoomSpec
	Amenities     []string
	Title         string // base/fallback (usually RU)
	Description   string
	TitleKz       string
	TitleRu       string
	TitleEn       string
	DescriptionKz string
	DescriptionRu string
	DescriptionEn string
	Contact       string
	CoverURL      string
	Images        []string
	Documents     []string
	// ContractURL is the lease the landlord will sign, published up front. Rent
	// only — see migration 20251107009600 for why sales carry title papers in
	// Documents instead.
	ContractURL string
	Status      string
	ViewsCount  int
	// Reports is how many distinct people have reported this listing. Shown to
	// its owner before the auto-hide threshold, so the warning arrives while
	// there is still time to fix the photos.
	Reports       int
	ContactsCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ExpiresAt     time.Time
	// GeoNodeID is the location the author picked in the cascade. Carried on the
	// listing so the edit form can restore that choice instead of making the
	// author re-pick their own address.
	GeoNodeID     string
	PromotedUntil *time.Time
	FeaturedUntil *time.Time
	BannerUntil   *time.Time
	// Agent attribution: set when the listing's owner is a registered
	// real-estate agent, so cards and the listing page can badge it and link to
	// the agent's public page.
	AgentID   string
	AgentName string
}

// TitleIn returns the title in the requested language, falling back to RU, then
// any non-empty language, then the base title — so a card never shows blank.
// CurrencySymbol returns the price symbol for the listing's currency: the ruble
// for Russian listings, the tenge otherwise (the default).
func (l Listing) CurrencySymbol() string {
	if l.Currency == "RUB" {
		return "₽"
	}
	return "₸"
}

func (l Listing) TitleIn(lang string) string {
	switch lang {
	case LangKZ:
		if l.TitleKz != "" {
			return l.TitleKz
		}
	case LangEN:
		if l.TitleEn != "" {
			return l.TitleEn
		}
	}
	for _, c := range []string{l.TitleRu, l.TitleKz, l.TitleEn, l.Title} {
		if c != "" {
			return c
		}
	}
	return l.Title
}

// DescriptionIn returns the description in the requested language, with the same
// fallback chain as TitleIn.
func (l Listing) DescriptionIn(lang string) string {
	switch lang {
	case LangKZ:
		if l.DescriptionKz != "" {
			return l.DescriptionKz
		}
	case LangEN:
		if l.DescriptionEn != "" {
			return l.DescriptionEn
		}
	}
	for _, c := range []string{l.DescriptionRu, l.DescriptionKz, l.DescriptionEn, l.Description} {
		if c != "" {
			return c
		}
	}
	return l.Description
}

// ByAgent reports whether this listing was posted by a registered agent.
func (l Listing) ByAgent() bool { return l.AgentID != "" }

// Expired reports whether the listing's free window has ended.
func (l Listing) Expired() bool { return l.ExpiresAt.Before(time.Now()) }

// DaysLeft is the whole days until expiry (rounded up, min 0).
func (l Listing) DaysLeft() int {
	d := time.Until(l.ExpiresAt)
	if d <= 0 {
		return 0
	}
	return int((d + 24*time.Hour - time.Nanosecond) / (24 * time.Hour))
}

// Promoted reports whether the listing is currently boosted to the top.
func (l Listing) Promoted() bool { return l.PromotedUntil != nil && l.PromotedUntil.After(time.Now()) }

// Featured reports whether the listing is currently visually highlighted.
func (l Listing) Featured() bool { return l.FeaturedUntil != nil && l.FeaturedUntil.After(time.Now()) }

// Banner reports whether the listing currently holds a paid sidebar banner slot
// on the real-estate page.
func (l Listing) Banner() bool { return l.BannerUntil != nil && l.BannerUntil.After(time.Now()) }

// listingBannerPrice returns the banner price (tenge) for 1..7 days. Priced
// above Top (299) and highlight (499) since the slot is always on screen, with
// a volume discount as the period grows.
func listingBannerPrice(days int) int64 { return bannerPriceVal(days) }

// BannerDays is the selectable banner period (days), for the form.
func BannerDays() []int { return []int{1, 2, 3, 4, 5, 6, 7} }

// BannerPrice exposes the price to templates.
func BannerPrice(days int) int64 { return listingBannerPrice(days) }

// Address renders the street part of a listing — "мкр Самал-2, ул. Абая, 15" —
// skipping whatever the author left blank. The city and region are rendered
// separately by Location, so this stays purely the within-city part.
func (l Listing) Address() string {
	parts := []string{}
	for _, p := range []string{l.Microdistrict, l.Street, l.House} {
		if p = strings.TrimSpace(p); p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ", ")
}

// HasPin reports whether the listing carries coordinates for a map marker.
func (l Listing) HasPin() bool { return l.Lat != nil && l.Lng != nil }

// Location renders the place parts that are set, most specific first.
func (l Listing) Location() string {
	parts := []string{}
	for _, p := range []string{l.District, l.City, l.Region, l.Country} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
