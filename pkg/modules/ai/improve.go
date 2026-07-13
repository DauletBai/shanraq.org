package ai

import (
	"context"
	"strings"
)

var improveLangName = map[string]string{
	"kz": "Kazakh (Қазақ тілі)",
	"ru": "Russian (Русский язык)",
	"en": "English",
}

// Improve rewrites a draft in the given language into a clearer, respectful yet
// professional journalistic tone, preserving meaning and Markdown. It returns
// the improved text. ErrDisabled is returned when the assistant is off.
func (m *Module) Improve(ctx context.Context, lang, draft string) (string, error) {
	if !m.Enabled() {
		return "", ErrDisabled
	}
	if strings.TrimSpace(draft) == "" {
		return draft, nil
	}
	return m.completer.Complete(ctx, Request{
		Model:     m.editorModel,
		System:    improveSystem(lang),
		User:      draft,
		MaxTokens: m.maxTokens,
	})
}

func improveSystem(lang string) string {
	name := improveLangName[lang]
	if name == "" {
		name = "the same language as the input"
	}
	return `You are a supportive editor for an independent journalism platform where ordinary people — not only trained journalists — publish articles.

Rewrite the user's draft so it reads clearly and professionally, in a respectful, measured tone. Specifically:
- Keep the author's meaning, facts, and intent unchanged. Do not invent facts or add claims.
- Improve structure, clarity, grammar, and flow. Fix awkward phrasing.
- Replace insulting, inflammatory, or defamatory wording with neutral, factual phrasing. Where a claim needs a source, keep it but phrase it as an assertion the author should support.
- Preserve all Markdown formatting (headings, lists, quotes, links).
- Write in ` + name + `.
- Output ONLY the improved article text — no preamble, no notes, no explanation of your changes.`
}
