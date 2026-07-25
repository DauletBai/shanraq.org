package ai

import "testing"

func TestNormalizeSettings(t *testing.T) {
	base := AISettings{Provider: ProviderOpenAI, EditorModel: "gpt-5", TranslateModel: "gpt-5-mini", MaxTokens: 4096}

	if _, err := normalizeSettings(base); err != nil {
		t.Fatalf("valid settings rejected: %v", err)
	}

	if _, err := normalizeSettings(AISettings{Provider: "grok", EditorModel: "x", TranslateModel: "y", MaxTokens: 1000}); err == nil {
		t.Error("unknown provider should be rejected")
	}

	if _, err := normalizeSettings(AISettings{Provider: ProviderAnthropic, EditorModel: "", TranslateModel: "y", MaxTokens: 1000}); err == nil {
		t.Error("empty editor model should be rejected")
	}

	// Clamp low and high.
	lo, err := normalizeSettings(AISettings{Provider: ProviderKimi, EditorModel: "a", TranslateModel: "b", MaxTokens: 10})
	if err != nil || lo.MaxTokens != 256 {
		t.Errorf("max_tokens floor: got %d err %v, want 256", lo.MaxTokens, err)
	}
	hi, err := normalizeSettings(AISettings{Provider: ProviderKimi, EditorModel: "a", TranslateModel: "b", MaxTokens: 999999})
	if err != nil || hi.MaxTokens != 32000 {
		t.Errorf("max_tokens ceiling: got %d err %v, want 32000", hi.MaxTokens, err)
	}

	// Whitespace is trimmed.
	got, err := normalizeSettings(AISettings{Provider: "  openai  ", EditorModel: "  gpt-5  ", TranslateModel: " m ", MaxTokens: 500})
	if err != nil || got.Provider != "openai" || got.EditorModel != "gpt-5" || got.TranslateModel != "m" {
		t.Errorf("trim: got %+v err %v", got, err)
	}
}

func TestAdminViewReflectsKeysAndActive(t *testing.T) {
	m := &Module{
		providers:  map[string]Completer{ProviderAnthropic: nil, ProviderOpenAI: nil},
		keyPresent: map[string]bool{ProviderAnthropic: true, ProviderOpenAI: true},
		settings:   &Settings{cur: AISettings{Enabled: true, Provider: ProviderOpenAI, EditorModel: "gpt-5", TranslateModel: "gpt-5-mini", MaxTokens: 2048}},
	}
	v := m.AdminView()

	if len(v.Providers) != 3 {
		t.Fatalf("want 3 providers in catalog order, got %d", len(v.Providers))
	}
	if v.Providers[0].Code != ProviderAnthropic || v.Providers[2].Code != ProviderKimi {
		t.Errorf("catalog order wrong: %+v", v.Providers)
	}
	// Kimi has no key here.
	if v.Providers[2].HasKey {
		t.Error("kimi should report no key")
	}
	// OpenAI is active and has a key.
	if !v.ActiveHasKey || v.Provider != ProviderOpenAI {
		t.Errorf("active provider/key wrong: active=%s hasKey=%v", v.Provider, v.ActiveHasKey)
	}
	for _, p := range v.Providers {
		if p.IsActive != (p.Code == ProviderOpenAI) {
			t.Errorf("IsActive wrong for %s", p.Code)
		}
	}
}

func TestClientsNilWhenDisabled(t *testing.T) {
	m := &Module{} // nothing enabled
	if c, _, _ := m.editorClient(); c != nil {
		t.Error("editorClient should be nil when disabled")
	}
	if c, _, _ := m.translateClient(); c != nil {
		t.Error("translateClient should be nil when disabled")
	}
	if c, _ := m.moderateClient(); c != nil {
		t.Error("moderateClient should be nil when disabled")
	}
}
