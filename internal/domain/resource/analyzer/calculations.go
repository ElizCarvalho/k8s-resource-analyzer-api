package analyzer

import (
	"fmt"
)

// Funções de cálculo de distribuição
func calculateCPUDistribution(usage float64) map[string]float64 {
	distribution := make(map[string]float64)
	usageMillicores := usage * 1000

	ranges := []struct {
		min, max float64
		label    string
	}{
		{0, 200, "0-200m"},
		{201, 400, "201-400m"},
		{401, 600, "401-600m"},
		{601, 800, "601-800m"},
		{801, 1000, "801-1000m"},
	}

	for _, r := range ranges {
		if usageMillicores >= r.min && usageMillicores <= r.max {
			distribution[r.label] = 100.0 // 100% dos pods nesta faixa
		} else {
			distribution[r.label] = 0.0
		}
	}

	return distribution
}

// Funções de consulta
func buildCPUHistoricalQuery(namespace, deployment string) string {
	return fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod=~"%s-.*"}[5m])) * 1000`, namespace, deployment)
}

func buildMemoryHistoricalQuery(namespace, deployment string) string {
	return fmt.Sprintf(`sum(container_memory_working_set_bytes{namespace="%s",pod=~"%s-.*"}) / (1024 * 1024)`, namespace, deployment)
}
