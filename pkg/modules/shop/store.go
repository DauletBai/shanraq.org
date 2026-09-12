package shop

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shanraq.org/pkg/site"
)

// Product statuses. A product moves left to right and can move back: draft
// while it is being written, announced when the page may be read but nothing
// can be bought, selling once an acquirer is connected, paused to stop sales
// without taking the page down.
const (
	StatusDraft     = "draft"
	StatusAnnounced = "announced"
	StatusSelling   = "selling"
	StatusPaused    = "paused"
)

// Statuses is the order the admin panel offers them in.
var Statuses = []string{StatusDraft, StatusAnnounced, StatusSelling, StatusPaused}

// ErrNotFound is returned for a slug that is not a product (or not a public
// one, for a reader who is not staff).
var ErrNotFound = errors.New("shop: no such product")

// Product is one thing the site sells, with its page copy in all three
// languages. Text is stored per language rather than translated on the fly: a
// price list and a promise about what the buyer receives are not something to
// hand to a machine translator.
type Product struct {
	ID         uuid.UUID
	Slug       string
	Kind       string
	Status     string
	Currency   string
	Edition    string
	CoverURL   string
	PreviewURL string
	Title      map[string]string
	Summary    map[string]string
	Body       map[string]string
	Position   int
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Plans are the packages this product is sold in, cheapest first. Loaded
	// with the product: a product page without its prices is not a product page.
	Plans []Plan
}

// TitleIn returns the title in lang, falling back the way the rest of the site
// does: the reader's language, then Russian, then anything that is filled in.
func (p Product) TitleIn(lang string) string   { return pick(p.Title, lang) }
func (p Product) SummaryIn(lang string) string { return pick(p.Summary, lang) }
func (p Product) BodyIn(lang string) string    { return pick(p.Body, lang) }

func pick(m map[string]string, lang string) string {
	if v := strings.TrimSpace(m[lang]); v != "" {
		return v
	}
	if v := strings.TrimSpace(m[site.LangRU]); v != "" {
		return v
	}
	for _, l := range site.Langs {
		if v := strings.TrimSpace(m[l]); v != "" {
			return v
		}
	}
	return ""
}

// Public reports whether a reader who is not staff may see this product.
func (p Product) Public() bool { return p.Status != StatusDraft }

// OnSale reports whether the buy action should be offered. It is the only
// place that decides it, so a page and a checkout can never disagree: selling,
// and at least one package with a price on it. A package priced at nothing is
// not a gift, it is a price nobody has set yet.
func (p Product) OnSale() bool {
	if p.Status != StatusSelling {
		return false
	}
	for _, pl := range p.Plans {
		if pl.Buyable() {
			return true
		}
	}
	return false
}

// From is the lowest price among the packages, for the card in the list.
// Zero means no package is priced yet.
func (p Product) From() int64 {
	var min int64
	for _, pl := range p.Plans {
		if pl.Buyable() && (min == 0 || pl.Price < min) {
			min = pl.Price
		}
	}
	return min
}

// Plan is one package a product is sold in: the book, or the book with the
// code. What differs is the price and the list of what the buyer downloads.
type Plan struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Code      string
	Price     int64
	Position  int
	Title     map[string]string
	Includes  map[string]string
}

// Buyable reports whether this package has a price to charge.
func (p Plan) Buyable() bool { return p.Price > 0 }

func (p Plan) TitleIn(lang string) string    { return pick(p.Title, lang) }
func (p Plan) IncludesIn(lang string) string { return pick(p.Includes, lang) }

// Store reads and writes products and the queue of people waiting for them.
type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

const productCols = `id, slug, kind, status, currency, edition, cover_url, preview_url,
	title_kz, title_ru, title_en, summary_kz, summary_ru, summary_en,
	body_kz, body_ru, body_en, position, created_at, updated_at`

func scanProduct(row pgx.Row) (Product, error) {
	var p Product
	var tkz, tru, ten, skz, sru, sen, bkz, bru, ben string
	err := row.Scan(&p.ID, &p.Slug, &p.Kind, &p.Status, &p.Currency, &p.Edition,
		&p.CoverURL, &p.PreviewURL, &tkz, &tru, &ten, &skz, &sru, &sen, &bkz, &bru, &ben,
		&p.Position, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return Product{}, err
	}
	p.Title = map[string]string{site.LangKZ: tkz, site.LangRU: tru, site.LangEN: ten}
	p.Summary = map[string]string{site.LangKZ: skz, site.LangRU: sru, site.LangEN: sen}
	p.Body = map[string]string{site.LangKZ: bkz, site.LangRU: bru, site.LangEN: ben}
	return p, nil
}

// List returns the products, newest position first. withDrafts is for staff.
func (s *Store) List(ctx context.Context, withDrafts bool) ([]Product, error) {
	q := `SELECT ` + productCols + ` FROM shop_products`
	if !withDrafts {
		q += ` WHERE status <> 'draft'`
	}
	q += ` ORDER BY position, created_at`
	rows, err := s.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()
	var out []Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	ids := make([]uuid.UUID, 0, len(out))
	for _, p := range out {
		ids = append(ids, p.ID)
	}
	plans, err := s.plansFor(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Plans = plans[out[i].ID]
	}
	return out, nil
}

