// Package collector fornece funcionalidades para coleta de métricas de recursos Kubernetes.
// Este pacote implementa a camada de acesso a dados, abstraindo a complexidade de
// interação com diferentes fontes de métricas (Kubernetes API, Prometheus/Mimir).
package collector

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// Collector define a interface para coleta de métricas de recursos Kubernetes.
// Esta interface segue o princípio de inversão de dependência do SOLID,
// permitindo diferentes implementações para coleta de métricas.
type Collector interface {
	// GetDeploymentMetrics retorna métricas atuais de um deployment.
	// Coleta dados em tempo real do cluster Kubernetes, incluindo:
	// - Utilização de CPU em milicores
	// - Utilização de memória em Mi
	// - Número de pods em execução
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - namespace: Namespace do Kubernetes
	//   - deployment: Nome do deployment
	//
	// Retorna:
	//   - K8sMetrics: Métricas atuais do deployment
	//   - error: Erro em caso de falha na coleta
	GetDeploymentMetrics(ctx context.Context, namespace, deployment string) (*types.K8sMetrics, error)

	// GetDeploymentConfig retorna configurações de um deployment.
	// Obtém dados de configuração do Kubernetes, incluindo:
	// - Limites e requests de recursos
	// - Configurações de HPA
	// - Número de réplicas
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - namespace: Namespace do Kubernetes
	//   - deployment: Nome do deployment
	//
	// Retorna:
	//   - K8sDeploymentConfig: Configurações do deployment
	//   - error: Erro em caso de falha na obtenção
	GetDeploymentConfig(ctx context.Context, namespace, deployment string) (*types.K8sDeploymentConfig, error)

	// Query executa uma query pontual no sistema de métricas.
	// Utiliza Prometheus/Mimir para consultas instantâneas.
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - query: Query PromQL
	//
	// Retorna:
	//   - QueryResult: Resultado da query
	//   - error: Erro em caso de falha na consulta
	Query(ctx context.Context, query string) (*types.QueryResult, error)

	// QueryRange executa uma query com range de tempo.
	// Obtém série temporal de métricas do Prometheus/Mimir.
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - query: Query PromQL
	//   - start: Início do período
	//   - end: Fim do período
	//   - step: Intervalo entre pontos
	//
	// Retorna:
	//   - QueryRangeResult: Série temporal de métricas
	//   - error: Erro em caso de falha na consulta
	QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error)
}
