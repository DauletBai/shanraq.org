package notifier

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

// Mailer sends e-mail notifications.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
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

var _ interface {
	shanraq.Module
	shanraq.InitializerModule
} = (*Module)(nil)

type smtpSender struct {
	cfg config.SMTPConfig
}

func (s *smtpSender) Send(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	// The From header keeps any display name ("Shanraq <no-reply@…>"), but the
	// SMTP envelope sender (MAIL FROM) must be a bare address.
	msg := buildMessage(s.cfg.From, to, subject, body)

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

func buildMessage(from, to, subject, body string) []byte {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
	}
	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body + "\r\n")
}
