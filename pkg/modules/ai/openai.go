package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// openaiCompleter implements Completer against any OpenAI-compatible Chat
// Completions endpoint. Both OpenAI (ChatGPT) and Moonshot (Kimi K3) speak this
// same protocol, so one client covers both — only the base URL and model differ,
// which is exactly what makes switching providers a config change, not a code
// change. Implemented over net/http to avoid pulling in a second vendor SDK.
type openaiCompleter struct {
	apiKey  string
	baseURL string // e.g. https://api.openai.com/v1
	http    *http.Client

	// Pacing learned from refusals. The account this runs against allows three
	// requests a minute, and the provider's own Retry-After says one second —
	// which is true of the moment and false of the window, so following it
	// burns every attempt inside twenty seconds. Rather than trust the header
	// or hard-code a number that would be wrong for a better account, the
	// client starts unpaced and slows down when it is told to.
	mu      sync.Mutex
	spacing time.Duration
	next    time.Time
}

func newOpenAICompleter(apiKey, baseURL string) *openaiCompleter {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &openaiCompleter{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		// A ceiling, not the working limit: each request sets its own deadline
		// from how much text it asked for (see requestTimeout).
		http: &http.Client{Timeout: 20 * time.Minute},
	}
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens,omitempty"`
	Messages  []openaiMessage `json:"messages"`
	Thinking  *thinkingParam  `json:"thinking,omitempty"`
	// Effort is Kimi K3's depth dial. K3 always reasons and cannot be told to
	// stop, but it can be told how hard — and translation and moderation are not
	// what deep reasoning is for. Only K3 accepts the field.
	Effort string `json:"reasoning_effort,omitempty"`
}

// lowEffort asks K3 to think briefly. Its default is "max", which on a
// one-sentence job spent three and a half thousand tokens deliberating over
// word choice. Models that do not take the parameter get an empty string and
// the field is omitted.
func lowEffort(model string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "kimi-k3") {
		return "low"
	}
	return ""
}

// thinkingParam switches off a Kimi model's deliberation.
type thinkingParam struct {
	Type string `json:"type"`
}

// noThinking returns the parameter that turns deliberation off, for models that
// both support switching it off and have no use for it here.
//
// Kimi's k2.x models think by default, and reasoning is billed as output and
// counted against max_tokens. Measured on a one-sentence translation: 3,459
// reasoning tokens against 35 tokens of answer. Two consequences, and the
// second is worse than the first.
//
// The bill: a hundredfold on every call, for jobs — translating a sentence,
// returning a JSON verdict — where there is nothing to deliberate about.
//
// The correctness: reasoning shares the max_tokens budget with the answer, and
// a title is translated with a cap of 512. Thinking would have consumed the
// whole allowance and left an empty or truncated title, which the caller would
// then have saved as the translation.
//
// Only the models that document a disable switch get one. kimi-k3 always thinks
// and the docs say not to send the parameter; kimi-k2.7-code errors on it.
func noThinking(model string) *thinkingParam {
	m := strings.ToLower(strings.TrimSpace(model))
	if strings.HasPrefix(m, "kimi-k2.6") || strings.HasPrefix(m, "kimi-k2.5") {
		return &thinkingParam{Type: "disabled"}
	}
	return nil
}

