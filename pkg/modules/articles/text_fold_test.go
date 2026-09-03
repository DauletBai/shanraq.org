package articles

import (
	"strings"
	"testing"
)

func TestFoldAnswers(t *testing.T) {
	lesson := "## Задание\n\nНапишите функцию.\n\n## Ответы\n\n1. Потому что.\n2. Оттого.\n\n## Источники\n\n- [go.dev](https://go.dev)\n"
	html, toc := RenderMarkdownTOC(lesson)
	s := string(html)
	if !strings.Contains(s, `<details class="fold">`) {
		t.Fatalf("ответы не свёрнуты: %s", s)
	}
	if !strings.Contains(s, "Показать ответы") {
		t.Error("нет подписи на русском")
	}
	if strings.Index(s, `id="sec-2"`) > strings.Index(s, "<details") {
		t.Error("якорь оглавления попал внутрь свёртки")
	}
	if strings.Contains(s[strings.Index(s, "<details"):strings.Index(s, "</details>")], "Источники") {
		t.Error("в свёртку затянуло следующий раздел")
	}
	if len(toc) != 3 {
		t.Errorf("оглавление сбилось: %d", len(toc))
	}
	// An ordinary article that happens to answer something is left alone.
	plain, _ := RenderMarkdownTOC("## Вопросы\n\nтекст\n\n## Ответы\n\nтекст\n")
	if strings.Contains(string(plain), "<details") {
		t.Error("свёрнута обычная статья, а не урок")
	}
	// Kazakh and English lessons fold with their own label.
	for src, want := range map[string]string{
		"## Тапсырма\n\nЖазыңыз.\n\n## Жауаптар\n\n1. Себебі.\n":  "Жауаптарды көрсету",
		"## Exercise\n\nWrite it.\n\n## Answers\n\n1. Because.\n": "Show the answers",
	} {
		got, _ := RenderMarkdownTOC(src)
		if !strings.Contains(string(got), want) {
			t.Errorf("нет подписи %q", want)
		}
	}
}
