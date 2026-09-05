// Package blog is everything the blog knows about a text, kept apart from the
// server that shows it. Only the names starting with a capital letter leave
// this folder.
package blog

import (
	"fmt"
	"unicode/utf8"
)

// Article is one text.
type Article struct {
	Slug  string
	Title string
	Words int
	Body  string
	Lang  string
}

// ReadingTime rounds up: an unfinished minute still counts, and exactly 400
// words must give 2 rather than 3.
func (a Article) ReadingTime() int {
	return (a.Words + 199) / 200
}

// Letters counts letters rather than bytes, which is what the runes lesson was
// for: a Kazakh letter takes two bytes and would otherwise be counted twice.
func (a Article) Letters() int {
	return utf8.RuneCountInString(a.Title)
}

func (a Article) String() string {
	return fmt.Sprintf("%s — %d мин", a.Title, a.ReadingTime())
}
