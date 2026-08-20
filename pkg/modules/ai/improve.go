package ai

import "context"

// Check runs an arbitrary system+user prompt and returns the raw reply. The
// articles module uses it for the publication rules check; keeping the prompt
// on the caller's side means the rules live with the rules, not with the
// transport.
func (m *Module) Check(ctx context.Context, system, user string, maxTokens int) (string, error) {
	// The rules check is the same kind of call as comment screening: a short
	// verdict, run on everything. It belongs on the moderation tier, not on the
	// model that writes columns.
	c, model := m.moderateClient()
	if c == nil {
		return "", ErrDisabled
	}
	if maxTokens <= 0 {
		maxTokens = 2000
	}
	return c.Complete(ctx, Request{Model: model, System: system, User: user, MaxTokens: maxTokens})
}
