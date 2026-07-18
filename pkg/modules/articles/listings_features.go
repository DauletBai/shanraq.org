package articles

// RoomSpec is one room of a listing: its kind, floor area (m²) and an optional note.
type RoomSpec struct {
	Type string  `json:"type"`
	Area float64 `json:"area"`
	Note string  `json:"note,omitempty"`
}

// maxRoomSpecs caps how many rooms a listing may detail.
const maxRoomSpecs = 20

// RoomArea returns the floor area of the first room of the given type (0 if
// none) — used to show "зал 20 м²", "кухня 7 м²" in the facts summary.
func (l Listing) RoomArea(roomType string) float64 {
	for _, r := range l.RoomSpecs {
		if r.Type == roomType {
			return r.Area
		}
	}
	return 0
}

// RoomCount returns how many rooms of the given type a listing has — used to
// show "спальни 2", "санузел 1".
func (l Listing) RoomCount(roomType string) int {
	n := 0
	for _, r := range l.RoomSpecs {
		if r.Type == roomType {
			n++
		}
	}
	return n
}

// roomTypeKeys is the ordered set of room kinds offered in the form. Labels
// live in i18n as "re.rt_<key>".
var roomTypeKeys = []string{
	"living", "bedroom", "kitchen", "bathroom", "wc", "hallway", "balcony", "loggia", "other",
}

var roomTypeSet = func() map[string]bool {
	m := make(map[string]bool, len(roomTypeKeys))
	for _, k := range roomTypeKeys {
		m[k] = true
	}
	return m
}()

// RoomTypeKeys exposes room kinds to templates (for the form select).
func RoomTypeKeys() []string { return roomTypeKeys }

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
