package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

// Kimi's k2.x models deliberate by default, and the reasoning is billed as
// output and counted against max_tokens. A one-sentence translation measured
// 3,459 reasoning tokens against 35 of answer — a hundredfold on the bill, and,
// because a title is translated with a 512-token cap, an empty title saved as
// the translation.
func TestKimiThinkingIsSwitchedOffWhereItCanBe(t *testing.T) {
	cases := []struct {
		model string
		want  string // "" = the parameter must not be sent at all
	}{
		{"kimi-k2.6", "disabled"},
		{"kimi-k2.5", "disabled"},
		{"KIMI-K2.6", "disabled"}, // the operator types it by hand
		{"  kimi-k2.6  ", "disabled"},
		// Always thinks, and the documentation says not to send the parameter.
		{"kimi-k3", ""},
		// Errors when asked to stop thinking.
		{"kimi-k2.7-code", ""},
		{"kimi-k2.7-code-highspeed", ""},
		// Not Kimi at all: the field is not part of the OpenAI API.
		{"gpt-5.6-luna", ""},
		{"claude-haiku-4-5", ""},
	}
	for _, c := range cases {
		t.Run(c.model, func(t *testing.T) {
			var sent openaiRequest
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(body, &sent)
				// Also assert the raw shape: an empty struct would serialise as
				// {"type":""} and be read as a request to think.
				if c.want == "" && strings.Contains(string(body), `"thinking"`) {
					t.Errorf("thinking was sent to %s: %s", c.model, body)
				}
				_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"},"finish_reason":"stop"}]}`))
			}))
			defer srv.Close()

			c2 := newOpenAICompleter("sk-test", srv.URL)
			if _, err := c2.Complete(context.Background(), Request{Model: c.model, User: "hi"}); err != nil {
				t.Fatalf("Complete: %v", err)
			}
			got := ""
			if sent.Thinking != nil {
				got = sent.Thinking.Type
			}
			if got != c.want {
				t.Errorf("thinking = %q, want %q", got, c.want)
			}
		})
	}
}

// Half a translation saved as the translation is worse than a failed job: the
// job retries, the bad text does not announce itself.
func TestTruncatedRepliesFailInsteadOfBeingSaved(t *testing.T) {
	cases := []struct {
		name, body, wantErr string
	}{
		{
			"nothing left after reasoning",
			`{"choices":[{"message":{"content":""},"finish_reason":"length"}],
			  "usage":{"completion_tokens":512,"completion_tokens_details":{"reasoning_tokens":512}}}`,
			"spent on reasoning",
		},
		{
			"cut off mid-answer",
			`{"choices":[{"message":{"content":"Тақыр"},"finish_reason":"length"}]}`,
			"before finishing",
		},
		{
			"answered nothing at all",
			`{"choices":[{"message":{"content":"   "},"finish_reason":"stop"}]}`,
			"no content",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(c.body))
			}))
			defer srv.Close()

			out, err := newOpenAICompleter("sk-test", srv.URL).
				Complete(context.Background(), Request{Model: "kimi-k2.6", User: "hi", MaxTokens: 512})
			if err == nil {
				t.Fatalf("a truncated reply was returned as an answer: %q", out)
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Errorf("error %q does not mention %q", err, c.wantErr)
			}
			if out != "" {
				t.Errorf("partial text was handed back anyway: %q", out)
			}
		})
	}
}

// The first article translation asked for nine thousand tokens and the client
// hung up at two minutes while the model was still writing. A flat timeout is
// the same mistake as a flat token cap: fine for a headline, fatal for an
// article.
func TestRequestTimeoutScalesWithTheAnswer(t *testing.T) {
	cases := []struct {
		name      string
		maxTokens int
		want      time.Duration
	}{
		{"заголовок получает минимум", 512, 124 * time.Second},
		{"статье хватает на длинный ответ", 9517, 724 * time.Second},
		{"потолок не превышается", 100000, 15 * time.Minute},
		{"нулевой запрос не даёт нулевого срока", 0, 90 * time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := requestTimeout(c.maxTokens); got != c.want {
				t.Errorf("requestTimeout(%d) = %v, ожидалось %v", c.maxTokens, got, c.want)
			}
		})
	}
	// Короткий вызов не должен получать столько же времени, сколько статья:
	// зависший запрос обязан отваливаться быстро.
	if requestTimeout(512) >= requestTimeout(9517) {
		t.Error("короткий запрос ждёт не меньше длинного — срок не зависит от объёма")
	}
}

// The request has to honour the caller's context, not just the client's own
// ceiling — otherwise a shutdown or a cancelled job leaves calls hanging.
func TestCompleteHonoursTheCallerDeadline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second) // дольше, чем срок вызывающей стороны
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := newOpenAICompleter("sk-test", srv.URL).
		Complete(ctx, Request{Model: "kimi-k2.6", User: "hi", MaxTokens: 512})
	if err == nil {
		t.Fatal("вызов пережил отменённый контекст")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("возврат занял %v — контекст вызывающей стороны не соблюдён", elapsed)
	}
}

// K3 always reasons and cannot be switched off, but it takes a depth dial that
// K2.6 does not. Its default is "max" — three and a half thousand tokens of
// deliberation over one sentence — and translation is not what that is for.
func TestOnlyK3GetsTheEffortDial(t *testing.T) {
	for model, want := range map[string]string{
		"kimi-k3":        "low",
		"KIMI-K3":        "low",
		"kimi-k2.6":      "", // не поддерживает поле — отправка была бы ошибкой
		"kimi-k2.7-code": "",
		"gpt-5.6-luna":   "",
	} {
		if got := lowEffort(model); got != want {
			t.Errorf("lowEffort(%q) = %q, ожидалось %q", model, got, want)
		}
	}
}
