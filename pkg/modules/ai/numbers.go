package ai

import (
	"regexp"
	"strconv"
	"strings"
)

// Comparing the numbers of a text with the numbers of its translation.
//
// This lives here, not in the articles module, because two callers need the
// same answer: the translator, which retries a paragraph whose figures came
// back wrong, and the editor, which shows the author what to look at. One
// implementation, so the warning an author sees is the rule the machine
// enforced.
//
// Numbers are compared as values, never as spellings. A translation is expected
// to write 2 101 as 2,101 and 98,5% as 98.5%; it is also entitled to write
// "сорок триллионов" as "forty trillion" or "двадцать три раза" as "23 есе".
// The first version of this check knew only digits, so every article where the
// author spelled a number out was reported as damaged — and the noise buried a
// real fault: "сорок триллионов" came back as "43 триллион", in Kazakh and in
// English both.

// numberWord maps a spelled-out numeral to its value.
//
// Whole forms, not prefixes. The first attempt matched on stems, reasoning that
// over-matching was the safe direction — and it was wrong in both directions at
// once. A stem that claims too much invents values on the translation side
// (reported as fabricated) and demands them on the original side (reported as
// lost). Russian "стоит" read as "сто", Kazakh "оның" read as "он", and the
// check flagged twenty-six paragraphs out of seventy-one on a translation whose
// only numeric fault was in one of them.
//
// So the forms are listed. It is a long table and it will miss an inflection
// now and then; missing one costs nothing, because a numeral the parser does
// not recognise simply is not compared.
type numberWord struct {
	word  string
	value int64
}

