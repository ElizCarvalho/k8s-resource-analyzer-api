package collector

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// K8sClient define a interface para o cliente Kubernetes
type K8sClient interface {
	GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error)
	GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error)
	CheckConnection(ctx context.Context) error
}

// MimirClient define a interface para o cliente Mimir
type MimirClient interface {
	Query(ctx context.Context, query string) (*types.QueryResult, error)
	QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error)
	CheckConnection(ctx context.Context) error
}

// K8sMimirCollector implementa a interface Collector usando K8s e Mimir
type K8sMimirCollector struct {
	K8sClient   K8sClient
	MimirClient MimirClient
}

// NewK8sMimirCollector cria uma nova instância do K8sMimirCollector
func NewK8sMimirCollector(k8sClient K8sClient, mimirClient MimirClient) *K8sMimirCollector {
	return &K8sMimirCollector{
		K8sClient:   k8sClient,
		MimirClient: mimirClient,
	}
}

// GetDeploymentMetrics retorna métricas atuais de um deployment
func (c *K8sMimirCollector) GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error) {
	return c.K8sClient.GetDeploymentMetrics(ctx, namespace, name)
}

// GetDeploymentConfig retorna configurações de um deployment
func (c *K8sMimirCollector) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	return c.K8sClient.GetDeploymentConfig(ctx, namespace, name)
}

// Query executa uma query pontual
func (c *K8sMimirCollector) Query(ctx context.Context, query string) (*types.QueryResult, error) {
	return c.MimirClient.Query(ctx, query)
}

// QueryRange executa uma query com range de tempo
func (c *K8sMimirCollector) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	return c.MimirClient.QueryRange(ctx, query, start, end, step)
}
