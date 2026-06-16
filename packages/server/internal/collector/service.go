package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	prom "github.com/Gabriel-Schiestl/sre-agent/packages/server/pkg/prometheus"
)

type Collector struct {
	client *prom.Client
}

func New(client *prom.Client) *Collector {
	return &Collector{client: client}
}

func (c *Collector) Collect(ctx context.Context, payload CollectPayload) (*PrometheusData, error) {
	duration := fmt.Sprintf("%.0fs", payload.EndTime.Sub(payload.StartTime).Seconds())
	step := 30 * time.Second

	var services []ServiceMetrics
	for _, ms := range payload.Microservices {
		if ms.PrometheusJobLabel() == nil {
			continue
		}
		metrics := c.collectService(ctx, ms.Name(), ms.PrometheusJobLabel(), ms.KubernetesNamespace(), payload.StartTime, payload.EndTime, step, duration)
		services = append(services, metrics)
	}

	return &PrometheusData{Services: services}, nil
}

func (c *Collector) collectService(
	ctx context.Context,
	name string,
	jobLabel *string,
	namespace *string,
	start, end time.Time,
	step time.Duration,
	duration string,
) ServiceMetrics {
	metrics := ServiceMetrics{MicroserviceName: name}

	if namespace != nil {
		// Kubernetes service: cAdvisor + kube-state-metrics
		queries := kubernetesQueries(*jobLabel, *namespace, duration)

		if series, err := c.client.QueryRange(ctx, queries["cpu"], start, end, step); err == nil && len(series) > 0 {
			metrics.CPURateCores = toPoints(series[0].Points)
		} else if err != nil {
			log.Printf("collector: %s cpu query: %v", name, err)
		}

		if series, err := c.client.QueryRange(ctx, queries["memory"], start, end, step); err == nil && len(series) > 0 {
			metrics.MemoryBytes = toPoints(series[0].Points)
		} else if err != nil {
			log.Printf("collector: %s memory query: %v", name, err)
		}

		if samples, err := c.client.Query(ctx, queries["restarts"], end); err == nil {
			metrics.RestartsDelta = int(sumSamples(samples))
		} else {
			log.Printf("collector: %s restarts query: %v", name, err)
		}

		if samples, err := c.client.Query(ctx, queries["oom"], end); err == nil {
			metrics.OOMKilledCount = len(samples)
		} else {
			log.Printf("collector: %s oom query: %v", name, err)
		}
	} else {
		// Non-Kubernetes service: generic process_* metrics
		queries := processQueries(*jobLabel)

		if series, err := c.client.QueryRange(ctx, queries["cpu"], start, end, step); err == nil && len(series) > 0 {
			metrics.ProcessCPURate = toPoints(series[0].Points)
		} else if err != nil {
			log.Printf("collector: %s process cpu query: %v", name, err)
		}

		if series, err := c.client.QueryRange(ctx, queries["memory"], start, end, step); err == nil && len(series) > 0 {
			metrics.ProcessMemBytes = toPoints(series[0].Points)
		} else if err != nil {
			log.Printf("collector: %s process memory query: %v", name, err)
		}
	}

	return metrics
}

func toPoints(src []prom.TimeseriesPoint) []TimeseriesPoint {
	out := make([]TimeseriesPoint, len(src))
	for i, p := range src {
		out[i] = TimeseriesPoint{Timestamp: p.Timestamp, Value: p.Value}
	}
	return out
}

func sumSamples(samples []prom.Sample) float64 {
	var total float64
	for _, s := range samples {
		total += s.Value
	}
	return total
}
