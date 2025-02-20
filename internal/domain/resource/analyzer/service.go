package analyzer

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/errors"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/collector"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/pricing"
)

// Service implementa a interface ResourceAnalyzer
type Service struct {
	metricsCollector collector.Collector
	pricingClient    *pricing.Client
}

// NewService cria uma nova instância do Service
func NewService(metricsCollector collector.Collector, pricingClient *pricing.Client) *Service {
	return &Service{
		metricsCollector: metricsCollector,
		pricingClient:    pricingClient,
	}
}

// GetMetrics retorna as métricas atuais e históricas de um deployment
func (s *Service) GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error) {
	logger.Info("Iniciando coleta de métricas",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", deployment),
		logger.NewField("period", period),
	)

	response := &types.MetricsResponse{
		Current: &types.CurrentMetrics{
			Deployment: struct {
				Config struct {
					CPU struct {
						Request float64 `json:"request"`
						Limit   float64 `json:"limit"`
					} `json:"cpu"`
					Memory struct {
						Request float64 `json:"request"`
						Limit   float64 `json:"limit"`
					} `json:"memory"`
					HPA struct {
						MinReplicas int     `json:"minReplicas"`
						MaxReplicas int     `json:"maxReplicas"`
						TargetCPU   float64 `json:"targetCPU"`
					} `json:"hpa"`
				} `json:"config"`
			}{},
			Analysis: struct {
				CPU struct {
					Distribution map[string]float64 `json:"distribution"`
					Alerts       struct {
						HighCPU    int `json:"highCPU"`
						NearLimit  int `json:"nearLimit"`
						HighMemory int `json:"highMemory"`
					} `json:"alerts"`
					Usage struct {
						Current struct {
							Average float64 `json:"average"`
							Peak    float64 `json:"peak"`
						} `json:"current"`
						Historical struct {
							Average float64 `json:"average"`
							Peak    float64 `json:"peak"`
						} `json:"historical"`
					} `json:"usage"`
				} `json:"cpu"`
				Memory struct {
					Usage struct {
						Current struct {
							Average float64 `json:"average"`
							Peak    float64 `json:"peak"`
						} `json:"current"`
						Historical struct {
							Average float64 `json:"average"`
							Peak    float64 `json:"peak"`
						} `json:"historical"`
					} `json:"usage"`
				} `json:"memory"`
			}{},
			CPU:    &types.ResourceMetrics{},
			Memory: &types.ResourceMetrics{},
			Pods:   &types.PodMetrics{},
		},
		Historical: &types.HistoricalMetrics{
			CPU:    []*types.ResourceMetrics{},
			Memory: []*types.ResourceMetrics{},
			Pods:   []*types.PodMetrics{},
		},
		Metadata: struct {
			Analysis struct {
				Timestamp  string   `json:"timestamp"`
				Period     string   `json:"period"`
				Cluster    string   `json:"cluster"`
				Sources    []string `json:"sources"`
				Confidence struct {
					CPU    float64 `json:"cpu"`
					Memory float64 `json:"memory"`
					Pods   float64 `json:"pods"`
				} `json:"confidence"`
			} `json:"analysis"`
		}{},
	}

	// Inicializa o mapa de distribuição
	response.Current.Analysis.CPU.Distribution = make(map[string]float64)

	// Obtém métricas atuais
	logger.Info("Coletando métricas atuais")
	k8sMetrics, err := s.metricsCollector.GetDeploymentMetrics(ctx, namespace, deployment)
	if err != nil {
		logger.Error("Erro ao obter métricas do deployment", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", deployment),
		)
		return nil, fmt.Errorf("erro ao obter métricas do deployment: %w", err)
	}

	// Obtém configuração do deployment
	logger.Info("Coletando configuração do deployment")
	config, err := s.metricsCollector.GetDeploymentConfig(ctx, namespace, deployment)
	if err != nil {
		logger.Error("Erro ao obter configuração do deployment", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", deployment),
		)
		return nil, fmt.Errorf("erro ao obter configuração do deployment: %w", err)
	}

	// Configura a resposta com os dados atuais
	// Os valores já estão nas unidades corretas
	response.Current.Deployment.Config.CPU.Request = config.CPU.Request
	response.Current.Deployment.Config.CPU.Limit = config.CPU.Limit
	response.Current.Deployment.Config.Memory.Request = config.Memory.Request
	response.Current.Deployment.Config.Memory.Limit = config.Memory.Limit
	response.Current.Deployment.Config.HPA.MinReplicas = config.Pods.MinReplicas
	response.Current.Deployment.Config.HPA.MaxReplicas = config.Pods.MaxReplicas
	response.Current.Deployment.Config.HPA.TargetCPU = config.Pods.TargetCPU * 100 // Converte para percentual

	// Configura análise de CPU (valores já em milicores)
	response.Current.Analysis.CPU.Usage.Current.Average = k8sMetrics.CPU.Average
	response.Current.Analysis.CPU.Usage.Current.Peak = k8sMetrics.CPU.Peak
	response.Current.Analysis.Memory.Usage.Current.Average = k8sMetrics.Memory.Average
	response.Current.Analysis.Memory.Usage.Current.Peak = k8sMetrics.Memory.Peak

	// Configura distribuição de CPU
	cpuUsage := k8sMetrics.CPU.Usage // Já está em milicores
	response.Current.Analysis.CPU.Distribution = calculateCPUDistribution(cpuUsage)

	// Configura alertas
	response.Current.Analysis.CPU.Alerts.HighCPU = countPodsAboveCPU(k8sMetrics, 999.0)
	response.Current.Analysis.CPU.Alerts.NearLimit = countPodsAboveCPU(k8sMetrics, 900.0)
	response.Current.Analysis.CPU.Alerts.HighMemory = countPodsAboveMemory(k8sMetrics, 800.0)

	// Obtém métricas históricas
	end := time.Now()
	start := end.Add(-period)
	step := 5 * time.Minute

	logger.Info("Coletando métricas históricas",
		logger.NewField("start", start),
		logger.NewField("end", end),
		logger.NewField("step", step),
	)

	// Query para CPU
	cpuQuery := buildCPUHistoricalQuery(namespace, deployment)
	logger.Info("Executando query de CPU histórica",
		logger.NewField("query", cpuQuery),
	)
	cpuResult, err := s.metricsCollector.QueryRange(ctx, cpuQuery, start, end, step)
	if err != nil {
		logger.Error("Erro ao obter métricas históricas de CPU", err,
			logger.NewField("query", cpuQuery),
		)
		return nil, fmt.Errorf("erro ao obter métricas históricas de CPU: %w", err)
	}

	// Query para memória
	memoryQuery := buildMemoryHistoricalQuery(namespace, deployment)
	logger.Info("Executando query de memória histórica",
		logger.NewField("query", memoryQuery),
	)
	memoryResult, err := s.metricsCollector.QueryRange(ctx, memoryQuery, start, end, step)
	if err != nil {
		logger.Error("Erro ao obter métricas históricas de memória", err,
			logger.NewField("query", memoryQuery),
		)
		return nil, fmt.Errorf("erro ao obter métricas históricas de memória: %w", err)
	}

	// Configura métricas históricas
	if len(cpuResult.Values) > 0 {
		var sum, peak float64
		for _, v := range cpuResult.Values {
			// Converte de cores para milicores
			cpuValue := v.Value * 1000 // Já está em cores, precisamos converter para milicores
			sum += cpuValue
			if cpuValue > peak {
				peak = cpuValue
			}
		}
		response.Current.Analysis.CPU.Usage.Historical.Average = sum / float64(len(cpuResult.Values))
		response.Current.Analysis.CPU.Usage.Historical.Peak = peak
	}

	if len(memoryResult.Values) > 0 {
		var sum, peak float64
		for _, v := range memoryResult.Values {
			// Converte de bytes para Mi
			memValue := v.Value / (1024 * 1024) // Converte bytes para Mi
			sum += memValue
			if memValue > peak {
				peak = memValue
			}
		}
		response.Current.Analysis.Memory.Usage.Historical.Average = sum / float64(len(memoryResult.Values))
		response.Current.Analysis.Memory.Usage.Historical.Peak = peak
	}

	// Configura metadados
	response.Metadata.Analysis.Timestamp = time.Now().Format(time.RFC3339)
	response.Metadata.Analysis.Period = period.String()
	response.Metadata.Analysis.Cluster = config.ClusterName
	response.Metadata.Analysis.Sources = []string{"kubernetes", "prometheus"}
	response.Metadata.Analysis.Confidence.CPU = 95.0
	response.Metadata.Analysis.Confidence.Memory = 95.0
	response.Metadata.Analysis.Confidence.Pods = 100.0

	// Calcula recomendações
	logger.Info("Calculando recomendações")
	recommendations := s.CalculateRecommendations(response.Current, response.Historical)

	// Calcula custos
	logger.Info("Calculando custos")
	costs, err := s.CalculateCosts(ctx, response.Current, recommendations)
	if err != nil {
		logger.Error("Erro ao calcular custos", err)
		return nil, errors.NewInvalidMetricsError("costs", "erro ao calcular custos")
	}

	// Atualiza a resposta com as recomendações e custos
	response.Analysis = recommendations
	response.Costs = costs

	logger.Info("Análise concluída com sucesso",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", deployment),
	)

	return response, nil
}

