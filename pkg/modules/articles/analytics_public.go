package articles

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// The public audience page.
//
// A publication that asks other people for verifiable numbers has to publish
// its own, including the ones that flatter nobody. The figures here are the
// same ones the admin panel shows, drawn from the same table by the same code —
// there is no second, friendlier set.
//
// What is deliberately absent: any cut that could point at a person. Individual
// addresses are never recorded, only page kinds, and country counts are shown
// without crossing them with anything. A country is not a reader; a country
// crossed with a page and a day can be.

// publicStatsTTL is how long an assembled page is reused. The underlying rows
// change once a day, so a shorter cache would only cost queries.
const publicStatsTTL = time.Hour

// PublicStats backs /analytics.
type PublicStats struct {
	Traffic TrafficChart
	Base
	Guests GuestAnalytics

	// One total that lives outside the counter: how much there is to read.
	// The stock figures beside it -- listings, subscribers, authors, comments --
	// were removed: at this size they say more about how small the place is
	// than about whether it is worth reading, which is what the page is for.
	Articles int64

	Since   string // the day the counter started measuring
	Updated string // when these figures were assembled
}

type statsCache struct {
	mu   sync.Mutex
	at   time.Time
	byLg map[string]PublicStats
}

var publicStats = statsCache{byLg: map[string]PublicStats{}}

// handlePublicStats renders the audience page.
func (m *Module) handlePublicStats(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)

	publicStats.mu.Lock()
	fresh := time.Since(publicStats.at) < publicStatsTTL
	cached, ok := publicStats.byLg[lang]
	publicStats.mu.Unlock()

	page := cached
	if !fresh || !ok {
		page = m.buildPublicStats(r.Context(), lang)
		publicStats.mu.Lock()
		if !fresh {
			// A stale set is dropped whole: keeping one language from the old
			// hour beside another from the new one would put two different
			// days on the same page.
			publicStats.byLg = map[string]PublicStats{}
			publicStats.at = time.Now()
		}
		publicStats.byLg[lang] = page
		publicStats.mu.Unlock()
	}

	page.Base = m.base(r, T(lang, "stats.title"), lang)
	page.Desc = T(lang, "stats.desc")
	// Outside the cache: the switches are part of the address, so a cached page
	// would serve one reader's choice of audience to the next one.
	page.Traffic = m.trafficChartFrom(r)
	m.render(w, "stats", page)
}

// buildPublicStats assembles the figures. Every part degrades to zero rather
// than failing the page: a missing count is better than a blank page that says
// nothing at all.
func (m *Module) buildPublicStats(ctx context.Context, lang string) PublicStats {
	p := PublicStats{Guests: m.guestAnalytics(ctx, lang)}
	p.Updated = time.Now().Format("02.01.2006")

	count := func(q string) int64 {
		var n int64
		if err := m.rt.DB.QueryRow(ctx, q).Scan(&n); err != nil {
			m.rt.Logger.Warn("public stats count", zap.String("q", q), zap.Error(err))
			return 0
		}
		return n
	}
	p.Articles = count(`SELECT count(*) FROM articles WHERE status = 'published'`)

	var since time.Time
	if err := m.rt.DB.QueryRow(ctx,
		`SELECT min(day) FROM analytics_daily`).Scan(&since); err == nil && !since.IsZero() {
		p.Since = since.Format("02.01.2006")
	}
	return p
}
