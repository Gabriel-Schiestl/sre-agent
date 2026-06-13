package services

import "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"

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
}

// Runner is the interface that the runner module must implement.
type Runner interface {
	Process(payload RunPayload) (AggregatedData, error)
}

// Analyst is the interface that the analyst module must implement.
type Analyst interface {
	Analyze(payload AnalysisPayload) (*types.Diagnosis, error)
}
