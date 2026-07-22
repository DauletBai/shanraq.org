package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Advertiser is a company account in the self-serve ad cabinet. One per user
// (the responsible person) for the MVP.
type Advertiser struct {
	ID           string
	CompanyName  string
	BIN          string
	LegalForm    string
	Address      string
	Website      string
	Industry     string
	ContactName  string
	ContactRole  string
	ContactPhone string
	ContactEmail string
}

// AdOrder is one booked ad placement: a creative, a zone, a date range and an
// optional exclusive hold. Captured as pending_payment until a payment provider
// is wired; once paid it flips to active and is served automatically.
type AdOrder struct {
	ID            string
	Title         string
	Body          string
	ImageURL      string
	TargetURL     string
	CTA           string
	Placement     string   // legacy single placement (kept for old rows)
	Surfaces      []string // the surfaces this order covers (new model)
	GeoRegion     string
	Rubric        string
	Lang          string
	Exclusive     bool
	StartsAt      string // YYYY-MM-DD
	EndsAt        string
	DurationDays  int
	Price         int64
	PaymentMethod string // kaspi | card | invoice
	Status        string
	Created       string
}

// adPayMethods is the validated payment-method set (zones/durations live in
// ads_rates.go so the rate card is the single source of truth).
var adPayMethods = map[string]bool{"kaspi": true, "card": true, "invoice": true}

// AdStore persists advertiser accounts and their ad orders.
type AdStore struct{ db *pgxpool.Pool }

func NewAdStore(db *pgxpool.Pool) *AdStore { return &AdStore{db: db} }

// ByOwner loads the advertiser owned by a user, or nil if none yet.
func (s *AdStore) ByOwner(ctx context.Context, ownerID uuid.UUID) (*Advertiser, error) {
	var a Advertiser
	err := s.db.QueryRow(ctx, `
		SELECT id, company_name, bin, legal_form, address, website, industry,
		       contact_name, contact_role, contact_phone, contact_email
		FROM advertisers WHERE owner_id = $1`, ownerID).Scan(
		&a.ID, &a.CompanyName, &a.BIN, &a.LegalForm, &a.Address, &a.Website, &a.Industry,
		&a.ContactName, &a.ContactRole, &a.ContactPhone, &a.ContactEmail)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("advertiser by owner: %w", err)
	}
	return &a, nil
}

// Save inserts or updates the caller's company (one per owner).
func (s *AdStore) Save(ctx context.Context, ownerID uuid.UUID, a Advertiser) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO advertisers (owner_id, company_name, bin, legal_form, address, website, industry,
		                         contact_name, contact_role, contact_phone, contact_email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (owner_id) DO UPDATE SET
			company_name=EXCLUDED.company_name, bin=EXCLUDED.bin, legal_form=EXCLUDED.legal_form,
			address=EXCLUDED.address, website=EXCLUDED.website, industry=EXCLUDED.industry,
			contact_name=EXCLUDED.contact_name, contact_role=EXCLUDED.contact_role,
			contact_phone=EXCLUDED.contact_phone, contact_email=EXCLUDED.contact_email, updated_at=NOW()`,
		ownerID, a.CompanyName, a.BIN, a.LegalForm, a.Address, a.Website, a.Industry,
		a.ContactName, a.ContactRole, a.ContactPhone, a.ContactEmail)
	if err != nil {
		return fmt.Errorf("save advertiser: %w", err)
	}
	return nil
}

// CreateOrder books a new ad placement (pending_payment).
func (s *AdStore) CreateOrder(ctx context.Context, advertiserID string, o AdOrder) error {
	// placement carries the first surface for backward compatibility with the
	// old column's NOT NULL/CHECK; surfaces is the authoritative set.
	placement := "home"
	if len(o.Surfaces) > 0 {
		if strings.HasPrefix(o.Surfaces[0], surfaceRubricPfx) {
			placement = "rubric"
		} else {
			placement = o.Surfaces[0]
		}
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO ad_orders (advertiser_id, title, body, image_url, target_url, cta,
		                       placement, surfaces, geo_region, lang, exclusive,
		                       starts_at, ends_at, duration_days, price, payment_method)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12::date,$13::date,$14,$15,$16)`,
		advertiserID, o.Title, o.Body, o.ImageURL, o.TargetURL, o.CTA,
		placement, o.Surfaces, o.GeoRegion, o.Lang, o.Exclusive,
		o.StartsAt, o.EndsAt, o.DurationDays, o.Price, o.PaymentMethod)
	if err != nil {
		return fmt.Errorf("create ad order: %w", err)
	}
	return nil
}

