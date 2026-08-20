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
		// Разделители разрядов и десятичная запятая меняются при локализации —
		// это правильный перевод, а не поломка.
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

	// Пустые поля — не повод для претензий: переводить ещё нечего.
	if got := compareTranslation("", "что-то"); got != nil {
		t.Errorf("пустой оригинал дал замечания: %+v", got)
	}
	if got := compareTranslation(src, ""); got != nil {
		t.Errorf("пустой перевод дал замечания: %+v", got)
	}
}

// Первая версия проверки кричала на каждом переводе — и на всех трёх языках
// сразу. Автор, увидев четыре предупреждения на исправном тексте, перестаёт
// читать предупреждения вообще; ложная тревога опаснее молчания, потому что
// она обесценивает и настоящую.
//
// Все случаи ниже взяты из продакшена, из статьи про корь.
func TestTranslationChecksStayQuietOnHonestTranslations(t *testing.T) {
	cases := []struct{ name, src, dst string }{
		{
			// Разбор: «99% in 2019, 99% in 2022». Первая версия считала любую
			// запятую между цифрами разделителем разрядов и склеивала это в
			// число 201999 — которого, разумеется, не было в оригинале.
			"процент и год через запятую",
			"Охват 99% в 2019 году, 99% в 2022 году и 93% в 2023 году.",
			"Coverage was 99% in 2019, 99% in 2022 and 93% in 2023.",
		},
		{
			// Одно и то же число может встретиться в оригинале четыре раза, а
			// в переводе три: язык иначе строит фразу. Это разница изложения,
			// а не факта.
			"число повторено разное число раз",
			"28 147 случаев. Из этих 28 147 — большинство дети. Итого 28 147.",
			"28,147 cases, mostly children. That is 28,147 in total.",
		},
		{
			// Число прописью — стилистическое решение переводчика, и хорошее.
			"число прописью",
			"Госдолг США перевалил за $40 триллионов.",
			"U.S. national debt crossed the forty-trillion-dollar mark.",
		},
		{
			// «123 тысячи» и «123,000» — одно и то же, записанное по-разному.
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

// А это — то, ради чего проверка существует. В казахской версии статьи про корь
// модель выбросила целое предложение: «Таджикистан за 2024 год отчитался о нуле
// — находясь между Киргизией с её 14 380 случаями и Узбекистаном с 20 940».
// Прочитать текст автор не мог, а увидеть, что число 14 380 исчезло, — мог.
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
	// 20 940 в переводе осталось — жаловаться не на что.
	if strings.Contains(detail, "20 940") {
		t.Errorf("названо число, которое на месте: %q", detail)
	}
}

// Числа перечисляются так, как их написал автор: он будет искать «14 380» в
// своём тексте, а не «14380». И перечень не бесконечен — потерянный абзац
// уносит с собой десятки цифр, а предупреждение должно оставаться читаемым.
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

// Пропажа и выдумка — не одно и то же. Потерянное предложение оставляет дыру,
// которую читатель может заметить; выдуманное число читается как факт и цитируется
// как факт. Обе версии статьи про корь превратили «сорок триллионов» в «сорок три»
// — в тексте, где в списке источников стояло сорок.
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

// Встречная проверка не должна ругаться на честный перевод — иначе она разделит
// судьбу первой версии, на которую перестали смотреть.
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
