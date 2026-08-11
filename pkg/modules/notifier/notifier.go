package notifier

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"sort"
	"strings"

	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

// Mailer sends e-mail notifications.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

// HeaderMailer is the optional extension of Mailer for senders that need extra
// RFC 5322 headers — currently only the digest, which must carry
// List-Unsubscribe to stay out of Gmail's spam folder. Callers type-assert for
// it and fall back to plain Send, so the interface stays one method wide for
// every other caller.
type HeaderMailer interface {
	SendWithHeaders(ctx context.Context, to, subject, body string, headers map[string]string) error
}

// Module wires an SMTP-backed mailer based on configuration.
type Module struct {
	sender Mailer
	logger *zap.Logger
	cfg    config.SMTPConfig
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "notifier" }

func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.logger = rt.Logger
	m.cfg = rt.Config.Notifications.SMTP
	if m.cfg.Host == "" || m.cfg.From == "" {
		rt.Logger.Info("notifier: smtp disabled (host/from not configured)")
		return nil
	}
	if m.cfg.Port == 0 {
		m.cfg.Port = 587
	}
	m.sender = &smtpSender{cfg: m.cfg}
	rt.Logger.Info("notifier: smtp configured", zap.String("host", m.cfg.Host), zap.Int("port", m.cfg.Port))
	return nil
}

// Sender returns the configured mailer or nil if smtp is disabled.
func (m *Module) Sender() Mailer { return m.sender }

// Send allows the notifier module itself to satisfy the Mailer interface.
func (m *Module) Send(ctx context.Context, to, subject, body string) error {
	if m.sender == nil {
		return errors.New("mailer not configured")
	}
	return m.sender.Send(ctx, to, subject, body)
}

// SendWithHeaders satisfies HeaderMailer, forwarding to the SMTP sender.
func (m *Module) SendWithHeaders(ctx context.Context, to, subject, body string, headers map[string]string) error {
	s, ok := m.sender.(HeaderMailer)
	if !ok {
		return m.Send(ctx, to, subject, body)
	}
	return s.SendWithHeaders(ctx, to, subject, body, headers)
}

var _ interface {
	shanraq.Module
	shanraq.InitializerModule
} = (*Module)(nil)

type smtpSender struct {
	cfg config.SMTPConfig
}

func (s *smtpSender) Send(ctx context.Context, to, subject, body string) error {
	return s.SendWithHeaders(ctx, to, subject, body, nil)
}

func (s *smtpSender) SendWithHeaders(ctx context.Context, to, subject, body string, extra map[string]string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	// The From header keeps any display name ("Shanraq.org <no-reply@…>"), but
	// the SMTP envelope sender (MAIL FROM) must be a bare address.
	msg := buildMessage(s.cfg.From, to, subject, body, extra)

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	return smtp.SendMail(addr, auth, envelopeSender(s.cfg.From), []string{to}, msg)
}

// envelopeSender extracts the bare e-mail address from a possibly display-named
// From value ("Name <addr>" → "addr"), for use as the SMTP MAIL FROM.
func envelopeSender(from string) string {
	if a, err := mail.ParseAddress(from); err == nil {
		return a.Address
	}
	return strings.TrimSpace(from)
}

func buildMessage(from, to, subject, body string, extra map[string]string) []byte {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", scrubHeader(to)),
		// Non-ASCII subjects ("Shanraq.org: обзор недели") are not legal raw in
		// a header. Q-encoding them is what the spec asks for, and it also
		// makes header injection through a subject impossible.
		fmt.Sprintf("Subject: %s", mime.QEncoding.Encode("utf-8", subject)),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
	}
	// Sorted so the message is byte-identical for the same input — otherwise
	// map order would make the output untestable.
	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		headers = append(headers, fmt.Sprintf("%s: %s", scrubHeader(k), scrubHeader(extra[k])))
	}
	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body + "\r\n")
}

// scrubHeader strips CR/LF so a header value can never start a new header.
func scrubHeader(v string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(v)
}