// ListOrders returns an advertiser's orders, newest first.
func (s *AdStore) ListOrders(ctx context.Context, advertiserID string) ([]AdOrder, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, title, placement, COALESCE(rubric,''), exclusive,
		       to_char(starts_at,'DD.MM.YYYY'), to_char(ends_at,'DD.MM.YYYY'),
		       duration_days, price, payment_method, status, to_char(created_at,'DD.MM.YYYY')
		FROM ad_orders WHERE advertiser_id = $1 ORDER BY created_at DESC`, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("list ad orders: %w", err)
	}
	defer rows.Close()
	var out []AdOrder
	for rows.Next() {
		var o AdOrder
		if err := rows.Scan(&o.ID, &o.Title, &o.Placement, &o.Rubric, &o.Exclusive,
			&o.StartsAt, &o.EndsAt, &o.DurationDays, &o.Price, &o.PaymentMethod,
			&o.Status, &o.Created); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// SlotsTaken counts how many of a zone's rotation slots are already booked for
// any day in [start,end]. An "all" package occupies a slot in every zone, and an
// exclusive booking takes the whole rotation. Asking about the "all" package is
// deliberately conservative: it competes with every zone, so we never oversell.
func (s *AdStore) SlotsTaken(ctx context.Context, surface, start, end string) (int, error) {
	var taken int
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN exclusive THEN $4 ELSE 1 END), 0)::int
		FROM ad_orders
		WHERE status IN ('pending_payment','active')
		  AND starts_at <= $2::date AND ends_at >= $1::date
		  AND surfaces @> ARRAY[$3]`,
		start, end, surface, adSlotCapacity).Scan(&taken)
	if err != nil {
		return 0, fmt.Errorf("slots taken: %w", err)
	}
	return taken, nil
}

// DayFree is one calendar day with the number of still-free rotation slots.
type DayFree struct {
	Date string `json:"date"` // YYYY-MM-DD
	Free int    `json:"free"`
}