var numberWords = []numberWord{
	// Russian, nominative and the oblique forms that actually turn up in prose.
	{"ноль", 0}, {"нуля", 0},
	{"один", 1}, {"одна", 1}, {"одно", 1}, {"одного", 1}, {"одну", 1}, {"одним", 1},
	{"два", 2}, {"две", 2}, {"двух", 2}, {"двум", 2}, {"двумя", 2},
	{"три", 3}, {"трёх", 3}, {"трех", 3}, {"трём", 3}, {"трем", 3}, {"тремя", 3},
	{"четыре", 4}, {"четырёх", 4}, {"четырех", 4}, {"четырьмя", 4},
	{"пять", 5}, {"пяти", 5}, {"пятью", 5},
	{"шесть", 6}, {"шести", 6}, {"шестью", 6},
	{"семь", 7}, {"семи", 7}, {"семью", 7},
	{"восемь", 8}, {"восьми", 8},
	{"девять", 9}, {"девяти", 9},
	{"десять", 10}, {"десяти", 10},
	{"одиннадцать", 11}, {"одиннадцати", 11},
	{"двенадцать", 12}, {"двенадцати", 12},
	{"тринадцать", 13}, {"тринадцати", 13},
	{"четырнадцать", 14}, {"четырнадцати", 14},
	{"пятнадцать", 15}, {"пятнадцати", 15},
	{"шестнадцать", 16}, {"шестнадцати", 16},
	{"семнадцать", 17}, {"семнадцати", 17},
	{"восемнадцать", 18}, {"восемнадцати", 18},
	{"девятнадцать", 19}, {"девятнадцати", 19},
	{"двадцать", 20}, {"двадцати", 20},
	{"тридцать", 30}, {"тридцати", 30},
	{"сорок", 40}, {"сорока", 40},
	{"пятьдесят", 50}, {"пятидесяти", 50},
	{"шестьдесят", 60}, {"шестидесяти", 60},
	{"семьдесят", 70}, {"семидесяти", 70},
	{"восемьдесят", 80}, {"восьмидесяти", 80},
	{"девяносто", 90}, {"девяноста", 90},
	{"сто", 100}, {"ста", 100},
	// A ratio Russian folds into one adjective while English and Kazakh spell it
	// out: "стопроцентное резервирование" against "100 percent reserves". Without
	// these the check reports the 100 as invented on every such pair.
	{"стопроцентный", 100}, {"стопроцентное", 100}, {"стопроцентная", 100},
	{"стопроцентного", 100}, {"стопроцентной", 100}, {"стопроцентных", 100},
	{"стопроцентном", 100}, {"стопроцентную", 100},
	{"двести", 200}, {"двухсот", 200},
	{"триста", 300}, {"трёхсот", 300}, {"трехсот", 300},
	{"четыреста", 400}, {"четырёхсот", 400},
	{"пятьсот", 500}, {"пятисот", 500},
	{"шестьсот", 600}, {"шестисот", 600},
	{"семьсот", 700}, {"семисот", 700},
	{"восемьсот", 800}, {"восьмисот", 800},
	{"девятьсот", 900}, {"девятисот", 900},
	{"тысяча", 1000}, {"тысячи", 1000}, {"тысяч", 1000}, {"тысячу", 1000}, {"тысячам", 1000},
	{"миллион", 1000000}, {"миллиона", 1000000}, {"миллионов", 1000000}, {"миллионам", 1000000},
	{"миллиард", 1000000000}, {"миллиарда", 1000000000}, {"миллиардов", 1000000000},
	{"триллион", 1000000000000}, {"триллиона", 1000000000000}, {"триллионов", 1000000000000},

	// Kazakh. "Он" (ten) is left out on purpose: it is also the Russian pronoun
	// and the opening of half the Kazakh words that follow one, and it cost more
	// than it was worth.
	{"нөл", 0}, {"бір", 1}, {"екі", 2}, {"үш", 3}, {"төрт", 4}, {"бес", 5},
	{"алты", 6}, {"жеті", 7}, {"сегіз", 8}, {"тоғыз", 9},
	{"жиырма", 20}, {"отыз", 30}, {"қырық", 40}, {"елу", 50},
	{"алпыс", 60}, {"жетпіс", 70}, {"сексен", 80}, {"тоқсан", 90},
	{"жүз", 100}, {"жүзге", 100}, {"жүздеген", 100},
	{"мың", 1000}, {"мыңға", 1000}, {"мыңы", 1000}, {"мыңнан", 1000},
	{"мыңдаған", 1000}, {"мыңдық", 1000},

	// English.
	{"zero", 0}, {"one", 1}, {"two", 2}, {"three", 3}, {"four", 4}, {"five", 5},
	{"six", 6}, {"seven", 7}, {"eight", 8}, {"nine", 9}, {"ten", 10},
	{"eleven", 11}, {"twelve", 12}, {"thirteen", 13}, {"fourteen", 14},
	{"fifteen", 15}, {"sixteen", 16}, {"seventeen", 17}, {"eighteen", 18},
	{"nineteen", 19}, {"twenty", 20}, {"thirty", 30}, {"forty", 40},
	{"fifty", 50}, {"sixty", 60}, {"seventy", 70}, {"eighty", 80}, {"ninety", 90},
	{"hundred", 100}, {"hundreds", 100},
	{"thousand", 1000}, {"thousands", 1000},
	{"million", 1000000}, {"millions", 1000000},
	{"billion", 1000000000}, {"billions", 1000000000},
	{"trillion", 1000000000000}, {"trillions", 1000000000000},
}

// wordValues indexes the table for exact lookup.
var wordValues = func() map[string]int64 {
	m := make(map[string]int64, len(numberWords))
	for _, w := range numberWords {
		m[w.word] = w.value
	}
	return m
}()

// numberToken matches a number written in digits. A separator only joins digits
// when it groups exactly three of them: without that rule "99% in 2019, 99% in
// 2022" reads as the single number 201999, and every English translation is
// reported as damaged.
//
// The grouped form carries an optional fractional tail, and it has to. English
// is the one language here that uses both separators in the same number, so a
// glacier area written "1744,8" in Russian comes back as "1,744.8". Without the
// tail the match stopped at the group and yielded 1744 against the original's
// 17448 -- the figure was reported lost and invented at once, on a translation
// that was in fact correct. Every English article carrying a grouped decimal hit
// this, which is why the false alarms looked like a quirk of one language.
//
// The tail cannot swallow a sentence boundary: the separator must be followed
// immediately by digits, so "cost 22,100. 5 people came" still parses as two
// numbers.
var numberToken = regexp.MustCompile(`\d{1,3}(?:[\s,.\x{00a0} ]\d{3})+(?:[.,]\d+)?|\d+(?:[.,]\d+)?`)

