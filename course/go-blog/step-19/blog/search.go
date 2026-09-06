package blog

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

// Found is one article the search turned up, together with the piece of text
// the words were found in.
type Found struct {
	Article  Article
	Fragment template.HTML
}

// The markers snippet() wraps a match in. They travel through Go as text and
// only the pieces between them are turned into real tags.
const (
	markOpen  = "<mark>"
	markClose = "</mark>"
)

// highlighted escapes the fragment and then puts the two mark tags back. The
// body of an article came from a form, so handing the whole fragment to the
// template as safe HTML would hand a reader's script tag to every visitor.
// Everything is escaped; only our own markers become tags again.
func highlighted(s string) template.HTML {
	var b strings.Builder
	for _, part := range strings.Split(s, markOpen) {
		if inner := strings.SplitN(part, markClose, 2); len(inner) == 2 {
			b.WriteString(markOpen)
			b.WriteString(html.EscapeString(inner[0]))
			b.WriteString(markClose)
			b.WriteString(html.EscapeString(inner[1]))
			continue
		}
		b.WriteString(html.EscapeString(part))
	}
	return template.HTML(b.String())
}

// ftsQuery turns what a person typed into a query FTS5 will accept. Every word
// goes in quotes, so AND, OR and a star inside stay letters rather than
// becoming commands, and a star outside makes it match the beginning of a word.
//
// A parameter protects the text of the SQL; it does not protect the search
// query, because FTS5 parses that itself.
func ftsQuery(s string) string {
	var parts []string
	for _, w := range strings.Fields(s) {
		w = strings.ReplaceAll(w, `"`, `""`)
		parts = append(parts, `"`+w+`"*`)
	}
	return strings.Join(parts, " AND ")
}

// Search hands back what matches, best first. bm25 gives the title ten times
// the weight of the body, and its numbers are negative — the smaller, the
// better a match — so the order is ascending.
func (s *Store) Search(q string) ([]Found, error) {
	fts := ftsQuery(q)
	if fts == "" {
		return nil, nil
	}

	rows, err := s.db.Query(`select a.slug, a.title, a.words, a.lang, a.body, a.updated_at,
			snippet(articles_fts, 1, '<mark>', '</mark>', '…', 8)
		from articles_fts
		join articles a on a.id = articles_fts.rowid
		where articles_fts match ?
		order by bm25(articles_fts, 10.0, 1.0)`, fts)
	if err != nil {
		return nil, fmt.Errorf("іздеу %q: %w", q, err)
	}
	defer rows.Close()

	var out []Found
	for rows.Next() {
		var f Found
		var fragment string
		if err := rows.Scan(&f.Article.Slug, &f.Article.Title, &f.Article.Words,
			&f.Article.Lang, &f.Article.Body, &f.Article.UpdatedAt, &fragment); err != nil {
			return nil, err
		}
		f.Fragment = highlighted(fragment)
		out = append(out, f)
	}
	return out, rows.Err()
}
