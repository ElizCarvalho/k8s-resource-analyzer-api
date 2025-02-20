package analyzer

import (
	"testing"
)

func TestCalculateCPUDistribution(t *testing.T) {
	tests := []struct {
		name string
		cpu  float64
		want map[string]float64
	}{
		{
			name: "CPU em 150m",
			cpu:  0.150,
			want: map[string]float64{
				"0-200m":    100.0,
				"201-400m":  0.0,
				"401-600m":  0.0,
				"601-800m":  0.0,
				"801-1000m": 0.0,
			},
		},
		{
			name: "CPU em 350m",
			cpu:  0.350,
			want: map[string]float64{
				"0-200m":    0.0,
				"201-400m":  100.0,
				"401-600m":  0.0,
				"601-800m":  0.0,
				"801-1000m": 0.0,
			},
		},
		{
			name: "CPU em 950m",
			cpu:  0.950,
			want: map[string]float64{
				"0-200m":    0.0,
				"201-400m":  0.0,
				"401-600m":  0.0,
				"601-800m":  0.0,
				"801-1000m": 100.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCPUDistribution(tt.cpu)
			if len(got) != len(tt.want) {
				t.Errorf("calculateCPUDistribution() returned %d ranges, expected %d", len(got), len(tt.want))
				for k, v := range tt.want {
					t.Errorf("calculateCPUDistribution() for range %s = %v, expected %v", k, got[k], v)
				}
			}
		})
	}
}

func TestBuildCPUHistoricalQuery(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		deployment string
		want       string
	}{
		{
			name:       "query básica",
			namespace:  "default",
			deployment: "nginx",
			want:       `sum(rate(container_cpu_usage_seconds_total{namespace="default",pod=~"nginx-.*"}[5m])) * 1000`,
		},
		{
			name:       "namespace e deployment com caracteres especiais",
			namespace:  "prod-env",
			deployment: "web-app",
			want:       `sum(rate(container_cpu_usage_seconds_total{namespace="prod-env",pod=~"web-app-.*"}[5m])) * 1000`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCPUHistoricalQuery(tt.namespace, tt.deployment)
			if got != tt.want {
				t.Errorf("buildCPUHistoricalQuery() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestBuildMemoryHistoricalQuery(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		deployment string
		want       string
	}{
		{
			name:       "query básica",
			namespace:  "default",
			deployment: "nginx",
			want:       `sum(container_memory_working_set_bytes{namespace="default",pod=~"nginx-.*"}) / (1024 * 1024)`,
		},
		{
			name:       "namespace e deployment com caracteres especiais",
			namespace:  "prod-env",
			deployment: "web-app",
			want:       `sum(container_memory_working_set_bytes{namespace="prod-env",pod=~"web-app-.*"}) / (1024 * 1024)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMemoryHistoricalQuery(tt.namespace, tt.deployment)
			if got != tt.want {
				t.Errorf("buildMemoryHistoricalQuery() = %v, expected %v", got, tt.want)
			}
		})
	}
}
