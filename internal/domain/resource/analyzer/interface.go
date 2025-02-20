package analyzer

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// ResourceAnalyzer define a interface para análise de recursos
type ResourceAnalyzer interface {
	GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error)
	GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error)
	AnalyzeResources(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis
	CalculateCosts(ctx context.Context, current *types.CurrentMetrics, analysis *types.ResourceRecommendationAnalysis) (*types.CostAnalysis, error)
	GenerateAlerts(current *types.CurrentMetrics, historical *types.HistoricalMetrics) []types.Alert
}
