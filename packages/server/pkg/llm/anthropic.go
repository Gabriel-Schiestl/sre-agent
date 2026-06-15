package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	inner *anthropic.Client
}

func NewAnthropic(apiKey string) Client {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicClient{inner: &c}
}

func (c *AnthropicClient) Complete(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	resp, err := c.inner.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{{
			Text: systemPrompt,
		}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("llm.Complete: %w", err)
	}

	for _, block := range resp.Content {
		if v, ok := block.AsAny().(anthropic.TextBlock); ok {
			return v.Text, nil
		}
	}
	return "", fmt.Errorf("llm.Complete: no text block in response")
}
