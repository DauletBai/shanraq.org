package articles

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ListingsPage backs the real-estate marketplace grid.
type ListingsPage struct {
	Base
	Listings   []*Listing
	ActiveDeal string
	ActiveType string
	// Search form state (echoed back so the panel stays filled).
	Query      string
	PriceMin   int64
	PriceMax   int64
	RoomsMin   int
	RegionText string
	GeoNodeID  string
	Searching  bool // any filter beyond deal/type is active → open the panel
	Count      int
}

// ListingFormPage backs the submission form.
type ListingFormPage struct {
	Base
	Values ListingInput
	Error  string
}

// ListingViewPage backs a single listing.
type ListingViewPage struct {
	Base
	L          *Listing
	Subscribed bool
}

func (m *Module) handleListings(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	q := r.URL.Query()
	deal := q.Get("deal")
	ptype := q.Get("type")

	f := ListingFilter{Deal: deal, PropertyType: ptype, Limit: 30}
	f.PriceMin, _ = strconv.ParseInt(digitsOnly(q.Get("pmin")), 10, 64)
	f.PriceMax, _ = strconv.ParseInt(digitsOnly(q.Get("pmax")), 10, 64)
	f.RoomsMin, _ = strconv.Atoi(digitsOnly(q.Get("rooms")))
	f.Query = strings.TrimSpace(q.Get("q"))
	f.RegionText = strings.TrimSpace(q.Get("region"))
	if gid, err := uuid.Parse(strings.TrimSpace(q.Get("geo"))); err == nil {
		f.GeoNodeID = &gid
	}

	items, err := m.listings.List(r.Context(), f)
	if err != nil {
		m.rt.Logger.Error("listings list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := ListingsPage{Base: m.base(r, T(lang, "nav.realestate"), lang)}
	page.ActiveCat = "realestate"
	page.Listings = items
	page.Count = len(items)
	if isDealType(deal) {
		page.ActiveDeal = deal
	}
	if isPropertyType(ptype) {
		page.ActiveType = ptype
	}
	page.Query = f.Query
	page.PriceMin = f.PriceMin
	page.PriceMax = f.PriceMax
	page.RoomsMin = f.RoomsMin
	page.RegionText = f.RegionText
	if f.GeoNodeID != nil {
		page.GeoNodeID = f.GeoNodeID.String()
	}
	page.Searching = f.Query != "" || f.PriceMin > 0 || f.PriceMax > 0 || f.RoomsMin > 0 || f.RegionText != "" || f.GeoNodeID != nil
	m.render(w, "listings", page)
}

func (m *Module) handleListingNew(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	if _, ok := m.authorID(r); !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	page := ListingFormPage{Base: m.base(r, T(lang, "re.new_title"), lang)}
	page.ActiveCat = "realestate"
	page.Values = ListingInput{DealType: "sale", PropertyType: "apartment", Country: countryDefault(lang)}
	m.render(w, "listing_new", page)
}

func (m *Module) handleListingCreate(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	in := parseListingForm(r)

	// Resolve the selected location node into the denormalized address fields.
	if in.GeoNodeID != nil {
		anc, err := m.geo.Ancestry(r.Context(), *in.GeoNodeID, lang)
		if err != nil || len(anc) == 0 {
			in.GeoNodeID = nil
		} else {
			in.Country, in.Region, in.City, in.Village = "", "", "", ""
			for _, n := range anc {
				switch n.Level {
				case 0:
					in.Country = n.Name
				case 1:
					in.Region = n.Name
				case 2:
					in.City = n.Name
				default:
					in.Village = n.Name
				}
			}
		}
	}

	if in.Title == "" || in.Contact == "" {
		page := ListingFormPage{Base: m.base(r, T(lang, "re.new_title"), lang)}
		page.ActiveCat = "realestate"
		page.Values = in
		page.Error = T(lang, "re.err_required")
		m.render(w, "listing_new", page)
		return
	}

	id, err := m.listings.Create(r.Context(), authorID, in)
	if err != nil {
		m.rt.Logger.Error("create listing", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listings/"+id.String(), http.StatusSeeOther)
}

func (m *Module) handleListingView(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	l, err := m.listings.GetByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	page := ListingViewPage{Base: m.base(r, l.Title, lang)}
	page.ActiveCat = "realestate"
	page.L = l
	m.render(w, "listing_view", page)
}

func parseListingForm(r *http.Request) ListingInput {
	deal := r.FormValue("deal_type")
	if !isDealType(deal) {
		deal = "sale"
	}
	ptype := r.FormValue("property_type")
	if !isPropertyType(ptype) {
		ptype = "apartment"
	}
	price, _ := strconv.ParseInt(digitsOnly(r.FormValue("price")), 10, 64)
	area, _ := strconv.ParseFloat(strings.Replace(strings.TrimSpace(r.FormValue("area")), ",", ".", 1), 64)
	rooms, _ := strconv.Atoi(digitsOnly(r.FormValue("rooms")))
	var geoID *uuid.UUID
	if gid, err := uuid.Parse(strings.TrimSpace(r.FormValue("geo_node_id"))); err == nil {
		geoID = &gid
	}
	return ListingInput{
		DealType:     deal,
		PropertyType: ptype,
		Country:      strings.TrimSpace(r.FormValue("country")),
		Region:       strings.TrimSpace(r.FormValue("region")),
		City:         strings.TrimSpace(r.FormValue("city")),
		Village:      strings.TrimSpace(r.FormValue("village")),
		Price:        price,
		Area:         area,
		Rooms:        rooms,
		Title:        strings.TrimSpace(r.FormValue("title")),
		Description:  strings.TrimSpace(r.FormValue("description")),
		Contact:      strings.TrimSpace(r.FormValue("contact")),
		Cover:        strings.TrimSpace(r.FormValue("cover_url")),
		GeoNodeID:    geoID,
	}
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func countryDefault(lang string) string {
	switch lang {
	case LangKZ:
		return "Қазақстан"
	case LangEN:
		return "Kazakhstan"
	default:
		return "Казахстан"
	}
}

// money formats an integer amount with thin thousands separators.
func money(v int64) string {
	s := strconv.FormatInt(v, 10)
	n := len(s)
	if n <= 3 {
		return s
	}
	var b strings.Builder
	pre := n % 3
	if pre > 0 {
		b.WriteString(s[:pre])
	}
	for i := pre; i < n; i += 3 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}
