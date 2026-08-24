package articles

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"
)

// Reader-reported typos, and the machinery that decides what to do with one.
//
// The whole design answers a single question: how do you let a stranger edit a
// published article without letting a stranger edit a published article. The
// answer is that nobody edits anything — the reader points, and the server makes
// one mechanical substitution it can describe, justify and undo.
//
// Four tiers, cheapest and most certain first. A claim has to survive all of
// them, and each one can only ever narrow what the next is allowed to do.
//
//  1. Structure. The sentence must be findable in the article, the word must be
//     inside that sentence, and it must not sit in a code block, a link target
//     or a heading. This is arithmetic on strings: free, certain, and it throws
//     out the malformed and the mischievous before anything costs money.
//
//  2. Mechanics. One large class of typos in Kazakh and Russian needs no
//     judgement at all: a Latin letter that looks like its Cyrillic twin,
//     usually from a keyboard switched mid-word. "Aлматы" with a Latin A is not
//     a matter of opinion, and the repair is a table lookup. These are settled
//     here and never reach the model.
//
//  3. Judgement. What is left is the part that genuinely needs reading —
//     agreement, case, homophones, a name spelled two ways. That goes to the
//     model, which answers with one word and nothing else. A proposed change is
//     then put back to the model as something to refute, and survives only if
//     the second reading agrees. Refusals get no second pass: nothing changes on
//     a refusal, so there is nothing to double-check.
//
//  4. Application. Bounded edit distance, a character whitelist, one occurrence
//     replaced, the whole claim written down, and the author told what changed
//     in their text.
//
// Established spell-checking was weighed against this and is not used as the
// judge. A Hunspell dictionary cannot reject a correctly-spelled word that is
// wrong in its sentence, so it cannot be trusted to refuse; and the Kazakh
// dictionaries are thin enough that using one to constrain fixes would block
// legitimate corrections in the language that needs them most. Its one sound
// contribution — edit distance as a sanity bound — is in tier 4.

// Correction statuses.
const (
	// CorrPending is recorded but undecided: the checker was unreachable.
	CorrPending = "pending"
	// CorrApplied means the word was replaced in the article.
	CorrApplied = "applied"
	// CorrRejected means the checker found nothing to fix.
	CorrRejected = "rejected"
	// CorrFailed means the claim did not survive the server's own checks.
	CorrFailed = "failed"
)

// Limits on what one reader may send. Generous for a person reading carefully,
// tight enough that an automated client cannot use the checker as a free model.
const (
	// correctionsPerDay is how many claims one account may send in a day.
	correctionsPerDay = 20
	// correctionWordMax is the longest word a claim may point at. A typo is a
	// word; anything longer is a sentence in disguise.
	correctionWordMax = 48
	// correctionSentenceMax bounds the sentence field. Long enough for a real
	// sentence in three languages, short enough not to be a paste of the piece.
	correctionSentenceMax = 600
	// correctionChapterMax bounds the chapter field.
	correctionChapterMax = 200
	// correctionMinOverlap is the share of the reader's sentence that must be
	// found around a candidate word for that word to be the one they meant.
	correctionMinOverlap = 0.6
)

// ErrCorrectionRate is returned when a reader has sent too many claims today.
var ErrCorrectionRate = errors.New("correction rate exceeded")

// Correction is one reader's claim and its outcome.
type Correction struct {
	ID        uuid.UUID
	ArticleID uuid.UUID
	Lang      string
	Reporter  uuid.UUID
	Chapter   string
	Sentence  string
	Word      string
	Status    string
	Fixed     string
	Reason    string
	CreatedAt time.Time
}

// correctionSite is where in the body a fix would land.
type correctionSite struct {
	// Start and End bound the word inside the body, in bytes.
	Start, End int
	// Score is how much of the reader's sentence was found around it.
	Score float64
	// Para is the paragraph the word sits in, for the model to read.
	Para string
}