// BySlug loads one product. A draft is found only for staff, so an unfinished
// page cannot be read by guessing its address.
func (s *Store) BySlug(ctx context.Context, slug string, withDrafts bool) (Product, error) {
	q := `SELECT ` + productCols + ` FROM shop_products WHERE slug = $1`
	if !withDrafts {
		q += ` AND status <> 'draft'`
	}
	p, err := scanProduct(s.db.QueryRow(ctx, q, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("get product: %w", err)
	}
	plans, err := s.plansFor(ctx, []uuid.UUID{p.ID})
	if err != nil {
		return Product{}, err
	}
	p.Plans = plans[p.ID]
	return p, nil
}

// Save writes the editable fields of a product. The slug identifies it and is
// not editable: it is the address the page is linked from.
func (s *Store) Save(ctx context.Context, p Product) error {
	_, err := s.db.Exec(ctx, `
		UPDATE shop_products SET
			status = $2, edition = $3, cover_url = $4, preview_url = $5,
			title_kz = $6, title_ru = $7, title_en = $8,
			summary_kz = $9, summary_ru = $10, summary_en = $11,
			body_kz = $12, body_ru = $13, body_en = $14,
			updated_at = NOW()
		 WHERE slug = $1`,
		p.Slug, p.Status, p.Edition, p.CoverURL, p.PreviewURL,
		p.Title[site.LangKZ], p.Title[site.LangRU], p.Title[site.LangEN],
		p.Summary[site.LangKZ], p.Summary[site.LangRU], p.Summary[site.LangEN],
		p.Body[site.LangKZ], p.Body[site.LangRU], p.Body[site.LangEN])
	if err != nil {
		return fmt.Errorf("save product: %w", err)
	}
	return nil
}

const planCols = `id, product_id, code, price, position,
	title_kz, title_ru, title_en, includes_kz, includes_ru, includes_en`

func scanPlan(row pgx.Row) (Plan, error) {
	var pl Plan
	var tkz, tru, ten, ikz, iru, ien string
	err := row.Scan(&pl.ID, &pl.ProductID, &pl.Code, &pl.Price, &pl.Position,
		&tkz, &tru, &ten, &ikz, &iru, &ien)
	if err != nil {
		return Plan{}, err
	}
	pl.Title = map[string]string{site.LangKZ: tkz, site.LangRU: tru, site.LangEN: ten}
	pl.Includes = map[string]string{site.LangKZ: ikz, site.LangRU: iru, site.LangEN: ien}
	return pl, nil
}

// plansFor loads the packages of several products in one query: a list of
// products would otherwise ask the database once per card.
func (s *Store) plansFor(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]Plan, error) {
	out := map[uuid.UUID][]Plan{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := s.db.Query(ctx, `SELECT `+planCols+`
		FROM shop_plans WHERE product_id = ANY($1) ORDER BY position, price`, ids)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		pl, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		out[pl.ProductID] = append(out[pl.ProductID], pl)
	}
	return out, rows.Err()
}

// SavePlan writes the editable fields of one package.
func (s *Store) SavePlan(ctx context.Context, pl Plan) error {
	_, err := s.db.Exec(ctx, `
		UPDATE shop_plans SET
			price = $3, title_kz = $4, title_ru = $5, title_en = $6,
			includes_kz = $7, includes_ru = $8, includes_en = $9, updated_at = NOW()
		 WHERE product_id = $1 AND code = $2`,
		pl.ProductID, pl.Code, pl.Price,
		pl.Title[site.LangKZ], pl.Title[site.LangRU], pl.Title[site.LangEN],
		pl.Includes[site.LangKZ], pl.Includes[site.LangRU], pl.Includes[site.LangEN])
	if err != nil {
		return fmt.Errorf("save plan: %w", err)
	}
	return nil
}

// AddLead records that someone wants to be told when this product is ready.
// Asking twice is not an error: the second time is the same promise, and
// telling a reader "you are already on the list" is friendlier than an error
// page. Returns false when the address was already there.
func (s *Store) AddLead(ctx context.Context, productID uuid.UUID, email, lang string) (added bool, err error) {
	tag, err := s.db.Exec(ctx, `
		INSERT INTO shop_leads (product_id, email, lang) VALUES ($1, lower($2), $3)
		ON CONFLICT (product_id, email) DO NOTHING`, productID, strings.TrimSpace(email), lang)
	if err != nil {
		return false, fmt.Errorf("add lead: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// LeadCount is how many people are waiting for a product, for the admin panel.
func (s *Store) LeadCount(ctx context.Context, productID uuid.UUID) (int, error) {
	var n int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM shop_leads WHERE product_id = $1`, productID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count leads: %w", err)
	}
	return n, nil
}
