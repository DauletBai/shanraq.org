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
		http:    &http.Client{Timeout: 120 * time.Second},
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
}

type openaiResponse struct {
	Choices []struct {
		Message openaiMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *openaiCompleter) Complete(ctx context.Context, req Request) (string, error) {
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	msgs := make([]openaiMessage, 0, 2)
	if strings.TrimSpace(req.System) != "" {
		msgs = append(msgs, openaiMessage{Role: "system", Content: req.System})
	}
	msgs = append(msgs, openaiMessage{Role: "user", Content: req.User})

	body, err := json.Marshal(openaiRequest{Model: req.Model, MaxTokens: maxTokens, Messages: msgs})
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
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}
