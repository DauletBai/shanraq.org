package articles

import (
	"context"
	"errors"
	"math"
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
	Query        string
	PriceMin     int64
	PriceMax     int64
	RoomsMin     int
	RegionText   string
	GeoNodeID    string
	SelAmenities map[string]bool // selected amenity filters (for re-checking boxes)
	Searching    bool            // any filter beyond deal/type is active → open the panel
	Count        int
	BannerAds    []*Listing    // paid banner slots shown in the real-estate sidebar
	Facets       ListingFacets // active-listing counts per deal/type, for filter badges
	Reported     string        // "hidden" flash after a report auto-hid a listing
}

// ListingFormPage backs the submission form.
type ListingFormPage struct {
	Base
	Values ListingInput
	Error  string
	// EditID is set when the form is editing an existing listing rather than
	// creating one; it switches the action URL and the submit label.
	EditID string
}

// ListingViewPage backs a single listing.
type ListingViewPage struct {
	Base
	L             *Listing
	Owner         bool
	Subscribed    bool
	IsFavorite    bool
	Reported      bool   // just submitted a report (thank-you flash)
	NeedVerify    bool   // tried to report/act without a verified email
	CanReport     bool   // logged-in and not the owner
	ShowContact   bool   // reveal the full seller contact
	MaskedContact string // partly-hidden contact shown before reveal
	ViewsCount    int
}

