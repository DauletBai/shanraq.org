package syndicate

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestTelegramPayloadRoundTrip(t *testing.T) {
	id := uuid.New()
	raw, err := TelegramPayload(id)
	if err != nil {
		t.Fatalf("TelegramPayload: %v", err)
	}
	var p TelegramJobPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatalf("payload not valid JSON: %v", err)
	}
	if p.ArticleID != id.String() {
		t.Errorf("article id = %q, want %q", p.ArticleID, id.String())
	}
}

func TestBuildTelegramMessage(t *testing.T) {
	msg := buildTelegramMessage("Мир & <прогресс>", "Итоги <дня>", "https://shanraq.org/read/x?a=1&b=2")
	// HTML special characters from the DATA must be escaped.
	if !strings.Contains(msg, "&amp;") || !strings.Contains(msg, "&lt;прогресс&gt;") {
		t.Errorf("title not HTML-escaped: %q", msg)
	}
	// Structural markup and the fixed CTA remain.
	if !strings.Contains(msg, "📰 <b>") || !strings.Contains(msg, "Оқу · Читать") {
		t.Errorf("message structure missing: %q", msg)
	}
	// An empty summary is omitted (no stray blank block).
	noSummary := buildTelegramMessage("T", "  ", "u")
	if strings.Count(noSummary, "\n\n") != 1 { // only the CTA separator
		t.Errorf("empty summary should be skipped: %q", noSummary)
	}
}
