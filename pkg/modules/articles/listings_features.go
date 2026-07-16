package articles

// amenityKeys is the ordered set of listing amenities offered as checkboxes on
// the submission form. Labels live in i18n as "re.am_<key>".
var amenityKeys = []string{
	"air_conditioner", "pool", "parking", "garage", "furniture", "fridge",
	"washer", "internet", "tv", "security", "elevator", "heating",
	"hot_water", "plastic_windows", "playground", "gas",
}

var amenitySet = func() map[string]bool {
	m := make(map[string]bool, len(amenityKeys))
	for _, k := range amenityKeys {
		m[k] = true
	}
	return m
}()

// AmenityKeys exposes the amenity list to templates (for rendering checkboxes).
func AmenityKeys() []string { return amenityKeys }

// HasAmenity reports whether a listing carries a given amenity (template helper).
func (l Listing) HasAmenity(key string) bool { return containsStr(l.Amenities, key) }

// HasAmenity lets the submission form re-check boxes after a validation error.
func (in ListingInput) HasAmenity(key string) bool { return containsStr(in.Amenities, key) }

func containsStr(ss []string, key string) bool {
	for _, s := range ss {
		if s == key {
			return true
		}
	}
	return false
}
