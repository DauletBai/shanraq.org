package articles

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// FavoritesPage backs the reader's "saved" page — bookmarked articles and listings.
type FavoritesPage struct {
	Base
	Articles []FeedItem
	Listings []*Listing
}

// toFeedItems projects hydrated articles into feed cards for the current language.
func toFeedItems(arts []*Article, lang string) []FeedItem {
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
		name, aiAuthor := authorDisplay(a)
		items = append(items, FeedItem{
			Slug:           a.Slug,
			Title:          tr.Title,
			Summary:        summary,
			AuthorName:     name,
			AIAuthor:       aiAuthor,
			ServedLang:     served,
			Category:       a.Category,
			Subcategory:    a.Subcategory,
			CoverURL:       a.CoverURL,
			Published:      a.PublishedAt,
			Views:          a.ViewsCount,
			Score:          a.Score,
			IsAI:           tr.Source == "ai",
			AvailableLangs: a.AvailableLangs(),
		})
	}
	return items
}

// handleFavorites renders the current user's saved articles and listings.
func (m *Module) handleFavorites(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	page := FavoritesPage{Base: m.base(r, T(lang, "fav.title"), lang)}
	if arts, err := m.store.ListFavorited(r.Context(), uid); err != nil {
		m.rt.Logger.Error("favorited articles", zap.Error(err))
	} else {
		page.Articles = toFeedItems(arts, lang)
	}
	if ls, err := m.listings.ListFavorited(r.Context(), uid); err != nil {
		m.rt.Logger.Error("favorited listings", zap.Error(err))
	} else {
		page.Listings = ls
	}
	page.SidebarNews = m.latestNews(r, lang, 6)
	m.render(w, "favorites", page)
}

// handleFavoriteToggle flips a bookmark for the current user and returns to the
// page they came from (server-rendered, no JavaScript required).
func (m *Module) handleFavoriteToggle(w http.ResponseWriter, r *http.Request) {
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	itemType := chi.URLParam(r, "type")
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil || !favTypes[itemType] {
		http.NotFound(w, r)
		return
	}
	if _, err := m.favs.Toggle(r.Context(), uid, itemType, itemID); err != nil {
		m.rt.Logger.Error("toggle favorite", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	back := r.Header.Get("Referer")
	if back == "" {
		back = "/favorites"
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}