type openaiResponse struct {
	Choices []struct {
		Message      openaiMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		CompletionTokens int `json:"completion_tokens"`
		Details          struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// requestTimeout gives a request as long as its answer plausibly needs.
//
// A flat two-minute timeout was fine while everything asked for was a headline
// or a JSON verdict. The first full-length article translation asked for nine
// thousand tokens, and the client hung up at 120 seconds while the model was
// still writing — the job failed, retried, and failed again the same way.
//
// Generation runs on the order of tens of tokens a second, so the allowance
// scales with the tokens requested, with a floor for short calls and a ceiling
// so a stuck request cannot hold a worker forever.
func requestTimeout(maxTokens int) time.Duration {
	d := 90*time.Second + time.Duration(maxTokens/15)*time.Second
	if d > 15*time.Minute {
		d = 15 * time.Minute
	}
	return d
}

// rateLimitAttempts is how many times a request refused for rate limiting is
// tried again.
//
// Translating by paragraph turned one request per article into seventy, and the
// provider noticed: a run of a single article hit "429 Too Many Requests"
// repeatedly, and an unretried 429 fails the whole translation over a queue
// that would have cleared in ten seconds. Waiting is the correct response to
// being asked to wait.
const rateLimitAttempts = 6

// rateLimitStep is how much a refusal slows the client down, and the floor for
// how long it waits before trying again. It is a fifth of a minute because the
// limits that bite here are per-minute request counts.
const rateLimitStep = 12 * time.Second

func (c *openaiCompleter) Complete(ctx context.Context, req Request) (string, error) {
	for attempt := 1; ; attempt++ {
		if err := c.waitTurn(ctx); err != nil {
			return "", err
		}
		out, err := c.complete(ctx, req)
		var limited rateLimited
		if !errors.As(err, &limited) {
			return out, err
		}
		if attempt == rateLimitAttempts {
			return "", err
		}
		wait := c.slowDown(limited.retryAfter)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(wait):
		}
	}
}

// waitTurn holds the request until the pace allows it. Requests are serialised
// deliberately: the same account that caps requests per minute caps concurrency
// at one, and two callers racing only produces refusals.
func (c *openaiCompleter) waitTurn(ctx context.Context) error {
	c.mu.Lock()
	spacing := c.spacing
	if spacing == 0 {
		c.mu.Unlock()
		return nil
	}
	now := time.Now()
	wait := time.Duration(0)
	if c.next.After(now) {
		wait = c.next.Sub(now)
	}
	c.next = now.Add(wait + spacing)
	c.mu.Unlock()

	if wait <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(wait):
		return nil
	}
}

// slowDown widens the pace after a refusal and says how long to wait now.
func (c *openaiCompleter) slowDown(retryAfter time.Duration) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.spacing < rateLimitStep {
		c.spacing = rateLimitStep
	} else if c.spacing < time.Minute {
		c.spacing += rateLimitStep
	}
	wait := c.spacing
	if retryAfter > wait {
		wait = retryAfter
	}
	c.next = time.Now().Add(wait + c.spacing)
	return wait
}

// rateLimited marks a refusal that will pass on its own.
type rateLimited struct {
	retryAfter time.Duration
	body       string
}

func (e rateLimited) Error() string {
	return "openai: rate limited (" + strings.TrimSpace(e.body) + ")"
}

func (c *openaiCompleter) complete(ctx context.Context, req Request) (string, error) {
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	ctx, cancel := context.WithTimeout(ctx, requestTimeout(maxTokens))
	defer cancel()

	msgs := make([]openaiMessage, 0, 2)
	if strings.TrimSpace(req.System) != "" {
		msgs = append(msgs, openaiMessage{Role: "system", Content: req.System})
	}
	msgs = append(msgs, openaiMessage{Role: "user", Content: req.User})

	body, err := json.Marshal(openaiRequest{
		Model: req.Model, MaxTokens: maxTokens, Messages: msgs,
		Thinking: noThinking(req.Model), Effort: lowEffort(req.Model),
	})
	if err != nil {
		return "", fmt.Errorf("openai encode: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("openai call: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", fmt.Errorf("openai read: %w", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		var after time.Duration
		if v := resp.Header.Get("Retry-After"); v != "" {
			if secs, convErr := strconv.Atoi(strings.TrimSpace(v)); convErr == nil && secs > 0 {
				after = time.Duration(secs) * time.Second
			}
		}
		return "", rateLimited{retryAfter: after, body: string(raw)}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out openaiResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("openai decode: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("openai: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("openai: empty response")
	}
	choice := out.Choices[0]
	text := strings.TrimSpace(choice.Message.Content)

	// A truncated answer is worse than none: the caller would save half a
	// translation as the translation. Fail instead, and say what ran out —
	// a reasoning model can spend the whole allowance before writing a word.
	if choice.FinishReason == "length" {
		if r := out.Usage.Details.ReasoningTokens; r > 0 {
			return "", fmt.Errorf("openai: hit the %d-token cap with %d of them spent on reasoning", maxTokens, r)
		}
		return "", fmt.Errorf("openai: hit the %d-token cap before finishing", maxTokens)
	}
	if text == "" {
		return "", fmt.Errorf("openai: the reply had no content")
	}
	return text, nil
}
