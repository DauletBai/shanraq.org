// Package sms delivers short text messages (one-time verification codes) through
// a Kazakhstani SMS aggregator. The provider is chosen by config so the app is
// never locked to one gateway: if one aggregator's onboarding stalls, switching
// is a single environment variable. With no provider configured the package
// returns a nil Client and the caller falls back to dev-logging the code.
package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config selects and credentials the SMS gateway. A single generic credential
// block serves both providers: Mobizon authenticates with APIKey; SMSC.kz with
// Login+Password. From is the alphanumeric sender name registered with the
// operators (optional; the aggregator's default is used when empty).
type Config struct {
	Provider string `mapstructure:"provider"` // "" (disabled) | "mobizon" | "smsc"
	From     string `mapstructure:"from"`
	APIKey   string `mapstructure:"api_key"`  // mobizon
	Login    string `mapstructure:"login"`    // smsc
	Password string `mapstructure:"password"` // smsc
	BaseURL  string `mapstructure:"base_url"` // optional override (tests / custom host)
	// Channel routes the message off the default SMS rail (smsc only): "telegram"
	// delivers verification codes through Telegram (tg=1) — no operator sender
	// name, far cheaper. Empty = plain SMS.
	Channel string `mapstructure:"channel"` // "" | "telegram"
}

// Client sends SMS via the configured provider. Its SendSMS method satisfies the
// auth module's SMSSender interface.
type Client struct {
	cfg  Config
	http *http.Client
}

// New builds a Client from config. It returns (nil, false) when SMS is not
// configured — the caller then dev-logs codes instead of sending. It returns an
// error when a provider is named but its credentials are incomplete, so a
// half-configured production deploy fails loudly rather than silently dropping
// verification codes.
func New(cfg Config) (*Client, bool, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "", "none", "off":
		return nil, false, nil
	case "mobizon":
		if cfg.APIKey == "" {
			return nil, false, fmt.Errorf("sms: mobizon provider needs SHANRAQ_SMS_API_KEY")
		}
	case "smsc":
		if cfg.Login == "" || cfg.Password == "" {
			return nil, false, fmt.Errorf("sms: smsc provider needs SHANRAQ_SMS_LOGIN and SHANRAQ_SMS_PASSWORD")
		}
	default:
		return nil, false, fmt.Errorf("sms: unknown provider %q (want mobizon|smsc)", cfg.Provider)
	}
	cfg.Provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	return &Client{cfg: cfg, http: &http.Client{Timeout: 12 * time.Second}}, true, nil
}

// SendSMS delivers text to phone. The phone is reduced to bare digits (no '+'),
// which both Kazakhstani gateways expect for local 77XXXXXXXXXX numbers.
func (c *Client) SendSMS(ctx context.Context, phone, text string) error {
	to := digitsOnly(phone)
	if to == "" {
		return fmt.Errorf("sms: empty recipient")
	}
	switch c.cfg.Provider {
	case "mobizon":
		return c.sendMobizon(ctx, to, text)
	case "smsc":
		return c.sendSMSC(ctx, to, text)
	default:
		return fmt.Errorf("sms: provider %q not wired", c.cfg.Provider)
	}
}

// digitsOnly strips '+' and separators, leaving the bare MSISDN.
func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// --- Mobizon (api.mobizon.kz) ---

func (c *Client) sendMobizon(ctx context.Context, to, text string) error {
	base := c.cfg.BaseURL
	if base == "" {
		base = "https://api.mobizon.kz"
	}
	q := url.Values{}
	q.Set("output", "json")
	q.Set("api", "v1")
	q.Set("apiKey", c.cfg.APIKey)
	q.Set("recipient", to)
	q.Set("text", text)
	if c.cfg.From != "" {
		q.Set("from", c.cfg.From)
	}
	endpoint := strings.TrimRight(base, "/") + "/service/message/sendSmsMessage?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sms mobizon: %w", err)
	}
	defer resp.Body.Close()
	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("sms mobizon: decode response: %w", err)
	}
	if out.Code != 0 {
		return fmt.Errorf("sms mobizon: code %d: %s", out.Code, out.Message)
	}
	return nil
}

// --- SMSC.kz ---

func (c *Client) sendSMSC(ctx context.Context, to, text string) error {
	base := c.cfg.BaseURL
	if base == "" {
		base = "https://smsc.kz"
	}
	q := url.Values{}
	q.Set("login", c.cfg.Login)
	q.Set("psw", c.cfg.Password)
	q.Set("phones", to)
	q.Set("mes", text)
	q.Set("fmt", "3") // JSON response
	q.Set("charset", "utf-8")
	if strings.EqualFold(c.cfg.Channel, "telegram") || strings.EqualFold(c.cfg.Channel, "tg") {
		q.Set("tg", "1") // deliver the code via Telegram instead of SMS (no sender name, cheaper)
	}
	if c.cfg.From != "" {
		q.Set("sender", c.cfg.From)
	}
	endpoint := strings.TrimRight(base, "/") + "/sys/send.php?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sms smsc: %w", err)
	}
	defer resp.Body.Close()
	var out struct {
		ID        int    `json:"id"`
		Cnt       int    `json:"cnt"`
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("sms smsc: decode response: %w", err)
	}
	if out.Error != "" {
		return fmt.Errorf("sms smsc: error %d: %s", out.ErrorCode, out.Error)
	}
	return nil
}