// GetTrends retorna as tendências de uso de recursos
func (s *Service) GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error) {
	logger.Info("Iniciando análise de tendências",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", deployment),
		logger.NewField("period", period),
	)

	// Obtém métricas
	metricsResponse, err := s.GetMetrics(ctx, namespace, deployment, period)
	if err != nil {
		logger.Error("Erro ao obter métricas para análise de tendências", err)
		return nil, errors.NewInvalidMetricsError("trends", "erro ao obter métricas para análise de tendências")
	}

	// Calcula tendências
	logger.Info("Calculando tendências")
	response := &types.TrendsResponse{
		CPU: &types.TrendMetrics{
			Trend:      calculateUtilizationTrend(metricsResponse.Historical.CPU),
			Confidence: 0.8,
			Period:     period.String(),
		},
		Memory: &types.TrendMetrics{
			Trend:      calculateUtilizationTrend(metricsResponse.Historical.Memory),
			Confidence: 0.8,
			Period:     period.String(),
		},
		Pods: &types.TrendMetrics{
			Trend:      calculatePodsUtilizationTrend(metricsResponse.Historical.Pods),
			Confidence: 0.9,
			Period:     period.String(),
		},
	}

	logger.Info("Análise de tendências concluída",
		logger.NewField("cpu_trend", response.CPU.Trend),
		logger.NewField("memory_trend", response.Memory.Trend),
		logger.NewField("pods_trend", response.Pods.Trend),
	)

	return response, nil
}

