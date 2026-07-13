package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// claudeCompleter implements Completer against the Anthropic Messages API.
type claudeCompleter struct {
	client anthropic.Client
}

func newClaudeCompleter(apiKey string) *claudeCompleter {
	return &claudeCompleter{client: anthropic.NewClient(option.WithAPIKey(apiKey))}
}

func (c *claudeCompleter) Complete(ctx context.Context, req Request) (string, error) {
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(req.Model),
		MaxTokens: int64(maxTokens),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.User)),
		},
	}
	if strings.TrimSpace(req.System) != "" {
		params.System = []anthropic.TextBlockParam{{Text: req.System}}
	}

	resp, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("claude complete: %w", err)
	}

	var b strings.Builder
	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.TextBlock); ok {
			b.WriteString(text.Text)
		}
	}
	return strings.TrimSpace(b.String()), nil
}
