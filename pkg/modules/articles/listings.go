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
	ID           string
	AuthorEmail  string
	DealType     string
	PropertyType string
	Country      string
	Region       string
	City         string
	Village      string
	Price        int64
	Area         float64
	Rooms        int
	Title        string
	Description  string
	Contact      string
	CoverURL     string
	Status       string
	CreatedAt    time.Time
}

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
