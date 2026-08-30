package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serveThrough(t *testing.T, h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	conditionalGet(h).ServeHTTP(w, req)
	return w
}

func TestConditionalGetTagsHTML(t *testing.T) {
	page := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html>hello</html>"))
	}
	w := serveThrough(t, page, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	tag := w.Header().Get("ETag")
	if tag == "" {
		t.Fatal("no ETag on an HTML response")
	}
	if w.Body.String() != "<html>hello</html>" {
		t.Errorf("body altered: %q", w.Body.String())
	}
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache (store but revalidate)", cc)
	}

	// The whole point: the same page, asked for again, costs no bytes.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", tag)
	w2 := serveThrough(t, page, req)
	if w2.Code != http.StatusNotModified {
		t.Fatalf("repeat request = %d, want 304", w2.Code)
	}
	if w2.Body.Len() != 0 {
		t.Errorf("304 carried a body of %d bytes", w2.Body.Len())
	}
}

// A changed page must produce a changed tag, or crawlers would never see edits.
func TestConditionalGetTagFollowsContent(t *testing.T) {
	tagOf := func(body string) string {
		h := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(body))
		}
		return serveThrough(t, h, httptest.NewRequest(http.MethodGet, "/", nil)).Header().Get("ETag")
	}
	if tagOf("<html>a</html>") == tagOf("<html>b</html>") {
		t.Error("different pages share an ETag")
	}
	// Two separate requests for the same body, named so it is plain that this
	// is a determinism check and not an expression compared with itself.
	firstServe, secondServe := tagOf("<html>a</html>"), tagOf("<html>a</html>")
	if firstServe != secondServe {
		t.Errorf("the same page produced two different ETags: %q and %q", firstServe, secondServe)
	}
}

func TestConditionalGetIgnoresNonHTML(t *testing.T) {
	cases := []struct {
		name   string
		ct     string
		status int
	}{
		{"xml sitemap", "application/xml; charset=utf-8", http.StatusOK},
		{"redirect", "text/html", http.StatusSeeOther},
		{"not found", "text/html", http.StatusNotFound},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", c.ct)
				w.WriteHeader(c.status)
				_, _ = w.Write([]byte("body"))
			}
			w := serveThrough(t, h, httptest.NewRequest(http.MethodGet, "/", nil))
			if w.Header().Get("ETag") != "" {
				t.Errorf("%s should not be tagged", c.name)
			}
			if w.Code != c.status {
				t.Errorf("status = %d, want %d", w.Code, c.status)
			}
		})
	}
}

// A POST must pass through untouched: buffering a form submit to hash it would
// be pointless and would break streaming handlers.
func TestConditionalGetSkipsWrites(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html>ok</html>"))
	}
	w := serveThrough(t, h, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")))
	if w.Header().Get("ETag") != "" {
		t.Error("POST response was tagged")
	}
}

func TestMatchesETag(t *testing.T) {
	const tag = `"abc123"`
	for _, ok := range []string{tag, "*", `"other", "abc123"`, `W/"abc123"`} {
		if !matchesETag(ok, tag) {
			t.Errorf("%q should match", ok)
		}
	}
	for _, bad := range []string{"", `"nope"`, `"abc12"`, `"abc1234"`} {
		if matchesETag(bad, tag) {
			t.Errorf("%q should not match", bad)
		}
	}
}
