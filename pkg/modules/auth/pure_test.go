package auth

import (
	"encoding/base64"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// ---- token / request helpers ----

func TestGenerateSecureToken(t *testing.T) {
	a, err := generateSecureToken(32)
	if err != nil {
		t.Fatalf("generateSecureToken: %v", err)
	}
	raw, err := base64.RawURLEncoding.DecodeString(a)
	if err != nil {
		t.Fatalf("token is not base64url: %v", err)
	}
	if len(raw) != 32 {
		t.Errorf("decoded token = %d bytes, want 32", len(raw))
	}
	b, _ := generateSecureToken(32)
	if a == b {
		t.Error("two tokens should differ (entropy)")
	}
}

func TestHashHelpers(t *testing.T) {
	h := hashToken("hello")
	if len(h) != 64 { // sha256 hex
		t.Errorf("hashToken length = %d, want 64", len(h))
	}
	if hashToken("hello") != h {
		t.Error("hashToken must be deterministic")
	}
	if hashToken("hellO") == h {
		t.Error("different input must hash differently")
	}
	if k := hashKey("1.2.3.4"); len(k) != 16 { // first 8 bytes hex
		t.Errorf("hashKey length = %d, want 16", len(k))
	}
}

func TestIsJSONRequest(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if !isJSONRequest(r) {
		t.Error("application/json should be detected")
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if isJSONRequest(r) {
		t.Error("form content-type is not JSON")
	}
}

func TestClientIdentifier(t *testing.T) {
	// RemoteAddr host is used when there is no forwarded header.
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.9:5555"
	if got := clientIdentifier(r); got != "203.0.113.9" {
		t.Errorf("clientIdentifier(RemoteAddr) = %q, want 203.0.113.9", got)
	}
	// The first X-Forwarded-For hop wins when present (documents current
	// behavior — this header must only be trusted behind a sanitising proxy).
	r.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")
	if got := clientIdentifier(r); got != "198.51.100.7" {
		t.Errorf("clientIdentifier(XFF) = %q, want 198.51.100.7", got)
	}
}

// ---- in-memory rate limiter ----

func TestMemoryRateLimiterBurstAndKeys(t *testing.T) {
	rl := newMemoryRateLimiter(map[string]rateLimitRule{
		// One token per hour, burst of 2: two immediate calls, then blocked.
		"otp": {limit: rate.Every(time.Hour), burst: 2},
	})
	if !rl.Allow("otp", "user-a") || !rl.Allow("otp", "user-a") {
		t.Fatal("first two calls within burst must be allowed")
	}
	if rl.Allow("otp", "user-a") {
		t.Error("third call must be denied (burst exhausted)")
	}
	// A different key has its own bucket.
	if !rl.Allow("otp", "user-b") {
		t.Error("a different key must not be affected")
	}
	// An empty key collapses to a single shared "anonymous" bucket.
	if !rl.Allow("otp", "") {
		t.Error("first anonymous call allowed")
	}
}

func TestMemoryRateLimiterDefaultsAndUnlimited(t *testing.T) {
	// Unknown action falls back to the "default" rule (burst 20).
	rl := newMemoryRateLimiter(nil)
	for i := 0; i < 20; i++ {
		if !rl.Allow("totally-unknown-action", "k") {
			t.Fatalf("call %d under default burst should be allowed", i+1)
		}
	}
	if rl.Allow("totally-unknown-action", "k") {
		t.Error("21st call should exceed the default burst")
	}
	// A rule with a non-positive limit means unlimited.
	rl2 := newMemoryRateLimiter(map[string]rateLimitRule{"free": {limit: 0, burst: 0}})
	for i := 0; i < 100; i++ {
		if !rl2.Allow("free", "k") {
			t.Fatal("a zero-limit rule must always allow")
		}
	}
}

func TestDefaultRateLimitRules(t *testing.T) {
	rules := defaultRateLimitRules()
	for _, action := range []string{"default", "signin", "signup", "password_reset", "mfa_verify"} {
		if _, ok := rules[action]; !ok {
			t.Errorf("default rules missing %q", action)
		}
	}
	if rules["signup"].burst != 3 {
		t.Errorf("signup burst = %d, want 3", rules["signup"].burst)
	}
}

// ---- TOTP (pure parts: secret + otpauth URI) ----

func TestTOTPSecretAndURI(t *testing.T) {
	p := NewTOTPProvider(nil, "Shanraq", nil)
	u := User{ID: uuid.New(), Email: "author@example.kz"}

	secret, uri, err := p.generateSecret(u)
	if err != nil {
		t.Fatalf("generateSecret: %v", err)
	}
	if secret == "" {
		t.Error("secret must not be empty")
	}
	if !strings.HasPrefix(uri, "otpauth://totp/") {
		t.Errorf("uri should be an otpauth URI, got %q", uri)
	}

	formatted := p.formatURI(u, secret)
	if !strings.HasPrefix(formatted, "otpauth://totp/Shanraq:") {
		t.Errorf("formatURI prefix wrong: %q", formatted)
	}
	if !strings.Contains(formatted, "secret="+secret) {
		t.Error("formatURI must embed the secret")
	}
	if !strings.Contains(formatted, "issuer=Shanraq") {
		t.Error("formatURI must embed the issuer")
	}
	// Falls back to the user id when no email is set.
	anon := User{ID: uuid.New()}
	if !strings.Contains(p.formatURI(anon, secret), url.QueryEscape(anon.ID.String())) {
		t.Error("formatURI should use the id when email is blank")
	}
}

// ---- person name normalization ----

func TestNormalizePersonName(t *testing.T) {
	if NormalizePersonName("  Даулет  ") != "Даулет" {
		t.Error("NormalizePersonName should trim surrounding whitespace")
	}
	if NormalizePersonName("") != "" {
		t.Error("empty stays empty")
	}
}