// MyListingsPage backs the author's own-listings management view.
type MyListingsPage struct {
	Base
	Listings []*Listing
	Credit   int // referral promotion-day credit, enables the free-promote button
	Saved    string
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
		// A price range is read in the currency of the place being searched, so
		// the country filter decides which listings the range can even match.
		f.Currency = listingCurrency("", m.geoCountry(r, gid))
	}
	// Amenity filters: keep only known keys, de-duplicated.
	selAmen := map[string]bool{}
	for _, a := range q["am"] {
		if amenitySet[a] && !selAmen[a] {
			selAmen[a] = true
			f.Amenities = append(f.Amenities, a)
		}
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
	if facets, ferr := m.listings.Facets(r.Context()); ferr != nil {
		m.rt.Logger.Warn("listings facets", zap.Error(ferr)) // badges degrade to 0, page still renders
	} else {
		page.Facets = facets
	}
	if bads, berr := m.listings.BannerListings(r.Context(), 2); berr != nil {
		m.rt.Logger.Warn("banner listings", zap.Error(berr)) // sidebar falls back to the promo card
	} else {
		page.BannerAds = bads
	}
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
	page.SelAmenities = selAmen
	if f.GeoNodeID != nil {
		page.GeoNodeID = f.GeoNodeID.String()
	}
	page.Searching = f.Query != "" || f.PriceMin > 0 || f.PriceMax > 0 || f.RoomsMin > 0 || f.RegionText != "" || f.GeoNodeID != nil || len(f.Amenities) > 0
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

// listingFormFail re-renders the submission form with everything the author
// typed and a message naming what to fix. Every rejection goes through here, so
// no path can quietly drop a field on the way back — which is how the location,
// and with it the currency, used to disappear between two submits.
func (m *Module) listingFormFail(w http.ResponseWriter, r *http.Request, lang, editID string, in ListingInput, msg string) {
	title := "re.new_title"
	if editID != "" {
		title = "re.edit_title"
	}
	page := ListingFormPage{Base: m.base(r, T(lang, title), lang)}
	page.ActiveCat = "realestate"
	page.Values = in
	page.EditID = editID
	page.Error = msg
	m.render(w, "listing_new", page)
}

// geoCountry returns the ISO country code a location node sits under, or "" if
// it cannot be resolved.
func (m *Module) geoCountry(r *http.Request, node uuid.UUID) string {
	anc, err := m.geo.Ancestry(r.Context(), node, LangRU)
	if err != nil {
		return ""
	}
	for _, n := range anc {
		if n.Country != "" {
			return n.Country
		}
	}
	return ""
}

// resolveListingLocation rewrites the denormalized address fields from the
// location the author picked and returns its ISO country code. A node that no
// longer resolves is cleared rather than half-applied, so the caller's
// "location is required" check catches it.
func (m *Module) resolveListingLocation(r *http.Request, in *ListingInput, lang string) string {
	if in.GeoNodeID == nil {
		return ""
	}
	anc, err := m.geo.Ancestry(r.Context(), *in.GeoNodeID, lang)
	if err != nil || len(anc) == 0 {
		in.GeoNodeID = nil
		return ""
	}
	countryCode := ""
	in.Country, in.Region, in.City, in.District = "", "", "", ""
	for _, n := range anc {
		if n.Country != "" {
			countryCode = n.Country
		}
		// By kind, never by depth. Depth assumes every country nests the same
		// number of times, and none do: a federal city (Москва, Санкт-Петербург,
		// Алматы, Астана, Шымкент) sits exactly where an oblast sits, so depth
		// labelled the city an oblast and its district a city.
		switch n.Kind {
		case "country":
			in.Country = n.Name
		case "region", "oblast", "krai", "republic", "okrug":
			in.Region = n.Name
		case "city", "town", "village":
			in.City = n.Name
		case "district":
			in.District = n.Name
		}
	}
	return countryCode
}

// validateListing returns the localized message for the first problem with a
// submission, or "" when it is good to save. Shared by create and edit so an
// edited listing can never be saved in a state a new one would be refused in.
func validateListing(in ListingInput, countryCode, lang string) string {
	// The location is required, and the price is why. A listing with no country
	// has no currency to follow, so the price silently defaults to tenge — which
	// is how a Moscow rent ended up priced in the wrong money. Demanding the
	// address closes that hole at the source, and a property advertisement
	// without one was never worth publishing anyway.
	if in.GeoNodeID == nil {
		return T(lang, "re.err_location_required")
	}
	if !in.NoFilters {
		return T(lang, "re.err_no_filters")
	}
	// Kazakh title is mandatory (the flagship trilingual rule) — except for
	// Russian listings, where only Russian and English are required.
	kzRequired := countryCode != "RU" && in.Currency != "RUB"
	if (kzRequired && in.TitleKz == "") || in.TitleRu == "" || in.TitleEn == "" || in.Contact == "" {
		return T(lang, "re.err_required")
	}
	// Language sanity: English must be Latin (no Cyrillic), Russian must be
	// Cyrillic, and Kazakh may be either — Kazakh is transitioning from Cyrillic
	// to a Latin alphabet, so both scripts are valid. This catches the common
	// mistake of pasting one language into every tab.
	if !isLatinText(in.TitleEn) || !isCyrillicText(in.TitleRu) || (in.TitleKz != "" && !hasLetters(in.TitleKz)) {
		return T(lang, "re.err_lang_script")
	}
	return ""
}

func (m *Module) handleListingCreate(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	authorID, ok := m.authorID(r)
	if !ok {
		// Session expired (or never authenticated) — explain why on the login
		// page instead of a silent bounce that drops the filled form.
		http.Redirect(w, r, "/studio/login?reason=session_expired", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		m.listingFormFail(w, r, lang, "", ListingInput{}, T(lang, "re.err_bad_form"))
		return
	}

	in := parseListingForm(r)

	// Listing submission gate (staged launch): open / invite-only / closed.
	if ok, msg := m.gateReason(r, SvcListingSubmit, lang); !ok {
		m.listingFormFail(w, r, lang, "", in, msg)
		return
	}
	// Posting requires a verified email (blocks throwaway-account spam).
	if !m.auth.IsEmailVerified(r.Context(), authorID) {
		m.listingFormFail(w, r, lang, "", in, T(lang, "re.err_verify_email"))
		return
	}

	countryCode := m.resolveListingLocation(r, &in, lang)
	in.Currency = listingCurrency(in.Currency, countryCode)
	if msg := validateListing(in, countryCode, lang); msg != "" {
		m.listingFormFail(w, r, lang, "", in, msg)
		return
	}

	id, err := m.listings.Create(r.Context(), authorID, in)
	if err != nil {
		m.rt.Logger.Error("create listing", zap.Error(err))
		// Re-render the form with everything the user typed so a transient save
		// error never costs them their work; tell them plainly what happened.
		m.listingFormFail(w, r, lang, "", in, T(lang, "re.err_save_failed"))
		return
	}
	// A real listing is the rewardable action: if this author was invited,
	// their referrer earns promotion credit now. Best-effort — a reward failure
	// must not fail the listing.
	if _, err := m.refs.Qualify(r.Context(), authorID); err != nil {
		m.rt.Logger.Warn("qualify referral", zap.Error(err))
	}
	http.Redirect(w, r, "/listings/"+id.String(), http.StatusSeeOther)
}

// handleListingEdit shows the submission form filled with an existing listing.
// Scoped to the owner: another account's id is a 404, not a peek at the form.
func (m *Module) handleListingEdit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	l, _, ok := m.ownedListing(w, r)
	if !ok {
		return
	}
	page := ListingFormPage{Base: m.base(r, T(lang, "re.edit_title"), lang)}
	page.ActiveCat = "realestate"
	page.Values = listingToInput(l)
	page.EditID = l.ID
	m.render(w, "listing_new", page)
}

// handleListingUpdate saves an edit. It runs the identical validation as a new
// submission, so a listing cannot be edited into a state it could not be posted
// in — an edited price still has to carry the currency of a real country.
func (m *Module) handleListingUpdate(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	l, authorID, ok := m.ownedListing(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		m.listingFormFail(w, r, lang, l.ID, listingToInput(l), T(lang, "re.err_bad_form"))
		return
	}
	in := parseListingForm(r)
	countryCode := m.resolveListingLocation(r, &in, lang)
	in.Currency = listingCurrency(in.Currency, countryCode)
	if msg := validateListing(in, countryCode, lang); msg != "" {
		m.listingFormFail(w, r, lang, l.ID, in, msg)
		return
	}
	id, err := uuid.Parse(l.ID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := m.listings.Update(r.Context(), id, authorID, in); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("update listing", zap.Error(err))
		m.listingFormFail(w, r, lang, l.ID, in, T(lang, "re.err_save_failed"))
		return
	}
	http.Redirect(w, r, "/listings/"+l.ID, http.StatusSeeOther)
}

// handleListingDelete removes the author's own listing for good.
func (m *Module) handleListingDelete(w http.ResponseWriter, r *http.Request) {
	m.listingAction(w, r, m.listings.Delete)
}

// ownedListing loads the listing named in the URL and confirms the caller owns
// it, writing the redirect or 404 itself when they do not.
func (m *Module) ownedListing(w http.ResponseWriter, r *http.Request) (*Listing, uuid.UUID, bool) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login?reason=session_expired", http.StatusSeeOther)
		return nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return nil, uuid.Nil, false
	}
	l, err := m.listings.GetByID(r.Context(), id)
	if err != nil || l.AuthorID != authorID.String() {
		http.NotFound(w, r)
		return nil, uuid.Nil, false
	}
	return l, authorID, true
}

