package articles

import (
	"regexp"
	"strings"
)

// Mechanical comparison of a translation against its original.
//
// The guide used to tell authors to read the translation and correct it. Then
// the obvious question: how does someone who writes in Kazakh check the English
// version? They cannot read it, and advice they cannot act on is not advice.
//
// But a great deal is checkable without knowing a word of the language. Numbers
// must survive translation unchanged. So must links, headings, and the shape of
// a table. A version half the length of its original has lost something. None
// of that requires reading — it requires counting, which is work for the
// machine, not the author.
//
// These checks find mechanical damage, not bad language. A fluent translation
// that says the wrong thing passes all of them, and the guide says so.

var (
	mdHeading   = regexp.MustCompile(`(?m)^#{2,6}\s`)
	mdLink      = regexp.MustCompile(`\]\(https?://`)
	mdTableRow  = regexp.MustCompile(`(?m)^\s*\|`)
	numberToken = regexp.MustCompile(`\d[\d  ,.\x{00a0}]*\d|\d`)
	notDigit    = regexp.MustCompile(`\D`)
)

// TranslationIssue is one discrepancy, ready for the editor to show.
type TranslationIssue struct {
	Key  string // i18n key under "tcheck."
	Have int
	Want int
}

// numbersIn returns every number in the text, reduced to its digits.
//
// Separators are deliberately discarded: a translation is expected to write
// 2 101 as 2,101 and 98,5% as 98.5% — that is correct localisation, not damage.
// What must not change is which digits are there.
func numbersIn(s string) map[string]int {
	out := map[string]int{}
	for _, m := range numberToken.FindAllString(s, -1) {
		if d := notDigit.ReplaceAllString(m, ""); d != "" {
			out[d]++
		}
	}
	return out
}

// compareTranslation lists what the translation lost or gained.
func compareTranslation(src, dst string) []TranslationIssue {
	if strings.TrimSpace(src) == "" || strings.TrimSpace(dst) == "" {
		return nil
	}
	var out []TranslationIssue
	count := func(key string, re *regexp.Regexp) {
		a, b := len(re.FindAllString(src, -1)), len(re.FindAllString(dst, -1))
		if a != b {
			out = append(out, TranslationIssue{Key: key, Have: b, Want: a})
		}
	}
	count("headings", mdHeading)
	count("links", mdLink)
	count("table_rows", mdTableRow)

	// Numbers are the check that matters most: a changed figure is a factual
	// error, and the one kind of error a reader of the translation would never
	// suspect.
	want, have := numbersIn(src), numbersIn(dst)
	missing := 0
	for n, k := range want {
		if have[n] < k {
			missing += k - have[n]
		}
	}
	if missing > 0 {
		out = append(out, TranslationIssue{Key: "numbers", Have: len(have), Want: len(want)})
	}

	// Length is the coarse net that catches a truncated or abandoned pass.
	sl, dl := len([]rune(src)), len([]rune(dst))
	if dl*10 < sl*6 || dl*10 > sl*18 {
		out = append(out, TranslationIssue{Key: "length", Have: dl, Want: sl})
	}
	return out
}
