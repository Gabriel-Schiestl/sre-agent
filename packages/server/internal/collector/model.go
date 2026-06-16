package collector

import (
	"time"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/types"
)

type CollectPayload struct {
	Microservices []*types.Microservice
	StartTime     time.Time
	EndTime       time.Time
}

type PrometheusData struct {
	Services []ServiceMetrics
}

type ServiceMetrics struct {
	MicroserviceName string

	// Kubernetes metrics — populated when KubernetesNamespace is set on the microservice.
	CPURateCores   []TimeseriesPoint
	MemoryBytes    []TimeseriesPoint
	RestartsDelta  int
	OOMKilledCount int

	// Process metrics — populated when only PrometheusJobLabel is set (non-Kubernetes).
	ProcessCPURate  []TimeseriesPoint
	ProcessMemBytes []TimeseriesPoint
}

type TimeseriesPoint struct {
	Timestamp time.Time
	Value     float64
}