// listingToInput turns a stored listing back into form values, so the edit form
// opens showing what is actually published rather than a blank slate.
func listingToInput(l *Listing) ListingInput {
	in := ListingInput{
		DealType: l.DealType, PropertyType: l.PropertyType,
		Country: l.Country, Region: l.Region, City: l.City, District: l.District,
		Microdistrict: l.Microdistrict, Street: l.Street, House: l.House,
		Lat: l.Lat, Lng: l.Lng,
		Price: l.Price, Currency: l.Currency, Area: l.Area, Rooms: l.Rooms,
		Title: l.Title, Description: l.Description, Contact: l.Contact, Cover: l.CoverURL,
		TitleKz: l.TitleKz, TitleRu: l.TitleRu, TitleEn: l.TitleEn,
		DescriptionKz: l.DescriptionKz, DescriptionRu: l.DescriptionRu, DescriptionEn: l.DescriptionEn,
		Images: l.Images, Documents: l.Documents, ContractURL: l.ContractURL,
		LandArea: l.LandArea, Amenities: l.Amenities, RoomSpecs: l.RoomSpecs,
		// Already attested once when the listing was posted; re-ticking the box
		// on every price fix would be a ritual, not a check.
		NoFilters: true,
	}
	if id, err := uuid.Parse(l.GeoNodeID); err == nil {
		in.GeoNodeID = &id
	}
	return in
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
	// Count a view — but not the owner's own visits, and not a crawler's: a
	// seller reads this number as interest in their flat, so it has to mean
	// people. See the same guard on articles in handlers.go.
	if !m.isListingOwner(r, l) && botLabel(r.UserAgent()) == "" {
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
	page := ListingViewPage{Base: m.base(r, l.TitleIn(lang), lang)}
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
	page.NeedVerify = r.URL.Query().Get("notice") == "verify"
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
	// Reporting requires a verified email, so throwaway accounts can't brigade.
	if !m.auth.IsEmailVerified(r.Context(), uid) {
		http.Redirect(w, r, "/listings/"+id.String()+"?lang="+lang+"&notice=verify", http.StatusSeeOther)
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
	page.Saved = r.URL.Query().Get("ok")
	if c, err := m.refs.Balance(r.Context(), authorID); err == nil {
		page.Credit = c
	}
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

// serviceGate blocks a paid action when its service is not fully available
// (maintenance or off), sending the user back where a localized notice explains
// why. Free actions (extend, referral-funded promotion) never call this.
func (m *Module) serviceGate(w http.ResponseWriter, r *http.Request, code, back string) bool {
	if m.flags.Available(code) {
		return true
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
	return false
}

func (m *Module) handleListingPromote(w http.ResponseWriter, r *http.Request) {
	if !m.serviceGate(w, r, SvcListingPromo, "/listings/my") {
		return
	}
	m.listingAction(w, r, m.listings.Promote)
}

func (m *Module) handleListingFeature(w http.ResponseWriter, r *http.Request) {
	if !m.serviceGate(w, r, SvcListingPromo, "/listings/my") {
		return
	}
	m.listingAction(w, r, m.listings.Feature)
}

// handleListingBanner buys the sidebar banner slot for 1..7 days.
func (m *Module) handleListingBanner(w http.ResponseWriter, r *http.Request) {
	if !m.serviceGate(w, r, SvcListingPromo, "/listings/my") {
		return
	}
	days, _ := strconv.Atoi(digitsOnly(r.FormValue("days")))
	if days < 1 || days > 7 {
		days = 1
	}
	m.listingAction(w, r, func(ctx context.Context, id, author uuid.UUID) error {
		return m.listings.Banner(ctx, id, author, days)
	})
}

// maxArea caps a single area field. The largest apartments and houses in the
// country are far below this; anything above it is a typo or an attempt to game
// the sort order, and both should be refused rather than stored.
const maxArea = 100000.0

// parseArea reads a "123,4" or "123.4" area field into a sane number.
//
// strconv.ParseFloat accepts "NaN", "Inf" and "-1" quite happily, and the old
// code discarded the error, so those landed in the database and then in the
// filters — where NaN compares false against everything and quietly removes a
// listing from every search. Out-of-range values become zero, which the form
// already treats as "not specified".
// mediaPath keeps only URLs that point at our own uploads.
//
// The form fields carry paths returned by /media/upload, but nothing stopped a
// hand-made POST from putting any absolute URL there — the comment promised
// "/media/..." while the code accepted whatever arrived. That let a listing
// embed a remote image, which loads from a third-party server on every view: a
// tracking pixel for anyone who opens the listing, and a link that can be
// swapped for something else after moderation has passed.
func mediaPath(raw string) (string, bool) {
	u := strings.TrimSpace(raw)
	// Reject anything with a scheme or host, and any traversal attempt, before
	// the prefix test — "/media/../etc" starts with /media/ too.
	if u == "" || !strings.HasPrefix(u, "/media/") || strings.Contains(u, "..") ||
		strings.Contains(u, "//") || strings.ContainsAny(u, " \t\r\n\\") {
		return "", false
	}
	return u, true
}

func parseArea(raw string) float64 {
	v, err := strconv.ParseFloat(strings.Replace(strings.TrimSpace(raw), ",", ".", 1), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > maxArea {
		return 0
	}
	return v
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
	area := parseArea(r.FormValue("area"))
	landArea := parseArea(r.FormValue("land_area"))
	rooms, _ := strconv.Atoi(digitsOnly(r.FormValue("rooms")))

	// Amenity checkboxes — keep only recognized keys.
	var amenities []string
	for _, a := range r.Form["amenity"] {
		if amenitySet[a] {
			amenities = append(amenities, a)
		}
	}

	// Per-room breakdown — parallel arrays room_type / room_area / room_note.
	var roomSpecs []RoomSpec
	types := r.Form["room_type"]
	areas := r.Form["room_area"]
	notes := r.Form["room_note"]
	for i, t := range types {
		if !roomTypeSet[t] || len(roomSpecs) >= maxRoomSpecs {
			continue
		}
		var ar float64
		if i < len(areas) {
			ar = parseArea(areas[i])
		}
		var note string
		if i < len(notes) {
			note = strings.TrimSpace(notes[i])
		}
		if ar <= 0 && note == "" { // an empty row — skip
			continue
		}
		roomSpecs = append(roomSpecs, RoomSpec{Type: t, Area: ar, Note: note})
	}
	// Coordinates come from the author dragging a pin; keep them only if they
	// are real numbers inside the world's range, so a broken client cannot
	// write nonsense into the map.
	var lat, lng *float64
	if a, err1 := strconv.ParseFloat(strings.TrimSpace(r.FormValue("lat")), 64); err1 == nil {
		if o, err2 := strconv.ParseFloat(strings.TrimSpace(r.FormValue("lng")), 64); err2 == nil {
			if a >= -90 && a <= 90 && o >= -180 && o <= 180 {
				lat, lng = &a, &o
			}
		}
	}

	var geoID *uuid.UUID
	if gid, err := uuid.Parse(strings.TrimSpace(r.FormValue("geo_node_id"))); err == nil {
		geoID = &gid
	}

	// Up to maxListingPhotos photo URLs (each an already-uploaded /media/... path).
	images := make([]string, 0, maxListingPhotos)
	for _, raw := range r.Form["image"] {
		u, ok := mediaPath(raw)
		if !ok {
			continue
		}
		images = append(images, u)
		if len(images) >= maxListingPhotos {
			break
		}
	}
	cover, _ := mediaPath(r.FormValue("cover_url"))
	if cover == "" && len(images) > 0 {
		cover = images[0]
	}

	// Up to maxListingDocs document URLs (PDF plans/passports or image schemes,
	// already uploaded via /media/upload-doc).
	documents := make([]string, 0, maxListingDocs)
	for _, raw := range r.Form["document"] {
		u, ok := mediaPath(raw)
		if !ok {
			continue
		}
		documents = append(documents, u)
		if len(documents) >= maxListingDocs {
			break
		}
	}

	// The lease the landlord is ready to sign — one file, and only for rent. A
	// sale contract is notarial and settled at the deal, so a published draft
	// would inform nobody; a buyer's papers belong in `documents` above.
	contract := ""
	if deal == "rent" {
		contract, _ = mediaPath(r.FormValue("contract"))
	}

	return ListingInput{
		DealType:      deal,
		PropertyType:  ptype,
		Country:       strings.TrimSpace(r.FormValue("country")),
		Region:        strings.TrimSpace(r.FormValue("region")),
		City:          strings.TrimSpace(r.FormValue("city")),
		District:      strings.TrimSpace(r.FormValue("district")),
		Microdistrict: clip(strings.TrimSpace(r.FormValue("microdistrict")), 60),
		Street:        clip(strings.TrimSpace(r.FormValue("street")), 80),
		House:         clip(strings.TrimSpace(r.FormValue("house")), 20),
		Lat:           lat,
		Lng:           lng,
		Price:         price,
		Area:          area,
		Rooms:         rooms,
		Title:         strings.TrimSpace(r.FormValue("title_ru")), // base/fallback
		Description:   strings.TrimSpace(r.FormValue("description_ru")),
		TitleKz:       strings.TrimSpace(r.FormValue("title_kz")),
		TitleRu:       strings.TrimSpace(r.FormValue("title_ru")),
		TitleEn:       strings.TrimSpace(r.FormValue("title_en")),
		DescriptionKz: strings.TrimSpace(r.FormValue("description_kz")),
		DescriptionRu: strings.TrimSpace(r.FormValue("description_ru")),
		DescriptionEn: strings.TrimSpace(r.FormValue("description_en")),
		Contact:       strings.TrimSpace(r.FormValue("contact")),
		Cover:         cover,
		Images:        images,
		Documents:     documents,
		ContractURL:   contract,
		LandArea:      landArea,
		Amenities:     amenities,
		RoomSpecs:     roomSpecs,
		NoFilters:     r.FormValue("no_filters") == "on",
		GeoNodeID:     geoID,
		Currency:      strings.ToUpper(strings.TrimSpace(r.FormValue("currency"))),
	}
}

// listingCurrencies are the currencies a listing may be priced in. Anything
// outside this set is not a currency this site knows how to display.
var listingCurrencies = map[string]bool{"KZT": true, "RUB": true}

// listingCurrency decides what a price is denominated in.
//
// The author's explicit choice wins, because it is the symbol printed next to
// the price box while they type: a server that silently overrules what they can
// see is exactly how a ruble price ends up stored as tenge. The country is the
// fallback for a submission that carries no choice at all — an old cached form,
// or a browser with JavaScript off — and the tenge is the last resort.
func listingCurrency(picked, countryCode string) string {
	if c := strings.ToUpper(strings.TrimSpace(picked)); listingCurrencies[c] {
		return c
	}
	if countryCode == "RU" {
		return "RUB"
	}
	return "KZT"
}

// isLatinText reports whether s carries Latin letters and no Cyrillic — used to
// keep the English tab from holding Russian/Kazakh text.
func isLatinText(s string) bool {
	hasLatin := false
	for _, r := range s {
		if r >= 0x0400 && r <= 0x04FF { // any Cyrillic letter disqualifies
			return false
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			hasLatin = true
		}
	}
	return hasLatin
}

// isCyrillicText reports whether s carries at least one Cyrillic letter — used
// to keep the Kazakh/Russian tabs from holding Latin (English/transliterated)
// text.
func isCyrillicText(s string) bool {
	for _, r := range s {
		if r >= 0x0400 && r <= 0x04FF {
			return true
		}
	}
	return false
}

// hasLetters reports whether s carries at least one Cyrillic or Latin letter —
// used for the Kazakh tab, which may be written in either script.
func hasLetters(s string) bool {
	for _, r := range s {
		if (r >= 0x0400 && r <= 0x04FF) || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			return true
		}
	}
	return false
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
