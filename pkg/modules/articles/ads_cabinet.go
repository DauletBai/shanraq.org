package articles

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

// AdOrder is one ad-placement request. Captured as pending_payment until a
// payment provider is wired (Phase 2). Billing is intentionally not live yet.
type AdOrder struct {
	ID            string
	Title         string
	Body          string
	ImageURL      string
	TargetURL     string
	CTA           string
	Placement     string // all | articles | listings
	GeoRegion     string
	Rubric        string
	DurationDays  int
	Price         int64
	PaymentMethod string // kaspi | card | invoice
	Status        string
	Created       string
}

// adPlacements / adDurations / adPayMethods are the validated option sets.
var adPlacements = map[string]bool{"all": true, "articles": true, "listings": true}
var adDurations = map[int]bool{7: true, 30: true}
var adPayMethods = map[string]bool{"kaspi": true, "card": true, "invoice": true}

// adFlatPrice returns the flat placement price (tenge) for a duration. Charm
// pricing, founding-friendly; real billing applies these once Kaspi is live.
func adFlatPrice(days int) int64 {
	switch days {
	case 30:
		return 19900
	default:
		return 6900
	}
}

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

// CreateOrder stores a new ad-placement request (pending_payment).
func (s *AdStore) CreateOrder(ctx context.Context, advertiserID string, o AdOrder) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO ad_orders (advertiser_id, title, body, image_url, target_url, cta,
		                       placement, geo_region, rubric, duration_days, price, payment_method)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		advertiserID, o.Title, o.Body, o.ImageURL, o.TargetURL, o.CTA,
		o.Placement, o.GeoRegion, o.Rubric, o.DurationDays, o.Price, o.PaymentMethod)
	if err != nil {
		return fmt.Errorf("create ad order: %w", err)
	}
	return nil
}

// ListOrders returns an advertiser's orders, newest first.
func (s *AdStore) ListOrders(ctx context.Context, advertiserID string) ([]AdOrder, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, title, placement, geo_region, rubric, duration_days, price, payment_method, status,
		       to_char(created_at, 'DD.MM.YYYY')
		FROM ad_orders WHERE advertiser_id = $1 ORDER BY created_at DESC`, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("list ad orders: %w", err)
	}
	defer rows.Close()
	var out []AdOrder
	for rows.Next() {
		var o AdOrder
		if err := rows.Scan(&o.ID, &o.Title, &o.Placement, &o.GeoRegion, &o.Rubric,
			&o.DurationDays, &o.Price, &o.PaymentMethod, &o.Status, &o.Created); err != nil {
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
	Saved      string // "company" | "order" flash
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
	days := 7
	if r.FormValue("duration") == "30" {
		days = 30
	}
	o := AdOrder{
		Title:         strings.TrimSpace(r.FormValue("title")),
		Body:          strings.TrimSpace(r.FormValue("body")),
		ImageURL:      strings.TrimSpace(r.FormValue("image_url")),
		TargetURL:     strings.TrimSpace(r.FormValue("target_url")),
		CTA:           strings.TrimSpace(r.FormValue("cta")),
		Placement:     r.FormValue("placement"),
		GeoRegion:     strings.TrimSpace(r.FormValue("geo_region")),
		Rubric:        strings.TrimSpace(r.FormValue("rubric")),
		DurationDays:  days,
		PaymentMethod: r.FormValue("payment_method"),
		Price:         adFlatPrice(days),
	}
	if !adPlacements[o.Placement] {
		o.Placement = "all"
	}
	if !adPayMethods[o.PaymentMethod] {
		o.PaymentMethod = "kaspi"
	}
	if o.Rubric != "" && !IsCategory(o.Rubric) && o.Rubric != "realestate" {
		o.Rubric = ""
	}
	if o.Title == "" || o.TargetURL == "" {
		page := AdvertisePage{Base: m.base(r, T(lang, "adv.title"), lang), Advertiser: adv, Error: T(lang, "adv.err_order")}
		if orders, oerr := m.ads.ListOrders(r.Context(), adv.ID); oerr == nil {
			page.Orders = orders
		}
		m.render(w, "advertise", page)
		return
	}
	if err := m.ads.CreateOrder(r.Context(), adv.ID, o); err != nil {
		m.rt.Logger.Error("create ad order", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/advertise?saved=order", http.StatusSeeOther)
}
