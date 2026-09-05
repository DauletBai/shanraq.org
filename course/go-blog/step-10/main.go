// Step 10 — after the lesson "Templates".
//
// The HTML left the code. The handlers now hand data to a template and say
// which one; the markup lives in templates/ where an editor can see it. The
// page is assembled in a buffer first, so a template that fails halfway does
// not leave the reader with half a page and a 200 that cannot be taken back.
package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

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

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		render(w, http.StatusOK, "list.html", map[string]any{"Articles": store.All()})
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
