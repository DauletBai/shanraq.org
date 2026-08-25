package articles

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// The archive: one address per day of publishing.
//
// The date at the top of every page was the last thing in the strip that led
// nowhere. Made into a link it answers a question readers ask of any
// publication — what came out today — and gives search engines a second
// dimension to index the archive along, beside places and rubrics.
//
// Only days that have something on them get an address in the sitemap. A page
// for a day nobody published on is a real URL a crawler can reach, and a few
// hundred of them would say nothing except that the site is mostly empty.

// archivePageSize is how many articles one day's page holds. A day rarely fills
// it, and when one does the rest is a page away.
const archivePageSize = 24

// ArchiveDay backs /archive/{YYYY-MM-DD}.
type ArchiveDay struct {
	Base
	Desc string
	// Human is the day written out: "24 августа 2026".
	Human string
	// Weekday is its name, shown beside the date.
	Weekday string
	// ISO is the machine form, used in links and in the dateline.
	ISO   string
	Posts []FeedItem
	// PrevURL and NextURL walk to the neighbouring days that actually have
	// something on them, so the reader is never sent into a run of empty pages.
	PrevURL string
	NextURL string
	// IsToday marks the current day, which is the one the strip links to.
	IsToday bool
}

// handleArchive renders one day of publishing.
func (m *Module) handleArchive(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	raw := chi.URLParam(r, "date")

	day, err := time.ParseInLocation("2006-01-02", raw, siteLoc())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	today := siteDay(siteNow())
	// A day in the future has nothing on it and never will have retroactively;
	// answering 404 keeps a crawler from walking the calendar forwards for ever.
	if day.After(today) {
		http.NotFound(w, r)
		return
	}

	arts, err := m.store.ListByDay(r.Context(), day, archivePageSize, 0, m.addressedTo(r))
	if err != nil {
		m.rt.Logger.Error("archive day", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	page := ArchiveDay{
		Human:   localizedDate(lang, day),
		Weekday: wxWeekday(day.Weekday(), lang),
		ISO:     day.Format("2006-01-02"),
		IsToday: day.Equal(today),
	}
	page.Posts = m.withOrgs(r.Context(), arts, feedItems(arts, lang))

	if prev, err := m.store.NeighbourDay(r.Context(), day, false); err == nil && !prev.IsZero() {
		page.PrevURL = "/archive/" + prev.Format("2006-01-02") + "?lang=" + lang
	}
	if next, err := m.store.NeighbourDay(r.Context(), day, true); err == nil && !next.IsZero() {
		page.NextURL = "/archive/" + next.Format("2006-01-02") + "?lang=" + lang
	}

	title := localizedDate(lang, day)
	page.Base = m.base(r, T(lang, "arch.title")+" — "+title, lang)
	page.Desc = T(lang, "arch.desc_pre") + " " + title + "."
	page.Base.CanonURL = canonURL("/archive/"+page.ISO, "", lang)
	page.Base.LangLinks = langLinks("/archive/"+page.ISO, "")
	// A day with nothing on it is a legitimate address — the link from the strip
	// lands on it on a quiet day — but it has no business in the index claiming
	// to be about that day.
	if len(page.Posts) == 0 {
		page.NoIndex = true
	}
	m.render(w, "archive", page)
}
