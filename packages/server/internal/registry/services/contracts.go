package services

import (
	"context"
	"time"

	collectorpkg "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/collector"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

// RunPayload is the input to the runner module.
type RunPayload struct {
	Run           *types.TestRun
	Suite         *types.Suite
	Microservices []*types.Microservice
	JTLContent    []byte
}

// AggregatedData is the output from the runner module.
type AggregatedData struct {
	TotalRequests   int
	ErrorRate       float64
	LatencyP50Ms    float64
	LatencyP90Ms    float64
	LatencyP99Ms    float64
	ErrorsByType    []ErrorGroup
	EndpointMetrics []EndpointMetric
	Timeline        []TimelinePoint
	StartTime       time.Time
	EndTime         time.Time
}

type ErrorGroup struct {
	ResponseCode   string
	FailureMessage string
	Count          int
}

type EndpointMetric struct {
	Label      string
	P50Ms      float64
	P90Ms      float64
	P99Ms      float64
	ErrorRate  float64
	TotalCalls int
}

type TimelinePoint struct {
	WindowStart int64
	ErrorCount  int
	TotalCount  int
	P99Ms       float64
}

// AnalysisPayload is the input to the analyst module.
type AnalysisPayload struct {
	Run           *types.TestRun
	Suite         *types.Suite
	Microservices []*types.Microservice
	Data          AggregatedData
	Prometheus    *collectorpkg.PrometheusData // nil when Prometheus is not configured
}

// CollectPayload is the input to the collector module.
type CollectPayload = collectorpkg.CollectPayload

// Runner is the interface that the runner module must implement.
type Runner interface {
	Process(payload RunPayload) (AggregatedData, error)
}

// Analyst is the interface that the analyst module must implement.
type Analyst interface {
	Analyze(payload AnalysisPayload) (*types.Diagnosis, error)
}

// Collector is the interface that the collector module must implement.
// It is optional — a nil Collector means Prometheus is not configured.
type Collector interface {
	Collect(ctx context.Context, payload collectorpkg.CollectPayload) (*collectorpkg.PrometheusData, error)
}
