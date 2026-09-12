// Package shop sells the site's own digital goods: the Go book first, and
// whatever the house publishes after it.
//
// It is deliberately not part of the publishing module. A product is not an
// article -- it has a price, an edition, a file to deliver and a buyer to
// remember -- and the one thing they share, the page frame, now comes from
// pkg/site.
//
// Money is the payments module's business: the shop records what is sold and
// to whom, and hands the charge to an acquirer it knows nothing about. Until
// one is connected, a product can be announced and its free fragment read, and
// nothing can be bought.
package shop

import (
	"context"
	"embed"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/payments"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/site"
)

//go:embed templates/*.html
var templateFiles embed.FS

// Module is the shop's product surface.
type Module struct {
	rt       *shanraq.Runtime
	store    *Store
	auth     *auth.Module
	payments *payments.Module
}

// New returns the module. The auth module identifies the buyer and the staff
// who edit a product; payments is held for the checkout that follows once an
// acquirer is connected.
func New(authModule *auth.Module, payModule *payments.Module) *Module {
	return &Module{auth: authModule, payments: payModule}
}

func (m *Module) Name() string { return "shop" }

// Init opens the store and adds the shop's pages to the site's template set.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	rt.Site.Add(nil, templateFiles, "templates/*.html")
	return nil
}

// Routes registers the shop's browser surface. Everything here is a page or a
// form, so the whole group carries the site's cross-origin guard and loads the
// session: the admin pages need to know who is asking, and the reader's pages
// show the header as signed in.
func (m *Module) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(site.VerifyOrigin(m.logger()))
		if m.auth != nil {
			r.Use(m.auth.LoadSession)
		}
		r.Get("/shop", m.handleIndex)
		r.Get("/shop/{slug}", m.handleProduct)
		r.Post("/shop/{slug}/notify", m.handleNotify)

		r.Get("/admin/shop", m.handleAdmin)
		r.Post("/admin/shop/{slug}", m.handleAdminSave)
	})
}

func (m *Module) logger() *zap.Logger {
	if m.rt == nil {
		return nil
	}
	return m.rt.Logger
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
