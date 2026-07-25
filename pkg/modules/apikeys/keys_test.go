package apikeys

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	plain, prefix, hash, err := generateKey()
	if err != nil {
		t.Fatalf("generateKey: %v", err)
	}
	if !strings.HasPrefix(plain, "sk_") {
		t.Errorf("key should start with sk_, got %q", plain)
	}
	if len(prefix) != 10 || prefix != plain[:10] {
		t.Errorf("prefix = %q, want first 10 of key", prefix)
	}
	if len(hash) != 64 { // sha256 hex
		t.Errorf("hash length = %d, want 64", len(hash))
	}
	if hashKey(plain) != hash {
		t.Error("stored hash must equal hashKey(plain)")
	}
	other, _, _, _ := generateKey()
	if other == plain {
		t.Error("two generated keys must differ")
	}
}

func TestExtractAPIKey(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-API-Key", "  sk_abc  ")
	if got := extractAPIKey(r); got != "sk_abc" {
		t.Errorf("X-API-Key: got %q", got)
	}
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("Authorization", "ApiKey sk_xyz")
	if got := extractAPIKey(r2); got != "sk_xyz" {
		t.Errorf("Authorization ApiKey: got %q", got)
	}
	r3 := httptest.NewRequest("GET", "/", nil)
	r3.Header.Set("Authorization", "Bearer sk_no")
	if got := extractAPIKey(r3); got != "" {
		t.Errorf("Bearer is not an api key, got %q", got)
	}
	if extractAPIKey(httptest.NewRequest("GET", "/", nil)) != "" {
		t.Error("no header → empty")
	}
}
