package collector

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// Collector define a interface para coleta de métricas
type Collector interface {
	// GetDeploymentMetrics retorna métricas atuais de um deployment
	GetDeploymentMetrics(ctx context.Context, namespace, deployment string) (*types.K8sMetrics, error)

	// GetDeploymentConfig retorna configurações de um deployment
	GetDeploymentConfig(ctx context.Context, namespace, deployment string) (*types.K8sDeploymentConfig, error)

	// Query executa uma query pontual
	Query(ctx context.Context, query string) (*types.QueryResult, error)

	// QueryRange executa uma query com range de tempo
	QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error)
}
