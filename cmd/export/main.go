// Command export writes the platform's own published content to a directory of
// JSON files: articles with their translations, the prediction ledger, and the
// editable pages. Nothing else.
//
// It exists so a copy of the writing can live outside the country. The full
// backup cannot: it holds accounts, e-mail addresses, sellers' phone numbers
// and avatars, and Article 12(2) of the Law on Personal Data keeps those in a
// database inside Kazakhstan. This export carries none of them, so it may go
// anywhere — and it is also exactly the material a read-only mirror would need,
// which makes the two jobs one job.
//
// What is deliberately left out, and why it is a list rather than a filter:
// every table not named here is excluded by default, so a table added later
// cannot leak into the export by being forgotten. Accounts, sessions, tokens,
// listings (they carry sellers' telephone numbers), comments, votes, favourites,
// reputation, moderation, analytics, payments, advertisers and the job queue all
// stay behind.
//
// Author names are included as a display string because a mirror without
// bylines is not the same publication. Today they are two: the platform's own
// AI columnist and the owner. A third-party author would be personal data
// belonging to someone else, so the export refuses to run when one appears —
// better a failed job than a quiet export of somebody's name abroad.
//
//	DATABASE_URL=... EXPORT_DIR=/tmp/content go run ./cmd/export
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// knownAuthors are the bylines the export may carry abroad: the platform's own
// AI columnist and the owner. Anyone else is a third party whose name is not
// ours to move.
var knownAuthors = map[string]bool{"sana": true, "Daulet Baimurza": true}

type translation struct {
	Lang    string `json:"lang"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body_md"`
	Source  string `json:"source,omitempty"`
	Status  string `json:"status"`
}

type article struct {
	Slug         string        `json:"slug"`
	Author       string        `json:"author"`
	Category     string        `json:"category"`
	Subcategory  string        `json:"subcategory,omitempty"`
	OriginalLang string        `json:"original_lang"`
	CoverURL     string        `json:"cover_url,omitempty"`
	Indexable    bool          `json:"indexable"`
	PublishedAt  *time.Time    `json:"published_at,omitempty"`
	Translations []translation `json:"translations"`
}

type prediction struct {
	ArticleSlug string            `json:"article_slug,omitempty"`
	MadeOn      time.Time         `json:"made_on"`
	Horizon     *time.Time        `json:"horizon,omitempty"`
	Status      string            `json:"status"`
	ResolvedOn  *time.Time        `json:"resolved_on,omitempty"`
	SourceURL   string            `json:"source_url,omitempty"`
	Statement   map[string]string `json:"statement"`
	Verdict     map[string]string `json:"verdict,omitempty"`
}

type page struct {
	Key   string `json:"page_key"`
	Lang  string `json:"lang"`
	Title string `json:"title"`
	Body  string `json:"body_md"`
}

