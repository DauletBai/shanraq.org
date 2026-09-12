package shop

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"shanraq.org/pkg/site"
)

// The product and the queue of people waiting for it, against the real tables:
// a draft must be invisible to a reader, an address must not be recorded twice,
// and an edit must survive the round trip. Skipped unless SHANRAQ_TEST_DB names
// a test database.
func TestShopStore(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the shop store test")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	// Closing is registered as a cleanup rather than deferred: cleanups run
	// last-in-first-out, so the rows below are deleted while the pool is still
	// open. A deferred Close would shut it first and the deletes would quietly
	// do nothing, leaving rows that fail the next run.
	t.Cleanup(pool.Close)
	s := NewStore(pool)

	slug := "test-book-" + strings.ReplaceAll(t.Name(), "/", "-")
	if _, err := pool.Exec(ctx, `
		INSERT INTO shop_products (slug, kind, status, title_ru) VALUES ($1,'book','draft','Тест')
		ON CONFLICT (slug) DO UPDATE SET status = 'draft'`, slug); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO shop_plans (product_id, code, position, title_ru, includes_ru)
		SELECT id, 'book', 10, 'Книга', '- PDF и EPUB' FROM shop_products WHERE slug = $1
		ON CONFLICT (product_id, code) DO NOTHING`, slug); err != nil {
		t.Fatalf("seed package: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM shop_products WHERE slug = $1`, slug) })

	// A draft is staff-only: guessing the address must not open the page.
	if _, err := s.BySlug(ctx, slug, false); err == nil {
		t.Error("a draft product must not be readable by a guest")
	}
	p, err := s.BySlug(ctx, slug, true)
	if err != nil {
		t.Fatalf("staff read: %v", err)
	}

	// An edit round trip, including the price rule the panel enforces.
	p.Status = StatusSelling
	p.Edition = "1.0"
	p.Title[site.LangKZ] = "Кітап"
	p.Body[site.LangRU] = "Текст"
	if err := s.Save(ctx, p); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := s.BySlug(ctx, slug, false)
	if err != nil {
		t.Fatalf("read after save: %v", err)
	}
	if got.Edition != "1.0" || got.TitleIn(site.LangKZ) != "Кітап" {
		t.Errorf("saved product came back as %+v", got)
	}

	// The packages: a price set in the panel comes back on the product, and it
	// is what decides whether anything can be bought.
	if len(got.Plans) == 0 {
		t.Fatal("a product with no packages cannot be sold; the seed must create them")
	}
	if got.OnSale() {
		t.Error("nothing is priced yet, so the product must not be on sale")
	}
	pl := got.Plans[0]
	pl.Price = 12900
	if err := s.SavePlan(ctx, pl); err != nil {
		t.Fatalf("save plan: %v", err)
	}
	got, err = s.BySlug(ctx, slug, false)
	if err != nil {
		t.Fatalf("read after saving a package: %v", err)
	}
	if got.From() != 12900 || !got.OnSale() {
		t.Errorf("priced package did not put the product on sale: from=%d onsale=%v", got.From(), got.OnSale())
	}

	// The waiting list: the same address twice is one row and not an error.
	added, err := s.AddLead(ctx, got.ID, "Reader@Example.KZ", site.LangRU)
	if err != nil || !added {
		t.Fatalf("first lead: added=%v err=%v", added, err)
	}
	// The address is stored folded, so the same person in different case is the
	// same person — and does not get two letters.
	again, err := s.AddLead(ctx, got.ID, "reader@example.kz", site.LangRU)
	if err != nil {
		t.Fatalf("second lead: %v", err)
	}
	if again {
		t.Error("the same address must not be recorded twice")
	}
	n, err := s.LeadCount(ctx, got.ID)
	if err != nil || n != 1 {
		t.Errorf("lead count = %d (err %v), want 1", n, err)
	}

	// The list a reader sees never contains a draft.
	if _, err := pool.Exec(ctx, `UPDATE shop_products SET status='draft' WHERE slug=$1`, slug); err != nil {
		t.Fatalf("back to draft: %v", err)
	}
	items, err := s.List(ctx, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, it := range items {
		if it.Slug == slug {
			t.Error("a draft appeared in the public list")
		}
	}
}
