package articles

import (
	"regexp"
	"strings"

	"shanraq.org/pkg/modules/ai"
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
	blockSplit = regexp.MustCompile(`\n\s*\n`)
)

// TranslationIssue is one discrepancy, ready for the editor to show.
type TranslationIssue struct {
	Key    string // i18n key under "tcheck."
	Field  string // i18n key under "tfield."
	Have   int
	Want   int
	Detail string // the numbers themselves, when that is what is wrong
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
	// suspect. The comparison itself lives in the ai module, because the
	// translator enforces the same rule when it decides whether to ask for a
	// paragraph again — the warning shown here and the rule applied there must
	// be the same rule.
	//
	// Paragraph by paragraph whenever the two texts have the same shape, which
	// since translation runs a paragraph at a time is nearly always. Comparing
	// whole texts hides exactly the fault that prompted this: the article that
	// came back saying "43 триллион" where the original said forty passed a
	// whole-text check clean, because 43 appeared elsewhere in the piece — in
	// "сорок три тысячи заболевших", four paragraphs down.
	diff := comparePiecewise(src, dst)
	if len(diff.Missing) > 0 {
		out = append(out, TranslationIssue{Key: "numbers", Detail: ai.CapList(diff.Missing, 5)})
	}
	if len(diff.Invented) > 0 {
		out = append(out, TranslationIssue{Key: "invented", Detail: ai.CapList(diff.Invented, 5)})
	}

	// Length is the coarse net that catches a truncated or abandoned pass.
	sl, dl := len([]rune(src)), len([]rune(dst))
	if dl*10 < sl*6 || dl*10 > sl*18 {
		out = append(out, TranslationIssue{Key: "length", Have: dl, Want: sl})
	}
	return out
}

// comparePiecewise compares the figures paragraph against paragraph when the
// two texts still line up, and as whole texts when they do not.
//
// Lining up is what makes the check precise: a figure is then judged against
// the sentence it belongs to rather than against everything the article
// happens to mention. When the shapes differ the translation has already lost
// or gained a paragraph, and pairing them by position would report every
// paragraph after the first difference — so it falls back, and the length
// check is what speaks up.
func comparePiecewise(src, dst string) ai.NumberDiff {
	a := blockSplit.Split(strings.TrimSpace(src), -1)
	b := blockSplit.Split(strings.TrimSpace(dst), -1)
	if len(a) < 2 || !sameShape(a, b) {
		return ai.CompareNumbers(src, dst)
	}
	var out ai.NumberDiff
	touched := 0
	for i := range a {
		d := ai.CompareNumbers(a[i], b[i])
		if !d.Empty() {
			touched++
		}
		out.Missing = append(out.Missing, d.Missing...)
		out.Invented = append(out.Invented, d.Invented...)
	}
	// Matching kinds is a good test of alignment, not a perfect one: a text of
	// nothing but paragraphs keeps its shape however far it has shifted. So the
	// count of objections is the second guard. Real damage is a figure or two;
	// when a third of the paragraphs disagree, what is wrong is the pairing.
	if touched >= 3 && touched*3 > len(a) {
		return ai.CompareNumbers(src, dst)
	}
	return out
}

// sameShape reports whether two texts are built of the same run of blocks, so
// that the nth paragraph of one answers the nth paragraph of the other. It is
// the first of two alignment guards; the second is in comparePiecewise.
//
// An equal count is not enough. A translation that dropped one sentence and
// split another adds up to the same total while everything after the first
// change sits against the wrong original — and the article that prompted these
// checks did exactly that, so pairing by position reported half the piece.
// Comparing the kinds catches it: heading, table row and paragraph fall in a
// particular order, and a shifted text stops matching that order almost at
// once.
func sameShape(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if blockKind(a[i]) != blockKind(b[i]) {
			return false
		}
	}
	return true
}

func blockKind(b string) byte {
	switch t := strings.TrimSpace(b); {
	case t == "":
		return ' '
	case strings.HasPrefix(t, "#"):
		return 'H'
	case strings.HasPrefix(t, "|"):
		return 'T'
	case strings.HasPrefix(t, ">"):
		return 'Q'
	case strings.HasPrefix(t, "-"), strings.HasPrefix(t, "*"):
		return 'L'
	default:
		return 'P'
	}
}
