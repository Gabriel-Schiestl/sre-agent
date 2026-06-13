package runner

import (
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
)

type Processor struct {}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(payload services.RunPayload) (services.AggregatedData, error) {
	if len(payload.JTLContent) == 0 {
		return services.AggregatedData{}, fmt.Errorf("no content provided")
	}

	records, err := parse(payload.JTLContent)
	if err != nil {
		return services.AggregatedData{}, fmt.Errorf("failed to parse JTL content: %w", err)
	}

	result := services.AggregatedData{
		TotalRequests: len(records),
	}

	//TODO: implement actual aggregation logic to populate the rest of the fields in result based on records.

	return result, nil
}