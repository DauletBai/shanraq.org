package articles

import (
	"context"
	"errors"
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
	Reported   string // "hidden" flash after a report auto-hid a listing
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
	Owner         bool
	Subscribed    bool
	IsFavorite    bool
	Reported      bool   // just submitted a report (thank-you flash)
	CanReport     bool   // logged-in and not the owner
	ShowContact   bool   // reveal the full seller contact
	MaskedContact string // partly-hidden contact shown before reveal
	ViewsCount    int
}

// MyListingsPage backs the author's own-listings management view.
type MyListingsPage struct {
	Base
	Listings []*Listing
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
	page.SidebarNews = m.latestNews(r, lang, 6)
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
	page.Reported = q.Get("reported")
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

	if in.Title == "" || in.Contact == "" || !in.NoFilters {
		page := ListingFormPage{Base: m.base(r, T(lang, "re.new_title"), lang)}
		page.ActiveCat = "realestate"
		page.Values = in
		if !in.NoFilters {
			page.Error = T(lang, "re.err_no_filters")
		} else {
			page.Error = T(lang, "re.err_required")
		}
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
	// Count a view — but not the owner's own visits.
	if !m.isListingOwner(r, l) {
		if err := m.listings.RecordView(r.Context(), id); err == nil {
			l.ViewsCount++
		}
	}
	m.renderListingView(w, r, l, false)
}

// handleListingContact reveals the seller's contact and counts the reveal. The
// full number is only ever rendered in response to this POST, so it stays out
// of the crawlable page markup.
func (m *Module) handleListingContact(w http.ResponseWriter, r *http.Request) {
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
	if !m.isListingOwner(r, l) {
		if err := m.listings.RecordContact(r.Context(), id); err != nil {
			m.rt.Logger.Warn("record contact", zap.Error(err))
		}
	}
	m.renderListingView(w, r, l, true)
}

// isListingOwner reports whether the current session user owns the listing.
func (m *Module) isListingOwner(r *http.Request, l *Listing) bool {
	aid, ok := m.authorID(r)
	return ok && aid.String() == l.AuthorID
}

// renderListingView builds and renders a listing page. reveal (or ownership)
// shows the full contact; otherwise it is masked behind a "show contact" button.
func (m *Module) renderListingView(w http.ResponseWriter, r *http.Request, l *Listing, reveal bool) {
	lang := m.resolveLang(w, r)
	page := ListingViewPage{Base: m.base(r, l.Title, lang)}
	page.ActiveCat = "realestate"
	page.L = l
	page.MaskedContact = maskContact(l.Contact)
	page.ViewsCount = l.ViewsCount
	if authorID, ok := m.authorID(r); ok {
		if authorID.String() == l.AuthorID {
			page.Owner = true
		} else {
			page.CanReport = true
		}
		if lid, err := uuid.Parse(l.ID); err == nil {
			page.IsFavorite = m.favs.IsFavorite(r.Context(), authorID, "listing", lid)
		}
	}
	page.ShowContact = reveal || page.Owner
	page.Reported = r.URL.Query().Get("reported") == "ok"
	page.SidebarNews = m.latestNews(r, lang, 6)
	m.applyListingSEO(&page)
	m.render(w, "listing_view", page)
}

// maskContact hides the middle of a phone/handle, keeping a recognizable
// prefix and suffix (spaces preserved).
func maskContact(s string) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= 6 {
		return "••••"
	}
	out := make([]rune, len(r))
	for i, c := range r {
		if i < 5 || i >= len(r)-2 || c == ' ' {
			out[i] = c
		} else {
			out[i] = '•'
		}
	}
	return string(out)
}

// handleListingReport records a reader's report of a listing (mainly filtered,
// dimension-distorting photos), warns the seller, and auto-hides the listing
// once enough distinct users report it.
func (m *Module) handleListingReport(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
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
	if l.AuthorID == uid.String() { // can't report your own listing
		http.Redirect(w, r, "/listings/"+l.ID+"?lang="+lang, http.StatusSeeOther)
		return
	}

	count, hidden, err := m.listings.Report(r.Context(), id, uid, strings.TrimSpace(r.FormValue("reason")))
	if err != nil {
		m.rt.Logger.Error("listing report", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Warn the seller by email (best-effort): on the first report, and again
	// when reports cross the threshold and the listing is hidden.
	if m.mailer != nil && (count == 1 || hidden) {
		subject := T(lang, "re.report_mail_subject")
		body := T(lang, "re.report_mail_body")
		if hidden {
			body = T(lang, "re.report_mail_hidden")
		}
		if err := m.mailer.Send(r.Context(), l.AuthorEmail, subject, body+"\n\n"+l.Title); err != nil {
			m.rt.Logger.Warn("seller report email", zap.Error(err))
		}
	}

	if hidden { // detail page now 404s; land the reporter on the index with a note
		http.Redirect(w, r, "/listings?lang="+lang+"&reported=hidden", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/listings/"+l.ID+"?lang="+lang+"&reported=ok", http.StatusSeeOther)
}

func (m *Module) handleMyListings(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	items, err := m.listings.MyListings(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Error("my listings", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := MyListingsPage{Base: m.base(r, T(lang, "re.my_listings"), lang)}
	page.ActiveCat = "realestate"
	page.Listings = items
	m.render(w, "listing_my", page)
}

// listingAction runs an owner-only lifecycle mutation and returns to /listings/my.
func (m *Module) listingAction(w http.ResponseWriter, r *http.Request, fn func(context.Context, uuid.UUID, uuid.UUID) error) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := fn(r.Context(), id, authorID); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("listing action", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listings/my", http.StatusSeeOther)
}

func (m *Module) handleListingExtend(w http.ResponseWriter, r *http.Request) {
	m.listingAction(w, r, m.listings.Extend)
}

func (m *Module) handleListingPromote(w http.ResponseWriter, r *http.Request) {
	m.listingAction(w, r, m.listings.Promote)
}

func (m *Module) handleListingFeature(w http.ResponseWriter, r *http.Request) {
	m.listingAction(w, r, m.listings.Feature)
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
	landArea, _ := strconv.ParseFloat(strings.Replace(strings.TrimSpace(r.FormValue("land_area")), ",", ".", 1), 64)
	rooms, _ := strconv.Atoi(digitsOnly(r.FormValue("rooms")))

	// Amenity checkboxes — keep only recognized keys.
	var amenities []string
	for _, a := range r.Form["amenity"] {
		if amenitySet[a] {
			amenities = append(amenities, a)
		}
	}
	var geoID *uuid.UUID
	if gid, err := uuid.Parse(strings.TrimSpace(r.FormValue("geo_node_id"))); err == nil {
		geoID = &gid
	}

	// Up to maxListingPhotos photo URLs (each an already-uploaded /media/... path).
	images := make([]string, 0, maxListingPhotos)
	for _, u := range r.Form["image"] {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		images = append(images, u)
		if len(images) >= maxListingPhotos {
			break
		}
	}
	cover := strings.TrimSpace(r.FormValue("cover_url"))
	if cover == "" && len(images) > 0 {
		cover = images[0]
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
		Cover:        cover,
		Images:       images,
		LandArea:     landArea,
		Amenities:    amenities,
		NoFilters:    r.FormValue("no_filters") == "on",
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
