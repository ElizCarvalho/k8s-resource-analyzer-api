package analyzer

import (
	"math"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// determineOverallStatus determina o status geral da análise
func determineOverallStatus(cpu, memory *types.ResourceTypeAnalysis, pods *types.PodAnalysis) string {
	// Verifica se há algum problema crítico
	if cpu.Utilization > 90 || memory.Utilization > 90 {
		return "critical"
	}

	// Verifica se há algum alerta
	if cpu.Utilization > 75 || memory.Utilization > 75 {
		return "warning"
	}

	// Verifica se há subutilização
	if cpu.Utilization < 30 && memory.Utilization < 30 {
		return "underutilized"
	}

	// Se tudo estiver normal
	return "normal"
}

// calculateHistoricalAverage calcula a média histórica de uso
func calculateHistoricalAverage(metrics []*types.ResourceMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}

	var sum float64
	for _, m := range metrics {
		sum += m.Usage
	}
	return sum / float64(len(metrics))
}

// calculateHistoricalPeak calcula o pico histórico de uso
func calculateHistoricalPeak(metrics []*types.ResourceMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}

	peak := metrics[0].Usage
	for _, m := range metrics {
		if m.Usage > peak {
			peak = m.Usage
		}
	}
	return peak
}

// calculateHistoricalUtilization calcula a utilização histórica média
func calculateHistoricalUtilization(metrics []*types.ResourceMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}

	var sum float64
	for _, m := range metrics {
		sum += m.Utilization
	}
	return sum / float64(len(metrics))
}

// calculatePodsHistoricalUtilization calcula a utilização histórica média de pods
func calculatePodsHistoricalUtilization(metrics []*types.PodMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}

	var sum float64
	for _, m := range metrics {
		sum += float64(m.Running) / float64(m.Replicas)
	}
	return (sum / float64(len(metrics))) * 100
}

// calculateScalingEfficiency calcula a eficiência do escalonamento de pods
func calculateScalingEfficiency(current *types.PodMetrics, historical []*types.PodMetrics) float64 {
	if len(historical) == 0 {
		return 0
	}

	// Calcula a média de utilização
	avgUtilization := calculatePodsHistoricalUtilization(historical)

	// Calcula o desvio padrão da utilização
	var sumSquares float64
	for _, m := range historical {
		utilization := float64(m.Running) / float64(m.Replicas) * 100
		diff := utilization - avgUtilization
		sumSquares += diff * diff
	}
	stdDev := math.Sqrt(sumSquares / float64(len(historical)))

	// Quanto menor o desvio padrão, melhor a eficiência
	// Normaliza para um valor entre 0 e 100
	efficiency := 100 - (stdDev / avgUtilization * 100)
	if efficiency < 0 {
		efficiency = 0
	}
	if efficiency > 100 {
		efficiency = 100
	}

	return efficiency
}

// detectPattern detecta padrões nos dados históricos
func detectPattern(metrics []*types.ResourceMetrics) string {
	if len(metrics) < 2 {
		return "insufficient_data"
	}

	// Calcula a tendência
	trend := calculateUtilizationTrend(metrics)

	// Analisa o padrão com base na tendência
	switch {
	case trend > 10:
		return "increasing"
	case trend < -10:
		return "decreasing"
	case trend >= -5 && trend <= 5:
		return "stable"
	default:
		return "fluctuating"
	}
}

// detectPodsPattern detecta padrões no histórico de pods
func detectPodsPattern(metrics []*types.PodMetrics) string {
	if len(metrics) < 2 {
		return "insufficient_data"
	}

	// Calcula a tendência
	trend := calculatePodsUtilizationTrend(metrics)

	// Analisa o padrão com base na tendência
	switch {
	case trend > 10:
		return "scaling_up"
	case trend < -10:
		return "scaling_down"
	case trend >= -5 && trend <= 5:
		return "stable"
	default:
		return "fluctuating"
	}
}

