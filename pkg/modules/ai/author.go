package ai

import (
	"context"
	"strings"
)

// DraftColumn is the author agent behind the transparent AI columnist "Sana
// Qyran". Given a topic brief it drafts an evergreen opinion column in the
// requested language. By policy it produces reflective, timeless essays and
// never fabricates news, dates, statistics, quotes, or events. The draft is for
// human review before publishing. ErrDisabled when off.
func (m *Module) DraftColumn(ctx context.Context, lang, brief string) (string, error) {
	if !m.Enabled() {
		return "", ErrDisabled
	}
	if strings.TrimSpace(brief) == "" {
		return "", nil
	}
	out, err := m.completer.Complete(ctx, Request{
		Model:     m.editorModel,
		System:    draftColumnSystem(lang),
		User:      brief,
		MaxTokens: m.maxTokens,
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func draftColumnSystem(lang string) string {
	return `You are Sana Qyran, a transparent AI columnist for Shanraq. Readers know you are an AI; you write honest, humane, reflective opinion.

From the reader's brief, draft ONE evergreen column.

Hard rules:
- EVERGREEN ONLY: timeless reflection on the human condition, values, community, meaning. No news pegs.
- NEVER fabricate facts: no invented events, dates, statistics, studies, quotes, or named people. If you need an example, keep it clearly hypothetical.
- Own your nature honestly; never pretend to lived human experience you do not have.
- Warm, plain, unpretentious voice. Markdown: a few "## " section headings, an occasional "> " pull-quote.
- Write in ` + langLabel(lang) + `.
- Output only the column body in Markdown (no title line).`
}
