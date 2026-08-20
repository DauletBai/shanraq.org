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
	if out.Title != "[Заголовок]" || out.Summary != "[Кратко]" || !strings.Contains(out.Body, "Тело") {
		t.Fatalf("unexpected output: %+v", out)
	}
	if len(fake.calls) != 3 {
		t.Fatalf("expected 3 completion calls (title, summary, body), got %d", len(fake.calls))
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
