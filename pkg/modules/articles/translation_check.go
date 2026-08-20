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
	mdHeading  = regexp.MustCompile(`(?m)^#{2,6}\s`)
	mdLink     = regexp.MustCompile(`\]\(https?://`)
	mdTableRow = regexp.MustCompile(`(?m)^\s*\|`)
	// A separator only joins digits when it groups exactly three of them.
	// The first version treated any comma or space between digits as a
	// thousands mark, so "99% in 2019, 99% in 2022" became the single number
	// 201999 — and every English translation was reported as damaged.
	numberToken = regexp.MustCompile(`\d{1,3}(?:[\s,.\x{00a0} ]\d{3})+|\d+(?:[.,]\d+)?`)
	notDigit    = regexp.MustCompile(`\D`)
)

// TranslationIssue is one discrepancy, ready for the editor to show.
type TranslationIssue struct {
	Key    string // i18n key under "tcheck."
	Field  string // i18n key under "tfield."
	Have   int
	Want   int
	Detail string // the numbers themselves, when that is what is wrong
}

// numbersIn returns the numbers of a text in the order they appear, each as the
// text wrote it, along with the set of their digits.
//
// Separators are deliberately discarded from the set: a translation is expected
// to write 2 101 as 2,101 and 98,5% as 98.5% — that is correct localisation,
// not damage. What must not change is which digits are there.
//
// The written form is kept because that is what the author will search for. A
// warning naming "14380" sends them hunting through a text that says "14 380".
// The order is kept for the same reason: the author reads top to bottom, and a
// list that follows the text is a list they can walk.
//
// Numbers shorter than min digits are left out. Going missing and appearing
// from nowhere are not equally suspicious, so the two directions use different
// floors — see compareTranslation.
func numbersIn(s string, min int) ([]string, map[string]string) {
	var order []string
	digits := map[string]string{}
	for _, m := range numberToken.FindAllString(s, -1) {
		d := notDigit.ReplaceAllString(m, "")
		if len(d) < min {
			continue
		}
		if _, seen := digits[d]; !seen {
			digits[d] = m
			order = append(order, d)
		}
	}
	return order, digits
}

// survives reports whether a source number can be found in the translation,
// allowing for the one rewrite that changes the digits legitimately: a scale
// word traded for zeros. "123 тысячи" and "123,000" are the same number, and so
// are "40 миллионов" and "40,000,000".
//
// Only whole groups of three zeros count, which is what keeps this from
// swallowing real damage: 20 940 mistyped as 2 094 differs by one zero, not by
// a scale word, and is still reported.
func survives(digits string, have map[string]string) bool {
	for _, zeros := range []string{"", "000", "000000", "000000000", "000000000000"} {
		if _, ok := have[digits+zeros]; ok {
			return true
		}
		if trimmed := strings.TrimSuffix(digits, zeros); zeros != "" && trimmed != digits && len(trimmed) >= 3 {
			if _, ok := have[trimmed]; ok {
				return true
			}
		}
	}
	return false
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
	// Only numbers absent from the translation altogether are reported. A figure
	// used four times in Russian and three times in Kazakh is a difference of
	// phrasing, not of fact, and flagging it taught the author to ignore the
	// warnings — which is worse than not checking.
	//
	// And the numbers are named. "Some numbers do not match (68 / 68)" tells
	// nobody anything; "20 940, 14 380 are missing" is checked in five seconds
	// by someone who does not read the language.
	// Going missing is checked from three digits up. Shorter numbers are the
	// ones a good translator rewrites — "$40 триллионов" becomes "the
	// forty-trillion-dollar mark", "один больной" becomes "a single case" — and
	// a rule that cannot tell that apart from damage fires on every honest
	// translation.
	srcOrder, want := numbersIn(src, 3)
	_, have := numbersIn(dst, 3)
	var missing []string
	for _, digits := range srcOrder {
		if !survives(digits, have) {
			missing = append(missing, want[digits])
		}
	}
	if len(missing) > 0 {
		// The list follows the text, and stops while it is still readable: a
		// dropped paragraph takes dozens of figures with it, and a warning
		// naming all of them is a warning nobody reads.
		if len(missing) > 5 {
			missing = append(missing[:5:5], "…")
		}
		out = append(out, TranslationIssue{Key: "numbers", Detail: strings.Join(missing, ", ")})
	}

	// And the other direction: a number the translation states that the original
	// never did. This is the graver fault of the two — a lost sentence leaves a
	// gap a reader may notice, an invented figure reads as fact and is quoted as
	// one. It happened on the article that prompted these checks: "сорок
	// триллионов" came back as "43 триллион", twice, in a text whose own source
	// list said forty.
	//
	// Checked from two digits up, where losses are checked from three: a
	// translator has good reason to spell a small number out, and none to
	// produce one the original does not contain.
	dstOrder, got := numbersIn(dst, 2)
	_, base := numbersIn(src, 2)
	var invented []string
	for _, digits := range dstOrder {
		if !survives(digits, base) {
			invented = append(invented, got[digits])
		}
	}
	if len(invented) > 0 {
		if len(invented) > 5 {
			invented = append(invented[:5:5], "…")
		}
		out = append(out, TranslationIssue{Key: "invented", Detail: strings.Join(invented, ", ")})
	}

	// Length is the coarse net that catches a truncated or abandoned pass.
	sl, dl := len([]rune(src)), len([]rune(dst))
	if dl*10 < sl*6 || dl*10 > sl*18 {
		out = append(out, TranslationIssue{Key: "length", Have: dl, Want: sl})
	}
	return out
}
