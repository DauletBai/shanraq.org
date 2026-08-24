package articles

import (
	"strings"
	"testing"
)

// The question that prompted these checks: how does an author who writes in
// Kazakh verify the English version? They cannot read it — but they can be told
// that a number went missing, and that is work for the machine.
func TestTranslationChecksCatchMechanicalDamage(t *testing.T) {
	src := "## Заголовок\n\nВ 2024 году было 28 147 случаев, охват 98,5%.\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\nСм. [ВОЗ](https://who.int/x)."

	t.Run("честный перевод не даёт замечаний", func(t *testing.T) {
		// Thousands separators and the decimal comma change with localisation —
		// that is a correct translation, not a breakage.
		dst := "## Title\n\nIn 2024 there were 28,147 cases, coverage 98.5%.\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\nSee [WHO](https://who.int/x)."
		if got := compareTranslation(src, dst); len(got) != 0 {
			t.Errorf("на корректном переводе найдены замечания: %+v", got)
		}
	})

	cases := []struct {
		name, dst, wantKey string
	}{
		{"потерянная ссылка", strings.Replace(src, "[ВОЗ](https://who.int/x)", "ВОЗ", 1), "links"},
		{"пропавший заголовок", strings.Replace(src, "## Заголовок", "Заголовок", 1), "headings"},
		{"съеденная строка таблицы", strings.Replace(src, "| 1 | 2 |\n", "", 1), "table_rows"},
		{"искажённое число", strings.Replace(src, "28 147", "28 000", 1), "numbers"},
		{"обрезанный текст", "## Title", "length"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := compareTranslation(src, c.dst)
			for _, i := range got {
				if i.Key == c.wantKey {
					return
				}
			}
			t.Errorf("замечание %q не найдено, получено: %+v", c.wantKey, got)
		})
	}

	// Empty fields are no grounds for complaint: there is nothing to translate yet.
	if got := compareTranslation("", "что-то"); got != nil {
		t.Errorf("пустой оригинал дал замечания: %+v", got)
	}
	if got := compareTranslation(src, ""); got != nil {
		t.Errorf("пустой перевод дал замечания: %+v", got)
	}
}

// The first version of this check shouted on every translation — and in all
// three languages at once. An author who sees four warnings on a sound text
// stops reading warnings altogether; a false alarm is more dangerous than
// silence, because it devalues the real one too.
//
// Every case below is taken from production, from the measles article.
func TestTranslationChecksStayQuietOnHonestTranslations(t *testing.T) {
	cases := []struct{ name, src, dst string }{
		{
			// The text: "99% in 2019, 99% in 2022". The first version treated any comma
			// between digits as a thousands separator and glued this into the number
			// 201999 — which was, of course, nowhere in the original.
			"процент и год через запятую",
			"Охват 99% в 2019 году, 99% в 2022 году и 93% в 2023 году.",
			"Coverage was 99% in 2019, 99% in 2022 and 93% in 2023.",
		},
		{
			// The same number can appear four times in the original and three in the
			// translation: the language builds the sentence differently. That is a
			// difference of phrasing, not of fact.
			"число повторено разное число раз",
			"28 147 случаев. Из этих 28 147 — большинство дети. Итого 28 147.",
			"28,147 cases, mostly children. That is 28,147 in total.",
		},
		{
			// A number spelled out in words is the translator's stylistic choice, and a good one.
			"число прописью",
			"Госдолг США перевалил за $40 триллионов.",
			"U.S. national debt crossed the forty-trillion-dollar mark.",
		},
		{
			// "123 тысячи" and "123,000" are the same thing written two ways.
			"тысячи словом против разрядов",
			"123 тысячи случаев кори за два года.",
			"123,000 measles cases in two years.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, i := range compareTranslation(c.src, c.dst) {
				if i.Key == "numbers" {
					t.Errorf("ложная тревога по числам: пропало %q", i.Detail)
				}
			}
		})
	}
}

// And this is what the check exists for. In the Kazakh version of the measles
// article the model dropped a whole sentence: "Tajikistan reported zero for 2024
// — sitting between Kyrgyzstan with its 14,380 cases and Uzbekistan with
// 20,940". The author could not read the text, but could see 14 380 had gone.
func TestTranslationCheckNamesTheMissingNumber(t *testing.T) {
	src := "Узбекистан — 20 940 случаев. Таджикистан отчитался о нуле, " +
		"находясь между Киргизией с её 14 380 случаями и Узбекистаном с 20 940."
	dst := "Өзбекстан — 20 940 жағдай."

	var detail string
	for _, i := range compareTranslation(src, dst) {
		if i.Key == "numbers" {
			detail = i.Detail
		}
	}
	if !strings.Contains(detail, "14 380") {
		t.Fatalf("пропавшее число не названо, получено %q", detail)
	}
	// 20 940 survived the translation, so there is nothing to complain about.
	if strings.Contains(detail, "20 940") {
		t.Errorf("названо число, которое на месте: %q", detail)
	}
}