var notDigit = regexp.MustCompile(`\D`)

// numberValueFloor is the value below which numbers are not compared.
//
// Small numbers are the ones a translation legitimately rewrites out of
// existence: "один больной" becomes "a single case", "две недели" becomes
// "a fortnight". Nothing below ten is worth the false alarms, and no fabricated
// statistic hides there either.
const numberValueFloor = 10

// numbersFound is what a text says numerically: the set of values it states,
// and for each, the way it wrote it — which is what an author will search for.
// A warning naming "14380" sends them hunting through a text that says "14 380".
type numbersFound struct {
	order  []int64 // order of first appearance, for a list that follows the text
	values map[int64]string
}

func (n *numbersFound) add(v int64, written string) {
	if _, seen := n.values[v]; seen {
		return
	}
	n.values[v] = written
	n.order = append(n.order, v)
}

// wordValue reads a token as a spelled-out numeral, or reports that it is not one.
//
// An exact form is tried first, then the table as prefixes, because Kazakh
// glues its endings on — "миллионға", "мыңдаған" — and Russian declines. Loose
// matching is safe here and only here: a word never puts a number on the list
// of figures that must be found, it only permits one. Reading "стоит" as a
// hundred costs a hundred that would have been reported; reading it as nothing
// would cost a false alarm, which is worse.
func wordValue(tok string) (int64, bool) {
	if v, ok := wordValues[tok]; ok {
		return v, true
	}
	for _, w := range numberWords {
		stem := []rune(w.word)
		if len(stem) >= 5 && strings.HasPrefix(tok, w.word) && len([]rune(tok))-len(stem) <= 5 {
			return w.value, true
		}
	}
	return 0, false
}

// digitNumbers lists the numbers a text writes in digits, in the order they
// appear, each with the way the text wrote it — which is what an author will
// search for. A warning naming "14380" sends them hunting through a text that
// says "14 380".
//
// Only these are ever demanded of a translation. A number the original spells
// out is not required to come back in digits, nor the reverse: that choice
// belongs to the translator.
func digitNumbers(text string) numbersFound {
	out := numbersFound{values: map[int64]string{}}
	for _, m := range numberToken.FindAllString(text, -1) {
		d := notDigit.ReplaceAllString(m, "")
		if v, err := strconv.ParseInt(d, 10, 64); err == nil && v >= numberValueFloor {
			out.add(v, m)
		}
	}
	return out
}

// licensedValues is every value a text can be said to state, however it wrote
// it: in digits, in words, and in the combinations words make. Spelled-out
// numerals are composed the way they are spoken — adjacent numerals add
// ("двадцать три" is 23) and a scale word multiplies what came before ("сорок
// триллионов" is forty trillion) — and each part is kept alongside the whole,
// because a translation may render the same quantity at either grain.
//
// This set answers one question only: may the other text contain this figure?
// So it is built generously. Every value it holds is a warning not raised.
func licensedValues(text string) numbersFound {
	out := digitNumbers(text)
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !isLetter(r)
	})
	var run []int64
	var spelled []string
	flush := func() {
		for _, v := range composeNumerals(run) {
			if v >= numberValueFloor {
				out.add(v, strings.Join(spelled, " "))
			}
		}
		run, spelled = run[:0], spelled[:0]
	}
	for _, tok := range words {
		if v, ok := wordValue(tok); ok {
			run = append(run, v)
			spelled = append(spelled, tok)
			continue
		}
		flush()
	}
	flush()
	return out
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'ё' || r == 'Ё' ||
		(r >= 0x0400 && r <= 0x04FF) // the rest of Cyrillic, Kazakh letters included
}

// composeNumerals turns a run of spelled-out numerals into every value it could
// reasonably be stating. Generosity here is deliberate: an extra accepted value
// weakens the check slightly, while a missing one flags an honest translation.
func composeNumerals(run []int64) []int64 {
	if len(run) == 0 {
		return nil
	}
	seen := map[int64]bool{}
	var out []int64
	push := func(v int64) {
		if v > 0 && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}

	var total, current int64
	for _, v := range run {
		push(v)
		switch {
		case v >= 1000: // a scale word multiplies what stands before it
			if current == 0 {
				current = 1
			}
			current *= v
			total += current
			push(current)
			current = 0
		default:
			current += v
			push(current)
		}
	}
	push(total + current)
	return out
}

