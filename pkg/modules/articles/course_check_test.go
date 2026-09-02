package articles

import (
	"encoding/json"
	"net/http"
	"net/url"
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

// gofmt is the point, not an imitation of it: the reader must see their editor's
// own behaviour, and code that will not parse must be refused locally rather
// than sent to a paid reviewer.
func TestFormatSolution(t *testing.T) {
	messy := "package main\nimport \"fmt\"\nfunc main(){\nfmt.Println( \"x\" )\n}"
	out, err := formatSolution(messy)
	if err != nil {
		t.Fatalf("годный код не отформатирован: %v", err)
	}
	if !strings.Contains(out, "\tfmt.Println(\"x\")") {
		t.Errorf("отступ и пробелы не приведены к gofmt:\n%s", out)
	}
	if strings.Contains(out, "import \"fmt\"\nfunc") {
		t.Error("пустая строка между import и func не поставлена")
	}

	if _, err := formatSolution("package main\nfunc main() {"); err == nil {
		t.Error("непарсящийся код принят")
	} else if h := syntaxHint(err); h == "" || strings.Contains(h, ".go:") {
		t.Errorf("подсказка непригодна для новичка: %q", h)
	}
}

// The preview is rendered by the article pipeline, so it carries the same
// classes the lessons do and inherits their palette.
func TestHighlightGo(t *testing.T) {
	h := string(highlightGo("package main"))
	if !strings.Contains(h, `class="chroma"`) {
		t.Errorf("подсветка не применена: %.80s", h)
	}
	if !strings.Contains(h, "package") {
		t.Error("код потерян при подсветке")
	}
}

// The tidy endpoint end to end: a signed-in reader sends untidy code and gets
// it back gofmt'd and coloured, without any of it touching the reviewer or the
// three attempts.
func TestCourseFormatEndpoint(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("fmt@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = 'fmt@t.test'`)
	cookie := app.login("fmt@t.test", "Parol12345")
	_, slug := app.seedArticle(author, "published")

	send := func(code string) (int, checkResponse) {
		f := url.Values{}
		f.Set("solution", code)
		w := app.do(http.MethodPost, "/read/"+slug+"/format", f, withCookie(cookie))
		var res checkResponse
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		return w.Code, res
	}

	t.Run("кривой отступ выпрямляется", func(t *testing.T) {
		code, res := send("package main\nimport \"fmt\"\nfunc main(){\nfmt.Println( \"x\" )\n}")
		if code != http.StatusOK {
			t.Fatalf("код ответа %d: %s", code, res.Error)
		}
		if !strings.Contains(res.Code, "\tfmt.Println(\"x\")") {
			t.Errorf("не отформатировано:\n%s", res.Code)
		}
		if !strings.Contains(res.HTML, `class="chroma"`) {
			t.Error("подсветка не пришла")
		}
	})

	t.Run("непарсящийся код не уходит на проверку", func(t *testing.T) {
		code, res := send("package main\nfunc main() {")
		if code != http.StatusUnprocessableEntity {
			t.Fatalf("код ответа %d, ждали 422", code)
		}
		if res.Syntax == "" {
			t.Error("нет подсказки, где сломано")
		}
	})

	t.Run("без входа не отвечает", func(t *testing.T) {
		f := url.Values{}
		f.Set("solution", "package main")
		if w := app.do(http.MethodPost, "/read/"+slug+"/format", f); w.Code != http.StatusUnauthorized {
			t.Errorf("код ответа %d, ждали 401", w.Code)
		}
	})
}
