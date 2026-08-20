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
