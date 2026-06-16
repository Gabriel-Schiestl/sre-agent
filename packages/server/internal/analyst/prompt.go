package analyst

import (
	"fmt"
	"strings"

	collectorpkg "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/collector"
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
	sb.WriteString("\n## Aggregated Metrics (JMeter)\n")
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
		sb.WriteString("\n## Timeline (30s windows)\n")
		for _, tp := range d.Timeline {
			sb.WriteString(fmt.Sprintf("- t=%d: errors=%d/%d P99=%.1fms\n",
				tp.WindowStart, tp.ErrorCount, tp.TotalCount, tp.P99Ms))
		}
	}

	// Prometheus infrastructure metrics
	if p := payload.Prometheus; p != nil && len(p.Services) > 0 {
		sb.WriteString("\n## Infrastructure Metrics (Prometheus)\n")
		for _, svc := range p.Services {
			sb.WriteString(fmt.Sprintf("### %s\n", svc.MicroserviceName))

			// Kubernetes cAdvisor metrics
			if len(svc.CPURateCores) > 0 {
				avg := avgPoints(svc.CPURateCores)
				peak := maxPoints(svc.CPURateCores)
				sb.WriteString(fmt.Sprintf("- CPU usage: avg=%.3f cores, peak=%.3f cores\n", avg, peak))
			}
			if len(svc.MemoryBytes) > 0 {
				avg := avgPoints(svc.MemoryBytes)
				peak := maxPoints(svc.MemoryBytes)
				sb.WriteString(fmt.Sprintf("- Memory: avg=%.1f MB, peak=%.1f MB\n", avg/1e6, peak/1e6))
			}
			if svc.RestartsDelta > 0 {
				sb.WriteString(fmt.Sprintf("- Pod restarts during test: %d\n", svc.RestartsDelta))
			}
			if svc.OOMKilledCount > 0 {
				sb.WriteString(fmt.Sprintf("- OOMKilled events: %d\n", svc.OOMKilledCount))
			}

			// Non-Kubernetes process metrics
			if len(svc.ProcessCPURate) > 0 {
				avg := avgPoints(svc.ProcessCPURate)
				peak := maxPoints(svc.ProcessCPURate)
				sb.WriteString(fmt.Sprintf("- Process CPU rate: avg=%.3f, peak=%.3f\n", avg, peak))
			}
			if len(svc.ProcessMemBytes) > 0 {
				avg := avgPoints(svc.ProcessMemBytes)
				peak := maxPoints(svc.ProcessMemBytes)
				sb.WriteString(fmt.Sprintf("- Process memory: avg=%.1f MB, peak=%.1f MB\n", avg/1e6, peak/1e6))
			}
		}
	}

	sb.WriteString("\nAnalyze the data above and return the JSON diagnosis.")
	return sb.String()
}

func avgPoints(pts []collectorpkg.TimeseriesPoint) float64 {
	if len(pts) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range pts {
		sum += p.Value
	}
	return sum / float64(len(pts))
}

func maxPoints(pts []collectorpkg.TimeseriesPoint) float64 {
	if len(pts) == 0 {
		return 0
	}
	max := pts[0].Value
	for _, p := range pts[1:] {
		if p.Value > max {
			max = p.Value
		}
	}
	return max
}
