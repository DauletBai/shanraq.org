package ai

import (
	"strings"
	"testing"
)

// The first live translation came back with the link label still in the source
// language — `[ссылкой]` sitting inside a Kazakh sentence — because the prompt
// told the model to preserve the Markdown "exactly" and, with no room to think
// it over, it took that to cover the words inside the markup too.
//
// The distinction the prompt has to make: markup and URLs are structure, the
// words a reader sees are prose.
func TestTranslatePromptSeparatesMarkupFromProse(t *testing.T) {
	p := translateSystem("ru", "kz")

	for _, want := range []string{
		"link labels",   // the case that actually failed
		"[brackets]",    // named concretely, not left to inference
		"byte for byte", // what preservation applies to
		"(parentheses)", // URLs are not prose
		"code",          // nor is code
	} {
		if !strings.Contains(p, want) {
			t.Errorf("the prompt never mentions %q:\n%s", want, p)
		}
	}

	// The old wording is what caused it; it must not come back.
	if strings.Contains(p, "Markdown formatting exactly") {
		t.Error("the prompt still tells the model to preserve the formatting 'exactly', " +
			"which is what it read as 'do not touch the words inside it'")
	}

	// Both language names have to survive the format string — a %[2]s typo
	// would silently produce a prompt naming one language twice.
	if !strings.Contains(p, langFullName["ru"]) || !strings.Contains(p, langFullName["kz"]) {
		t.Errorf("the prompt does not name both languages:\n%s", p)
	}
	if strings.Count(p, langFullName["kz"]) < 2 {
		t.Error("the target language should be named where the prompt asks for idiomatic writing")
	}
}

// The first long article this site tried to translate failed three times: the
// body needed more room than the configured 4096-token ceiling, and a fixed cap
// is the wrong shape for a job whose output is as long as its input.
func TestTranslationBudgetFollowsTheText(t *testing.T) {
	cases := []struct {
		name  string
		chars int
		floor int
		want  int
	}{
		{"заголовок остаётся на своём минимуме", 45, 512, 512},
		{"короткое описание тоже", 450, 1024, 1024},
		{"статья получает по своему размеру", 12690, 4096, 9517},
		{"очень длинный текст упирается в потолок", 100000, 4096, maxTranslationTokens},
		{"пустая строка не даёт нулевой лимит", 0, 4096, 4096},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Cyrillic: count in runes, not bytes — otherwise the estimate is doubled.
			src := strings.Repeat("я", c.chars)
			if got := translationBudget(src, c.floor); got != c.want {
				t.Errorf("translationBudget(%d знаков, минимум %d) = %d, ожидалось %d",
					c.chars, c.floor, got, c.want)
			}
		})
	}

	// Bytes against runes: a Russian letter takes two bytes, and counting by string
	// length would give twice the limit needed.
	if translationBudget(strings.Repeat("я", 1000), 0) != translationBudget(strings.Repeat("a", 1000), 0) {
		t.Error("оценка зависит от алфавита — значит, считает байты, а не символы")
	}
}

// One of nine URLs came back swapped in the first real translation: the WHO
// fact sheet pointed at a Bangladeshi newspaper. A prompt cannot guarantee
// this and a reader cannot check it, so the code does.
func TestLinkTargetsAreRestoredFromTheSource(t *testing.T) {
	src := "См. [ВОЗ](https://who.int/measles) и [газету](https://tbsnews.net/story)."
	bad := "Қараңыз: [ДДҰ](https://tbsnews.net/story) және [газет](https://tbsnews.net/story)."

	got, ok := restoreLinkTargets(src, bad)
	if !ok {
		t.Fatal("восстановление не сработало на совпадающем числе ссылок")
	}
	if !strings.Contains(got, "[ДДҰ](https://who.int/measles)") {
		t.Errorf("первая ссылка не восстановлена:\n%s", got)
	}
	if !strings.Contains(got, "[газет](https://tbsnews.net/story)") {
		t.Errorf("вторая ссылка испорчена:\n%s", got)
	}
	// Captions are the translator's work and must not be touched.
	if strings.Contains(got, "ВОЗ") {
		t.Error("восстановление URL затронуло текст подписи")
	}

	// If the model lost or added a link, matching by position starts to lie — then it is
	// better not to touch anything at all.
	if _, ok := restoreLinkTargets(src, "Только [одна](https://example.com/x)."); ok {
		t.Error("при несовпадении числа ссылок восстановление должно отказываться")
	}
	// Text without links is not an error.
	if _, ok := restoreLinkTargets("без ссылок", "сілтемесіз"); !ok {
		t.Error("текст без ссылок не должен считаться сбоем")
	}
}

// Told to use the glossary but not translate it, the model translated it,
// echoed the marker, and appended the answer — a 461-character summary came
// back as 1,275 with the real translation at the end. The instruction is not
// enforceable; the cut is.
func TestOnlyWhatFollowsTheMarkerIsKept(t *testing.T) {
	cases := []struct{ name, reply, want string }{
		{
			"модель вернула весь промпт",
			"Аударылған контекст мәтіні.\n\n---TRANSLATE ONLY WHAT FOLLOWS, matching the terminology above---\n\nНағыз аударма.",
			"Нағыз аударма.",
		},
		{
			"маркер переформатирован",
			"контекст\nTRANSLATE ONLY WHAT FOLLOWS\nответ",
			"ответ",
		},
		{
			"послушная модель — ответ не трогаем",
			"Просто перевод без всякого маркера.",
			"Просто перевод без всякого маркера.",
		},
		{
			"маркер встречается дважды — берём последний",
			"a TRANSLATE ONLY WHAT FOLLOWS\nb\nTRANSLATE ONLY WHAT FOLLOWS\nитог",
			"итог",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := afterMarker(c.reply); got != c.want {
				t.Errorf("afterMarker() = %q, ожидалось %q", got, c.want)
			}
		})
	}
	// The prompt and the cutter have to agree on one and the same string.
	if !strings.Contains(withGlossary("глоссарий", "поле", "Kazakh"), translateMarker) {
		t.Error("промпт не содержит маркера, по которому режет afterMarker")
	}
}
