package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"myblog/blog"
)

// newTestStore opens a database in a folder the testing package removes on its
// own. Nothing outside the test is touched, and two tests never share a file.
func newTestStore(t *testing.T) *blog.Store {
	t.Helper()

	store, err := blog.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestReadHandler(t *testing.T) {
	store := newTestStore(t)
	if err := store.AddWithTags(t.Context(), blog.Article{
		Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz", Body: "Мәтін",
	}, []string{"дала"}); err != nil {
		t.Fatal(err)
	}

	handler := routes(store, t.TempDir())

	cases := []struct {
		name string
		path string
		code int
		body string
	}{
		{"мақала бар", "/read/dala", http.StatusOK, "Дала туралы"},
		{"мақала жоқ", "/read/joq", http.StatusNotFound, "joq"},
		{"тег беті", "/tag/дала", http.StatusOK, "Дала туралы"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// No server is started: a handler is an ordinary function, and the
			// recorder stands in for the network.
			req := httptest.NewRequest(http.MethodGet, c.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("код %d, күткеніміз %d", rec.Code, c.code)
			}
			if !strings.Contains(rec.Body.String(), c.body) {
				t.Errorf("жауапта %q жоқ", c.body)
			}
		})
	}
}

func TestStoreRoundTrip(t *testing.T) {
	store := newTestStore(t)

	want := blog.Article{Slug: "salem", Title: "Сәлем", Words: 150, Lang: "kz", Body: "Бірінші"}
	if err := store.AddWithTags(t.Context(), want, nil); err != nil {
		t.Fatal(err)
	}

	got, err := store.Get(t.Context(), "salem")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != want.Title || got.Words != want.Words {
		t.Errorf("алғанымыз %+v, күткеніміз %+v", got, want)
	}
}
