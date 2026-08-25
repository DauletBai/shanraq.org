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

	// Calendar is the month this day sits in, so a reader can jump straight to
	// any day of it instead of stepping through them one at a time.
	Calendar [][]ArchCell
	Weekdays []string
	// MonthName is the month being shown, PrevMonth and NextMonth its
	// neighbours. NextMonth is empty once it would run past today.
	MonthName string
	PrevMonth string
	NextMonth string
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

	// The month around this day. The reader asked for a particular date, so the
	// calendar opens on its month rather than on the current one.
	addressed := m.addressedTo(r)
	if has, herr := m.store.DaysWithArticlesInMonth(r.Context(), day.Year(), day.Month(), addressed); herr != nil {
		m.rt.Logger.Warn("archive calendar", zap.Error(herr))
	} else {
		page.Calendar = archiveCalendar(day, has, day, today)
		page.Weekdays = archWeekdays(lang)
		page.MonthName = fxMonthFull(day.Month(), lang) + " " + day.Format("2006")
		page.PrevMonth = day.AddDate(0, -1, 0).Format("2006-01-02")
		// Never offer a month that has not started: a calendar of days that
		// cannot have anything on them is a walk into 404s.
		if next := firstOfNextMonth(day); !next.After(today) {
			page.NextMonth = next.Format("2006-01-02")
		}
	}

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

// ArchCell is one square of the month calendar.
type ArchCell struct {
	// Day is the day of the month; zero for the blanks that pad the first and
	// last weeks.
	Day int
	// ISO addresses this day's page. Empty when the square is a pad or when
	// nothing was published — a square that leads to an empty page is worse
	// than one that plainly does not lead anywhere.
	ISO string
	// Has marks a day with something on it, Today the current day, Sel the day
	// being viewed.
	Has   bool
	Today bool
	Sel   bool
	// Ahead marks a day that has not happened yet.
	Ahead bool
}

// archiveCalendar lays out one month, Monday first.
//
// Monday first because that is how a calendar is read here; Go counts weeks
// from Sunday, and taking its numbering literally would shift every date one
// column left.
func archiveCalendar(view time.Time, has map[int]bool, sel, today time.Time) [][]ArchCell {
	first := time.Date(view.Year(), view.Month(), 1, 0, 0, 0, 0, siteLoc())
	// Monday is 0 here: Go gives Sunday 0, and a plain subtraction would put
	// Sunday in the first column instead of the last.
	lead := (int(first.Weekday()) + 6) % 7
	days := first.AddDate(0, 1, -1).Day()

	weeks := [][]ArchCell{}
	row := make([]ArchCell, 0, 7)
	for i := 0; i < lead; i++ {
		row = append(row, ArchCell{})
	}
	for d := 1; d <= days; d++ {
		day := time.Date(view.Year(), view.Month(), d, 0, 0, 0, 0, siteLoc())
		c := ArchCell{
			Day:   d,
			Has:   has[d],
			Today: day.Equal(today),
			Sel:   day.Equal(sel),
			Ahead: day.After(today),
		}
		if c.Has {
			c.ISO = day.Format("2006-01-02")
		}
		row = append(row, c)
		if len(row) == 7 {
			weeks = append(weeks, row)
			row = make([]ArchCell, 0, 7)
		}
	}
	for len(row) > 0 && len(row) < 7 {
		row = append(row, ArchCell{})
	}
	if len(row) == 7 {
		weeks = append(weeks, row)
	}
	return weeks
}

// archWeekdays are the column headings, Monday first.
func archWeekdays(lang string) []string {
	switch lang {
	case LangKZ:
		return []string{"дс", "сс", "ср", "бс", "жм", "сн", "жс"}
	case LangEN:
		return []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	return []string{"пн", "вт", "ср", "чт", "пт", "сб", "вс"}
}

// firstOfNextMonth is the first day of the month after this one, used to step
// the calendar forward without landing on the 31st of a 30-day month.
func firstOfNextMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, siteLoc()).AddDate(0, 1, 0)
}
