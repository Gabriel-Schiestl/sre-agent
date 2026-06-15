package llm

import "context"

type Client interface {
	Complete(ctx context.Context, systemPrompt, userMessage string) (string, error)
}
