package articles

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// SearchPage backs the search results page for either scope.
type SearchPage struct {
	Base
	Query    string
	Scope    string // "articles" | "listings"
	Articles []FeedItem
	Listings []*Listing
	Count    int
	Searched bool // a query was submitted
}

// handleSearch runs a full-text search over articles or listings, chosen by the
// "in" parameter (default articles).
func (m *Module) handleSearch(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	scope := r.URL.Query().Get("in")
	if scope != "listings" {
		scope = "articles"
	}

	page := SearchPage{Base: m.base(r, T(lang, "search.title"), lang)}
	page.Query = query
	page.Scope = scope
	page.SearchQuery = query
	page.SearchScope = scope
	page.Searched = query != ""

	if query != "" {
		if scope == "listings" {
			if ls, err := m.listings.Search(r.Context(), query); err != nil {
				m.rt.Logger.Error("search listings", zap.Error(err))
			} else {
				page.Listings = ls
			}
			page.Count = len(page.Listings)
		} else {
			if arts, err := m.store.Search(r.Context(), query); err != nil {
				m.rt.Logger.Error("search articles", zap.Error(err))
			} else {
				page.Articles = toFeedItems(arts, lang)
			}
			page.Count = len(page.Articles)
		}
	}
	page.SidebarNews = m.latestNews(r, lang, 6)
	m.render(w, "search", page)
}
