// Step 11 — after the lesson "Forms".
//
// The blog stops being read-only. A form sends a title, the handler checks it
// on the server — the browser's own check does not travel with the request —
// and answers with a redirect, so pressing F5 cannot file the same article
// twice.
package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"myblog/blog"
)

//go:embed templates/*.html
var files embed.FS

// Parsed once at start: a broken template stops the server now rather than
// surprising a reader later.
var tpl = template.Must(template.ParseFS(files, "templates/*.html"))

type recorder struct {
	http.ResponseWriter
	code int
}

func (rec *recorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s → %d, %s ішінде", r.Method, r.URL.Path, rec.code,
			time.Since(start).Round(time.Millisecond))
	})
}

func htmlType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// render builds the page in full before a single byte goes out.
func render(w http.ResponseWriter, code int, name string, data any) {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, name, data); err != nil {
		log.Println("үлгі:", err)
		http.Error(w, "серверде қате", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	if _, err := buf.WriteTo(w); err != nil {
		log.Println("жазу:", err)
	}
}

// draftBody stands in until the form learns to take the text as well — that
// is this lesson's exercise. Counting its words keeps the reading time honest
// instead of printing a zero.
const draftBody = "Мәтін әзірге жоқ."

// front is what list.html asks for. Draft carries back what the reader typed,
// because making someone retype eighty letters is rudeness, not validation.
type front struct {
	Articles []blog.Article
	Draft    string
	Err      string
}

func seed() *blog.Store {
	store := blog.NewStore()
	store.Add(blog.Article{Slug: "salem", Title: "Сәлем, әлем", Words: 150, Lang: "kz",
		Body: "Бірінші жазба. Блог осыдан басталады."})
	store.Add(blog.Article{Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz",
		Body: "Дала жазда сары, көктемде жасыл болады."})
	store.Add(blog.Article{Slug: "shanyraq", Title: "Шаңырақ деген не", Words: 1000, Lang: "kz",
		Body: "Шаңырақ — киіз үйдің төбесіндегі шеңбер."})
	return store
}

func routes(store *blog.Store) http.Handler {
	mux := http.NewServeMux()

	showList := func(w http.ResponseWriter, code int, p front) {
		p.Articles = store.All()
		render(w, code, "list.html", p)
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		showList(w, http.StatusOK, front{})
	})

	mux.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		title := strings.TrimSpace(r.FormValue("title"))

		switch {
		case title == "":
			showList(w, http.StatusBadRequest, front{Err: "тақырып бос"})
			return
		case utf8.RuneCountInString(title) > 80:
			showList(w, http.StatusBadRequest, front{Draft: title, Err: "тақырып тым ұзын"})
			return
		}

		// Making an address out of a Kazakh title is a job for a package we
		// have not met yet, so the number does for now.
		slug := fmt.Sprintf("maqala-%d", store.Len()+1)
		store.Add(blog.Article{Slug: slug, Title: title, Lang: "kz",
			Body: draftBody, Words: len(strings.Fields(draftBody))})

		// A handler that changed something does not draw a page: it sends the
		// browser to one that can be refreshed as often as the reader likes.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		render(w, http.StatusOK, "about.html", nil)
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(slug)
		if err != nil {
			render(w, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		}
		render(w, http.StatusOK, "article.html", map[string]any{"Article": a})
	})

	return mux
}

func main() {
	handler := logging(htmlType(routes(seed())))

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
