package ai

import (
	"strings"
	"testing"
)

// The first version of this comparison knew only digits, and so shouted on every
// honest translation. An author who sees four warnings on a sound text stops
// reading warnings altogether — and the false alarm turns out to be more dangerous
// than silence, because it devalues the real one.
//
// Every case below is from production, from the measles article.
func TestCompareNumbersStaysQuietOnHonestTranslations(t *testing.T) {
	cases := []struct{ name, src, dst string }{
		{
			// "99% in 2019, 99% in 2022". The first version treated any comma between
			// digits as a thousands separator and glued this into the number 201999,
			// which was of course nowhere in the original.
			"процент и год через запятую",
			"Охват 99% в 2019 году, 99% в 2022 году и 93% в 2023 году.",
			"Coverage was 99% in 2019, 99% in 2022 and 93% in 2023.",
		},
		{
			"разделители разрядов и десятичная запятая",
			"Было 28 147 случаев при охвате 98,5%.",
			"There were 28,147 cases at 98.5% coverage.",
		},
		{
			// The same number appears four times in the original and three times in the
			// translation: that is the language building the sentence, not the meaning breaking.
			"число повторено разное число раз",
			"28 147 случаев. Из этих 28 147 — дети. Итого 28 147.",
			"28,147 cases, mostly children. That is 28,147 in total.",
		},
		{
			"цифры оригинала прописью в переводе",
			"Госдолг США перевалил за $40 триллионов.",
			"U.S. national debt crossed the forty-trillion-dollar mark.",
		},
		{
			"прописью в оригинале, цифрами в переводе",
			"Заболеваемость выросла в двадцать три раза.",
			"Аурушаңдық 23 есе өсті.",
		},
		{
			"масштабное слово вместо нулей",
			"123 тысячи случаев кори за два года.",
			"123,000 measles cases in two years.",
		},
		{
			"казахские окончания на числительных",
			"Заболеваемость — 1 367 случаев на миллион жителей.",
			"Аурушаңдық — миллион тұрғынға 1 367 жағдай.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if d := CompareNumbers(c.src, c.dst); !d.Empty() {
				t.Errorf("ложная тревога: пропало %v, выдумано %v", d.Missing, d.Invented)
			}
		})
	}
}

// What all of this is for. Kimi k2.6, asked four times to translate one sentence
// containing "сорок триллионов", answered 40, 43, 40, 43 — that is not a rare slip
// to be noted in a log, it is a tossed coin.
func TestCompareNumbersCatchesInventedFigures(t *testing.T) {
	d := CompareNumbers(
		"Государственный долг США перешёл отметку в сорок триллионов долларов.",
		"АҚШ мемлекеттік қарызы 43 триллион доллар шегінен асты.")
	if !strings.Contains(strings.Join(d.Invented, " "), "43") {
		t.Fatalf("подменённое число не найдено: %+v", d)
	}
}

// A dropped sentence takes its numbers with it, and that is the only thing an
// author who does not read the language can notice unaided.
func TestCompareNumbersCatchesLostFigures(t *testing.T) {
	d := CompareNumbers(
		"Узбекистан — 20 940 случаев. Таджикистан отчитался о нуле, находясь между "+
			"Киргизией с её 14 380 случаями и Узбекистаном с 20 940.",
		"Өзбекстан — 20 940 жағдай.")
	joined := strings.Join(d.Missing, " ")
	if !strings.Contains(joined, "14 380") {
		t.Fatalf("пропавшее число не названо: %+v", d)
	}
	if strings.Contains(joined, "20 940") {
		t.Errorf("названо число, которое на месте: %+v", d)
	}
}

// Rescaling is allowed only by whole groups of three zeros — otherwise the rule
// invented for "123 тысячи" against "123,000" would start letting real damage
// through.
func TestRescalingDoesNotHideDamage(t *testing.T) {
	d := CompareNumbers("В Узбекистане 20 940 случаев.", "Өзбекстанда 2 094 жағдай.")
	if d.Empty() {
		t.Fatal("потерянный ноль прошёл незамеченным")
	}
}

// Numerals are read as they are spoken: neighbours add up, and a scale word
// multiplies what stands before it.
func TestNumeralsCompose(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{"двадцать три", 23},
		{"сорок триллионов", 40000000000000},
		{"forty-three thousand", 43000},
		{"пятнадцать тысяч", 15000},
		{"қырық үш мың", 43000},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			if _, ok := licensedValues(c.text).values[c.want]; !ok {
				t.Errorf("%q не прочитано как %d, разобрано: %v", c.text, c.want, licensedValues(c.text).values)
			}
		})
	}
}

