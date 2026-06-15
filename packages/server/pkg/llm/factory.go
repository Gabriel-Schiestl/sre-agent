package llm

import "fmt"

// New returns a Client based on the provider name.
// Valid providers: "anthropic" (default), "ollama".
func New(provider, anthropicKey, ollamaURL, ollamaModel string) (Client, error) {
	switch provider {
	case "ollama":
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}
		if ollamaModel == "" {
			ollamaModel = "llama3.2"
		}
		return NewOllama(ollamaURL, ollamaModel), nil
	case "anthropic", "":
		if anthropicKey == "" {
			return nil, fmt.Errorf("llm: ANTHROPIC_API_KEY is required when LLM_PROVIDER=anthropic")
		}
		return NewAnthropic(anthropicKey), nil
	default:
		return nil, fmt.Errorf("llm: unknown provider %q (valid: anthropic, ollama)", provider)
	}
}
