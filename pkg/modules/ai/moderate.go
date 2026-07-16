package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ModerationVerdict is the structured result of screening a piece of user text.
type ModerationVerdict struct {
	Action     string  `json:"action"`     // "allow" | "flag"
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
	if !m.Enabled() {
		return ModerationVerdict{}, ErrDisabled
	}
	if strings.TrimSpace(text) == "" {
		return ModerationVerdict{Action: "allow"}, nil
	}
	raw, err := m.completer.Complete(ctx, Request{
		Model:     m.moderateModel(),
		System:    moderateSystem(kind),
		User:      text,
		MaxTokens: 256,
	})
	if err != nil {
		return ModerationVerdict{}, err
	}
	return parseModerationVerdict(raw)
}

// moderateModel picks the cheap/fast tier for moderation (Haiku by default).
func (m *Module) moderateModel() string {
	if strings.TrimSpace(m.translateModel) != "" {
		return m.translateModel
	}
	return "claude-haiku-4-5"
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
		v.Reason = ""
	}
	return v, nil
}

func moderateSystem(kind string) string {
	target := "user comment"
	if kind == "listing" {
		target = "real-estate listing text"
	}
	return `You are a content-safety moderator for an independent Kazakhstani publishing and classifieds platform (content is in Kazakh, Russian, or English).

Judge the following ` + target + ` and decide whether it should be hidden pending human review. Flag it only when it clearly contains:
- spam, scams, or purely promotional junk;
- insults, harassment, threats, or hate speech toward a person or group;
- doxxing or sharing another person's private data without consent;
- calls to violence or plainly illegal activity.

Do NOT flag ordinary criticism, strong opinions, political disagreement, or blunt language — this is a free-speech platform and over-blocking is worse than under-blocking. When unsure, allow.

Respond with ONLY a JSON object and nothing else:
{"action":"allow"|"flag","reason":"<short reason in the text's language, empty when allow>","confidence":<number 0..1>}`
}