// Numbers are listed as the author wrote them: they will look for "14 380" in
// their own text, not "14380". And the list is not endless — a lost paragraph
// takes dozens of figures with it, and the warning has to stay readable.
func TestTranslationCheckQuotesNumbersAsWrittenAndCapsTheList(t *testing.T) {
	var src strings.Builder
	src.WriteString("Было 14 380 случаев. ")
	for i := 0; i < 12; i++ {
		src.WriteString("Ещё ")
		src.WriteString(strings.Repeat("7", 1+i%3))
		src.WriteString(strings.Repeat("0", 3+i/3))
		src.WriteString(" штук. ")
	}
	var detail string
	for _, i := range compareTranslation(src.String(), "Ничего из этого.") {
		if i.Key == "numbers" {
			detail = i.Detail
		}
	}
	if !strings.Contains(detail, "14 380") {
		t.Errorf("число названо не так, как написано в тексте: %q", detail)
	}
	if n := strings.Count(detail, ",") + 1; n > 6 {
		t.Errorf("перечень не обрезан, %d позиций: %q", n, detail)
	}
	if !strings.HasSuffix(detail, "…") {
		t.Errorf("обрезанный перечень не помечен многоточием: %q", detail)
	}
}

// A loss and an invention are not the same thing. A dropped sentence leaves a
// hole the reader may notice; an invented number reads as fact and gets quoted
// as fact. Both versions of the measles article turned "forty trillion" into
// "forty-three" — in a text whose source list said forty.
func TestTranslationCheckCatchesInventedNumbers(t *testing.T) {
	got := compareTranslation(
		"Государственный долг США перешёл отметку в сорок триллионов долларов.",
		"АҚШ мемлекеттік қарызы 43 триллион доллар шекарасын асты.")
	var detail string
	for _, i := range got {
		if i.Key == "invented" {
			detail = i.Detail
		}
	}
	if !strings.Contains(detail, "43") {
		t.Fatalf("выдуманное число не найдено, получено %+v", got)
	}
}

// The counter-check must not complain about an honest translation — or it will
// share the fate of the first version, which people stopped looking at.
func TestInventedNumbersStayQuietOnHonestTranslations(t *testing.T) {
	cases := []struct{ name, src, dst string }{
		{"те же числа, другая запись", "Было 28 147 случаев и 98,5% охвата.",
			"There were 28,147 cases and 98.5% coverage."},
		{"масштабное слово раскрыто", "123 тысячи случаев.", "123,000 cases."},
		{"перевод пишет меньше чисел", "В 2023 году было 15 111 случаев.",
			"Fifteen thousand cases that year."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, i := range compareTranslation(c.src, c.dst) {
				if i.Key == "invented" {
					t.Errorf("ложная тревога: %q", i.Detail)
				}
			}
		})
	}
}

// The comparison runs paragraph against paragraph, and that is not a nicety. The
// article this started with passed a whole-text check cleanly: the translation
// wrote "43 триллион" where the original said "сорок", and the number 43 turned
// up four paragraphs below in "forty-three thousand cases". Against the whole
// text the evidence dissolves; against its own paragraph it does not.
func TestNumbersAreJudgedAgainstTheirOwnParagraph(t *testing.T) {
	src := "Госдолг США перешёл отметку в сорок триллионов долларов.\n\n" +
		"## Число\n\n" +
		"Сорок три тысячи заболевших за два года."
	dst := "АҚШ қарызы 43 триллион доллардан асты.\n\n" +
		"## Сан\n\n" +
		"Екі жылда қырық үш мың науқас."

	var detail string
	for _, i := range compareTranslation(src, dst) {
		if i.Key == "invented" {
			detail = i.Detail
		}
	}
	if !strings.Contains(detail, "43") {
		t.Fatalf("подменённое число не найдено, получено %q", detail)
	}
}

// The same number of paragraphs does not mean they are the same paragraphs. A
// translation that loses one sentence and splits another in two gives the same
// count — and matching by position starts complaining about everything after
// the shift. Then it is honester to compare the texts whole: cruder, but true.
func TestShiftedParagraphsFallBackToWholeText(t *testing.T) {
	src := "В 2019 году было 13 326 случаев.\n\n" +
		"Таджикистан отчитался о нуле при 14 380 у соседа.\n\n" +
		"В 2023 году — 15 111 случаев.\n\n" +
		"А в 2024-м — 28 147."
	// The second paragraph is dropped and the last split in two: the count adds up, the content has moved.
	dst := "2019 жылы 13 326 жағдай болды.\n\n" +
		"2023 жылы — 15 111 жағдай.\n\n" +
		"Ал 2024 жылы —\n\n" +
		"28 147 жағдай."

	got := compareTranslation(src, dst)
	var missing string
	for _, i := range got {
		if i.Key == "invented" {
			t.Errorf("сдвиг абзацев выдан за выдуманные числа: %q", i.Detail)
		}
		if i.Key == "numbers" {
			missing = i.Detail
		}
	}
	if !strings.Contains(missing, "14 380") {
		t.Errorf("потерянное число не названо: %q", missing)
	}
	if strings.Contains(missing, "28 147") || strings.Contains(missing, "15 111") {
		t.Errorf("уцелевшие числа названы потерянными: %q", missing)
	}
}