// formatValue writes a value the way the article's own text does — grouped in
// threes with a space, which is how Russian and Kazakh print thousands.
func formatValue(v int64) string {
	s := strconv.FormatInt(v, 10)
	if len(s) <= 4 {
		return s
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// matches reports whether a value is stated by a text, allowing the one rewrite
// that legitimately changes it: a scale word traded for zeros, so that "123
// тысячи" answers "123,000". Only whole groups of three zeros count, which is
// what keeps 20 940 mistyped as 2 094 from passing as a rescaling.
func matches(v int64, in numbersFound) bool {
	if _, ok := in.values[v]; ok {
		return true
	}
	for _, scale := range []int64{1000, 1000000, 1000000000, 1000000000000} {
		if _, ok := in.values[v*scale]; ok {
			return true
		}
		if v%scale == 0 {
			if _, ok := in.values[v/scale]; ok {
				return true
			}
		}
	}
	return false
}

// linkTarget matches the address half of a Markdown link and a bare URL. The
// visible label is left alone: it is prose, and "Tengrinews, 17 July" carries a
// date a translation must keep.
var linkTarget = regexp.MustCompile(`\]\([^)]*\)|https?://\S+`)

// stripLinkTargets removes URLs before the figures are read.
//
// A digit inside an address is not a claim about the world; it is part of the
// address. An article translated properly links to the same outlet's own
// English page, and that page has a different id -- so comparing the digits
// reported one number lost and another invented on a translation where nothing
// at all had happened to the text.
func stripLinkTargets(s string) string {
	if !strings.Contains(s, "](") && !strings.Contains(s, "http") {
		return s
	}
	return linkTarget.ReplaceAllString(s, " ")
}

// NumberDiff is what a translation did to the figures of its original.
type NumberDiff struct {
	Missing  []string // stated by the original, absent from the translation
	Invented []string // stated by the translation, absent from the original
}

// Empty reports whether the figures came through untouched.
func (d NumberDiff) Empty() bool { return len(d.Missing) == 0 && len(d.Invented) == 0 }

// CompareNumbers compares the figures of a text and its translation.
//
// Each direction asks a different question, and they are not symmetrical. What
// must be found is only what the original wrote in digits; what may appear is
// anything the original states by any means. That asymmetry is the whole trick:
// it lets "двадцать три раза" come back as "23 есе" without complaint, while
// "сорок триллионов" returning as "43 триллион" is caught.
//
// The second direction matters more. A dropped sentence leaves a gap a reader
// may notice; a number that was never in the original reads as fact and is
// quoted as one. Kimi k2.6, asked four times to translate one sentence
// containing "сорок триллионов", answered 40, 43, 40, 43 — so this is not a
// rare accident to be logged, it is a coin toss to be caught.
func CompareNumbers(src, dst string) NumberDiff {
	var d NumberDiff
	if strings.TrimSpace(src) == "" || strings.TrimSpace(dst) == "" {
		return d
	}
	src, dst = stripLinkTargets(src), stripLinkTargets(dst)
	srcDigits, srcAll := digitNumbers(src), licensedValues(src)
	dstDigits, dstAll := digitNumbers(dst), licensedValues(dst)

	for _, v := range srcDigits.order {
		if !matches(v, dstAll) {
			d.Missing = append(d.Missing, srcDigits.values[v])
		}
	}
	for _, v := range dstDigits.order {
		if !matches(v, srcAll) {
			d.Invented = append(d.Invented, dstDigits.values[v])
		}
	}
	return d
}

// CapList trims a list of figures to something a person will actually read. A
// dropped paragraph takes dozens of numbers with it, and a warning naming all
// of them is a warning nobody reads.
func CapList(items []string, max int) string {
	if len(items) > max {
		items = append(items[:max:max], "…")
	}
	return strings.Join(items, ", ")
}
