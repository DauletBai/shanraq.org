// Step 3 — after the lesson "A map: key and value".
//
// The blog stops being one article and becomes a list. Both of the last two
// lessons earn their place here: the slice keeps the order the articles are
// shown in, and the maps find a title and a length by the article's address.
// Neither could do the other's job -- which is the whole reason Go has both.
package main

import (
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

func readingTime(words int) int {
	return (words + 199) / 200
}

// longer keeps the titles with more than n letters. From the slices lesson.
func longer(titles []string, n int) []string {
	var out []string
	for _, t := range titles {
		if utf8.RuneCountInString(t) > n {
			out = append(out, t)
		}
	}
	return out
}

func main() {
	// The order of the front page lives in the slice.
	slugs := []string{"salem", "dala", "shanyraq"}

	// What is known about each article lives in the maps.
	titles := map[string]string{
		"salem":    "Сәлем, әлем",
		"dala":     "Дала туралы",
		"shanyraq": "Шаңырақ деген не",
	}
	words := map[string]int{
		"salem":    150,
		"dala":     400,
		"shanyraq": 1000,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Менің блогым")
		fmt.Fprintln(w, "-------------")

		for i, slug := range slugs {
			title, ok := titles[slug]
			if !ok {
				continue
			}
			fmt.Fprintf(w, "%d. %s — %d мин\n", i+1, title, readingTime(words[slug]))
		}

		all := make([]string, 0, len(slugs))
		for _, slug := range slugs {
			all = append(all, titles[slug])
		}
		fmt.Fprintln(w, "\nұзын тақырыптар:", longer(all, 11))
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
