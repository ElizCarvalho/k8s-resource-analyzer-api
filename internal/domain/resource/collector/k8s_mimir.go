package collector

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
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
	logger.Info("Collecting deployment metrics",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", name),
	)
	metrics, err := c.K8sClient.GetDeploymentMetrics(ctx, namespace, name)
	if err != nil {
		logger.Error("Failed to collect deployment metrics", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
		return nil, err
	}
	return metrics, nil
}

// GetDeploymentConfig retorna configurações de um deployment
func (c *K8sMimirCollector) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	logger.Info("Collecting deployment configuration",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", name),
	)
	config, err := c.K8sClient.GetDeploymentConfig(ctx, namespace, name)
	if err != nil {
		logger.Error("Failed to collect deployment configuration", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
		return nil, err
	}
	return config, nil
}

// Query executa uma query pontual
func (c *K8sMimirCollector) Query(ctx context.Context, query string) (*types.QueryResult, error) {
	logger.Info("Executing instant query",
		logger.NewField("query", query),
	)
	result, err := c.MimirClient.Query(ctx, query)
	if err != nil {
		logger.Error("Failed to execute instant query", err,
			logger.NewField("query", query),
		)
		return nil, err
	}
	return result, nil
}

// QueryRange executa uma query com range de tempo
func (c *K8sMimirCollector) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	logger.Info("Executing range query",
		logger.NewField("query", query),
		logger.NewField("start", start),
		logger.NewField("end", end),
		logger.NewField("step", step),
	)
	result, err := c.MimirClient.QueryRange(ctx, query, start, end, step)
	if err != nil {
		logger.Error("Failed to execute range query", err,
			logger.NewField("query", query),
		)
		return nil, err
	}
	return result, nil
}
