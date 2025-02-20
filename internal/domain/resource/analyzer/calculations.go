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

func calculateMemoryDistribution(usageMi int64) map[string]float64 {
	distribution := make(map[string]float64)
	ranges := []struct {
		min, max int64
		label    string
	}{
		{0, 256, "0-256Mi"},
		{257, 512, "257-512Mi"},
		{513, 1024, "513-1024Mi"},
		{1025, 2048, "1025-2048Mi"},
		{2049, 4096, "2049-4096Mi"},
	}

	for _, r := range ranges {
		if usageMi >= r.min && usageMi <= r.max {
			distribution[r.label] = 100.0 // 100% dos pods nesta faixa
		} else {
			distribution[r.label] = 0.0
		}
	}
	return distribution
}

// Funções de cálculo de utilização
func calculateUtilization(usage, request int64) float64 {
	if request == 0 {
		return 0
	}
	result := (float64(usage) / float64(request)) * 100
	if result < 0 || result > 100 {
		return 0
	}
	return result
}

func calculatePodUtilization(running, total int) float64 {
	if total <= 0 || running < 0 || running > total {
		return 0
	}
	result := float64(running) / float64(total) * 100
	if result < 0 || result > 100 {
		return 0
	}
	return result
}

// Funções de consulta
func buildCPUHistoricalQuery(namespace, deployment string) string {
	return fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod=~"%s-.*"}[5m])) * 1000`, namespace, deployment)
}

func buildMemoryHistoricalQuery(namespace, deployment string) string {
	return fmt.Sprintf(`sum(container_memory_working_set_bytes{namespace="%s",pod=~"%s-.*"}) / (1024 * 1024)`, namespace, deployment)
}
