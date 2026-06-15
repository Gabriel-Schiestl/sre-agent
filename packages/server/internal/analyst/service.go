package analyst

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/pkg/llm"
)

type llmResponse struct {
	ErrorPlan   []types.ErrorCategory `json:"errorPlan"`
	Bottlenecks []types.Bottleneck    `json:"bottlenecks"`
	NextSteps   []string              `json:"nextSteps"`
}

type Analyst struct {
	llm *llm.Client
}

func New(llmClient *llm.Client) *Analyst {
	return &Analyst{llm: llmClient}
}

func (a *Analyst) Analyze(payload services.AnalysisPayload) (*types.Diagnosis, error) {
	userMsg := buildUserMessage(payload)

	raw, err := a.llm.Complete(context.Background(), systemPrompt, userMsg)
	if err != nil {
		return nil, fmt.Errorf("analyst.Analyze: %w", err)
	}

	var parsed llmResponse
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("analyst.Analyze parse response: %w", err)
	}

	return types.NewDiagnosis(
		payload.Run.ID(),
		parsed.ErrorPlan,
		parsed.Bottlenecks,
		parsed.NextSteps,
		raw,
	), nil
}
