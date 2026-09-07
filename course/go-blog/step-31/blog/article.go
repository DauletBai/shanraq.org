// Package blog is everything the blog knows about a text, kept apart from the
// server that shows it. Only the names starting with a capital letter leave
// this folder.
package blog

import (
	"fmt"
	"time"
	"unicode/utf8"
)

// Article is one text.
type Article struct {
	Slug  string
	Title string
	Words int
	Body  string
	Lang  string
	// UpdatedAt is the zero time until someone edits the article, and holds
	// the moment the database wrote, not the one the server guessed.
	UpdatedAt time.Time
	// AuthorID is zero for the articles that were here before there were
	// people. Nobody owns those, so nobody may change them.
	AuthorID int64
	// Cover is a file name inside the uploads folder, empty when there is
	// none. Not a path: where uploads live is a setting.
	Cover string
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

// months, because Go has no Kazakh month names: package time knows the English
// ones only.
var months = [...]string{"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
	"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан"}

// Updated is the date as a person reads it. The formatting lives here rather
// than in the template: a template shows, it does not calculate.
func (a Article) Updated() string {
	if a.UpdatedAt.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d %s %d", a.UpdatedAt.Day(),
		months[a.UpdatedAt.Month()-1], a.UpdatedAt.Year())
}
