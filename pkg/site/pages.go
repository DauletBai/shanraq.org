package site

import (
	"context"
	"sync"
	"time"
)

// Page is one address worth offering to a search engine: where it is and when
// its content last changed. A zero time means "no better answer than now",
// which is what a standing page reports.
type Page struct {
	Path    string
	Updated time.Time
}

// Pages is how a module tells the site which of its addresses are public.
//
// The sitemap and the IndexNow submissions are built in one place, but the
// addresses are not knowable there: the shop knows which products are on show,
// the weather knows which towns have a forecast. Each module registers a source
// during Init; whoever builds the sitemap asks for all of them at request time,
// so a product published this morning is in the sitemap this morning.
type Pages struct {
	mu      sync.RWMutex
	sources []func(context.Context) []Page
}

// NewPages returns an empty registry.
func NewPages() *Pages { return &Pages{} }

// Add registers a source of public addresses. A source is called on every
// sitemap request, so it should read from a store rather than do real work.
func (p *Pages) Add(src func(context.Context) []Page) {
	if src == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sources = append(p.sources, src)
}

// All collects every registered address. A source that fails says so by
// returning nothing: a sitemap missing one module's pages is a smaller loss
// than a sitemap that fails to render.
func (p *Pages) All(ctx context.Context) []Page {
	p.mu.RLock()
	srcs := make([]func(context.Context) []Page, len(p.sources))
	copy(srcs, p.sources)
	p.mu.RUnlock()

	var out []Page
	seen := map[string]bool{}
	for _, src := range srcs {
		for _, pg := range src(ctx) {
			if pg.Path == "" || seen[pg.Path] {
				continue
			}
			seen[pg.Path] = true
			out = append(out, pg)
		}
	}
	return out
}
