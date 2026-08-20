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
	if len(fake.calls) != 3 {
		t.Fatalf("expected 3 completion calls (title, summary, body), got %d", len(fake.calls))
	}

	// The body is translated first and unframed: it is the text that settles
	// which words this article uses.
	if strings.Contains(fake.calls[0].User, "TRANSLATE ONLY WHAT FOLLOWS") {
		t.Error("the body was wrapped in context it does not need")
	}
	if !strings.Contains(fake.calls[0].User, "Тело") {
		t.Errorf("the first call was not the body: %q", fake.calls[0].User)
	}

	// The headline and the summary follow, each shown the finished translation
	// so a term chosen once is used everywhere. Showing them the Russian source
	// did not help: the question is not what the article is about, but which
	// word to use in the target language — and that word is now already chosen.
	for i, want := range []string{"Заголовок", "Кратко"} {
		u := fake.calls[i+1].User
		if !strings.Contains(u, "already translated into") {
			t.Errorf("call %d was not given the finished translation to match: %q", i+1, u)
		}
		if !strings.Contains(u, out.Body[:20]) {
			t.Errorf("call %d does not carry the translated body: %q", i+1, u)
		}
		if !strings.Contains(u, "TRANSLATE ONLY WHAT FOLLOWS") {
			t.Errorf("call %d does not mark which part to translate", i+1)
		}
		if !strings.HasSuffix(strings.TrimSpace(u), want) {
			t.Errorf("call %d does not end with the field itself: %q", i+1, u)
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
