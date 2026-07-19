package articles

import (
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	Name     string
	Initials string // avatar fallback for human authors (e.g. "АН")
	IsAI     bool
	Posts    []FeedItem
	// Public analytics showcase.
	Count int
	Views int64
	Score int
	ByCat []CatCount
}

// handleAuthor renders any author's public profile — name, karma, and a grid of
// their published articles (good for readers discovering more and for internal
// SEO links). The AI columnist keeps its friendly slug and richer bio/badge.
func (m *Module) handleAuthor(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	raw := chi.URLParam(r, "id")

	if raw == SanaSlug || raw == SanaAuthorID {
		m.renderAuthor(w, r, lang, SanaAuthorID, SanaName, true)
		return
	}
	uid, err := uuid.Parse(raw)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	first, last, _ := m.auth.AuthorIdentity(r.Context(), uid)
	m.renderAuthor(w, r, lang, uid.String(), strings.TrimSpace(first+" "+last), false)
}

// renderAuthor builds and renders a profile for the given author.
func (m *Module) renderAuthor(w http.ResponseWriter, r *http.Request, lang, authorID, name string, isAI bool) {
	arts, err := m.store.ListPublishedByAuthor(r.Context(), authorID, 60)
	if err != nil {
		m.rt.Logger.Error("author articles", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// Fall back to the byline of the author's own articles when the identity
	// lookup gave nothing; a truly unknown/nameless author is a 404.
	if strings.TrimSpace(name) == "" && len(arts) > 0 {
		name = arts[0].AuthorName()
	}
	if strings.TrimSpace(name) == "" {
		http.NotFound(w, r)
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
			AuthorName: name, AuthorID: authorID, AIAuthor: isAI, ServedLang: served,
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

	page := AuthorPage{Base: m.base(r, name, lang)}
	page.Name = name
	page.Initials = initials(name)
	page.IsAI = isAI
	if isAI {
		page.Desc = T(lang, "author.sana_bio")
	}
	page.Posts = items
	page.Count = len(items)
	page.Views = views
	page.Score = score
	page.ByCat = byCat
	m.render(w, "author", page)
}

// initials returns up to two uppercased leading letters of a name ("Айгерим
// Нурланова" -> "АН"), used as the avatar fallback for human authors.
func initials(name string) string {
	out := ""
	for i, f := range strings.Fields(name) {
		if i >= 2 {
			break
		}
		rs := []rune(f)
		if len(rs) > 0 {
			out += strings.ToUpper(string(rs[0]))
		}
	}
	return out
}
