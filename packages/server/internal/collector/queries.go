package collector

import "fmt"

// kubernetesQueries builds PromQL queries using cAdvisor and kube-state-metrics labels.
// containerName maps to the Prometheus `container` label (usually the deployment/app name).
func kubernetesQueries(containerName, namespace, duration string) map[string]string {
	return map[string]string{
		"cpu": fmt.Sprintf(
			`rate(container_cpu_usage_seconds_total{namespace=%q,container=%q,container!=""}[1m])`,
			namespace, containerName,
		),
		"memory": fmt.Sprintf(
			`container_memory_working_set_bytes{namespace=%q,container=%q,container!=""}`,
			namespace, containerName,
		),
		"restarts": fmt.Sprintf(
			`increase(kube_pod_container_status_restarts_total{namespace=%q,container=%q}[%s])`,
			namespace, containerName, duration,
		),
		"oom": fmt.Sprintf(
			`kube_pod_container_status_last_terminated_reason{namespace=%q,container=%q,reason="OOMKilled"}`,
			namespace, containerName,
		),
	}
}

// processQueries builds PromQL queries using generic process_* metrics exposed by most runtimes.
// jobLabel maps to the Prometheus `job` label configured for the scrape target.
func processQueries(jobLabel string) map[string]string {
	return map[string]string{
		"cpu":    fmt.Sprintf(`rate(process_cpu_seconds_total{job=%q}[1m])`, jobLabel),
		"memory": fmt.Sprintf(`process_resident_memory_bytes{job=%q}`, jobLabel),
	}
}
