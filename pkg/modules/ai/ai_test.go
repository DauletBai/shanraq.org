package ai

import (
	"context"
	"strings"
	"testing"
)

// fakeCompleter records requests and returns scripted output.
type fakeCompleter struct {
	calls []Request
	reply func(Request) string
}

func (f *fakeCompleter) Complete(_ context.Context, req Request) (string, error) {
	f.calls = append(f.calls, req)
	if f.reply != nil {
		return f.reply(req), nil
	}
	return "TR:" + req.User, nil
}

func TestModuleDisabledByDefault(t *testing.T) {
	m := New()
	if m.Enabled() {
		t.Fatal("module should be disabled before Init/config")
	}
	if _, err := m.Check(context.Background(), "sys", "текст", 100); err != ErrDisabled {
		t.Fatalf("expected ErrDisabled, got %v", err)
	}
}

func TestTranslateContent(t *testing.T) {
	fake := &fakeCompleter{reply: func(r Request) string { return "[" + r.User + "]" }}
	m := New()
	m.setCompleter(fake)

	out, err := m.translateContent(context.Background(), "ru", "kz", content{
		Title:   "Заголовок",
		Summary: "Кратко",
		Body:    "## Тело\n\nтекст",
	})
	if err != nil {
		t.Fatalf("translateContent: %v", err)
	}
	if !strings.Contains(out.Title, "Заголовок") || !strings.Contains(out.Summary, "Кратко") ||
		!strings.Contains(out.Body, "Тело") {
		t.Fatalf("unexpected output: %+v", out)
	}
	// Two paragraphs of body, then the headline and the summary.
	if len(fake.calls) != 4 {
		t.Fatalf("expected 4 completion calls (2 paragraphs, title, summary), got %d", len(fake.calls))
	}
	// The body keeps its shape: as many paragraphs out as went in.
	if got := len(blockSplit.Split(out.Body, -1)); got != 2 {
		t.Errorf("body came back as %d paragraphs, want 2: %q", got, out.Body)
	}

	// The first paragraph goes unframed: it is the text that settles which
	// words this article uses.
	if strings.Contains(fake.calls[0].User, translateMarker) {
		t.Error("the first paragraph was wrapped in context it does not need")
	}
	if !strings.Contains(fake.calls[0].User, "Тело") {
		t.Errorf("the first call was not the opening paragraph: %q", fake.calls[0].User)
	}

	// Everything after it is shown the wording already chosen, so a term
	// settled in the first paragraph is the term used in the last. Showing the
	// Russian source instead did not help: the question is not what the article
	// is about but which word to use in the target language — and that word is
	// by then already picked.
	for i := 1; i < len(fake.calls); i++ {
		u := fake.calls[i].User
		if !strings.Contains(u, "already translated into") {
			t.Errorf("call %d was not given the finished translation to match: %q", i, u)
		}
		if !strings.Contains(u, translateMarker) {
			t.Errorf("call %d does not mark which part to translate", i)
		}
	}
	for i, want := range []string{"Заголовок", "Кратко"} {
		if u := strings.TrimSpace(fake.calls[len(fake.calls)-2+i].User); !strings.HasSuffix(u, want) {
			t.Errorf("call for %q does not end with the field itself: %q", want, u)
		}
	}
	// The system prompt must name both source and target languages.
	sys := fake.calls[0].System
	if !strings.Contains(sys, "Russian") || !strings.Contains(sys, "Kazakh") {
		t.Fatalf("system prompt missing language names: %q", sys)
	}
	if fake.calls[0].Model != m.translateModel {
		t.Fatalf("translation should use translate model, got %q", fake.calls[0].Model)
	}
}