type manifest struct {
	GeneratedAt time.Time `json:"generated_at"`
	Articles    int       `json:"articles"`
	Predictions int       `json:"predictions"`
	Pages       int       `json:"pages"`
	Contains    string    `json:"contains"`
	Excludes    string    `json:"excludes"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	dir := os.Getenv("EXPORT_DIR")
	if dir == "" {
		log.Fatal("EXPORT_DIR is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("export dir: %v", err)
	}

	arts, err := readArticles(ctx, pool)
	if err != nil {
		log.Fatalf("articles: %v", err)
	}
	preds, err := readPredictions(ctx, pool)
	if err != nil {
		log.Fatalf("predictions: %v", err)
	}
	pages, err := readPages(ctx, pool)
	if err != nil {
		log.Fatalf("pages: %v", err)
	}

	for name, v := range map[string]any{
		"articles.json":    arts,
		"predictions.json": preds,
		"pages.json":       pages,
		"MANIFEST.json": manifest{
			GeneratedAt: time.Now().UTC(),
			Articles:    len(arts),
			Predictions: len(preds),
			Pages:       len(pages),
			Contains:    "published articles with translations, the prediction ledger, editable pages",
			Excludes:    "accounts, sessions, listings, comments, votes, moderation, analytics, payments — everything carrying personal data",
		},
	} {
		if err := writeJSON(filepath.Join(dir, name), v); err != nil {
			log.Fatalf("write %s: %v", name, err)
		}
	}
	log.Printf("exported %d articles, %d predictions, %d pages to %s", len(arts), len(preds), len(pages), dir)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// readArticles returns published articles with every translation attached. The
// author arrives as a display name and never as an id: an id is a handle on a
// person, and this file leaves the country.
func readArticles(ctx context.Context, pool *pgxpool.Pool) ([]article, error) {
	rows, err := pool.Query(ctx, `
		SELECT a.slug,
		       COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), split_part(u.email, '@', 1)),
		       COALESCE(a.category, ''), COALESCE(a.subcategory, ''), a.original_lang,
		       COALESCE(a.cover_url, ''), a.indexable, a.published_at
		  FROM articles a JOIN auth_users u ON u.id = a.author_id
		 WHERE a.status = 'published'
		 ORDER BY a.published_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []article
	bySlug := map[string]int{}
	for rows.Next() {
		var a article
		if err := rows.Scan(&a.Slug, &a.Author, &a.Category, &a.Subcategory,
			&a.OriginalLang, &a.CoverURL, &a.Indexable, &a.PublishedAt); err != nil {
			return nil, err
		}
		if !knownAuthors[a.Author] {
			return nil, fmt.Errorf("article %q is by %q, who is not one of the bylines this export may carry abroad; "+
				"add them to knownAuthors only once it is settled that their name may travel", a.Slug, a.Author)
		}
		bySlug[a.Slug] = len(out)
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	trs, err := pool.Query(ctx, `
		SELECT a.slug, t.lang, t.title, COALESCE(t.summary, ''), t.body_md, COALESCE(t.source, ''), t.status
		  FROM article_translations t JOIN articles a ON a.id = t.article_id
		 WHERE a.status = 'published'
		 ORDER BY a.slug, t.lang`)
	if err != nil {
		return nil, err
	}
	defer trs.Close()
	for trs.Next() {
		var slug string
		var t translation
		if err := trs.Scan(&slug, &t.Lang, &t.Title, &t.Summary, &t.Body, &t.Source, &t.Status); err != nil {
			return nil, err
		}
		if i, ok := bySlug[slug]; ok {
			out[i].Translations = append(out[i].Translations, t)
		}
	}
	return out, trs.Err()
}

// readPredictions returns the ledger with its texts folded into maps, keyed by
// language, and the article referenced by slug rather than by id.
func readPredictions(ctx context.Context, pool *pgxpool.Pool) ([]prediction, error) {
	rows, err := pool.Query(ctx, `
		SELECT COALESCE(a.slug, ''), p.made_on, p.horizon, p.status, p.resolved_on, COALESCE(p.source_url, ''),
		       COALESCE(jsonb_object_agg(x.lang, x.statement) FILTER (WHERE x.lang IS NOT NULL), '{}'::jsonb),
		       COALESCE(jsonb_object_agg(x.lang, x.verdict) FILTER (WHERE x.lang IS NOT NULL AND x.verdict <> ''), '{}'::jsonb)
		  FROM predictions p
		  LEFT JOIN articles a ON a.id = p.article_id
		  LEFT JOIN prediction_texts x ON x.prediction_id = p.id
		 GROUP BY a.slug, p.id, p.made_on, p.horizon, p.status, p.resolved_on, p.source_url
		 ORDER BY p.made_on`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []prediction{}
	for rows.Next() {
		var p prediction
		if err := rows.Scan(&p.ArticleSlug, &p.MadeOn, &p.Horizon, &p.Status, &p.ResolvedOn,
			&p.SourceURL, &p.Statement, &p.Verdict); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// readPages returns the editable pages. updated_by is a user id and stays home.
func readPages(ctx context.Context, pool *pgxpool.Pool) ([]page, error) {
	rows, err := pool.Query(ctx, `
		SELECT page_key, lang, COALESCE(title, ''), COALESCE(body_md, '')
		  FROM content_pages ORDER BY page_key, lang`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []page{}
	for rows.Next() {
		var p page
		if err := rows.Scan(&p.Key, &p.Lang, &p.Title, &p.Body); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
