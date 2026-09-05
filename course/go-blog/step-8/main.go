// Step 8 — after the lesson "Routes".
//
// The hand-written checks are gone. ServeMux now holds the address list, and
// it answers 404 and 405 by itself; the handlers only do their own work. The
// slug is no longer cut out of the path by hand — the pattern names it.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"myblog/blog"
)

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

func plain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func main() {
	store := seed()
	mux := http.NewServeMux()

	// "/{$}" is the front page only. Without the {$} the pattern "/" would
	// catch every address that no other pattern claims, and the 404 below
	// would never happen.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		plain(w)
		fmt.Fprintln(w, "Менің блогым —", store.Len(), "мақала")
		fmt.Fprintln(w, "-------------")
		for i, a := range store.All() {
			fmt.Fprintf(w, "%d. %v → /read/%s\n", i+1, a, a.Slug)
		}
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		plain(w)
		fmt.Fprintln(w, "Бұл блогты Go үйреніп жүрген адам жазады.")
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		plain(w)
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

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
