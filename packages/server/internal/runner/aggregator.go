package runner

import (
	"slices"
	"sort"
	"time"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
)

func aggregateData(records []JTLRecord) services.AggregatedData {
	errorRate := calculateErrorRate(records)
	p50Ms := calculatePercentile(records, 50)
	p90Ms := calculatePercentile(records, 90)
	p99Ms := calculatePercentile(records, 99)
	errorsGrouped := groupErrors(records)
	endpointMetrics := calculateEndpointMetrics(records)
	timelines := calculateTimelines(records)
	startTime, endTime := extractTimeRange(records)

	return services.AggregatedData{
		TotalRequests:   len(records),
		ErrorRate:       errorRate,
		LatencyP50Ms:    p50Ms,
		LatencyP90Ms:    p90Ms,
		LatencyP99Ms:    p99Ms,
		ErrorsByType:    errorsGrouped,
		EndpointMetrics: endpointMetrics,
		Timeline:        timelines,
		StartTime:       startTime,
		EndTime:         endTime,
	}
}

func extractTimeRange(records []JTLRecord) (start, end time.Time) {
	if len(records) == 0 {
		now := time.Now()
		return now, now
	}
	minTS, maxTS := records[0].TimeStamp, records[0].TimeStamp
	for _, r := range records[1:] {
		if r.TimeStamp < minTS {
			minTS = r.TimeStamp
		}
		if r.TimeStamp > maxTS {
			maxTS = r.TimeStamp
		}
	}
	return time.UnixMilli(minTS), time.UnixMilli(maxTS)
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

func calculateEndpointMetrics(records []JTLRecord) []services.EndpointMetric {
	groupedEndpoints := make(map[string][]JTLRecord)
	for _, rec := range records {
		if _, ok := groupedEndpoints[rec.Label]; !ok {
			groupedEndpoints[rec.Label] = []JTLRecord{}
		}
		groupedEndpoints[rec.Label] = append(groupedEndpoints[rec.Label], rec)
	}

	results := make([]services.EndpointMetric, len(groupedEndpoints))
	for _, records := range groupedEndpoints {
		results = append(results, aggregateEndpointData(records))
	}

	return results
}

func aggregateEndpointData(records []JTLRecord) (res services.EndpointMetric) {
	res.Label = records[0].Label
	res.TotalCalls = len(records)
	res.ErrorRate = calculateErrorRate(records)

	res.P50Ms = calculatePercentile(records, 50)
	res.P90Ms = calculatePercentile(records, 90)
	res.P99Ms = calculatePercentile(records, 99)

	return
}

func calculateTimelines(records []JTLRecord) []services.TimelinePoint {
	if len(records) == 0 {
		return nil
	}

	const windowSizeMs = int64(30_000)

	startTime := records[0].TimeStamp
	windows := make(map[int64][]JTLRecord)

	for _, rec := range records {
		idx := (rec.TimeStamp - startTime) / windowSizeMs
		windows[idx] = append(windows[idx], rec)
	}

	keys := make([]int64, 0, len(windows))
	for k := range windows {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	result := make([]services.TimelinePoint, 0, len(windows))
	for _, k := range keys {
		recs := windows[k]
		errorCount := 0
		for _, r := range recs {
			if !r.Success {
				errorCount++
			}
		}
		result = append(result, services.TimelinePoint{
			WindowStart: startTime + k*windowSizeMs,
			ErrorCount:  errorCount,
			TotalCount:  len(recs),
			P99Ms:       calculatePercentile(recs, 99),
		})
	}

	return result
}