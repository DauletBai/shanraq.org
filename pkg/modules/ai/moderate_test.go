package ai

import (
	"context"
	"strings"
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

// Скрытая реклама — отдельное правило, а не разновидность спама: спам виден
// сразу, а отзыв, написанный ради продажи, выглядит как мнение читателя. Ради
// него код правила и появился.
func TestModeratePromptCoversEveryPublishedRule(t *testing.T) {
	sys := moderateSystem("comment")
	for _, code := range []string{
		"spam", "scam", "hidden_ad", "prohibited_goods", "abuse",
		"hatred", "defamation", "personal_data", "adult", "illegal",
	} {
		if !strings.Contains(sys, code+":") {
			t.Errorf("правило %q не описано в промпте", code)
		}
	}
	// Объявления добавляют своё правило поверх общих.
	if !strings.Contains(moderateSystem("listing"), "misleading:") {
		t.Error("для объявлений нет правила о недостоверных данных")
	}
	if strings.Contains(moderateSystem("comment"), "misleading:") {
		t.Error("правило для объявлений попало в промпт для комментариев")
	}
}

// Перекос в сторону разрешения — не небрежность, а позиция площадки, и он
// должен быть в промпте написан, а не подразумеваться.
func TestModeratePromptProtectsCriticism(t *testing.T) {
	sys := moderateSystem("comment")
	for _, phrase := range []string{
		"over-blocking is worse than under-blocking",
		"however harsh or blunt",
		"When unsure, allow",
	} {
		if !strings.Contains(sys, phrase) {
			t.Errorf("в промпте нет защиты критики: %q", phrase)
		}
	}
}

// Скрытая причина — это не причина. Площадка обещает обжалование, а обжаловать
// «модели не понравилось» нельзя.
func TestModerationVerdictCarriesTheRule(t *testing.T) {
	v, err := parseModerationVerdict(
		`{"action":"flag","rule":"hidden_ad","reason":"скрытая реклама клиники","confidence":0.8}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if v.Rule != "hidden_ad" {
		t.Errorf("код правила потерян: %+v", v)
	}
	// А на «разрешено» — ничего лишнего.
	v, err = parseModerationVerdict(`{"action":"allow","rule":"spam","reason":"на всякий случай"}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if v.Rule != "" || v.Reason != "" {
		t.Errorf("разрешённый текст унёс с собой причину: %+v", v)
	}
}
