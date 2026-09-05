// Step 9 — after the lesson "Handlers and middleware".
//
// Two things that every page needed are no longer written on every page. The
// log line and the Content-Type header moved into wrappers, so the handlers
// went back to doing only their own work — and adding a tenth page adds
// nothing to either of them.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"myblog/blog"
)

// recorder remembers the status code, because ResponseWriter does not keep it
// and the log line has nothing to print without it.
type recorder struct {
	http.ResponseWriter
	code int
}

func (rec *recorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}

// logging writes one line per request: what was asked, what was answered, how
// long it took.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s → %d, %s ішінде", r.Method, r.URL.Path, rec.code,
			time.Since(start).Round(time.Millisecond))
	})
}

// plainText sets the header once for every page. It has to run before the
// handler writes anything, which is exactly what a wrapper does.
func plainText(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
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
		fmt.Fprintln(w, "Менің блогым —", store.Len(), "мақала")
		fmt.Fprintln(w, "-------------")
		for i, a := range store.All() {
			fmt.Fprintf(w, "%d. %v → /read/%s\n", i+1, a, a.Slug)
		}
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Бұл блогты Go үйреніп жүрген адам жазады.")
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "мақала табылмады:", slug)
			return
		}
		fmt.Fprintln(w, a.Title)
		fmt.Fprintln(w, strings.Repeat("=", a.Letters()))
		fmt.Fprintln(w, a.Body)
		fmt.Fprintln(w, "оқу уақыты:", a.ReadingTime(), "мин")
	})

	return mux
}

func main() {
	// The outermost wrapper is written first and runs first: the log covers
	// the header wrapper as well, so a request lost inside it is still logged.
	handler := logging(plainText(routes(seed())))

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