// AvailabilityByDay returns, for every day in [from,to], how many rotation
// slots of the zone are still free — this feeds the booking calendar.
func (s *AdStore) AvailabilityByDay(ctx context.Context, surface, from, to string) ([]DayFree, error) {
	rows, err := s.db.Query(ctx, `
		SELECT to_char(d, 'YYYY-MM-DD'),
		       $4 - COALESCE(SUM(CASE WHEN o.exclusive THEN $4 ELSE 1 END), 0)::int AS free
		FROM generate_series($1::date, $2::date, '1 day') d
		LEFT JOIN ad_orders o
		       ON o.status IN ('pending_payment','active')
		      AND o.starts_at <= d::date AND o.ends_at >= d::date
		      AND o.surfaces @> ARRAY[$3]
		GROUP BY d ORDER BY d`, from, to, surface, adSlotCapacity)
	if err != nil {
		return nil, fmt.Errorf("availability by day: %w", err)
	}
	defer rows.Close()
	out := []DayFree{}
	for rows.Next() {
		var d DayFree
		if err := rows.Scan(&d.Date, &d.Free); err != nil {
			return nil, err
		}
		if d.Free < 0 {
			d.Free = 0
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ActiveBySurface returns creatives currently running on a surface, for serving.
func (s *AdStore) ActiveBySurface(ctx context.Context, surface, lang string, limit int) ([]AdOrder, error) {
	if limit <= 0 || limit > 10 {
		limit = adSlotCapacity
	}
	rows, err := s.db.Query(ctx, `
		SELECT title, COALESCE(body,''), COALESCE(image_url,''), COALESCE(target_url,''), COALESCE(cta,'')
		FROM ad_orders
		WHERE status = 'active'
		  AND starts_at <= CURRENT_DATE AND ends_at >= CURRENT_DATE
		  AND surfaces @> ARRAY[$1]
		  AND COALESCE(lang,'') IN ('', $2)
		ORDER BY exclusive DESC, created_at DESC
		LIMIT $3`, surface, lang, limit)
	if err != nil {
		return nil, fmt.Errorf("active ads: %w", err)
	}
	defer rows.Close()
	var out []AdOrder
	for rows.Next() {
		var o AdOrder
		if err := rows.Scan(&o.Title, &o.Body, &o.ImageURL, &o.TargetURL, &o.CTA); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// ---------- handlers ----------

// AdvertisePage backs the advertiser cabinet.
type AdvertisePage struct {
	Base
	Advertiser *Advertiser
	Orders     []AdOrder
	Avail      []SurfaceAvail // free rotation slots per surface for the previewed period
	Today      string         // YYYY-MM-DD, min bookable date
	Pricing    AdOrderPricing // last computed total (on validation failure)
	Saved      string         // "company" | "order" flash
	Error      string
}

func (m *Module) handleAdvertise(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	adv, err := m.ads.ByOwner(r.Context(), uid)
	if err != nil {
		m.rt.Logger.Error("advertiser load", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := AdvertisePage{Base: m.base(r, T(lang, "adv.title"), lang), Advertiser: adv}
	if adv != nil {
		if orders, oerr := m.ads.ListOrders(r.Context(), adv.ID); oerr == nil {
			page.Orders = orders
		} else {
			m.rt.Logger.Warn("advertiser orders", zap.Error(oerr))
		}
	}
	today := time.Now().Truncate(24 * time.Hour)
	page.Today = today.Format("2006-01-02")
	page.Avail = m.surfaceAvailability(r.Context(), page.Today, today.AddDate(0, 0, 29).Format("2006-01-02"))
	page.Saved = r.URL.Query().Get("saved")
	m.render(w, "advertise", page)
}

func (m *Module) handleAdvertiseCompany(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	a := Advertiser{
		CompanyName:  strings.TrimSpace(r.FormValue("company_name")),
		BIN:          strings.TrimSpace(r.FormValue("bin")),
		LegalForm:    strings.TrimSpace(r.FormValue("legal_form")),
		Address:      strings.TrimSpace(r.FormValue("address")),
		Website:      strings.TrimSpace(r.FormValue("website")),
		Industry:     strings.TrimSpace(r.FormValue("industry")),
		ContactName:  strings.TrimSpace(r.FormValue("contact_name")),
		ContactRole:  strings.TrimSpace(r.FormValue("contact_role")),
		ContactPhone: strings.TrimSpace(r.FormValue("contact_phone")),
		ContactEmail: strings.TrimSpace(r.FormValue("contact_email")),
	}
	if a.CompanyName == "" || a.ContactName == "" || a.ContactPhone == "" {
		adv, _ := m.ads.ByOwner(r.Context(), uid)
		page := AdvertisePage{Base: m.base(r, T(lang, "adv.title"), lang), Advertiser: adv, Error: T(lang, "adv.err_company")}
		if adv == nil {
			page.Advertiser = &a // echo entered values back into the form
		}
		m.render(w, "advertise", page)
		return
	}
	if err := m.ads.Save(r.Context(), uid, a); err != nil {
		m.rt.Logger.Error("advertiser save", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/advertise?saved=company", http.StatusSeeOther)
}

func (m *Module) handleAdvertiseOrder(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	adv, err := m.ads.ByOwner(r.Context(), uid)
	if err != nil || adv == nil {
		http.Redirect(w, r, "/advertise", http.StatusSeeOther) // must register a company first
		return
	}
	_ = r.ParseForm()

	days, _ := strconv.Atoi(digitsOnly(r.FormValue("duration")))
	if !adDurationSet[days] {
		days = adDurationDays[0]
	}
	// Selected surfaces (checkboxes). Keep only real ones, de-duplicated.
	seen := map[string]bool{}
	var surfaces []string
	for _, sfc := range r.Form["surface"] {
		if isAdSurface(sfc) && !seen[sfc] {
			seen[sfc] = true
			surfaces = append(surfaces, sfc)
		}
	}
	exclusive := r.FormValue("exclusive") == "on"

	// Start date: today by default, never in the past.
	today := time.Now().Truncate(24 * time.Hour)
	start := today
	if v, perr := time.Parse("2006-01-02", strings.TrimSpace(r.FormValue("starts_at"))); perr == nil {
		start = v
	}
	if start.Before(today) {
		start = today
	}
	end := start.AddDate(0, 0, days-1)

	pricing := AdOrderTotal(surfaces, days)
	o := AdOrder{
		Title:         strings.TrimSpace(r.FormValue("title")),
		Body:          strings.TrimSpace(r.FormValue("body")),
		ImageURL:      strings.TrimSpace(r.FormValue("image_url")),
		TargetURL:     strings.TrimSpace(r.FormValue("target_url")),
		CTA:           strings.TrimSpace(r.FormValue("cta")),
		Surfaces:      surfaces,
		GeoRegion:     strings.TrimSpace(r.FormValue("geo_region")),
		Lang:          r.FormValue("target_lang"),
		Exclusive:     exclusive,
		StartsAt:      start.Format("2006-01-02"),
		EndsAt:        end.Format("2006-01-02"),
		DurationDays:  days,
		PaymentMethod: r.FormValue("payment_method"),
		Price:         pricing.Total,
	}
	if !adPayMethods[o.PaymentMethod] {
		o.PaymentMethod = "kaspi"
	}
	if !IsLang(o.Lang) {
		o.Lang = ""
	}

	fail := func(msg string) {
		page := AdvertisePage{Base: m.base(r, T(lang, "adv.title"), lang), Advertiser: adv, Error: msg}
		if orders, oerr := m.ads.ListOrders(r.Context(), adv.ID); oerr == nil {
			page.Orders = orders
		}
		page.Avail = m.surfaceAvailability(r.Context(), o.StartsAt, o.EndsAt)
		page.Today = today.Format("2006-01-02")
		page.Pricing = pricing
		m.render(w, "advertise", page)
	}

	if o.Title == "" || o.TargetURL == "" {
		fail(T(lang, "adv.err_order"))
		return
	}
	if len(surfaces) == 0 {
		fail(T(lang, "adv.err_no_surface"))
		return
	}

	// Auto-check the booking: every selected surface must still have room for
	// those dates. A booking needs a slot on each surface it covers.
	need := 1
	if exclusive {
		need = adSlotCapacity
	}
	for _, sfc := range surfaces {
		taken, aerr := m.ads.SlotsTaken(r.Context(), sfc, o.StartsAt, o.EndsAt)
		if aerr != nil {
			m.rt.Logger.Error("slots taken", zap.Error(aerr))
			fail(T(lang, "adv.err_order"))
			return
		}
		if taken+need > adSlotCapacity {
			fail(T(lang, "adv.err_busy"))
			return
		}
	}

	if err := m.ads.CreateOrder(r.Context(), adv.ID, o); err != nil {
		m.rt.Logger.Error("create ad order", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/advertise?saved=order", http.StatusSeeOther)
}

// SurfaceAvail is a surface's free-slot count for the previewed period.
type SurfaceAvail struct {
	Code string
	Free int
	Cap  int
}

// surfaceAvailability reports free slots per surface for a date range.
func (m *Module) surfaceAvailability(ctx context.Context, start, end string) []SurfaceAvail {
	all := AdSurfaces()
	out := make([]SurfaceAvail, 0, len(all))
	for _, sfc := range all {
		free := adSlotCapacity
		if taken, err := m.ads.SlotsTaken(ctx, sfc.Code, start, end); err == nil {
			free = adSlotCapacity - taken
			if free < 0 {
				free = 0
			}
		}
		out = append(out, SurfaceAvail{Code: sfc.Code, Free: free, Cap: adSlotCapacity})
	}
	return out
}

// handleAdsAvailability feeds the booking calendar: free slots per day for a
// zone, so the advertiser sees busy and free dates before choosing.
func (m *Module) handleAdsAvailability(w http.ResponseWriter, r *http.Request) {
	// The calendar reflects the surfaces currently ticked: a booking needs a
	// free slot on EVERY one, so a day's availability is the minimum across
	// them. With none ticked, fall back to the home surface as a hint.
	var surfaces []string
	for _, sfc := range strings.Split(r.URL.Query().Get("surfaces"), ",") {
		sfc = strings.TrimSpace(sfc)
		if isAdSurface(sfc) {
			surfaces = append(surfaces, sfc)
		}
	}
	if len(surfaces) == 0 {
		surfaces = []string{surfaceHome}
	}
	today := time.Now().Truncate(24 * time.Hour)
	from := today
	if v, err := time.Parse("2006-01-02", r.URL.Query().Get("from")); err == nil && !v.Before(today) {
		from = v
	}
	to := from.AddDate(0, 0, 62) // ~2 months of calendar

	minFree := map[string]int{}
	order := []string{}
	for i, sfc := range surfaces {
		days, err := m.ads.AvailabilityByDay(r.Context(), sfc, from.Format("2006-01-02"), to.Format("2006-01-02"))
		if err != nil {
			m.rt.Logger.Error("ads availability", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		for _, d := range days {
			if i == 0 {
				minFree[d.Date] = d.Free
				order = append(order, d.Date)
			} else if d.Free < minFree[d.Date] {
				minFree[d.Date] = d.Free
			}
		}
	}
	out := make([]DayFree, 0, len(order))
	for _, date := range order {
		out = append(out, DayFree{Date: date, Free: minFree[date]})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"capacity": adSlotCapacity,
		"days":     out,
	})
}
