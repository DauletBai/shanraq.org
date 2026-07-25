package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAICompleterRoundTrip(t *testing.T) {
	var gotAuth, gotPath, gotModel, gotSystem, gotUser string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		var req openaiRequest
		_ = json.Unmarshal(body, &req)
		gotModel = req.Model
		for _, m := range req.Messages {
			switch m.Role {
			case "system":
				gotSystem = m.Content
			case "user":
				gotUser = m.Content
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  привет  "}}]}`))
	}))
	defer srv.Close()

	c := newOpenAICompleter("sk-test-123", srv.URL)
	out, err := c.Complete(context.Background(), Request{Model: "gpt-5", System: "sys", User: "hi", MaxTokens: 100})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if out != "привет" {
		t.Errorf("content = %q, want trimmed %q", out, "привет")
	}
	if gotAuth != "Bearer sk-test-123" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotPath != "/chat/completions" {
		t.Errorf("path = %q, want /chat/completions", gotPath)
	}
	if gotModel != "gpt-5" || gotSystem != "sys" || gotUser != "hi" {
		t.Errorf("payload wrong: model=%q system=%q user=%q", gotModel, gotSystem, gotUser)
	}
}

func TestOpenAICompleterHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer srv.Close()

	c := newOpenAICompleter("nope", srv.URL)
	_, err := c.Complete(context.Background(), Request{Model: "gpt-5", User: "hi"})
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Errorf("expected a 401 error, got %v", err)
	}
}

func TestOpenAICompleterDefaultBaseURL(t *testing.T) {
	c := newOpenAICompleter("k", "")
	if c.baseURL != "https://api.openai.com/v1" {
		t.Errorf("default base URL = %q", c.baseURL)
	}
	// Trailing slash is trimmed so path joins cleanly.
	if c2 := newOpenAICompleter("k", "https://x/v1/"); c2.baseURL != "https://x/v1" {
		t.Errorf("base URL not trimmed: %q", c2.baseURL)
	}
}

func TestFirstNonEmptyAndEnvKey(t *testing.T) {
	if got := firstNonEmpty("", "  ", "x", "y"); got != "x" {
		t.Errorf("firstNonEmpty = %q, want x", got)
	}
	if got := firstNonEmpty("", "   "); got != "" {
		t.Errorf("firstNonEmpty all-blank = %q, want empty", got)
	}

	t.Setenv("AI_TEST_KEY_A", "")
	t.Setenv("AI_TEST_KEY_B", "secret")
	if got := envKey("AI_TEST_KEY_A", "AI_TEST_KEY_B"); got != "secret" {
		t.Errorf("envKey = %q, want secret (first non-empty)", got)
	}
	if got := envKey("AI_TEST_MISSING_X"); got != "" {
		t.Errorf("envKey missing = %q, want empty", got)
	}
}
