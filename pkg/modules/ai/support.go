package ai

import (
	"context"
	"strings"
)

// Answer is the support/consultant agent: it replies to a visitor's question
// about how the platform works (posting, pricing, listing lifecycle, disputes)
// grounded in a fixed knowledge base, in the visitor's own language. Escalation
// to a human is signalled to the caller by an empty reply. ErrDisabled when off.
func (m *Module) Answer(ctx context.Context, lang, question string) (string, error) {
	if !m.Enabled() {
		return "", ErrDisabled
	}
	if strings.TrimSpace(question) == "" {
		return "", nil
	}
	out, err := m.completer.Complete(ctx, Request{
		Model:     m.editorModel,
		System:    supportSystem(lang),
		User:      question,
		MaxTokens: m.maxTokens,
	})
	if err != nil {
		return "", err
	}
	// The model returns the ESCALATE token when the question is out of scope;
	// surface that to the caller as an empty reply (hand off to a human).
	if strings.EqualFold(strings.TrimSpace(out), "ESCALATE") {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

// supportKB is the ground truth the consultant may state. Keep it factual and in
// sync with the product; the model must not invent policies beyond it.
const supportKB = `Platform facts (the ONLY facts you may assert):
- Shanraq is an independent Kazakhstani portal: reader articles (KZ/RU/EN) plus a real-estate classifieds section.
- Publishing articles is free for registered subscribers who have agreed to the documents and tariffs. Optional paid services are the AI editor/translation/cover and listing promotion; banner advertising goes through an advertiser cabinet. Prices take effect only when paid billing launches, with at least 60 days' notice.
- Posting a listing is free; it stays active for 21 days (3 weeks), then the listing and all its data are permanently deleted. The owner is reminded 2 days before expiry and can extend (+21 days), raise it to the top once, or highlight it for 7 days.
- To post: register, open the Studio, and use "New article" or "New listing". Listings require an honest photo set — filtered/warped photos are forbidden and can be reported.
- Readers can up/down-vote articles; authors accumulate karma. Comments are moderated.
- Disputes between buyer and seller are settled directly between them; the platform only hosts listings and can hide ones that violate the rules after a report.`

func supportSystem(lang string) string {
	return `You are the support consultant for Shanraq, an independent Kazakhstani publishing and classifieds platform.

` + supportKB + `

Rules:
- Answer ONLY from the facts above. If the question is outside them or needs a human (account problems, payments, abuse reports, legal disputes), reply with the single token: ESCALATE
- Be brief, concrete, and friendly. No marketing fluff.
- Reply in ` + langLabel(lang) + `.`
}

// langLabel returns the human name of a language code, defaulting to Russian.
func langLabel(lang string) string {
	if name, ok := langFullName[lang]; ok {
		return name
	}
	return langFullName["ru"]
}
