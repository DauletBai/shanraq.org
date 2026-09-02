package articles

import (
	"strings"
	"testing"
)

// The exercise is read out of the lesson's own text, so a lesson edited without
// anyone remembering this code still marks readers against what they read.
func TestLessonExercise(t *testing.T) {
	body := `## Зачем это нужно

Текст.

## Задание

**Обязательное.** Напишите функцию.

- по желанию раз
- по желанию два

## Куда это встанет в блоге

Другое.`
	got := lessonExercise(body)
	if got == "" {
		t.Fatal("задание не найдено")
	}
	if want := "**Обязательное.** Напишите функцию."; got[:len(want)] != want {
		t.Errorf("захвачено не с начала задания: %.40q", got)
	}
	if strings.Contains(got, "Куда это встанет") {
		t.Error("захвачен следующий раздел")
	}
	for _, head := range []string{"## Тапсырма", "## Exercise"} {
		if lessonExercise("## A\n\nx\n\n"+head+"\n\ny\n\n## B\n\nz") != "y" {
			t.Errorf("не найдено по заголовку %q", head)
		}
	}
	if lessonExercise("## Только текст\n\nбез задания") != "" {
		t.Error("выдумано задание там, где его нет")
	}
}

// A model asked for bare JSON returns it wrapped, prefaced or both. The reader
// must not see a stack trace because of that.
func TestParseCheckVerdict(t *testing.T) {
	ok := []struct{ name, raw string }{
		{"чистый json", `{"passed": true, "note": "Работает."}`},
		{"в заборе", "```json\n{\"passed\": true, \"note\": \"Работает.\"}\n```"},
		{"с предисловием", "Вот разбор:\n{\"passed\": false, \"note\": \"Деление целое.\"}"},
	}
	for _, c := range ok {
		t.Run(c.name, func(t *testing.T) {
			v, err := parseCheckVerdict(c.raw)
			if err != nil {
				t.Fatalf("не разобрано: %v", err)
			}
			if v.Note == "" {
				t.Error("пустой комментарий принят")
			}
		})
	}
	bad := []struct{ name, raw string }{
		{"не json", "Мне кажется, всё хорошо."},
		{"без комментария", `{"passed": true, "note": "   "}`},
		{"пусто", ""},
	}
	for _, c := range bad {
		t.Run(c.name, func(t *testing.T) {
			if _, err := parseCheckVerdict(c.raw); err == nil {
				t.Error("принят непригодный ответ")
			}
		})
	}
}

// A fence around a pasted solution is the reader being tidy, not part of their
// program.
func TestUnfence(t *testing.T) {
	if got := unfence("```go\npackage main\n```"); got != "package main" {
		t.Errorf("забор не снят: %q", got)
	}
	if got := unfence("package main"); got != "package main" {
		t.Errorf("код без забора испорчен: %q", got)
	}
}
