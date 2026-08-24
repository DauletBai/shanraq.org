package ai

import (
	"context"
	"strings"
	"testing"
)

func TestParseProofVerdict(t *testing.T) {
	cases := []struct {
		name  string
		raw   string
		fix   bool
		fixed string
	}{
		{"plain", `{"fix":true,"fixed":"агентство","reason":""}`, true, "агентство"},
		{"fenced", "```json\n{\"fix\":true,\"fixed\":\"теңге\"}\n```", true, "теңге"},
		{"prose around it",
			`Sure! {"fix":false,"fixed":"","reason":"Слово написано верно."} Hope that helps.`,
			false, ""},
		// A model that answers "no" but fills in a word anyway must not have that
		// word used: the decision is the field that decides.
		{"refusal with a stray word", `{"fix":false,"fixed":"агентство","reason":"Верно."}`, false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := parseProofVerdict(c.raw)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if v.Fix != c.fix || v.Fixed != c.fixed {
				t.Errorf("got fix=%v fixed=%q, want %v %q", v.Fix, v.Fixed, c.fix, c.fixed)
			}
		})
	}
	if _, err := parseProofVerdict("I cannot help with that."); err == nil {
		t.Error("expected an error when there is no JSON at all")
	}
}

func TestProofreadDisabledWithoutKey(t *testing.T) {
	m := New()
	if _, err := m.Proofread(context.Background(), "ru", "para", "sentence", "word"); err != ErrDisabled {
		t.Errorf("got %v, want ErrDisabled", err)
	}
	if _, _, err := m.ProofreadRefute(context.Background(), "ru", "p", "s", "w", "f"); err != ErrDisabled {
		t.Errorf("refute: got %v, want ErrDisabled", err)
	}
}

func TestProofreadSendsTheContext(t *testing.T) {
	m := New()
	var seen Request
	m.setCompleter(&fakeCompleter{reply: func(r Request) string {
		seen = r
		return `{"fix":true,"fixed":"агентство","reason":""}`
	}})
	v, err := m.Proofread(context.Background(), "ru",
		"В прошлом году агенство отчиталось о росте.",
		"В прошлом году агенство отчиталось о росте.", "агенство")
	if err != nil {
		t.Fatalf("Proofread: %v", err)
	}
	if !v.Fix || v.Fixed != "агентство" {
		t.Errorf("verdict = %+v", v)
	}
	// The paragraph is what lets the model judge agreement and homophones, so it
	// has to actually be sent.
	for _, want := range []string{"PARAGRAPH:", "SENTENCE:", "WORD:", "агенство"} {
		if !strings.Contains(seen.User, want) {
			t.Errorf("request is missing %q:\n%s", want, seen.User)
		}
	}
	if !strings.Contains(seen.System, "Russian") {
		t.Error("the checker was not told which language it is reading")
	}
}

func TestProofreadRefuteBothWays(t *testing.T) {
	m := New()
	m.setCompleter(&fakeCompleter{reply: func(r Request) string {
		if strings.Contains(r.User, "агентство") {
			return `{"stands":true,"reason":""}`
		}
		return `{"stands":false,"reason":"Это имя собственное."}`
	}})
	stands, why, err := m.ProofreadRefute(context.Background(), "ru", "p", "s", "агенство", "агентство")
	if err != nil || !stands || why != "" {
		t.Errorf("a real fix did not survive: stands=%v why=%q err=%v", stands, why, err)
	}
	stands, why, err = m.ProofreadRefute(context.Background(), "ru", "p", "s", "Алматы", "Алмата")
	if err != nil {
		t.Fatalf("refute: %v", err)
	}
	if stands {
		t.Error("a proper name was allowed to be rewritten")
	}
	if why == "" {
		t.Error("a refusal came back with no reason for the reader")
	}
}

func TestLangName(t *testing.T) {
	for in, want := range map[string]string{
		"kz": "Kazakh", "kk": "Kazakh", "en": "English", "ru": "Russian", "": "Russian",
	} {
		if got := langName(in); got != want {
			t.Errorf("langName(%q) = %q, want %q", in, got, want)
		}
	}
}
