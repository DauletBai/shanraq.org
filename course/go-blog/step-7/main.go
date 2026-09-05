// Step 7 — after the lesson "HTTP for real".
//
// One handler still answers everything, but it now answers properly: the
// status code says what happened, Content-Type says what the body is, and an
// address that holds nothing gets a 404 instead of an empty page. Sorting the
// addresses out by hand is the part the next lesson takes away.
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

func list(w http.ResponseWriter, store *blog.Store) {
	fmt.Fprintln(w, "Менің блогым —", store.Len(), "мақала")
	fmt.Fprintln(w, "-------------")
	for i, a := range store.All() {
		fmt.Fprintf(w, "%d. %v → /read/%s\n", i+1, a, a.Slug)
	}
}

func read(w http.ResponseWriter, store *blog.Store, slug string) {
	a, err := store.Get(slug)
	if err != nil {
		// The store says what went wrong; the handler decides what the reader
		// sees. A missing article is a 404, not a broken server.
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "мақала табылмады:", slug)
		return
	}
	fmt.Fprintln(w, a.Title)
	fmt.Fprintln(w, strings.Repeat("=", a.Letters()))
	fmt.Fprintln(w, a.Body)
	fmt.Fprintln(w, "оқу уақыты:", a.ReadingTime(), "мин")
}

func main() {
	store := seed()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "тек GET")
			return
		}

		// Headers first, then the code, then the body: the order is one way
		// only, and a header set after WriteHeader never reaches the reader.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		switch {
		case r.URL.Path == "/":
			list(w, store)
		case strings.HasPrefix(r.URL.Path, "/read/"):
			read(w, store, strings.TrimPrefix(r.URL.Path, "/read/"))
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "мұндай бет жоқ:", r.URL.Path)
		}
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
