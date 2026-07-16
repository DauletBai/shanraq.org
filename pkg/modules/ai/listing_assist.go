package ai

import (
	"context"
	"strings"
)

// DescribeListing is the listing-assistant agent: given the structured facts of
// a property (deal/type, area, rooms, location, amenities) it drafts a clear,
// honest description in the seller's language for them to edit before posting.
// It never invents features that are not in the facts. ErrDisabled when off.
func (m *Module) DescribeListing(ctx context.Context, lang, facts string) (string, error) {
	if !m.Enabled() {
		return "", ErrDisabled
	}
	if strings.TrimSpace(facts) == "" {
		return "", nil
	}
	out, err := m.completer.Complete(ctx, Request{
		Model:     m.editorModel,
		System:    describeListingSystem(lang),
		User:      facts,
		MaxTokens: m.maxTokens,
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func describeListingSystem(lang string) string {
	return `You help a seller write the description for a real-estate listing on Shanraq, an independent Kazakhstani classifieds platform.

You are given the property's structured facts (deal type, property type, area, rooms, location, price, amenities). Write a clear, honest, easy-to-read description a buyer would trust.

Rules:
- Use ONLY the facts provided. Never invent rooms, amenities, condition, or nearby infrastructure that are not stated.
- No hype, no ALL-CAPS, no exclamation storms. Calm and factual.
- 2–4 short paragraphs. Lead with what matters most (rooms, area, location).
- Write in ` + langLabel(lang) + `.
- Output only the description text.`
}
