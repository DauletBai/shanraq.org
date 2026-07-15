package articles

import (
	"net/http"

	"go.uber.org/zap"
)

// Sana Qyran is the platform's AI columnist — an AI model that publishes
// clearly-labeled opinion. The identity is transparent by design: every byline
// carries an "AI opinion" badge and links to the author profile.
const (
	SanaAuthorID = "5a2a0000-0000-0000-0000-000000000001"
	SanaName     = "Сана Қыран"
	SanaSlug     = "sana"
)

// authorDisplay returns the byline name and whether the author is the AI
// columnist (so templates can show the "AI opinion" badge).
func authorDisplay(a *Article) (string, bool) {
	if a.AuthorID.String() == SanaAuthorID {
		return SanaName, true
	}
	return a.AuthorName(), false
}

// AuthorPage is the template context for an author profile.
type AuthorPage struct {
	Base
	Name  string
	IsAI  bool
	Posts []FeedItem
}

// handleAuthorSana renders the AI columnist's profile: a transparent bio plus
// a grid of published columns (good for readers and for internal SEO links).
func (m *Module) handleAuthorSana(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	arts, err := m.store.ListPublishedByAuthor(r.Context(), SanaAuthorID, 60)
	if err != nil {
		m.rt.Logger.Error("author articles", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	items := make([]FeedItem, 0, len(arts))
	for _, a := range arts {
		tr, served := a.Translation(lang)
		if tr == nil {
			continue
		}
		summary := tr.Summary
		if summary == "" {
			summary = excerpt(stripMD(tr.BodyMD), 170)
		}
		items = append(items, FeedItem{
			Slug: a.Slug, Title: tr.Title, Summary: summary,
			AuthorName: SanaName, AIAuthor: true, ServedLang: served,
			Category: a.Category, Subcategory: a.Subcategory, CoverURL: a.CoverURL,
			Published: a.PublishedAt, Views: a.ViewsCount, Score: a.Score,
			AvailableLangs: a.AvailableLangs(),
		})
	}
	page := AuthorPage{Base: m.base(r, SanaName, lang)}
	page.Name = SanaName
	page.IsAI = true
	page.Desc = T(lang, "author.sana_bio")
	page.Posts = items
	m.render(w, "author", page)
}
