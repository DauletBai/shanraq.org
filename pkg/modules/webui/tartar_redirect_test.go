package webui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestLegacyTartarExcerptRedirectsToFramedPages(t *testing.T) {
	module := newTestModule(t)
	router := chi.NewRouter()
	module.Routes(router)
	for old, current := range map[string]string{
		"/static/shop/tartar-i-ya-sample/":               "/book/tartar-i-ya",
		"/static/shop/tartar-i-ya-sample/chapter-2.html": "/book/tartar-i-ya/chapter-2",
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, old, nil))
		if response.Code != http.StatusMovedPermanently || response.Header().Get("Location") != current {
			t.Errorf("%s: status=%d location=%q, want 301 to %q",
				old, response.Code, response.Header().Get("Location"), current)
		}
	}
}