// detectSeasonality detecta sazonalidade nos dados
func detectSeasonality(metrics []*types.ResourceMetrics) *types.Seasonality {
	if len(metrics) < 24 { // Precisa de pelo menos 24 pontos para detectar sazonalidade
		return &types.Seasonality{
			Pattern: "insufficient_data",
			Period:  "unknown",
		}
	}

	// Implementação simplificada - apenas verifica padrões diários
	// TODO: Implementar análise mais sofisticada (ex: FFT, autocorrelação)
	return &types.Seasonality{
		Pattern: "daily",
		Period:  "24h",
	}
}

// detectPodsSeasonality detecta sazonalidade no uso de pods
func detectPodsSeasonality(metrics []*types.PodMetrics) *types.Seasonality {
	if len(metrics) < 24 {
		return &types.Seasonality{
			Pattern: "insufficient_data",
			Period:  "unknown",
		}
	}

	// Implementação simplificada
	return &types.Seasonality{
		Pattern: "daily",
		Period:  "24h",
	}
}

// calculateUtilizationTrend calcula a tendência de utilização
func calculateUtilizationTrend(metrics []*types.ResourceMetrics) float64 {
	if len(metrics) < 2 {
		return 0
	}

	// Calcula a média móvel para suavizar flutuações
	windowSize := 3
	if len(metrics) < windowSize {
		windowSize = len(metrics)
	}

	var recentAvg, oldAvg float64
	for i := 0; i < windowSize; i++ {
		recentAvg += metrics[len(metrics)-1-i].Utilization
		oldAvg += metrics[i].Utilization
	}
	recentAvg /= float64(windowSize)
	oldAvg /= float64(windowSize)

	// Calcula a variação percentual
	if oldAvg == 0 {
		return 0
	}
	return ((recentAvg - oldAvg) / oldAvg) * 100
}

// calculatePodsUtilizationTrend calcula a tendência de utilização de pods
func calculatePodsUtilizationTrend(metrics []*types.PodMetrics) float64 {
	if len(metrics) < 2 {
		return 0
	}

	// Calcula a média móvel para suavizar flutuações
	windowSize := 3
	if len(metrics) < windowSize {
		windowSize = len(metrics)
	}

	var recentAvg, oldAvg float64
	for i := 0; i < windowSize; i++ {
		recentAvg += metrics[len(metrics)-1-i].Utilization
		oldAvg += metrics[i].Utilization
	}
	recentAvg /= float64(windowSize)
	oldAvg /= float64(windowSize)

	// Calcula a variação percentual
	if oldAvg == 0 {
		return 0
	}
	return ((recentAvg - oldAvg) / oldAvg) * 100
}

// generateCPURecommendation gera recomendações para CPU
func generateCPURecommendation(current *types.ResourceMetrics, historical []*types.ResourceMetrics) *types.Recommendation {
	if len(historical) == 0 {
		return &types.Recommendation{
			Current:   current.Usage,
			Suggested: current.Usage,
			Reason:    "insufficient_data",
		}
	}

	// Calcula valores estatísticos
	peak := calculateHistoricalPeak(historical)
	buffer := 0.2 // 20% de buffer

	// Calcula valor sugerido em milicores
	suggested := peak * (1 + buffer)

	// Determina a razão
	var reason string
	switch {
	case suggested < current.Usage*0.7:
		reason = "overprovisioned"
	case suggested > current.Usage*1.3:
		reason = "underprovisioned"
	default:
		reason = "optimal"
	}

	return &types.Recommendation{
		Current:   current.Usage,
		Suggested: suggested,
		Reason:    reason,
	}
}

// generateMemoryRecommendation gera recomendações para memória
func generateMemoryRecommendation(current *types.ResourceMetrics, historical []*types.ResourceMetrics) *types.Recommendation {
	if len(historical) == 0 {
		return &types.Recommendation{
			Current:   current.Usage,
			Suggested: current.Usage,
			Reason:    "insufficient_data",
		}
	}

	// Calcula valores estatísticos
	peak := calculateHistoricalPeak(historical)
	buffer := 0.3 // 30% de buffer para memória

	// Calcula valor sugerido em Mi
	suggested := peak * (1 + buffer)

	// Determina a razão
	var reason string
	switch {
	case suggested < current.Usage*0.7:
		reason = "overprovisioned"
	case suggested > current.Usage*1.3:
		reason = "underprovisioned"
	default:
		reason = "optimal"
	}

	return &types.Recommendation{
		Current:   current.Usage,
		Suggested: suggested,
		Reason:    reason,
	}
}

