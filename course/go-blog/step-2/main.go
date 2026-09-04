// Step 2 — after the lesson "Strings and runes".
//
// The page now describes one article: a title, how many letters it holds and
// how long it takes to read. Three lessons' worth of language, and nothing in
// here that has not been taught: variables, a function, an if, and counting
// letters rather than bytes.
package main

import (
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

// readingTime is the exercise from the lesson on functions. Two hundred words
// a minute, rounded up: 199 is added before the division because an unfinished
// minute still counts, and exactly 400 words must give 2 rather than 3.
func readingTime(words int) int {
	return (words + 199) / 200
}

func main() {
	title := "Шаңырақ деген не"
	words := 400

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, title)
		fmt.Fprintln(w, "әріп саны:", utf8.RuneCountInString(title))
		fmt.Fprintln(w, "оқу уақыты:", readingTime(words), "мин")

		if readingTime(words) > 5 {
			fmt.Fprintln(w, "ұзын мақала")
		} else {
			fmt.Fprintln(w, "қысқа мақала")
		}
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
