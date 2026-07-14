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
	deal := r.URL.Query().Get("deal")
	ptype := r.URL.Query().Get("type")

	items, err := m.listings.List(r.Context(), deal, ptype, 30)
	if err != nil {
		m.rt.Logger.Error("listings list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := ListingsPage{Base: m.base(r, T(lang, "nav.realestate"), lang)}
	page.ActiveCat = "realestate"
	page.Listings = items
	if isDealType(deal) {
		page.ActiveDeal = deal
	}
	if isPropertyType(ptype) {
		page.ActiveType = ptype
	}
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
