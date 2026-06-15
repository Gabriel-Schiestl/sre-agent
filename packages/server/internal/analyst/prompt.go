package analyst

import (
	"fmt"
	"strings"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
)

const systemPrompt = `You are an SRE (Site Reliability Engineering) expert specialized in analyzing stress test results.
You receive aggregated performance data from a JMeter stress test along with context about the microservices under test.
Your task is to diagnose performance issues, identify bottlenecks, and provide actionable next steps.

You must respond with a single valid JSON object matching exactly this schema (no markdown, no extra text):
{
  "errorPlan": [
    {
      "category": "string — short name for the error category",
      "description": "string — what this error means and why it occurs",
      "occurrences": <integer>,
      "affectedEndpoints": ["string"],
      "severity": "low | medium | high | critical"
    }
  ],
  "bottlenecks": [
    {
      "microservice": "string — name of the microservice",
      "confidence": "low | medium | high",
      "hypotheses": [
        {
          "title": "string — concise hypothesis title",
          "evidence": "string — data points supporting this hypothesis",
          "priority": <integer starting at 1>
        }
      ]
    }
  ],
  "nextSteps": ["string — actionable recommendation"]
}`

func buildUserMessage(payload services.AnalysisPayload) string {
	var sb strings.Builder

	// Test profile
	sb.WriteString("## Test Profile\n")
	sb.WriteString(fmt.Sprintf("- Run name: %s\n", payload.Run.Name()))
	sb.WriteString(fmt.Sprintf("- Suite: %s — %s\n", payload.Suite.Name(), payload.Suite.Description()))
	sb.WriteString(fmt.Sprintf("- Virtual users: %d\n", payload.Run.VirtualUsers()))
	sb.WriteString(fmt.Sprintf("- Duration: %d seconds\n", payload.Run.DurationSeconds()))
	if n := payload.Run.Notes(); n != "" {
		sb.WriteString(fmt.Sprintf("- Notes: %s\n", n))
	}

	// Microservices
	sb.WriteString("\n## Microservices Under Test\n")
	for _, ms := range payload.Microservices {
		sb.WriteString(fmt.Sprintf("### %s (%s)\n", ms.Name(), ms.Language()))
		sb.WriteString(fmt.Sprintf("- Description: %s\n", ms.Description()))
		sb.WriteString(fmt.Sprintf("- Main endpoints: %s\n", strings.Join(ms.MainEndpoints(), ", ")))
		sb.WriteString(fmt.Sprintf("- Resources: CPU %s / Memory %s\n", ms.CPULimit(), ms.MemoryLimit()))
		sb.WriteString(fmt.Sprintf("- SLOs: latency P99 ≤ %d ms, error rate ≤ %.2f%%\n",
			ms.SLOLatencyP99Ms(), ms.SLOErrorRatePct()))
	}

	// Aggregated metrics
	d := payload.Data
	sb.WriteString("\n## Aggregated Metrics\n")
	sb.WriteString(fmt.Sprintf("- Total requests: %d\n", d.TotalRequests))
	sb.WriteString(fmt.Sprintf("- Error rate: %.2f%%\n", d.ErrorRate*100))
	sb.WriteString(fmt.Sprintf("- Latency P50: %.1f ms\n", d.LatencyP50Ms))
	sb.WriteString(fmt.Sprintf("- Latency P90: %.1f ms\n", d.LatencyP90Ms))
	sb.WriteString(fmt.Sprintf("- Latency P99: %.1f ms\n", d.LatencyP99Ms))

	// Errors by type
	if len(d.ErrorsByType) > 0 {
		sb.WriteString("\n## Errors by Type\n")
		for _, eg := range d.ErrorsByType {
			sb.WriteString(fmt.Sprintf("- [%s] %s — %d occurrences\n",
				eg.ResponseCode, eg.FailureMessage, eg.Count))
		}
	}

	// Endpoint metrics
	if len(d.EndpointMetrics) > 0 {
		sb.WriteString("\n## Endpoint Metrics\n")
		for _, em := range d.EndpointMetrics {
			sb.WriteString(fmt.Sprintf("- %s: P50=%.1fms P90=%.1fms P99=%.1fms error=%.2f%% calls=%d\n",
				em.Label, em.P50Ms, em.P90Ms, em.P99Ms, em.ErrorRate*100, em.TotalCalls))
		}
	}

	// Timeline
	if len(d.Timeline) > 0 {
		sb.WriteString("\n## Timeline (sampled windows)\n")
		for _, tp := range d.Timeline {
			sb.WriteString(fmt.Sprintf("- t=%d: errors=%d/%d P99=%.1fms\n",
				tp.WindowStart, tp.ErrorCount, tp.TotalCount, tp.P99Ms))
		}
	}

	sb.WriteString("\nAnalyze the data above and return the JSON diagnosis.")
	return sb.String()
}