// generatePodsRecommendation gera recomendações para número de pods
func generatePodsRecommendation(current *types.PodMetrics, historical []*types.PodMetrics) *types.Recommendation {
	if len(historical) == 0 {
		return &types.Recommendation{
			Current:   float64(current.Replicas),
			Suggested: float64(current.Replicas),
			Reason:    "insufficient_data",
		}
	}

	// Calcula utilização média
	avgUtilization := calculatePodsHistoricalUtilization(historical)

	// Determina o número sugerido de pods
	suggested := float64(current.Replicas)
	if avgUtilization > 80 {
		suggested = float64(current.Replicas) * 1.2 // Aumenta 20%
	} else if avgUtilization < 50 {
		suggested = float64(current.Replicas) * 0.8 // Reduz 20%
	}

	// Arredonda para o inteiro mais próximo
	suggested = math.Round(suggested)

	// Determina a razão
	var reason string
	switch {
	case suggested < float64(current.Replicas):
		reason = "overprovisioned"
	case suggested > float64(current.Replicas):
		reason = "underprovisioned"
	default:
		reason = "optimal"
	}

	return &types.Recommendation{
		Current:   float64(current.Replicas),
		Suggested: suggested,
		Reason:    reason,
	}
}

// generateMinReplicasRecommendation gera recomendações para o mínimo de réplicas
func generateMinReplicasRecommendation(current *types.PodMetrics, historical []*types.PodMetrics) *types.Recommendation {
	if len(historical) == 0 {
		return &types.Recommendation{
			Current:   float64(current.MinReplicas),
			Suggested: float64(current.MinReplicas),
			Reason:    "insufficient_data",
		}
	}

	// Encontra o mínimo histórico de pods em uso
	minUsed := float64(current.Running)
	for _, m := range historical {
		if float64(m.Running) < minUsed {
			minUsed = float64(m.Running)
		}
	}

	// Adiciona um buffer de segurança
	suggested := math.Max(1, math.Round(minUsed*0.8)) // Pelo menos 1 pod

	// Determina a razão
	var reason string
	switch {
	case suggested < float64(current.MinReplicas):
		reason = "can_be_reduced"
	case suggested > float64(current.MinReplicas):
		reason = "should_be_increased"
	default:
		reason = "optimal"
	}

	return &types.Recommendation{
		Current:   float64(current.MinReplicas),
		Suggested: suggested,
		Reason:    reason,
	}
}

// generateMaxReplicasRecommendation gera recomendações para o máximo de réplicas
func generateMaxReplicasRecommendation(current *types.PodMetrics, historical []*types.PodMetrics) *types.Recommendation {
	if len(historical) == 0 {
		return &types.Recommendation{
			Current:   float64(current.MaxReplicas),
			Suggested: float64(current.MaxReplicas),
			Reason:    "insufficient_data",
		}
	}

	// Encontra o máximo histórico de pods em uso
	maxUsed := float64(current.Running)
	for _, m := range historical {
		if float64(m.Running) > maxUsed {
			maxUsed = float64(m.Running)
		}
	}

	// Adiciona um buffer para crescimento
	suggested := math.Round(maxUsed * 1.5) // 50% de margem para crescimento

	// Determina a razão
	var reason string
	switch {
	case suggested < float64(current.MaxReplicas)*0.7:
		reason = "can_be_reduced"
	case suggested > float64(current.MaxReplicas):
		reason = "should_be_increased"
	default:
		reason = "optimal"
	}

	return &types.Recommendation{
		Current:   float64(current.MaxReplicas),
		Suggested: suggested,
		Reason:    reason,
	}
}
