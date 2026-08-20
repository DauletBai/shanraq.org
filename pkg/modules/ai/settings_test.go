package ai

import (
	"context"
	"strings"
	"testing"
)

func TestNormalizeSettings(t *testing.T) {
	base := AISettings{Provider: ProviderOpenAI, TranslateModel: "gpt-5-mini", MaxTokens: 4096}

	if _, err := normalizeSettings(base); err != nil {
		t.Fatalf("valid settings rejected: %v", err)
	}

	if _, err := normalizeSettings(AISettings{Provider: "grok", TranslateModel: "y", MaxTokens: 1000}); err == nil {
		t.Error("unknown provider should be rejected")
	}

	if _, err := normalizeSettings(AISettings{Provider: ProviderAnthropic, TranslateModel: "", MaxTokens: 1000}); err == nil {
		t.Error("empty translation model should be rejected")
	}

	// The moderation model may be empty: that means "the translation model".
	if _, err := normalizeSettings(AISettings{Provider: ProviderAnthropic, TranslateModel: "y",
		ModerateModel: "", MaxTokens: 1000}); err != nil {
		t.Errorf("an empty moderation model is allowed and means the translation model: %v", err)
	}

	// Clamp low and high.
	lo, err := normalizeSettings(AISettings{Provider: ProviderKimi, TranslateModel: "b", MaxTokens: 10})
	if err != nil || lo.MaxTokens != 256 {
		t.Errorf("max_tokens floor: got %d err %v, want 256", lo.MaxTokens, err)
	}
	hi, err := normalizeSettings(AISettings{Provider: ProviderKimi, TranslateModel: "b", MaxTokens: 999999})
	if err != nil || hi.MaxTokens != 32000 {
		t.Errorf("max_tokens ceiling: got %d err %v, want 32000", hi.MaxTokens, err)
	}

	// Whitespace is trimmed.
	got, err := normalizeSettings(AISettings{Provider: "  openai  ", TranslateModel: " m ", MaxTokens: 500})
	if err != nil || got.Provider != "openai" || got.TranslateModel != "m" {
		t.Errorf("trim: got %+v err %v", got, err)
	}
}

func TestAdminViewReflectsKeysAndActive(t *testing.T) {
	m := &Module{
		providers:  map[string]Completer{ProviderAnthropic: nil, ProviderOpenAI: nil},
		keyPresent: map[string]bool{ProviderAnthropic: true, ProviderOpenAI: true},
		settings:   &Settings{cur: AISettings{Enabled: true, Provider: ProviderOpenAI, TranslateModel: "gpt-5-mini", MaxTokens: 2048}},
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
	if c, _, _ := m.translateClient(); c != nil {
		t.Error("translateClient should be nil when disabled")
	}
	if c, _ := m.moderateClient(); c != nil {
		t.Error("moderateClient should be nil when disabled")
	}
}

// Every provider in the picker must carry the model to drive it with. An empty
// field is not a neutral default: the admin switches provider, the field clears,
// and whatever is typed in by hand is a guess that saves without complaint and
// fails at the first call.
func TestEveryProviderNamesItsModels(t *testing.T) {
	if len(providerCatalog) == 0 {
		t.Fatal("the provider catalog is empty")
	}
	for _, p := range providerCatalog {
		if p.Translate == "" {
			t.Errorf("%s (%s) names no model", p.Label, p.Code)
		}
	}
}

// The moonshot-v1 series is switched off on 31 August 2026. It is the obvious
// pick for a cheap translation model and would work for a fortnight.
func TestNoRetiredModels(t *testing.T) {
	for _, p := range providerCatalog {
		for _, m := range []string{p.Translate} {
			if strings.HasPrefix(m, "moonshot-v1") {
				t.Errorf("%s uses %q, which is being retired on 2026-08-31", p.Code, m)
			}
		}
	}
}

// The gate and the assistant are different questions, and the switch that
// governs them used to be one: turning the pre-publication check off to let
// authors publish also turned off the trilingual translation that makes this a
// trilingual publication, and turning translation back on reinstated the gate.
func TestReviewCheckIsSeparateFromTheAssistant(t *testing.T) {
	cases := []struct {
		name                  string
		enabled, reviewCheck  bool
		wantAssistant, wantOn bool
	}{
		{"everything off", false, false, false, false},
		{"assistant on, no gate", true, false, true, false},
		{"assistant on, gate on", true, true, true, true},
		// The gate cannot stand without the model that mans it: asking for a
		// check from a disabled assistant would hold every article for a
		// checker that is never going to run.
		{"gate on, assistant off", false, true, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := &Module{settings: NewSettings(nil, AISettings{
				Enabled: c.enabled, ReviewCheck: c.reviewCheck,
			})}
			m.enabled = c.enabled
			m.completer = stubCompleter{} // a configured provider
			if got := m.Enabled(); got != c.wantAssistant {
				t.Errorf("Enabled() = %v, want %v", got, c.wantAssistant)
			}
			if got := m.ReviewCheckEnabled(); got != c.wantOn {
				t.Errorf("ReviewCheckEnabled() = %v, want %v", got, c.wantOn)
			}
		})
	}
}

// stubCompleter stands in for a configured provider: the switches under test
// are about policy, not about whether a request would succeed.
type stubCompleter struct{}

func (stubCompleter) Complete(context.Context, Request) (string, error) { return "", nil }

// The check that reads every article and the model that drafts a column are
// different jobs with different bills. One setting for both forces a choice
// between paying frontier prices to classify and writing columns with the
// cheap model.
func TestModerationModelFallsBackToTheCheapTier(t *testing.T) {
	cases := []struct {
		name                             string
		translate, moderate, wantModerat string
	}{
		{"unset follows translation", "kimi-k2.6", "", "kimi-k2.6"},
		{"named model wins", "kimi-k2.6", "kimi-k3", "kimi-k3"},
		{"blank is not a model", "claude-haiku-4-5", "   ", "claude-haiku-4-5"},
		{"nothing configured at all", "", "", "claude-haiku-4-5"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := &Module{
				enabled:        true,
				completer:      stubCompleter{},
				translateModel: c.translate,
				moderateModel:  c.moderate,
			}
			_, got := m.moderateClient()
			if got != c.wantModerat {
				t.Errorf("moderateClient model = %q, want %q", got, c.wantModerat)
			}
		})
	}
}