// Prefix matching once read the Russian "стоит" as "сто" and the Kazakh "оның"
// as "он", and the check flagged twenty-six paragraphs out of seventy-one on a
// translation with one error in it. Prefixes were kept only where the word is
// long enough for the match not to be a coincidence.
func TestOrdinaryWordsAreNotReadAsNumbers(t *testing.T) {
	for _, w := range []string{"стоит", "стоимость", "оның", "онда", "она", "простой", "семья", "стороны"} {
		if v, ok := wordValue(w); ok {
			t.Errorf("%q прочитано как число %d", w, v)
		}
	}
}

// Small numbers are legitimately rewritten in translation: "один больной" becomes
// "a single case", "две недели" becomes "a fortnight". Below ten there is nothing
// to catch and false alarms in plenty.
func TestSmallNumbersAreNotCompared(t *testing.T) {
	if d := CompareNumbers("Один больной заражает двенадцать человек.",
		"A single case infects a dozen people."); !d.Empty() {
		t.Errorf("мелкие числа дали замечания: %+v", d)
	}
}

// The list has to stay readable: a lost paragraph takes dozens of figures with it,
// and a warning that names them all is a warning nobody reads.
func TestCapListTrims(t *testing.T) {
	got := CapList([]string{"1", "2", "3", "4", "5", "6", "7"}, 5)
	if got != "1, 2, 3, 4, 5, …" {
		t.Errorf("перечень обрезан неверно: %q", got)
	}
	if got := CapList([]string{"1", "2"}, 5); got != "1, 2" {
		t.Errorf("короткий перечень тронут: %q", got)
	}
}

// Russian folds a ratio into a single adjective where the other two languages
// spell it out. Read literally, the translation looks like it invented a 100.
func TestPercentAdjectivesAreReadAsNumbers(t *testing.T) {
	pairs := [][2]string{
		{"стопроцентное резервирование", "100 percent reserves"},
		{"стопроцентное резервирование", "жүз пайыздық резервтеу"},
	}
	for _, p := range pairs {
		if d := CompareNumbers(p[0], p[1]); !d.Empty() {
			t.Errorf("CompareNumbers(%q, %q): потеряно %v, появилось %v", p[0], p[1], d.Missing, d.Invented)
		}
	}
}

// A translation links to the same outlet's own page in the target language, and
// that page has a different id in its address. Read literally, the comparison
// says one number was lost and another invented while the prose is identical.
func TestDigitsInsideLinksAreNotCompared(t *testing.T) {
	src := "Концерт перенесли — [Tengrinews](https://tengrinews.kz/kazakhstan_news/kontsert-perenesli-606228/)."
	dst := "The concert was postponed — [Tengrinews](https://en.tengrinews.kz/kazakhstan_news/concert-postponed-272924/)."
	if d := CompareNumbers(src, dst); !d.Empty() {
		t.Errorf("адреса ссылок сверены как числа: потеряно %v, появилось %v", d.Missing, d.Invented)
	}
	// A figure in the visible label is prose and must still be checked.
	if d := CompareNumbers(
		"[Tengrinews, 17 июля](https://a.kz/x-1.html): 31 000 зрителей",
		"[Tengrinews, 17 July](https://b.kz/y-2.html): 41 000 spectators"); d.Empty() {
		t.Error("подменённое число в тексте пропущено — сверка ослепла на весь абзац")
	}
}

// English is the only language here that uses a group separator and a decimal
// point in the same number, so "1744,8" comes back as "1,744.8". The parser used
// to stop at the group and read 1744 where the original said 17448, reporting
// the same figure lost and invented at once. Every English translation carrying
// such a number tripped it, which made the fault look like a quirk of one
// language rather than a rule that was one clause short.
func TestGroupedDecimalsSurviveTranslation(t *testing.T) {
	pairs := [][2]string{
		{"площадь сократилась с 1744,8 до 1032,1 км²", "the area shrank from 1,744.8 to 1,032.1 km²"},
		{"смета выросла до 22 100,5 тенге", "the estimate rose to 22,100.5 tenge"},
	}
	for _, p := range pairs {
		if d := CompareNumbers(p[0], p[1]); !d.Empty() {
			t.Errorf("CompareNumbers(%q, %q): missing %v, invented %v", p[0], p[1], d.Missing, d.Invented)
		}
	}
}

// The fractional tail must not reach across a sentence boundary: a separator
// counts only when digits follow it immediately.
func TestGroupedNumberDoesNotSwallowTheNextSentence(t *testing.T) {
	if d := CompareNumbers("смета 22 100 тенге. 15 человек ушли", "the estimate is 22,100 tenge. 15 people left"); !d.Empty() {
		t.Errorf("sentence boundary read as a decimal point: missing %v, invented %v", d.Missing, d.Invented)
	}
}
