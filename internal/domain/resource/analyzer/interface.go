// Package analyzer fornece funcionalidades para análise de recursos em clusters Kubernetes.
// Este pacote implementa a lógica de negócio para análise de métricas, custos e tendências,
// seguindo os princípios de FinOps para otimização de recursos em containers.
package analyzer

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// ResourceAnalyzer define a interface principal para análise de recursos em clusters Kubernetes.
// Esta interface segue o princípio de segregação de interfaces do SOLID, fornecendo
// métodos específicos para cada tipo de análise necessária.
type ResourceAnalyzer interface {
	// GetMetrics retorna métricas atuais e históricas de um deployment.
	// Coleta dados de utilização de CPU (em milicores), memória (em Mi) e pods,
	// incluindo configurações de HPA e limites de recursos.
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - namespace: Namespace do Kubernetes onde o deployment está localizado
	//   - deployment: Nome do deployment
	//   - period: Período de tempo para análise histórica
	//
	// Retorna:
	//   - MetricsResponse: Contém métricas atuais, históricas e metadados
	//   - error: Erro em caso de falha na coleta
	GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error)

	// GetTrends analisa tendências de utilização de recursos ao longo do tempo.
	// Calcula padrões de uso, sazonalidade e projeta tendências futuras.
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - namespace: Namespace do Kubernetes
	//   - deployment: Nome do deployment
	//   - period: Período para análise de tendências
	//
	// Retorna:
	//   - TrendsResponse: Contém análises de tendência para CPU, memória e pods
	//   - error: Erro em caso de falha na análise
	GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error)

	// AnalyzeResources realiza análise detalhada dos recursos atuais e históricos.
	// Avalia eficiência, identifica gargalos e sugere otimizações.
	//
	// Parâmetros:
	//   - current: Métricas atuais do deployment
	//   - historical: Histórico de métricas para análise comparativa
	//
	// Retorna:
	//   - ResourceAnalysis: Análise completa com recomendações
	AnalyzeResources(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis

	// CalculateCosts calcula custos atuais e projetados dos recursos.
	// Utiliza dados de preços de cloud providers para estimativas precisas.
	//
	// Parâmetros:
	//   - ctx: Contexto da requisição
	//   - current: Métricas atuais para cálculo de custos
	//   - analysis: Análise de recomendações para projeção de custos
	//
	// Retorna:
	//   - CostAnalysis: Análise detalhada de custos atuais e projetados
	//   - error: Erro em caso de falha no cálculo
	CalculateCosts(ctx context.Context, current *types.CurrentMetrics, analysis *types.ResourceRecommendationAnalysis) (*types.CostAnalysis, error)

	// GenerateAlerts produz alertas baseados em análise de métricas.
	// Identifica situações críticas e oportunidades de otimização.
	//
	// Parâmetros:
	//   - current: Métricas atuais para análise
	//   - historical: Histórico de métricas para contexto
	//
	// Retorna:
	//   - []Alert: Lista de alertas gerados
	GenerateAlerts(current *types.CurrentMetrics, historical *types.HistoricalMetrics) []types.Alert
}
