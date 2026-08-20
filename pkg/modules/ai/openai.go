package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (c *openaiCompleter) Complete(ctx context.Context, req Request) (string, error) {
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
		Model: req.Model, MaxTokens: maxTokens, Messages: msgs, Thinking: noThinking(req.Model),
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
