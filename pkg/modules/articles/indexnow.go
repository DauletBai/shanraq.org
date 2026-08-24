package articles

import (
	"context"
	"time"
)

// Telling search engines about the standing pages.
//
// A sitemap is an invitation to come by some day; IndexNow tells Bing, Yandex and
// everyone else who supports it right now. Articles have been going there since
// publication, but the standing pages — "Analytics", "Exchange rates", the rules,
// the tariffs — never went at all: nobody "publishes" them, so no occasion to send
// a submission ever arose.
//
// So the occasions are appointed here. After a start-up the site declares its whole
// standing set: a deployment changes the markup and the texts, and that is exactly
// when a search engine should come back. After that, once a day, the pages whose
// content changes daily on its own are submitted.

const (
	// indexNowSettle is the pause after start-up. The submission points at a key
	// file on our own domain, and sending it before the site answers at all means
	// being refused.
	indexNowSettle = 2 * time.Minute
	// indexNowEvery is how often to announce the pages that change daily.
	indexNowEvery = 24 * time.Hour
)

// indexNowDaily are the pages that change every day by themselves: the front
// page's feed, the exchange rates and the audience counters.
var indexNowDaily = []string{"/", "/rates", "/analytics"}

// runIndexNow announces the site's standing pages until the context is
// cancelled.
func (m *Module) runIndexNow(ctx context.Context) {
	if m.syndicate == nil {
		return
	}
	select {
	case <-ctx.Done():
		return
	case <-time.After(indexNowSettle):
	}

	all := append([]string{"/"}, publicPages...)
	m.syndicate.SubmitURLs(m.pageURLs(all), "постоянные страницы")

	t := time.NewTicker(indexNowEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.syndicate.SubmitURLs(m.pageURLs(indexNowDaily), "ежедневные страницы")
		}
	}
}

// pageURLs expands paths into full addresses in all three languages — exactly the
// ones the sitemap and the canonical carry. A submission for an address the site
// does not declare as its own gives a search engine nothing.
func (m *Module) pageURLs(paths []string) []string {
	site := m.rt.Config.PublicBase()
	out := make([]string, 0, len(paths)*len(Langs))
	for _, p := range paths {
		for _, lang := range Langs {
			out = append(out, site+canonURL(p, "", lang))
		}
	}
	return out
}