// AnalyzeResources realiza análise detalhada dos recursos
func (s *Service) AnalyzeResources(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis {
	// Análise de CPU
	cpuAnalysis := &types.ResourceTypeAnalysis{
		CurrentUsage:  current.CPU.Usage,
		HistoricalAvg: calculateHistoricalAverage(historical.CPU),
		Peak:          calculateHistoricalPeak(historical.CPU),
		Distribution:  map[string]float64{"normal": 100.0},
		Utilization:   current.CPU.Utilization,
		UtilizationTrend: &types.UtilizationTrend{
			Current:     current.CPU.Utilization,
			Historical:  calculateHistoricalUtilization(historical.CPU),
			Trend:       calculateUtilizationTrend(historical.CPU),
			Pattern:     detectPattern(historical.CPU),
			Seasonality: detectSeasonality(historical.CPU),
		},
		Recommendation: generateCPURecommendation(current.CPU, historical.CPU),
	}

	// Análise de memória
	memoryAnalysis := &types.ResourceTypeAnalysis{
		CurrentUsage:  current.Memory.Usage,
		HistoricalAvg: calculateHistoricalAverage(historical.Memory),
		Peak:          calculateHistoricalPeak(historical.Memory),
		Distribution:  map[string]float64{"normal": 100.0},
		Utilization:   current.Memory.Utilization,
		UtilizationTrend: &types.UtilizationTrend{
			Current:     current.Memory.Utilization,
			Historical:  calculateHistoricalUtilization(historical.Memory),
			Trend:       calculateUtilizationTrend(historical.Memory),
			Pattern:     detectPattern(historical.Memory),
			Seasonality: detectSeasonality(historical.Memory),
		},
		Recommendation: generateMemoryRecommendation(current.Memory, historical.Memory),
	}

	// Análise de pods
	podsAnalysis := &types.PodAnalysis{
		CurrentRunning:    current.Pods.Running,
		HistoricalAvg:     calculatePodsHistoricalUtilization(historical.Pods),
		Peak:              current.Pods.MaxReplicas,
		ScalingEfficiency: calculateScalingEfficiency(current.Pods, historical.Pods),
		UtilizationTrend: &types.UtilizationTrend{
			Current:     current.Pods.Utilization,
			Historical:  calculatePodsHistoricalUtilization(historical.Pods),
			Trend:       calculatePodsUtilizationTrend(historical.Pods),
			Pattern:     detectPodsPattern(historical.Pods),
			Seasonality: detectPodsSeasonality(historical.Pods),
		},
		Recommendations: struct {
			Current     *types.Recommendation `json:"current"`
			MinReplicas *types.Recommendation `json:"minReplicas"`
			MaxReplicas *types.Recommendation `json:"maxReplicas"`
		}{
			Current:     generatePodsRecommendation(current.Pods, historical.Pods),
			MinReplicas: generateMinReplicasRecommendation(current.Pods, historical.Pods),
			MaxReplicas: generateMaxReplicasRecommendation(current.Pods, historical.Pods),
		},
	}

	// Determina status geral
	status := determineOverallStatus(cpuAnalysis, memoryAnalysis, podsAnalysis)

	return &types.ResourceAnalysis{
		CPU:    cpuAnalysis,
		Memory: memoryAnalysis,
		Pods:   podsAnalysis,
		Status: status,
	}
}

// CalculateCosts calcula os custos dos recursos
func (s *Service) CalculateCosts(ctx context.Context, current *types.CurrentMetrics, analysis *types.ResourceRecommendationAnalysis) (*types.CostAnalysis, error) {
	// Obtém preços atuais
	prices, err := s.pricingClient.GetCurrentPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter preços: %w", err)
	}

	// Obtém taxa de câmbio USD -> BRL
	exchange, err := s.pricingClient.GetExchangeRate(ctx, "USD", "BRL")
	if err != nil {
		return nil, fmt.Errorf("erro ao obter taxa de câmbio: %w", err)
	}

	// Converte CPU de milicores para cores e memória de Mi para GB
	cpuCores := current.Deployment.Config.CPU.Request / 1000
	memoryGB := current.Deployment.Config.Memory.Request / 1024

	// Calcula custos por hora em USD e converte para BRL
	hourly := &types.ResourceCosts{
		CPU:    cpuCores * prices.CPU.PerCore * exchange.Rate,
		Memory: memoryGB * prices.Memory.PerGB * exchange.Rate,
	}
	hourly.Total = hourly.CPU + hourly.Memory

	// Calcula custos por dia (24 horas)
	daily := &types.ResourceCosts{
		CPU:    hourly.CPU * 24,
		Memory: hourly.Memory * 24,
		Total:  hourly.Total * 24,
	}

	// Calcula custos por mês (730 horas em média)
	monthly := &types.ResourceCosts{
		CPU:    hourly.CPU * 730,
		Memory: hourly.Memory * 730,
		Total:  hourly.Total * 730,
	}

	// Converte recomendações para unidades corretas
	recommendedCPUCores := analysis.CPU.Recommendation.Suggested / 1000
	recommendedMemoryGB := analysis.Memory.Recommendation.Suggested / 1024

	// Calcula custos recomendados em BRL
	recommendedHourly := &types.ResourceCosts{
		CPU:    recommendedCPUCores * prices.CPU.PerCore * exchange.Rate,
		Memory: recommendedMemoryGB * prices.Memory.PerGB * exchange.Rate,
	}
	recommendedHourly.Total = recommendedHourly.CPU + recommendedHourly.Memory

	recommendedDaily := &types.ResourceCosts{
		CPU:    recommendedHourly.CPU * 24,
		Memory: recommendedHourly.Memory * 24,
		Total:  recommendedHourly.Total * 24,
	}

	recommendedMonthly := &types.ResourceCosts{
		CPU:    recommendedHourly.CPU * 730,
		Memory: recommendedHourly.Memory * 730,
		Total:  recommendedHourly.Total * 730,
	}

	// Calcula economia potencial em BRL
	savings := &types.ResourceCosts{
		CPU:    monthly.CPU - recommendedMonthly.CPU,
		Memory: monthly.Memory - recommendedMonthly.Memory,
	}
	savings.Total = savings.CPU + savings.Memory

	return &types.CostAnalysis{
		Current: &types.CostData{
			Hourly:  hourly,
			Daily:   daily,
			Monthly: monthly,
		},
		Recommended: &types.CostData{
			Hourly:  recommendedHourly,
			Daily:   recommendedDaily,
			Monthly: recommendedMonthly,
		},
		Savings:  savings,
		Currency: "BRL",
		Exchange: &types.ExchangeInfo{
			Rate:         exchange.Rate,
			FromCurrency: exchange.FromCurrency,
			ToCurrency:   exchange.ToCurrency,
			UpdatedAt:    exchange.UpdatedAt,
		},
	}, nil
}

