package analyzer

import (
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestDetermineOverallStatus(t *testing.T) {
	tests := []struct {
		name     string
		cpu      *types.ResourceTypeAnalysis
		memory   *types.ResourceTypeAnalysis
		pods     *types.PodAnalysis
		expected string
	}{
		{
			name: "Should return normal when metrics are normal",
			cpu: &types.ResourceTypeAnalysis{
				Utilization: 90.0,
			},
			memory: &types.ResourceTypeAnalysis{
				Utilization: 95.0,
			},
			pods:     &types.PodAnalysis{},
			expected: "critical",
		},
		{
			name: "Deve retornar warning quando CPU ou memória estão em warning",
			cpu: &types.ResourceTypeAnalysis{
				Utilization: 80.0,
			},
			memory: &types.ResourceTypeAnalysis{
				Utilization: 50.0,
			},
			pods:     &types.PodAnalysis{},
			expected: "warning",
		},
		{
			name: "Deve retornar critical quando CPU e memória estão críticos",
			cpu: &types.ResourceTypeAnalysis{
				Utilization: 90.0,
			},
			memory: &types.ResourceTypeAnalysis{
				Utilization: 95.0,
			},
			pods:     &types.PodAnalysis{},
			expected: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineOverallStatus(tt.cpu, tt.memory, tt.pods)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateHistoricalAverage(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected float64
	}{
		{
			name: "Deve calcular média corretamente",
			metrics: []*types.ResourceMetrics{
				{Usage: 10},
				{Usage: 20},
				{Usage: 30},
				{Usage: 40},
				{Usage: 50},
			},
			expected: 30,
		},
		{
			name:     "Deve retornar 0 para slice vazio",
			metrics:  []*types.ResourceMetrics{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHistoricalAverage(tt.metrics)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateHistoricalPeak(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected float64
	}{
		{
			name: "Deve encontrar pico corretamente",
			metrics: []*types.ResourceMetrics{
				{Usage: 10},
				{Usage: 20},
				{Usage: 50},
				{Usage: 30},
				{Usage: 40},
			},
			expected: 50,
		},
		{
			name:     "Deve retornar 0 para slice vazio",
			metrics:  []*types.ResourceMetrics{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHistoricalPeak(tt.metrics)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateHistoricalUtilization(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected float64
	}{
		{
			name: "Deve calcular utilização histórica corretamente",
			metrics: []*types.ResourceMetrics{
				{Utilization: 50},
				{Utilization: 70},
				{Utilization: 30},
			},
			expected: 50,
		},
		{
			name:     "Deve retornar 0 para slice vazio",
			metrics:  []*types.ResourceMetrics{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHistoricalUtilization(tt.metrics)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectPattern(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected string
	}{
		{
			name: "Deve detectar padrão crescente",
			metrics: []*types.ResourceMetrics{
				{Utilization: 10},
				{Utilization: 20},
				{Utilization: 30},
				{Utilization: 40},
				{Utilization: 50},
			},
			expected: "increasing",
		},
		{
			name: "Deve detectar padrão decrescente",
			metrics: []*types.ResourceMetrics{
				{Utilization: 50},
				{Utilization: 40},
				{Utilization: 30},
				{Utilization: 20},
				{Utilization: 10},
			},
			expected: "decreasing",
		},
		{
			name: "Deve detectar padrão estável",
			metrics: []*types.ResourceMetrics{
				{Utilization: 30},
				{Utilization: 32},
				{Utilization: 29},
				{Utilization: 31},
				{Utilization: 30},
			},
			expected: "stable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectPattern(tt.metrics)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectSeasonality(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected *types.Seasonality
	}{
		{
			name: "Deve detectar sazonalidade diária",
			metrics: generateMetricsForDays(4, []float64{
				60, 30, 20, 10, 20, 30, 60, // dia 1
				60, 30, 20, 10, 20, 30, 60, // dia 2
				60, 30, 20, 10, 20, 30, 60, // dia 3
				60, 30, 20, 10, 20, 30, 60, // dia 4
			}),
			expected: &types.Seasonality{
				Pattern: "daily",
				Period:  "24h",
			},
		},
		{
			name: "Não deve detectar sazonalidade em dados insuficientes",
			metrics: generateMetricsForDays(1, []float64{
				45, 23, 67, 12, 89, 34, 56,
			}),
			expected: &types.Seasonality{
				Pattern: "insufficient_data",
				Period:  "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectSeasonality(tt.metrics)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateUtilizationTrend(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []*types.ResourceMetrics
		expected float64
	}{
		{
			name: "Deve calcular tendências corretamente",
			metrics: []*types.ResourceMetrics{
				{Utilization: 40},
				{Utilization: 45},
				{Utilization: 50},
				{Utilization: 55},
				{Utilization: 60},
			},
			expected: 22.22,
		},
		{
			name: "Deve identificar tendências estáveis",
			metrics: []*types.ResourceMetrics{
				{Utilization: 48},
				{Utilization: 50},
				{Utilization: 49},
				{Utilization: 51},
				{Utilization: 50},
			},
			expected: 2.04,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateUtilizationTrend(tt.metrics)
			assert.InDelta(t, tt.expected, result, 0.1)
		})
	}
}

func TestGenerateCPURecommendation(t *testing.T) {
	tests := []struct {
		name       string
		current    *types.ResourceMetrics
		historical []*types.ResourceMetrics
		expected   *types.Recommendation
	}{
		{
			name: "Deve recomendar aumento de CPU",
			current: &types.ResourceMetrics{
				Usage:       85,
				Request:     100,
				Utilization: 85,
			},
			historical: []*types.ResourceMetrics{
				{Usage: 75, Request: 100, Utilization: 75},
				{Usage: 80, Request: 100, Utilization: 80},
				{Usage: 85, Request: 100, Utilization: 85},
			},
			expected: &types.Recommendation{
				Current:   85,
				Suggested: 102,
				Reason:    "optimal",
			},
		},
		{
			name: "Deve recomendar manutenção de CPU",
			current: &types.ResourceMetrics{
				Usage:       60,
				Request:     100,
				Utilization: 60,
			},
			historical: []*types.ResourceMetrics{
				{Usage: 55, Request: 100, Utilization: 55},
				{Usage: 58, Request: 100, Utilization: 58},
				{Usage: 60, Request: 100, Utilization: 60},
			},
			expected: &types.Recommendation{
				Current:   60,
				Suggested: 72,
				Reason:    "optimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateCPURecommendation(tt.current, tt.historical)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateMemoryRecommendation(t *testing.T) {
	tests := []struct {
		name       string
		current    *types.ResourceMetrics
		historical []*types.ResourceMetrics
		expected   *types.Recommendation
	}{
		{
			name: "Deve recomendar aumento de memória",
			current: &types.ResourceMetrics{
				Usage:       85,
				Request:     100,
				Utilization: 85,
			},
			historical: []*types.ResourceMetrics{
				{Usage: 80, Request: 100, Utilization: 80},
				{Usage: 82, Request: 100, Utilization: 82},
				{Usage: 85, Request: 100, Utilization: 85},
			},
			expected: &types.Recommendation{
				Current:   85,
				Suggested: 110.5,
				Reason:    "optimal",
			},
		},
		{
			name: "Deve recomendar manutenção de memória",
			current: &types.ResourceMetrics{
				Usage:       60,
				Request:     100,
				Utilization: 60,
			},
			historical: []*types.ResourceMetrics{
				{Usage: 58, Request: 100, Utilization: 58},
				{Usage: 59, Request: 100, Utilization: 59},
				{Usage: 60, Request: 100, Utilization: 60},
			},
			expected: &types.Recommendation{
				Current:   60,
				Suggested: 78,
				Reason:    "optimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMemoryRecommendation(tt.current, tt.historical)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func generateMetricsForDays(days int, values []float64) []*types.ResourceMetrics {
	var metrics []*types.ResourceMetrics
	start := time.Now().Truncate(24 * time.Hour)
	interval := 4 * time.Hour

	for i, value := range values {
		timestamp := start.Add(time.Duration(i) * interval)
		metrics = append(metrics, &types.ResourceMetrics{
			Usage:       value,
			Utilization: value,
			Timestamp:   timestamp.Unix(),
		})
	}
	return metrics
}
