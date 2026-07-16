package ai

import (
	"context"
	"strings"
	"testing"
)

func TestSupportAnswerDisabled(t *testing.T) {
	m := New()
	if _, err := m.Answer(context.Background(), "ru", "как подать объявление?"); err != ErrDisabled {
		t.Fatalf("expected ErrDisabled, got %v", err)
	}
}

func TestSupportAnswerReplies(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string { return "  Откройте Студию и нажмите «Новое объявление».  " }}
	m := New()
	m.setCompleter(fake)

	got, err := m.Answer(context.Background(), "ru", "как подать объявление?")
	if err != nil {
		t.Fatalf("Answer: %v", err)
	}
	if strings.TrimSpace(got) != "Откройте Студию и нажмите «Новое объявление»." {
		t.Fatalf("unexpected reply: %q", got)
	}
	if !strings.Contains(fake.calls[0].System, "Russian") {
		t.Fatalf("system prompt should target Russian: %q", fake.calls[0].System)
	}
}

func TestSupportAnswerEscalates(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string { return "ESCALATE" }}
	m := New()
	m.setCompleter(fake)

	got, err := m.Answer(context.Background(), "en", "someone stole my money, refund me")
	if err != nil {
		t.Fatalf("Answer: %v", err)
	}
	if got != "" {
		t.Fatalf("escalation should yield empty reply, got %q", got)
	}
}

func TestDescribeListingUsesFactsOnly(t *testing.T) {
	fake := &fakeCompleter{reply: func(r Request) string { return "DESC of: " + r.User }}
	m := New()
	m.setCompleter(fake)

	got, err := m.DescribeListing(context.Background(), "ru", "2-комн., 54 м², Алматы")
	if err != nil {
		t.Fatalf("DescribeListing: %v", err)
	}
	if !strings.Contains(got, "54 м²") {
		t.Fatalf("description should be built from facts, got %q", got)
	}
	if !strings.Contains(fake.calls[0].System, "ONLY the facts provided") {
		t.Fatalf("system prompt must forbid inventing facts: %q", fake.calls[0].System)
	}
}

func TestDescribeListingEmptySkips(t *testing.T) {
	fake := &fakeCompleter{}
	m := New()
	m.setCompleter(fake)

	got, err := m.DescribeListing(context.Background(), "ru", "   ")
	if err != nil || got != "" {
		t.Fatalf("empty facts should return empty without error, got %q / %v", got, err)
	}
	if len(fake.calls) != 0 {
		t.Fatalf("empty facts must not call the model")
	}
}

func TestDraftColumnEvergreenPolicy(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string { return "## О доверии\n\nтекст" }}
	m := New()
	m.setCompleter(fake)

	got, err := m.DraftColumn(context.Background(), "ru", "о доверии между людьми")
	if err != nil {
		t.Fatalf("DraftColumn: %v", err)
	}
	if !strings.Contains(got, "## О доверии") {
		t.Fatalf("unexpected draft: %q", got)
	}
	sys := fake.calls[0].System
	if !strings.Contains(sys, "EVERGREEN ONLY") || !strings.Contains(sys, "NEVER fabricate") {
		t.Fatalf("author prompt must enforce evergreen + no-fabrication policy: %q", sys)
	}
	if fake.calls[0].Model != m.editorModel {
		t.Fatalf("column drafting should use the editor model, got %q", fake.calls[0].Model)
	}
}

func TestDraftColumnDisabled(t *testing.T) {
	m := New()
	if _, err := m.DraftColumn(context.Background(), "ru", "тема"); err != ErrDisabled {
		t.Fatalf("expected ErrDisabled, got %v", err)
	}
}
