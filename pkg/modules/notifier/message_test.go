package notifier

import (
	"mime"
	"strings"
	"testing"
)

func TestBuildMessage(t *testing.T) {
	m := string(buildMessage("noreply@shanraq.org", "user@example.kz", "Тема", "Текст письма", nil))
	for _, want := range []string{
		"From: noreply@shanraq.org",
		"To: user@example.kz",
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

// A Cyrillic subject is not legal raw in a header: it must go out as an RFC 2047
// encoded-word, or clients show mojibake and spam filters mark the message down.
func TestBuildMessageEncodesSubject(t *testing.T) {
	m := string(buildMessage("noreply@shanraq.org", "user@example.kz", "Тема", "тело", nil))
	if strings.Contains(m, "Subject: Тема") {
		t.Error("subject sent as raw UTF-8")
	}
	want := "Subject: " + mime.QEncoding.Encode("utf-8", "Тема")
	if !strings.Contains(m, want) {
		t.Errorf("subject not q-encoded: %q", m)
	}
	// A decoder must get the original text back.
	got, err := new(mime.WordDecoder).DecodeHeader(mime.QEncoding.Encode("utf-8", "Тема"))
	if err != nil || got != "Тема" {
		t.Errorf("round-trip failed: %q, %v", got, err)
	}
}

// The digest carries List-Unsubscribe; without it Gmail scores bulk mail as spam.
func TestBuildMessageExtraHeaders(t *testing.T) {
	m := string(buildMessage("noreply@shanraq.org", "user@example.kz", "Digest", "body", map[string]string{
		"List-Unsubscribe":      "<https://shanraq.org/unsubscribe?token=abc>",
		"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
	}))
	for _, want := range []string{
		"List-Unsubscribe: <https://shanraq.org/unsubscribe?token=abc>",
		"List-Unsubscribe-Post: List-Unsubscribe=One-Click",
	} {
		if !strings.Contains(m, want) {
			t.Errorf("missing %q in %q", want, m)
		}
	}
	// Extra headers must stay in the header block, above the blank line.
	head, _, ok := strings.Cut(m, "\r\n\r\n")
	if !ok || !strings.Contains(head, "List-Unsubscribe:") {
		t.Error("extra headers leaked into the body")
	}
}

// A newline in a header value would let a caller inject headers of its own. The
// text may survive inside the value it was smuggled in; what must not survive is
// a line of its own, because only that makes it a header.
func TestBuildMessageRejectsHeaderInjection(t *testing.T) {
	m := string(buildMessage("noreply@shanraq.org", "user@example.kz\r\nBcc: evil@example.com",
		"Subject", "body", map[string]string{"X-Test": "a\r\nBcc: also-evil@example.com"}))
	head, _, _ := strings.Cut(m, "\r\n\r\n")
	for _, line := range strings.Split(head, "\r\n") {
		if strings.HasPrefix(line, "Bcc:") {
			t.Errorf("header injection got through: %q", m)
		}
	}
	// Five fixed headers plus the one X-Test above — no sixth smuggled in.
	if n := strings.Count(head, "\r\n") + 1; n != 6 {
		t.Errorf("expected 6 header lines, got %d: %q", n, head)
	}
}
