package notifier

import (
	"strings"
	"testing"
)

func TestBuildMessage(t *testing.T) {
	m := string(buildMessage("noreply@shanraq.org", "user@example.kz", "Тема", "Текст письма"))
	for _, want := range []string{
		"From: noreply@shanraq.org",
		"To: user@example.kz",
		"Subject: Тема",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
	} {
		if !strings.Contains(m, want) {
			t.Errorf("message missing header %q", want)
		}
	}
	// Headers use CRLF and are separated from the body by a blank line.
	if !strings.Contains(m, "\r\n\r\nТекст письма\r\n") {
		t.Errorf("body separation wrong: %q", m)
	}
}
