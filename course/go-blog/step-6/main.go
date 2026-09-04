// Step 6 — after the lesson "Packages and modules".
//
// The same blog, now in two packages: what the blog knows lives in blog/, and
// this file only shows it. Nothing here reaches inside the store — it asks.
package main

import (
	"fmt"
	"log"
	"net/http"

	"myblog/blog"
)

func main() {
	store := blog.NewStore()
	store.Add(blog.Article{Slug: "salem", Title: "Сәлем, әлем", Words: 150, Lang: "kz"})
	store.Add(blog.Article{Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz"})
	store.Add(blog.Article{Slug: "shanyraq", Title: "Шаңырақ деген не", Words: 1000, Lang: "kz"})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Менің блогым —", store.Len(), "мақала")
		fmt.Fprintln(w, "-------------")
		for i, a := range store.All() {
			fmt.Fprintf(w, "%d. %v\n", i+1, a)
		}
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
