package jobs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"shanraq.org/pkg/shanraq"
)

// mount registers the module's routes on a fresh router.
func mount(m *Module) chi.Router {
	r := chi.NewRouter()
	m.Routes(r)
	return r
}

// The console is a browser page with a session cookie; /jobs wants an API key
// and a bearer token. The explorer therefore showed an error where the queue
// should be. The second mount exists so the console can be guarded by what a
// browser actually carries — and every one of its routes has to be behind that
// guard, not just the one somebody remembered.
func TestConsoleRoutesAllSitBehindTheGuard(t *testing.T) {
	var guarded int
	m := &Module{rt: &shanraq.Runtime{}}
	WithConsoleMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Stop here: reaching the handler would need a database, and what
			// is under test is that the guard is reached at all.
			guarded++
			w.WriteHeader(http.StatusTeapot)
		})
	})(m)
	r := mount(m)

	for _, c := range []struct{ method, path string }{
		{http.MethodGet, "/console/jobs"},
		{http.MethodPost, "/console/jobs"},
		{http.MethodPost, "/console/jobs/6f1c1e3e-0000-0000-0000-000000000000/retry"},
		{http.MethodPost, "/console/jobs/6f1c1e3e-0000-0000-0000-000000000000/cancel"},
	} {
		before := guarded
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
		if rec.Code != http.StatusTeapot {
			t.Errorf("%s %s answered %d without passing the guard", c.method, c.path, rec.Code)
		}
		if guarded != before+1 {
			t.Errorf("%s %s did not reach the guard", c.method, c.path)
		}
	}
}

// Fail closed: an unguarded console mount would let anyone enqueue work that
// spends money, so with no guard configured the door must not exist at all.
func TestConsoleRoutesAreAbsentWithoutAGuard(t *testing.T) {
	r := mount(&Module{rt: &shanraq.Runtime{}})
	for _, path := range []string{"/console/jobs", "/console/jobs/x/retry"} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s answered %d with no console guard configured; want it unregistered", path, rec.Code)
		}
	}
}
