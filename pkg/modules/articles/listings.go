package articles

import "time"

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
	Village       string
	Price         int64
	Area          float64
	LandArea      float64
	Rooms         int
	RoomSpecs     []RoomSpec
	Amenities     []string
	Title         string
	Description   string
	Contact       string
	CoverURL      string
	Images        []string
	Status        string
	ViewsCount    int
	ContactsCount int
	CreatedAt     time.Time
	ExpiresAt     time.Time
	PromotedUntil *time.Time
	FeaturedUntil *time.Time
	BannerUntil   *time.Time
}

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
func listingBannerPrice(days int) int64 {
	switch days {
	case 2:
		return 1890
	case 3:
		return 2690
	case 4:
		return 3390
	case 5:
		return 3990
	case 6:
		return 4490
	case 7:
		return 4990
	default:
		return 990
	}
}

// BannerDays is the selectable banner period (days), for the form.
func BannerDays() []int { return []int{1, 2, 3, 4, 5, 6, 7} }

// BannerPrice exposes the price to templates.
func BannerPrice(days int) int64 { return listingBannerPrice(days) }

// Location renders the place parts that are set, most specific first.
func (l Listing) Location() string {
	parts := []string{}
	for _, p := range []string{l.Village, l.City, l.Region, l.Country} {
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
