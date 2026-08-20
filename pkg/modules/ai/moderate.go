package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ModerationVerdict is the structured result of screening a piece of user text.
type ModerationVerdict struct {
	Action string `json:"action"` // "allow" | "flag"
	// Rule is which published rule was broken, by code. It exists so a hidden
	// comment can be answered with the rule it broke rather than with a mood:
	// the platform promises an appeal, and an appeal against "the model did not
	// like it" is not one.
	Rule       string  `json:"rule"`
	Reason     string  `json:"reason"`     // short reason when flagged (text's language)
	Confidence float64 `json:"confidence"` // 0..1
}

// Flagged reports whether the content should be hidden pending human review.
func (v ModerationVerdict) Flagged() bool { return v.Action == "flag" }

// Moderate screens a user comment or listing for spam, abuse, or policy
// violations and returns a verdict. It is deliberately conservative: ordinary
// criticism and strong opinions are allowed. ErrDisabled is returned when the
// assistant is off (no API key) — callers then fall back to human moderation.
func (m *Module) Moderate(ctx context.Context, kind, text string) (ModerationVerdict, error) {
	c, model := m.moderateClient()
	if c == nil {
		return ModerationVerdict{}, ErrDisabled
	}
	if strings.TrimSpace(text) == "" {
		return ModerationVerdict{Action: "allow"}, nil
	}
	raw, err := c.Complete(ctx, Request{
		Model:     model,
		System:    moderateSystem(kind),
		User:      text,
		MaxTokens: 256,
	})
	if err != nil {
		return ModerationVerdict{}, err
	}
	return parseModerationVerdict(raw)
}

// parseModerationVerdict extracts the verdict JSON, tolerating code fences or
// stray prose around it, and normalizes the action.
func parseModerationVerdict(raw string) (ModerationVerdict, error) {
	raw = strings.TrimSpace(raw)
	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			raw = raw[i : j+1]
		}
	}
	var v ModerationVerdict
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return ModerationVerdict{}, fmt.Errorf("parse moderation verdict: %w", err)
	}
	if v.Action != "flag" {
		v.Action = "allow"
		v.Rule = ""
		v.Reason = ""
	}
	return v, nil
}

func moderateSystem(kind string) string {
	target := "user comment"
	extra := ""
	if kind == "listing" {
		target = "real-estate listing text"
		extra = "- misleading: floor area, price, address or photographs that contradict the description, " +
			"or an agency presenting itself as the owner\n"
	}
	return `You are the content moderator for an independent Kazakhstani publishing and classifieds platform. The text may be in Kazakh, Russian or English.

Decide whether this ` + target + ` should be hidden pending human review. Hide it only when it clearly contains one of the following, and name the code:
- spam: unsolicited bulk promotion, link dumping, repeated advertising junk
- scam: fraud, fake earnings schemes, phishing, forged documents or diplomas offered for sale
- hidden_ad: promotion written to pass as a reader's own opinion or experience — an endorsement of a named business, product or service by someone who does not disclose a connection to it; an affiliate or referral link; a "review" whose real purpose is to sell
- prohibited_goods: drugs, weapons, counterfeit documents, and other goods or services that may not lawfully be advertised in Kazakhstan
- abuse: insults, harassment or threats directed at a person
- hatred: incitement to violence, or to enmity against a national, ethnic, religious or other group
- defamation: an accusation of a specific crime or wrongdoing against an identifiable person, stated as established fact
- personal_data: another person's private data — home address, telephone, documents — published without their consent
- adult: sexual content, or anything harmful to children
- illegal: other plainly unlawful content under the law of the Republic of Kazakhstan
` + extra + `
Allow everything else. This is a free-speech platform and over-blocking is worse than under-blocking. In particular, these are ALLOWED:
- criticism of the article, its author, officials, agencies or the state, however harsh or blunt;
- political disagreement, sarcasm, and strong opinion;
- saying that an institution lies or falsifies its figures — that is opinion, and it is defamation only when a named individual is accused of a specific crime as fact;
- quoting or describing something offensive in order to condemn, report or discuss it;
- a person publishing their own contact details;
- naming a business in passing, or answering someone's question about where to buy something — that is hidden_ad only when the text reads as an advertisement dressed as an opinion;
- rudeness and disrespectful tone on their own. Readers moderate tone here; the model does not.

When unsure, allow.

Respond with ONLY a JSON object and nothing else:
{"action":"allow"|"flag","rule":"<code from the list, empty when allow>","reason":"<short reason in the text's language, empty when allow>","confidence":<number 0..1>}`
}
