package ai

import (
	"context"
	"testing"
)

func TestModerateDisabledByDefault(t *testing.T) {
	m := New()
	if _, err := m.Moderate(context.Background(), "comment", "любой текст"); err != ErrDisabled {
		t.Fatalf("expected ErrDisabled when off, got %v", err)
	}
}

func TestModerateFlags(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string {
		return `{"action":"flag","reason":"спам","confidence":0.92}`
	}}
	m := New()
	m.setCompleter(fake)

	v, err := m.Moderate(context.Background(), "comment", "Купите дешево, пишите в личку!!!")
	if err != nil {
		t.Fatalf("Moderate: %v", err)
	}
	if !v.Flagged() || v.Reason != "спам" {
		t.Fatalf("expected flagged verdict, got %+v", v)
	}
	if fake.calls[0].Model != m.translateModel {
		t.Fatalf("moderation should use the cheap (translate) model, got %q", fake.calls[0].Model)
	}
	if fake.calls[0].MaxTokens != 256 {
		t.Fatalf("moderation should cap tokens tight, got %d", fake.calls[0].MaxTokens)
	}
}

func TestModerateAllowsOrdinaryText(t *testing.T) {
	fake := &fakeCompleter{reply: func(Request) string {
		// Model may wrap the JSON in prose or fences; parser must tolerate it.
		return "Here is my verdict:\n```json\n{\"action\":\"allow\",\"reason\":\"\",\"confidence\":0.1}\n```"
	}}
	m := New()
	m.setCompleter(fake)

	v, err := m.Moderate(context.Background(), "comment", "Не согласен с автором, статья слабая.")
	if err != nil {
		t.Fatalf("Moderate: %v", err)
	}
	if v.Flagged() {
		t.Fatalf("ordinary criticism must be allowed, got %+v", v)
	}
}

func TestModerateEmptyIsAllowedWithoutCall(t *testing.T) {
	fake := &fakeCompleter{}
	m := New()
	m.setCompleter(fake)

	v, err := m.Moderate(context.Background(), "comment", "   ")
	if err != nil {
		t.Fatalf("Moderate: %v", err)
	}
	if v.Flagged() {
		t.Fatalf("empty text should be allowed, got %+v", v)
	}
	if len(fake.calls) != 0 {
		t.Fatalf("empty text must not call the model, got %d calls", len(fake.calls))
	}
}

func TestParseModerationVerdictNormalizesUnknownAction(t *testing.T) {
	v, err := parseModerationVerdict(`{"action":"maybe","reason":"noise","confidence":0.5}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if v.Action != "allow" || v.Reason != "" {
		t.Fatalf("unknown action should normalize to allow, got %+v", v)
	}
}