// GenerateAlerts gera alertas baseados nas métricas
func (s *Service) GenerateAlerts(current *types.CurrentMetrics, historical *types.HistoricalMetrics) []types.Alert {
	return nil
}

// CalculateRecommendations calcula recomendações de recursos baseadas nas métricas
func (s *Service) CalculateRecommendations(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceRecommendationAnalysis {
	analysis := &types.ResourceRecommendationAnalysis{
		CPU: &types.ResourceRecommendation{
			Status: "insufficient_data",
		},
		Memory: &types.ResourceRecommendation{
			Status: "insufficient_data",
		},
		Pods: &types.PodRecommendation{
			Status: "insufficient_data",
		},
	}

	// Calcula recomendações de CPU
	cpuPeak := current.Analysis.CPU.Usage.Current.Peak
	cpuAvg := current.Analysis.CPU.Usage.Current.Average
	cpuRequest := current.Deployment.Config.CPU.Request

	if cpuPeak > 0 && cpuAvg > 0 {
		// Calcula a recomendação baseada no uso médio + 30% de buffer
		suggestedCPU := cpuAvg * 1.3

		// Ajusta para o pico se necessário
		if suggestedCPU < cpuPeak {
			suggestedCPU = cpuPeak * 1.1 // 10% de buffer para picos
		}

		// Arredonda para o próximo múltiplo de 100m
		suggestedCPU = math.Ceil(suggestedCPU/100) * 100

		analysis.CPU.Status = "optimized"
		analysis.CPU.Recommendation = &types.ResourceSuggestion{
			Current:   cpuRequest,
			Suggested: suggestedCPU,
			Action:    determineAction(cpuRequest, suggestedCPU),
		}
	}

	// Calcula recomendações de memória
	memPeak := current.Analysis.Memory.Usage.Current.Peak
	memAvg := current.Analysis.Memory.Usage.Current.Average
	memRequest := current.Deployment.Config.Memory.Request

	if memPeak > 0 && memAvg > 0 {
		// Calcula a recomendação baseada no uso médio + 40% de buffer
		suggestedMem := memAvg * 1.4

		// Ajusta para o pico se necessário
		if suggestedMem < memPeak {
			suggestedMem = memPeak * 1.2 // 20% de buffer para picos
		}

		// Arredonda para o próximo múltiplo de 128Mi
		suggestedMem = math.Ceil(suggestedMem/128) * 128

		analysis.Memory.Status = "optimized"
		analysis.Memory.Recommendation = &types.ResourceSuggestion{
			Current:   memRequest,
			Suggested: suggestedMem,
			Action:    determineAction(memRequest, suggestedMem),
		}
	}

	// Calcula recomendações de pods
	if current.Pods.Running > 0 {
		analysis.Pods.Status = "optimized"
		analysis.Pods.Recommendation = &types.PodCount{
			Current:   current.Pods.Running,
			Suggested: calculateSuggestedPods(current, historical),
			Min:       current.Deployment.Config.HPA.MinReplicas,
			Max:       current.Deployment.Config.HPA.MaxReplicas,
		}
	}

	return analysis
}

// determineAction determina a ação recomendada baseada nos valores atual e sugerido
func determineAction(current, suggested float64) string {
	diff := math.Abs(current - suggested)
	threshold := current * 0.1 // 10% de diferença

	if diff <= threshold {
		return "maintain"
	} else if suggested > current {
		return "increase"
	}
	return "decrease"
}

// calculateSuggestedPods calcula o número sugerido de pods
func calculateSuggestedPods(current *types.CurrentMetrics, historical *types.HistoricalMetrics) int {
	// Usa o número atual de pods como base
	suggested := current.Pods.Running

	// Ajusta baseado no uso de CPU e memória
	cpuUtilization := current.Analysis.CPU.Usage.Current.Average / current.Deployment.Config.CPU.Request
	memUtilization := current.Analysis.Memory.Usage.Current.Average / current.Deployment.Config.Memory.Request

	// Se ambos CPU e memória estão acima de 80%, sugere aumentar
	if cpuUtilization > 0.8 && memUtilization > 0.8 {
		suggested++
	}

	// Se ambos estão abaixo de 40%, sugere diminuir
	if cpuUtilization < 0.4 && memUtilization < 0.4 && suggested > 1 {
		suggested--
	}

	// Garante que está dentro dos limites do HPA
	if suggested < current.Deployment.Config.HPA.MinReplicas {
		suggested = current.Deployment.Config.HPA.MinReplicas
	}
	if suggested > current.Deployment.Config.HPA.MaxReplicas {
		suggested = current.Deployment.Config.HPA.MaxReplicas
	}

	return suggested
}

// Funções auxiliares

func countPodsAboveCPU(metrics *types.K8sMetrics, threshold float64) int {
	count := 0
	if metrics.CPU.Usage > threshold {
		count++
	}
	return count
}

func countPodsAboveMemory(metrics *types.K8sMetrics, threshold float64) int {
	count := 0
	if metrics.Memory.Usage > threshold {
		count++
	}
	return count
}
