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
			// Кириллица: считаем в рунах, а не в байтах — иначе оценка вдвое завышена.
			src := strings.Repeat("я", c.chars)
			if got := translationBudget(src, c.floor); got != c.want {
				t.Errorf("translationBudget(%d знаков, минимум %d) = %d, ожидалось %d",
					c.chars, c.floor, got, c.want)
			}
		})
	}

	// Байты против рун: русская буква занимает два байта, и подсчёт по длине
	// строки дал бы вдвое больший лимит, чем нужно.
	if translationBudget(strings.Repeat("я", 1000), 0) != translationBudget(strings.Repeat("a", 1000), 0) {
		t.Error("оценка зависит от алфавита — значит, считает байты, а не символы")
	}
}
