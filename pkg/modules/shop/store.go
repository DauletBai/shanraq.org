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
	Price      int64
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
// place that decides it, so a page and a checkout can never disagree.
func (p Product) OnSale() bool { return p.Status == StatusSelling && p.Price > 0 }

// Store reads and writes products and the queue of people waiting for them.
type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

const productCols = `id, slug, kind, status, price, currency, edition, cover_url, preview_url,
	title_kz, title_ru, title_en, summary_kz, summary_ru, summary_en,
	body_kz, body_ru, body_en, position, created_at, updated_at`

func scanProduct(row pgx.Row) (Product, error) {
	var p Product
	var tkz, tru, ten, skz, sru, sen, bkz, bru, ben string
	err := row.Scan(&p.ID, &p.Slug, &p.Kind, &p.Status, &p.Price, &p.Currency, &p.Edition,
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
	return out, rows.Err()
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
	return p, nil
}

// Save writes the editable fields of a product. The slug identifies it and is
// not editable: it is the address the page is linked from.
func (s *Store) Save(ctx context.Context, p Product) error {
	_, err := s.db.Exec(ctx, `
		UPDATE shop_products SET
			status = $2, price = $3, edition = $4, cover_url = $5, preview_url = $6,
			title_kz = $7, title_ru = $8, title_en = $9,
			summary_kz = $10, summary_ru = $11, summary_en = $12,
			body_kz = $13, body_ru = $14, body_en = $15,
			updated_at = NOW()
		 WHERE slug = $1`,
		p.Slug, p.Status, p.Price, p.Edition, p.CoverURL, p.PreviewURL,
		p.Title[site.LangKZ], p.Title[site.LangRU], p.Title[site.LangEN],
		p.Summary[site.LangKZ], p.Summary[site.LangRU], p.Summary[site.LangEN],
		p.Body[site.LangKZ], p.Body[site.LangRU], p.Body[site.LangEN])
	if err != nil {
		return fmt.Errorf("save product: %w", err)
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