// The fault that made all of this necessary: asked for a whole article in one
// request, the model returned fewer paragraphs than it was given, and said
// nothing about it. Paragraph by paragraph the count cannot drift — every block
// is asked for and every answer is put back where it came from.
func TestBodyKeepsEveryParagraph(t *testing.T) {
	fake := &fakeCompleter{reply: func(r Request) string { return "аударма" }}
	m := New()
	m.setCompleter(fake)

	var src []string
	for i := 0; i < 40; i++ {
		src = append(src, "Абзац номер такой-то, и в нём есть текст.")
	}
	body := strings.Join(src, "\n\n")

	out, err := m.translateContent(context.Background(), "ru", "kz", content{Body: body})
	if err != nil {
		t.Fatalf("translateContent: %v", err)
	}
	if got := len(blockSplit.Split(out.Body, -1)); got != 40 {
		t.Fatalf("40 paragraphs in, %d out", got)
	}
}

// A table rule and a horizontal rule have nothing to translate. Sending them to
// a model is a request that can only do harm, and it is one request per rule.
func TestBlocksWithoutWordsAreNotSent(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string { return "аударма" }}
	m := New()
	m.setCompleter(fake)

	out, err := m.translateContent(context.Background(), "ru", "kz", content{
		Body: "Текст\n\n|---|---|\n\n---\n\nЕщё текст",
	})
	if err != nil {
		t.Fatalf("translateContent: %v", err)
	}
	if len(fake.calls) != 2 {
		t.Errorf("expected 2 calls for the 2 blocks with words, got %d", len(fake.calls))
	}
	if !strings.Contains(out.Body, "|---|---|") || !strings.Contains(out.Body, "\n---\n") {
		t.Errorf("the rules did not come through unchanged: %q", out.Body)
	}
}

// A wrong figure is the one fault that survives every kind of review: it reads
// as fact, it is quoted as fact, and the author cannot catch it because the
// whole reason the site translates for them is that they do not read the
// language. So it is refused, not reported.
func TestTranslationThatChangesNumbersIsRefused(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string {
		return "АҚШ қарызы 43 триллион доллардан асты."
	}}
	m := New()
	m.setCompleter(fake)

	_, err := m.translateContent(context.Background(), "ru", "kz", content{
		Body: "Госдолг США перешёл отметку в сорок триллионов долларов.",
	})
	if err == nil {
		t.Fatal("a translation that invented a number was accepted")
	}
	if !strings.Contains(err.Error(), "43") {
		t.Errorf("the error does not name the offending figure: %v", err)
	}
	if len(fake.calls) != translateAttempts {
		t.Errorf("expected %d attempts, got %d", translateAttempts, len(fake.calls))
	}
	// Every retry after the first must say what went wrong; a blind retry only
	// re-rolls the same dice.
	if !strings.Contains(fake.calls[1].User, "43") {
		t.Errorf("the retry was not told which figure was wrong: %q", fake.calls[1].User)
	}
}

// One bad roll must not cost the article: the next attempt is told what was
// wrong and is usually right.
func TestTranslationRetriesUntilTheNumbersMatch(t *testing.T) {
	n := 0
	fake := &fakeCompleter{reply: func(Request) string {
		n++
		if n == 1 {
			return "АҚШ қарызы 43 триллион доллардан асты."
		}
		return "АҚШ қарызы қырық триллион доллардан асты."
	}}
	m := New()
	m.setCompleter(fake)

	out, err := m.translateContent(context.Background(), "ru", "kz", content{
		Body: "Госдолг США перешёл отметку в сорок триллионов долларов.",
	})
	if err != nil {
		t.Fatalf("translateContent: %v", err)
	}
	if strings.Contains(out.Body, "43") {
		t.Errorf("the bad attempt was kept: %q", out.Body)
	}
	if len(fake.calls) != 2 {
		t.Errorf("expected 2 calls (one wrong, one right), got %d", len(fake.calls))
	}
}

func TestTranslateContentSkipsEmptyFields(t *testing.T) {
	fake := &fakeCompleter{}
	m := New()
	m.setCompleter(fake)

	out, err := m.translateContent(context.Background(), "ru", "en", content{Body: "only body"})
	if err != nil {
		t.Fatalf("translateContent: %v", err)
	}
	if out.Title != "" || out.Summary != "" {
		t.Fatalf("empty fields should stay empty, got %+v", out)
	}
	if len(fake.calls) != 1 {
		t.Fatalf("expected 1 call (body only), got %d", len(fake.calls))
	}
}