// homoglyphs maps a Latin letter to the Cyrillic letter it is indistinguishable
// from on screen. A word that mixes the two scripts is almost always a keyboard
// that changed layout mid-word, and the repair does not need an opinion.
//
// Only the pairs that genuinely render alike are listed. Adding a pair that
// merely looks similar in one typeface would start "correcting" real Latin words.
var homoglyphs = map[rune]rune{
	'A': 'А', 'B': 'В', 'C': 'С', 'E': 'Е', 'H': 'Н', 'I': 'І', 'K': 'К',
	'M': 'М', 'O': 'О', 'P': 'Р', 'T': 'Т', 'X': 'Х', 'Y': 'У',
	'a': 'а', 'c': 'с', 'e': 'е', 'i': 'і', 'o': 'о', 'p': 'р', 'x': 'х',
	'y': 'у', 'h': 'һ',
}

// scriptOf tells the two alphabets apart. Digits, punctuation and the rest
// belong to neither and never decide a word's script.
func scriptOf(r rune) string {
	switch {
	case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
		return "latin"
	case unicode.Is(unicode.Cyrillic, r):
		return "cyrillic"
	}
	return ""
}

// repairHomoglyphs turns a mixed-script word into a single-script one, when that
// is what it plainly is. The second result is false when the word is not of that
// kind — pure Latin, pure Cyrillic, or a genuine mixture no table can settle.
func repairHomoglyphs(word string) (string, bool) {
	var latin, cyr int
	for _, r := range word {
		switch scriptOf(r) {
		case "latin":
			latin++
		case "cyrillic":
			cyr++
		}
	}
	// Both scripts must be present, and Cyrillic must be the majority: that is
	// what makes this a Cyrillic word with intruders rather than a Latin word.
	if latin == 0 || cyr == 0 || latin >= cyr {
		return "", false
	}
	var b strings.Builder
	for _, r := range word {
		if scriptOf(r) == "latin" {
			c, ok := homoglyphs[r]
			if !ok {
				// A Latin letter with no Cyrillic twin means this is not a
				// look-alike mix-up, and guessing would corrupt the word.
				return "", false
			}
			b.WriteRune(c)
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if out == word {
		return "", false
	}
	return out, true
}

// normalizeSpace collapses runs of whitespace and normalizes the text to NFC, so
// a sentence copied out of the rendered page compares equal to the same sentence
// in the source. Kazakh in particular arrives both composed and decomposed
// depending on the keyboard, and the two are byte-different and eye-identical.
func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}

// foldWord reduces a word to what two spellings of it have in common: NFC,
// lower case, no surrounding punctuation.
func foldWord(s string) string {
	s = strings.TrimFunc(norm.NFC.String(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.ToLower(s)
}

// wordSet splits text into the folded words it contains.
func wordSet(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.FieldsFunc(norm.NFC.String(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if w = strings.ToLower(w); w != "" {
			out[w] = true
		}
	}
	return out
}

// findCorrectionSite decides which occurrence of the word the reader meant.
//
// The sentence they pasted came out of the rendered page, so it will not match
// the Markdown source byte for byte: emphasis, links and non-breaking spaces all
// stand between the two. So instead of matching the sentence, every occurrence of
// the word is scored against the sentence that actually contains it, and the best
// one wins — provided it wins clearly enough.
//
// Scoring is recall-weighted. A reader who pastes the whole sentence and a reader
// who pastes the half of it they were looking at should both be understood, so
// having found all of their words matters more than having found only theirs.
//
// The one case scoring cannot settle is a word that occurs exactly once: there the
// word alone names the place, and the sentence has nothing left to disambiguate.
func findCorrectionSite(body, sentence, word string) (correctionSite, bool) {
	target := foldWord(word)
	if target == "" {
		return correctionSite{}, false
	}
	want := wordSet(sentence)
	if len(want) == 0 {
		return correctionSite{}, false
	}

	var open []span
	for _, at := range wordOccurrences(body, target) {
		if inCodeSpan(body, at.start) || inLinkTarget(body, at.start) || inHeading(body, at.start) {
			continue
		}
		open = append(open, at)
	}
	if len(open) == 0 {
		return correctionSite{}, false
	}

	best, found := correctionSite{Score: -1}, false
	for _, at := range open {
		got := wordSet(stripMD(sentenceAround(body, at.start, at.end)))
		score := recallWeightedF(want, got)
		if score > best.Score {
			best = correctionSite{
				Start: at.start, End: at.end, Score: score,
				Para: paragraphAt(body, at.start),
			}
			found = true
		}
	}
	if !found {
		return correctionSite{}, false
	}
	// A single unambiguous occurrence is accepted on the word alone: the reader
	// may have paraphrased, and there is nowhere else the fix could land.
	if best.Score < correctionMinOverlap && len(open) != 1 {
		return correctionSite{}, false
	}
	return best, true
}

// recallWeightedF scores how well a candidate sentence answers the reader's.
//
// It is the F-measure with recall weighted four to one. Precision still counts,
// or a long paragraph would beat the right sentence simply by containing more
// words; recall counts more, because a reader who quotes half a sentence has
// still identified it.
func recallWeightedF(want, got map[string]bool) float64 {
	if len(want) == 0 || len(got) == 0 {
		return 0
	}
	hit := 0
	for w := range want {
		if got[w] {
			hit++
		}
	}
	if hit == 0 {
		return 0
	}
	precision := float64(hit) / float64(len(got))
	recall := float64(hit) / float64(len(want))
	const beta2 = 4
	return (1 + beta2) * precision * recall / (beta2*precision + recall)
}

// sentenceAround returns the sentence holding a byte range: back to the previous
// terminator and on to the next one, and never across a blank line.
func sentenceAround(body string, start, end int) string {
	lo := 0
	for i := start - 1; i >= 0; i-- {
		switch body[i] {
		case '.', '!', '?', '\n':
			// A period inside a number or an abbreviation does not end a
			// sentence; a period followed by space or newline does.
			if body[i] == '\n' || i+1 >= len(body) || body[i+1] == ' ' || body[i+1] == '\n' {
				lo = i + 1
				i = -1
			}
		}
		if lo != 0 {
			break
		}
	}
	hi := len(body)
	for i := end; i < len(body); i++ {
		switch body[i] {
		case '.', '!', '?':
			if i+1 >= len(body) || body[i+1] == ' ' || body[i+1] == '\n' {
				hi = i + 1
				i = len(body)
			}
		case '\n':
			hi = i
			i = len(body)
		}
		if hi != len(body) {
			break
		}
	}
	lo, hi = clampRunes(body, lo, hi)
	return strings.TrimSpace(body[lo:hi])
}

// span is a byte range inside the body.
type span struct{ start, end int }

// wordOccurrences finds every whole-word occurrence of a folded word.
//
// Whole-word is checked by hand rather than with \b: Go's \b is ASCII-only, and
// on Cyrillic it fires between every letter and its neighbour, which would make
// "он" match inside "фон".
func wordOccurrences(body, target string) []span {
	var out []span
	runes := []rune(norm.NFC.String(body))
	// Byte offsets are what the caller needs, so walk the runes and keep the
	// offset alongside.
	offs := make([]int, len(runes)+1)
	b := 0
	for i, r := range runes {
		offs[i] = b
		b += len(string(r))
	}
	offs[len(runes)] = b

	i := 0
	for i < len(runes) {
		if !isWordRune(runes[i]) {
			i++
			continue
		}
		j := i
		for j < len(runes) && isWordRune(runes[j]) {
			j++
		}
		if strings.ToLower(string(runes[i:j])) == target {
			out = append(out, span{start: offs[i], end: offs[j]})
		}
		i = j
	}
	return out
}

// isWordRune says what counts as part of a word. The hyphen is deliberately not
// one: "жизненно-важный" is two words for this purpose, and a reader pointing at
// one half of it means that half.
func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\''
}

// clampRunes widens a byte range to the nearest rune boundaries inside the body.
func clampRunes(body string, lo, hi int) (int, int) {
	if lo < 0 {
		lo = 0
	}
	if hi > len(body) {
		hi = len(body)
	}
	for lo > 0 && lo < len(body) && !utf8Start(body[lo]) {
		lo--
	}
	for hi < len(body) && !utf8Start(body[hi]) {
		hi++
	}
	return lo, hi
}

// utf8Start reports whether a byte begins a rune.
func utf8Start(b byte) bool { return b&0xC0 != 0x80 }

// paragraphAt returns the paragraph containing an offset, for the model to read
// the word in context.
func paragraphAt(body string, at int) string {
	start := strings.LastIndex(body[:at], "\n\n")
	if start < 0 {
		start = 0
	} else {
		start += 2
	}
	end := strings.Index(body[at:], "\n\n")
	if end < 0 {
		end = len(body)
	} else {
		end += at
	}
	return normalizeSpace(stripMD(body[start:end]))
}

// inCodeSpan reports whether an offset sits inside code — a fenced block or an
// inline span. Correcting the spelling of an identifier breaks it.
func inCodeSpan(body string, at int) bool {
	// Fences: count the ones that open a line before the offset. An odd count
	// means the offset is inside a block that has not closed yet.
	fences := 0
	for i, line := range strings.Split(body[:at], "\n") {
		if strings.HasPrefix(strings.TrimLeft(line, " \t"), "```") {
			fences++
		}
		_ = i
	}
	if fences%2 == 1 {
		return true
	}
	// Inline: an odd number of backticks on the same line before the offset.
	lineStart := strings.LastIndexByte(body[:at], '\n') + 1
	return strings.Count(body[lineStart:at], "`")%2 == 1
}

// inLinkTarget reports whether an offset sits inside a link's URL. Correcting a
// URL's spelling breaks the link, silently.
func inLinkTarget(body string, at int) bool {
	open := strings.LastIndex(body[:at], "](")
	if open < 0 {
		return false
	}
	close := strings.IndexByte(body[open:], ')')
	return close < 0 || open+close > at
}

// inHeading reports whether an offset sits in a Markdown heading. Headings are
// anchors: the table of contents and every link into the article are built from
// them, and editing one moves a target other pages point at.
func inHeading(body string, at int) bool {
	lineStart := strings.LastIndexByte(body[:at], '\n') + 1
	return strings.HasPrefix(strings.TrimLeft(body[lineStart:], " \t"), "#")
}

// safeFix is the last gate before an author's text changes. It knows nothing
// about language and everything about damage: whatever the model returned, this
// decides whether it is small, plain and word-shaped enough to be a typo fix.
func safeFix(word, fixed string) error {
	word, fixed = norm.NFC.String(word), norm.NFC.String(fixed)
	if fixed == "" {
		return errors.New("fix is empty")
	}
	if fixed == word {
		return errors.New("fix is the word itself")
	}
	rf := []rune(fixed)
	rw := []rune(word)
	if len(rf) > 3*len(rw)+8 {
		return errors.New("fix is far longer than the word")
	}
	// A whitelist, not a blacklist: everything that could carry markup, break a
	// line or open a link is absent from it, and so is anything we did not think
	// of. Letters, digits, one hyphen's worth of punctuation, and at most one
	// space — a missing space between two words is a real typo.
	spaces := 0
	for _, r := range rf {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
		case r == '-' || r == '\'' || r == '’' || r == '.' || r == ',':
		case r == ' ':
			spaces++
		default:
			return fmt.Errorf("fix contains %q", r)
		}
	}
	if spaces > 1 {
		return errors.New("fix is more than two words")
	}
	// Edit distance is the one thing conventional spell-checking contributes
	// here, and it contributes it well: a typo fix is a small edit, and a model
	// that has decided to rewrite the word instead of correcting it fails this.
	limit := 3
	if l := len(rw) / 2; l > limit {
		limit = l
	}
	if d := levenshtein(strings.ToLower(word), strings.ToLower(fixed)); d > limit {
		return fmt.Errorf("fix is %d edits away, limit %d", d, limit)
	}
	return nil
}

// levenshtein is the plain two-row edit distance over runes.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	prev := make([]int, len(rb)+1)
	cur := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min3(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(rb)]
}

func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

// matchCase gives the fix the capitalisation the word had. A model told to keep
// it usually does; this makes it certain, because a lower-cased word at the start
// of a sentence is a new typo introduced while fixing one.
func matchCase(word, fixed string) string {
	rw, rf := []rune(word), []rune(fixed)
	if len(rw) == 0 || len(rf) == 0 {
		return fixed
	}
	allUpper := true
	for _, r := range rw {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			allUpper = false
			break
		}
	}
	switch {
	case allUpper && len(rw) > 1:
		return strings.ToUpper(fixed)
	case unicode.IsUpper(rw[0]):
		rf[0] = unicode.ToUpper(rf[0])
		return string(rf)
	case unicode.IsLower(rw[0]):
		rf[0] = unicode.ToLower(rf[0])
		return string(rf)
	}
	return fixed
}

// applyAt replaces the byte range with the fix. One occurrence, named by offset,
// so a word that appears elsewhere in the article is untouched.
func applyAt(body string, s correctionSite, fixed string) string {
	return body[:s.Start] + fixed + body[s.End:]
}
