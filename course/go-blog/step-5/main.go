// Step 5 — after the lesson "A struct: your own type".
//
// The three parallel maps of step 3 are gone. Everything the blog knows about
// one text now lives in one value, so adding an article is one line instead of
// three edits that had to be kept in step by hand.
package main

import (
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

// Article is one text. The fields start with a capital letter for a reason
// that becomes important in the lesson on packages.
type Article struct {
	Slug  string
	Title string
	Words int
	Lang  string
}

func readingTime(a Article) int {
	return (a.Words + 199) / 200
}

// longTitles keeps the titles with more than n letters. It still counts
// letters rather than bytes, which is what the runes lesson was for.
func longTitles(blog []Article, n int) []string {
	var out []string
	for _, a := range blog {
		if utf8.RuneCountInString(a.Title) > n {
			out = append(out, a.Title)
		}
	}
	return out
}

func main() {
	blog := []Article{
		{Slug: "salem", Title: "Сәлем, әлем", Words: 150, Lang: "kz"},
		{Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz"},
		{Slug: "shanyraq", Title: "Шаңырақ деген не", Words: 1000, Lang: "kz"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Менің блогым")
		fmt.Fprintln(w, "-------------")

		for i, a := range blog {
			fmt.Fprintf(w, "%d. %s — %d мин\n", i+1, a.Title, readingTime(a))
		}

		fmt.Fprintln(w, "\nұзын тақырыптар:", longTitles(blog, 11))
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
