package articles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFoldAnswers(t *testing.T) {
	lesson := "## Задание\n\nНапишите функцию.\n\n## Ответы\n\n1. Потому что.\n2. Оттого.\n\n## Источники\n\n- [go.dev](https://go.dev)\n"
	html, toc := renderLessonTOC(lesson)
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
	plain, _ := renderLessonTOC("## Вопросы\n\nтекст\n\n## Ответы\n\nтекст\n")
	if strings.Contains(string(plain), "<details") {
		t.Error("свёрнута обычная статья, а не урок")
	}
	// Kazakh and English lessons fold with their own label.
	for src, want := range map[string]string{
		"## Тапсырма\n\nЖазыңыз.\n\n## Жауаптар\n\n1. Себебі.\n":  "Жауаптарды көрсету",
		"## Exercise\n\nWrite it.\n\n## Answers\n\n1. Because.\n": "Show the answers",
	} {
		got, _ := renderLessonTOC(src)
		if !strings.Contains(string(got), want) {
			t.Errorf("нет подписи %q", want)
		}
	}
}

func TestRustNinthBatchAnswersStayBehindFold(t *testing.T) {
	lessons := []string{"41-iterators", "42-closures", "43-maps-sets", "44-modules", "45-generics"}
	for _, stem := range lessons {
		for _, locale := range []struct{ suffix, heading string }{
			{"", "После проверки"}, {"-kz", "Тексергеннен кейін"}, {"-en", "After checking"},
		} {
			t.Run(stem+locale.suffix, func(t *testing.T) {
				path := filepath.Join("..", "..", "..", "course", "lessons", "rust", fmt.Sprintf("%s%s.md", stem, locale.suffix))
				file, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				lines := strings.SplitN(string(file), "\n", 5)
				if len(lines) != 5 {
					t.Fatal("missing lesson body")
				}
				html, _ := renderLessonTOC(lines[4])
				page := string(html)
				start := strings.Index(page, `<details class="fold">`)
				end := strings.Index(page, "</details>")
				next := strings.Index(page, locale.heading)
				if start < 0 || end < start || next < end {
					t.Fatalf("answer fold or following navigation section is misplaced: fold=%d end=%d next=%d", start, end, next)
				}
			})
		}
	}
}
