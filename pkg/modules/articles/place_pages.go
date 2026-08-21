package articles

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Place pages: one address per place, one feed per address.
//
// The point is not filtering but reach. The same location data used as a filter
// narrows what a reader sees; used as a page it creates something that did not
// exist before — an address a search engine can index and a person can be sent
// to. Somebody googling "Качар ЖКХ" cannot find us today because there is no
// page about Kachar to find.
//
// It also answers the question of regional editions without founding any. A
// section for Kostanay oblast inside one publication shares the domain's
// standing, its comment rules and its readers; ten separate regional sites
// would each start from nothing, which is where this site was a month ago.

// placePageSize is how many articles one page of a place's feed holds.
const placePageSize = 21

// PlacePage backs /place/{slug}.
type PlacePage struct {
	Base
	PlaceName  string
	PlaceLabel string // "Качар, Костанайская область"
	Kind       string // country | region | city | town | village | district
	Slug       string
	Posts      []FeedItem
	Ancestors  []GeoNode // the way back up, for the breadcrumb
	Children   []GeoNode // places inside this one that a reader can descend into

	Page    int
	PrevURL string
	NextURL string
}

// handlePlace renders the feed of one place: what was published for it, and for
// anywhere inside it.
func (m *Module) handlePlace(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	slug := chi.URLParam(r, "slug")

	node, err := m.geo.BySlug(r.Context(), slug, lang)
	if err != nil {
		m.rt.Logger.Error("place by slug", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if node == nil {
		http.NotFound(w, r)
		return
	}
	id, err := uuid.Parse(node.ID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pageNo := 1
	if p, _ := strconv.Atoi(r.URL.Query().Get("page")); p > 1 {
		pageNo = p
	}
	offset := (pageNo - 1) * placePageSize

	// One row over the page size answers "is there a next page" without a
	// second count query.
	arts, err := m.store.ListForPlace(r.Context(), id, placePageSize+1, offset)
	if err != nil {
		m.rt.Logger.Error("place feed", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	hasNext := len(arts) > placePageSize
	if hasNext {
		arts = arts[:placePageSize]
	}

	page := PlacePage{Base: m.base(r, node.Name+" — "+T(lang, "place.title_suffix"), lang)}
	page.PlaceName = node.Name
	page.Kind = node.Kind
	page.Slug = slug
	page.Posts = m.withOrgs(r.Context(), arts, feedItems(arts, lang))
	page.Page = pageNo

	if label, err := m.geo.PlaceLabel(r.Context(), id, lang); err == nil {
		page.PlaceLabel = label
	}
	if chain, err := m.geo.Ancestry(r.Context(), id, lang); err == nil && len(chain) > 1 {
		// Drop the node itself: a breadcrumb points at where you can go, not
		// at where you already are.
		page.Ancestors = chain[:len(chain)-1]
	}
	if kids, err := m.geo.Children(r.Context(), id, lang); err == nil {
		page.Children = kids
	}

	base := "/place/" + slug + "?lang=" + lang
	if pageNo > 1 {
		page.PrevURL = base
		if pageNo > 2 {
			page.PrevURL = base + "&page=" + strconv.Itoa(pageNo-1)
		}
		// A page past the end of a small place is a real URL a crawler can
		// reach, and it must not join the index claiming to be about the place.
		if len(page.Posts) == 0 {
			page.NoIndex = true
		}
	}
	if hasNext {
		page.NextURL = base + "&page=" + strconv.Itoa(pageNo+1)
	}

	m.render(w, "place", page)
}
