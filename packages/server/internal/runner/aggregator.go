package runner

import (
	"sort"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
)

func aggregateData(records []JTLRecord) services.AggregatedData {
	errorRate := calculateErrorRate(records)
	p50Ms := calculatePercentile(records, 50)
	p90Ms := calculatePercentile(records, 90)
	p99Ms := calculatePercentile(records, 99)
	errorsGrouped := groupErrors(records)

	return services.AggregatedData{
		TotalRequests: len(records),
		ErrorRate: errorRate,
		LatencyP50Ms: p50Ms,
		LatencyP90Ms: p90Ms,
		LatencyP99Ms: p99Ms,
		ErrorsByType: errorsGrouped,
	}
}

func calculateErrorRate(records []JTLRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	errorCount := 0
	for _, record := range records {
		if !record.Success {
			errorCount++
		}
	}
	return float64(errorCount) / float64(len(records))
}

func calculatePercentile(records []JTLRecord, percentile int) float64 {
	if len(records) == 0 {
		return 0
	}
	latencies := make([]int, len(records))
	for i, record := range records {
		latencies[i] = record.Latency
	}

	sort.Ints(latencies)

	index := int(float64(len(latencies)-1) * float64(percentile) / 100)
	return float64(latencies[index])
}

func groupErrors(records []JTLRecord) []services.ErrorGroup {
	errorMap := make(map[string]map[string]int)
	for _, record := range records {
		if !record.Success {
			if _, exists := errorMap[record.ResponseCode]; !exists {
				errorMap[record.ResponseCode] = make(map[string]int)
			}
			errorMap[record.ResponseCode][record.FailureMessage]++
		}
	}

	var errorGroups []services.ErrorGroup
	for code, messages := range errorMap {
		for msg, count := range messages {
			errorGroups = append(errorGroups, services.ErrorGroup{
				ResponseCode:   code,
				FailureMessage: msg,
				Count:          count,
			})
		}
	}

	return errorGroups
}