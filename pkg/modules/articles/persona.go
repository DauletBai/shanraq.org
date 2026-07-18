package articles

import (
	"net/http"
	"sort"

	"go.uber.org/zap"
)

// AI Dake is the platform's AI columnist — an AI model that publishes
// clearly-labeled opinion. The identity is transparent by design: every byline
// carries an "AI opinion" badge and links to the author profile. The display
// name is intentionally the same Latin form in all three languages; only the
// surrounding "(AI)" descriptor is localized. (The internal slug/email stay
// "sana" for URL and history stability.)
const (
	SanaAuthorID = "5a2a0000-0000-0000-0000-000000000001"
	SanaName     = "AI Dake"
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

// CatCount is one rubric's article count for the profile breakdown.
type CatCount struct {
	Slug string
	N    int
}

// AuthorPage is the template context for an author profile.
type AuthorPage struct {
	Base
	Name  string
	IsAI  bool
	Posts []FeedItem
	// Public analytics showcase.
	Count int
	Views int64
	Score int
	ByCat []CatCount
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
	var views int64
	var score int
	catN := map[string]int{}
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
		views += a.ViewsCount
		score += a.Score
		catN[a.Category]++
	}
	byCat := make([]CatCount, 0, len(catN))
	for slug, n := range catN {
		byCat = append(byCat, CatCount{Slug: slug, N: n})
	}
	sort.Slice(byCat, func(i, j int) bool {
		if byCat[i].N != byCat[j].N {
			return byCat[i].N > byCat[j].N
		}
		return byCat[i].Slug < byCat[j].Slug
	})

	page := AuthorPage{Base: m.base(r, SanaName, lang)}
	page.Name = SanaName
	page.IsAI = true
	page.Desc = T(lang, "author.sana_bio")
	page.Posts = items
	page.Count = len(items)
	page.Views = views
	page.Score = score
	page.ByCat = byCat
	m.render(w, "author", page)
}
