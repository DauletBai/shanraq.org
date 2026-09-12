package articles

import (
	"net/http"

	"go.uber.org/zap"

	"shanraq.org/pkg/site"
)

// verifyOrigin guards this module's browser surface. The guard itself is the
// frame's (site.VerifyOrigin): every module that serves a form needs the same
// one, and a second copy is a second thing to get wrong.
func (m *Module) verifyOrigin(next http.Handler) http.Handler {
	var log *zap.Logger
	if m.rt != nil {
		log = m.rt.Logger
	}
	return site.VerifyOrigin(log)(next)
}

// sameOrigin reports whether a request came from our own pages.
func (m *Module) sameOrigin(r *http.Request) bool { return site.SameOrigin(r) }
