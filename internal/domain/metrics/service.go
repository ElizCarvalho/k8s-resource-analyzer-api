package metrics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/pricing"
)

// K8sClient define a interface para o cliente Kubernetes
type K8sClient interface {
	GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error)
	GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error)
}

// MimirClient define a interface para o cliente Mimir
type MimirClient interface {
	Query(ctx context.Context, query string) (*types.QueryResult, error)
	QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error)
}

// MetricsProvider define a interface para provedores de métricas
type MetricsProvider interface {
	// GetMetrics retorna as métricas atuais e históricas de um deployment
	GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error)

	// GetTrends retorna as tendências de uso de recursos de um deployment
	GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error)

	// AnalyzeResources realiza uma análise detalhada dos recursos
	AnalyzeResources(current *types.ResourceMetrics, historical []*types.ResourceMetrics) *types.ResourceAnalysis

	// CalculateCosts calcula os custos dos recursos
	CalculateCosts(ctx context.Context, current *types.ResourceMetrics, analysis *types.ResourceAnalysis) (*types.CostAnalysis, error)

	// GenerateAlerts gera alertas baseados nas métricas
	GenerateAlerts(current *types.ResourceMetrics, historical []*types.ResourceMetrics) []types.Alert
}

// Service implementa a interface MetricsProvider
type Service struct {
	k8sClient     K8sClient
	mimirClient   MimirClient
	pricingClient *pricing.Client
}

// NewService cria uma nova instância do serviço de métricas
func NewService(k8sClient K8sClient, mimirClient MimirClient, pricingClient *pricing.Client) *Service {
	return &Service{
		k8sClient:     k8sClient,
		mimirClient:   mimirClient,
		pricingClient: pricingClient,
	}
}

// GetCurrentMetrics retorna as métricas atuais de um deployment
func (s *Service) GetCurrentMetrics(ctx context.Context, params types.MetricsParams) (*types.ResourceMetrics, error) {
	// Obtém as métricas do Kubernetes
	k8sMetrics, err := s.k8sClient.GetDeploymentMetrics(ctx, params.Namespace, params.Deployment)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter métricas do Kubernetes: %w", err)
	}

	// Obtém a configuração do deployment
	deployConfig, err := s.k8sClient.GetDeploymentConfig(ctx, params.Namespace, params.Deployment)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter configuração do deployment: %w", err)
	}

	// Calcula as métricas atuais
	current := &types.CurrentMetrics{
		CPU: types.CPUMetrics{
			Average:     k8sMetrics.CPU.Average,
			Peak:        k8sMetrics.CPU.Peak,
			Usage:       k8sMetrics.CPU.Usage,
			Request:     deployConfig.CPU.Request,
			Limit:       deployConfig.CPU.Limit,
			Utilization: calculateUtilization(k8sMetrics.CPU.Usage, deployConfig.CPU.Request),
		},
		Memory: types.MemoryMetrics{
			Average:     k8sMetrics.Memory.Average,
			Peak:        k8sMetrics.Memory.Peak,
			Usage:       k8sMetrics.Memory.Usage,
			Request:     deployConfig.Memory.Request,
			Limit:       deployConfig.Memory.Limit,
			Utilization: calculateUtilization(k8sMetrics.Memory.Usage, deployConfig.Memory.Request),
		},
		Pods: types.PodMetrics{
			Running:     k8sMetrics.Pods.Running,
			Replicas:    deployConfig.Pods.Replicas,
			MinReplicas: deployConfig.Pods.MinReplicas,
			MaxReplicas: deployConfig.Pods.MaxReplicas,
		},
	}

	metadata := &types.MetricsMetadata{
		CollectedAt: time.Now(),
		TimeWindow:  params.Period,
	}

	return &types.ResourceMetrics{
		Current:  current,
		Metadata: metadata,
	}, nil
}

// GetHistoricalMetrics retorna as métricas históricas de um deployment
func (s *Service) GetHistoricalMetrics(ctx context.Context, params types.MetricsParams) (*types.HistoricalMetrics, error) {
	fmt.Printf("Obtendo métricas históricas para %s no namespace %s (período: %s)\n",
		params.Deployment, params.Namespace, params.Period)

	// Define o período de análise
	endTime := time.Now()
	startTime := endTime.Add(-params.Period)
	step := time.Hour // Intervalo de 1 hora entre pontos

	fmt.Printf("Período de análise: %s até %s (step: %s)\n",
		startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), step)

	// Queries para CPU
	cpuQuery := fmt.Sprintf(`avg(rate(container_cpu_usage_seconds_total{namespace="%s",pod=~"%s.*",container!="POD"}[5m]))`,
		params.Namespace, params.Deployment)
	cpuPeakQuery := fmt.Sprintf(`max(rate(container_cpu_usage_seconds_total{namespace="%s",pod=~"%s.*",container!="POD"}[5m]))`,
		params.Namespace, params.Deployment)

	fmt.Printf("Query CPU média: %s\n", cpuQuery)
	fmt.Printf("Query CPU pico: %s\n", cpuPeakQuery)

	// Queries para Memória (convertendo para GB)
	memoryQuery := fmt.Sprintf(`avg(container_memory_working_set_bytes{namespace="%s",pod=~"%s.*",container!="POD"}) / (1024 * 1024 * 1024)`,
		params.Namespace, params.Deployment)
	memoryPeakQuery := fmt.Sprintf(`max(container_memory_working_set_bytes{namespace="%s",pod=~"%s.*",container!="POD"}) / (1024 * 1024 * 1024)`,
		params.Namespace, params.Deployment)

	fmt.Printf("Query Memória média: %s\n", memoryQuery)
	fmt.Printf("Query Memória pico: %s\n", memoryPeakQuery)

	// Query para Pods
	podsQuery := fmt.Sprintf(`kube_deployment_status_replicas{namespace="%s",deployment="%s"}`,
		params.Namespace, params.Deployment)

	fmt.Printf("Query Pods: %s\n", podsQuery)

	// Obtém métricas históricas de CPU
	fmt.Println("Obtendo métricas históricas de CPU...")
	fmt.Printf("DEBUG: Enviando query para Mimir: %s\n", cpuQuery)
	cpuResult, err := s.mimirClient.QueryRange(ctx, cpuQuery, startTime, endTime, step)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao obter histórico de CPU: %v\n", err)
		return nil, fmt.Errorf("erro ao obter histórico de CPU: %w", err)
	}
	fmt.Printf("DEBUG: Obtidos %d pontos de dados de CPU média\n", len(cpuResult.Values))

	fmt.Printf("DEBUG: Enviando query para Mimir: %s\n", cpuPeakQuery)
	cpuPeakResult, err := s.mimirClient.QueryRange(ctx, cpuPeakQuery, startTime, endTime, step)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao obter picos de CPU: %v\n", err)
		return nil, fmt.Errorf("erro ao obter picos de CPU: %w", err)
	}
	fmt.Printf("DEBUG: Obtidos %d pontos de dados de CPU pico\n", len(cpuPeakResult.Values))

	// Obtém métricas históricas de Memória
	fmt.Println("Obtendo métricas históricas de Memória...")
	memoryResult, err := s.mimirClient.QueryRange(ctx, memoryQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter histórico de memória: %w", err)
	}
	fmt.Printf("Obtidos %d pontos de dados de Memória média\n", len(memoryResult.Values))

	memoryPeakResult, err := s.mimirClient.QueryRange(ctx, memoryPeakQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter picos de memória: %w", err)
	}
	fmt.Printf("Obtidos %d pontos de dados de Memória pico\n", len(memoryPeakResult.Values))

	// Obtém métricas históricas de Pods
	fmt.Println("Obtendo métricas históricas de Pods...")
	podsResult, err := s.mimirClient.QueryRange(ctx, podsQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter histórico de pods: %w", err)
	}
	fmt.Printf("Obtidos %d pontos de dados de Pods\n", len(podsResult.Values))

	// Obtém a configuração do deployment para requests/limits
	fmt.Println("Obtendo configuração do deployment...")
	deployConfig, err := s.k8sClient.GetDeploymentConfig(ctx, params.Namespace, params.Deployment)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter configuração do deployment: %w", err)
	}

	// Processa os resultados
	fmt.Println("Processando resultados...")
	cpuMetrics := make([]types.CPUMetrics, len(cpuResult.Values))
	memoryMetrics := make([]types.MemoryMetrics, len(memoryResult.Values))
	podMetrics := make([]types.PodMetrics, len(podsResult.Values))

	// Processa métricas de CPU
	for i, value := range cpuResult.Values {
		peakValue := cpuPeakResult.Values[i]
		cpuMetrics[i] = types.CPUMetrics{
			Average:     sanitizeFloat64(value.Value),
			Peak:        sanitizeFloat64(peakValue.Value),
			Request:     deployConfig.CPU.Request,
			Limit:       deployConfig.CPU.Limit,
			Utilization: calculateUtilization(value.Value, deployConfig.CPU.Request),
		}
	}

	// Processa métricas de Memória
	for i, value := range memoryResult.Values {
		peakValue := memoryPeakResult.Values[i]
		memoryMetrics[i] = types.MemoryMetrics{
			Average:     sanitizeFloat64(value.Value),
			Peak:        sanitizeFloat64(peakValue.Value),
			Request:     deployConfig.Memory.Request,
			Limit:       deployConfig.Memory.Limit,
			Utilization: calculateUtilization(value.Value, deployConfig.Memory.Request),
		}
	}

	// Processa métricas de Pods
	for i, value := range podsResult.Values {
		podMetrics[i] = types.PodMetrics{
			Running:     int(value.Value),
			Replicas:    deployConfig.Pods.Replicas,
			MinReplicas: deployConfig.Pods.MinReplicas,
			MaxReplicas: deployConfig.Pods.MaxReplicas,
		}
	}

	fmt.Printf("Métricas históricas processadas com sucesso: %d pontos de CPU, %d pontos de Memória, %d pontos de Pods\n",
		len(cpuMetrics), len(memoryMetrics), len(podMetrics))

	return &types.HistoricalMetrics{
		CPU:    cpuMetrics,
		Memory: memoryMetrics,
		Pods:   podMetrics,
		Period: params.Period,
	}, nil
}

// GetTrends retorna as tendências de uso de recursos de um deployment
func (s *Service) GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error) {
	params := types.MetricsParams{
		Namespace:  namespace,
		Deployment: deployment,
		Period:     period,
	}

	trends, err := s.getTrends(ctx, params)
	if err != nil {
		return nil, err
	}

	return &types.TrendsResponse{
		CPU:    trends,
		Memory: trends,
		Pods:   trends,
	}, nil
}

func (s *Service) getTrends(ctx context.Context, params types.MetricsParams) (*types.TrendMetrics, error) {
	// Obtém métricas históricas para análise de tendências
	historical, err := s.GetHistoricalMetrics(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter métricas históricas: %w", err)
	}

	// Calcula tendências de CPU
	cpuTrend := calculateTrend(historical.CPU, func(m types.CPUMetrics) float64 {
		return m.Usage
	})

	// Calcula tendências de memória
	memoryTrend := calculateTrend(historical.Memory, func(m types.MemoryMetrics) float64 {
		return m.Usage
	})

	// Calcula tendências de pods
	podsTrend := calculateTrend(historical.Pods, func(m types.PodMetrics) float64 {
		return float64(m.Running)
	})

	return &types.TrendMetrics{
		CPU:    cpuTrend,
		Memory: memoryTrend,
		Pods:   podsTrend,
	}, nil
}

// calculateTrend calcula a tendência de uma série de métricas
func calculateTrend[T any](metrics []T, getValue func(T) float64) types.TrendData {
	if len(metrics) < 2 {
		return types.TrendData{
			Trend:      0,
			Confidence: 0,
			Direction:  "stable",
		}
	}

	// Calcula a média móvel para suavizar os dados
	windowSize := 6 // 6 horas
	if len(metrics) < windowSize {
		windowSize = len(metrics)
	}

	smoothed := make([]float64, len(metrics)-windowSize+1)
	for i := 0; i <= len(metrics)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += getValue(metrics[i+j])
		}
		smoothed[i] = sum / float64(windowSize)
	}

	// Calcula a tendência linear
	n := float64(len(smoothed))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range smoothed {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calcula o coeficiente de inclinação
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Calcula o R² (coeficiente de determinação)
	meanY := sumY / n
	totalSS := 0.0
	residualSS := 0.0

	for i, y := range smoothed {
		x := float64(i)
		yPred := slope*x + (sumY-slope*sumX)/n
		totalSS += math.Pow(y-meanY, 2)
		residualSS += math.Pow(y-yPred, 2)
	}

	r2 := 1 - (residualSS / totalSS)

	// Determina a direção da tendência
	direction := "stable"
	if math.Abs(slope) > 0.01 {
		if slope > 0 {
			direction = "increase"
		} else {
			direction = "decrease"
		}
	}

	// Normaliza a tendência para porcentagem
	trendPerc := (slope * float64(len(smoothed))) / smoothed[0] * 100

	return types.TrendData{
		Trend:      trendPerc,
		Confidence: r2,
		Direction:  direction,
	}
}

// calculateUtilization calcula a porcentagem de utilização de um recurso
func calculateUtilization(usage, request float64) float64 {
	if request == 0 || usage == 0 {
		return 0
	}
	result := (usage / request) * 100
	if result < 0 || result > 1000 || result != result { // Verifica NaN
		return 0
	}
	return result
}

// sanitizeFloat64 trata valores NaN e Inf
func sanitizeFloat64(value float64) float64 {
	if value != value || value < 0 || value > 1e6 { // Verifica NaN e valores absurdos
		return 0
	}
	return value
}

// AnalyzeResources realiza uma análise detalhada dos recursos
func (s *Service) AnalyzeResources(current *types.ResourceMetrics, historical []*types.ResourceMetrics) *types.ResourceAnalysis {
	// Análise de CPU
	cpuAnalysis := analyzeCPU(current, historical)

	// Análise de Memória
	memAnalysis := analyzeMemory(current, historical)

	// Análise de Pods
	podAnalysis := analyzePods(current, historical)

	// Gera recomendações gerais
	recommendations := generateRecommendations(cpuAnalysis, memAnalysis, podAnalysis)

	// Calcula distribuição
	distribution := calculateDistribution(current, historical)

	// Identifica oportunidades de otimização
	optimizations := identifyOptimizations(cpuAnalysis, memAnalysis, podAnalysis)

	return &types.ResourceAnalysis{
		CPU:             cpuAnalysis,
		Memory:          memAnalysis,
		Pods:            podAnalysis,
		Recommendations: recommendations,
		Distribution:    distribution,
		OptimizationOps: optimizations,
	}
}

// CalculateCosts calcula os custos dos recursos
func (s *Service) CalculateCosts(ctx context.Context, current *types.ResourceMetrics, analysis *types.ResourceAnalysis) (*types.CostAnalysis, error) {
	// Obtém preços dos recursos
	pricing, err := s.pricingClient.GetResourcePricing(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter preços: %w", err)
	}

	// Obtém taxa de câmbio
	exchange, err := s.pricingClient.GetExchangeRate(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter taxa de câmbio: %w", err)
	}

	// Calcula custos atuais
	currentCPUCost := current.Current.CPU.Usage * pricing.CPU * 24 * 30 // Custo mensal
	currentMemCost := current.Current.Memory.Usage * pricing.Memory * 24 * 30

	// Calcula custos recomendados
	recommendedCPUCost := analysis.CPU.Recommended * pricing.CPU * 24 * 30
	recommendedMemCost := analysis.Memory.Recommended * pricing.Memory * 24 * 30

	// Calcula economias
	cpuSavings := currentCPUCost - recommendedCPUCost
	memSavings := currentMemCost - recommendedMemCost
	totalSavings := cpuSavings + memSavings

	// Calcula porcentagem de economia
	currentTotal := currentCPUCost + currentMemCost
	savingsPerc := (totalSavings / currentTotal) * 100

	return &types.CostAnalysis{
		Current: types.CostData{
			Monthly:    currentTotal,
			MonthlyBRL: currentTotal * exchange.Rate,
			Resources: types.ResourceCosts{
				CPU:    currentCPUCost,
				Memory: currentMemCost,
				Total:  currentTotal,
			},
			TotalPods:  current.Current.Pods.Replicas,
			TotalCores: current.Current.CPU.Usage,
			TotalMemGB: current.Current.Memory.Usage,
		},
		Recommended: types.CostData{
			Monthly:    recommendedCPUCost + recommendedMemCost,
			MonthlyBRL: (recommendedCPUCost + recommendedMemCost) * exchange.Rate,
			Resources: types.ResourceCosts{
				CPU:    recommendedCPUCost,
				Memory: recommendedMemCost,
				Total:  recommendedCPUCost + recommendedMemCost,
			},
			TotalPods:  analysis.Pods.Recommended,
			TotalCores: analysis.CPU.Recommended,
			TotalMemGB: analysis.Memory.Recommended,
		},
		Savings: types.ResourceCosts{
			CPU:    cpuSavings,
			Memory: memSavings,
			Total:  totalSavings,
		},
		SavingsPerc:  savingsPerc,
		Currency:     "USD",
		ExchangeRate: exchange.Rate,
	}, nil
}

// GenerateAlerts gera alertas baseados nas métricas
func (s *Service) GenerateAlerts(current *types.ResourceMetrics, historical []*types.ResourceMetrics) []types.Alert {
	alerts := []types.Alert{}

	// Alerta de CPU alta
	if current.Current.CPU.Utilization > 80 {
		alerts = append(alerts, types.Alert{
			Type:       "high_cpu",
			Severity:   "warning",
			Message:    fmt.Sprintf("Alta utilização de CPU: %.1f%%", current.Current.CPU.Utilization),
			Resource:   "cpu",
			CurrentVal: current.Current.CPU.Utilization,
			Threshold:  80,
		})
	}

	// Alerta de memória alta
	if current.Current.Memory.Utilization > 80 {
		alerts = append(alerts, types.Alert{
			Type:       "high_memory",
			Severity:   "warning",
			Message:    fmt.Sprintf("Alta utilização de memória: %.1f%%", current.Current.Memory.Utilization),
			Resource:   "memory",
			CurrentVal: current.Current.Memory.Utilization,
			Threshold:  80,
		})
	}

	// Alerta de pods no limite
	if current.Current.Pods.Running == current.Current.Pods.MaxReplicas {
		alerts = append(alerts, types.Alert{
			Type:       "max_pods",
			Severity:   "warning",
			Message:    "Número máximo de pods atingido",
			Resource:   "pods",
			CurrentVal: float64(current.Current.Pods.Running),
			Threshold:  float64(current.Current.Pods.MaxReplicas),
		})
	}

	return alerts
}

// analyzeCPU realiza análise detalhada de CPU
func analyzeCPU(current *types.ResourceMetrics, historical []*types.ResourceMetrics) types.CPUAnalysis {
	var totalCPU, peakCPU float64
	for _, m := range historical {
		totalCPU += m.Current.CPU.Usage
		if m.Current.CPU.Peak > peakCPU {
			peakCPU = m.Current.CPU.Peak
		}
	}
	avgCPU := totalCPU / float64(len(historical))

	// Calculando a distribuição
	distribution := make(map[string]int)
	for _, m := range historical {
		if m.Current.CPU.Usage < 0.2 {
			distribution["0-200m"]++
		} else if m.Current.CPU.Usage < 0.4 {
			distribution["200-400m"]++
		} else if m.Current.CPU.Usage < 0.6 {
			distribution["400-600m"]++
		} else if m.Current.CPU.Usage < 0.8 {
			distribution["600-800m"]++
		} else {
			distribution["800m+"]++
		}
	}

	// Convertendo mapa em slice de DistributionBucket
	buckets := []types.DistributionBucket{
		{Range: "0-200m", Count: distribution["0-200m"], StartVal: 0, EndVal: 0.2},
		{Range: "200-400m", Count: distribution["200-400m"], StartVal: 0.2, EndVal: 0.4},
		{Range: "400-600m", Count: distribution["400-600m"], StartVal: 0.4, EndVal: 0.6},
		{Range: "600-800m", Count: distribution["600-800m"], StartVal: 0.6, EndVal: 0.8},
		{Range: "800m+", Count: distribution["800m+"], StartVal: 0.8, EndVal: 1.0},
	}

	// Calculando percentuais
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	for i := range buckets {
		buckets[i].Percent = float64(buckets[i].Count) / float64(total) * 100
	}

	// Gerando recomendações
	recommendedRequest := avgCPU * 1.5 // 50% buffer
	if recommendedRequest > current.Current.CPU.Request {
		recommendedRequest = current.Current.CPU.Request
	}

	return types.CPUAnalysis{
		CurrentUsage:  current.Current.CPU.Usage,
		HistoricalAvg: avgCPU,
		Peak:          peakCPU,
		Distribution:  buckets,
		Recommendation: &types.Recommendation{
			Current:    current.Current.CPU.Request,
			Suggested:  recommendedRequest,
			Confidence: 0.95,
			Reason:     fmt.Sprintf("Baseado no uso histórico de %.2fm com pico de %.2fm", avgCPU*1000, peakCPU*1000),
			Priority:   1,
		},
		Recommended: recommendedRequest,
	}
}

// analyzeMemory realiza análise detalhada de memória
func analyzeMemory(current *types.ResourceMetrics, historical []*types.ResourceMetrics) types.MemoryAnalysis {
	var totalMem, peakMem float64
	for _, m := range historical {
		totalMem += m.Current.Memory.Usage
		if m.Current.Memory.Peak > peakMem {
			peakMem = m.Current.Memory.Peak
		}
	}
	avgMem := totalMem / float64(len(historical))

	// Calculando a distribuição
	distribution := make(map[string]int)
	for _, m := range historical {
		memGi := m.Current.Memory.Usage
		if memGi < 0.2 {
			distribution["0-200Mi"]++
		} else if memGi < 0.4 {
			distribution["200-400Mi"]++
		} else if memGi < 0.6 {
			distribution["400-600Mi"]++
		} else if memGi < 0.8 {
			distribution["600-800Mi"]++
		} else {
			distribution["800Mi+"]++
		}
	}

	// Convertendo mapa em slice de DistributionBucket
	buckets := []types.DistributionBucket{
		{Range: "0-200Mi", Count: distribution["0-200Mi"], StartVal: 0, EndVal: 0.2},
		{Range: "200-400Mi", Count: distribution["200-400Mi"], StartVal: 0.2, EndVal: 0.4},
		{Range: "400-600Mi", Count: distribution["400-600Mi"], StartVal: 0.4, EndVal: 0.6},
		{Range: "600-800Mi", Count: distribution["600-800Mi"], StartVal: 0.6, EndVal: 0.8},
		{Range: "800Mi+", Count: distribution["800Mi+"], StartVal: 0.8, EndVal: 1.0},
	}

	// Calculando percentuais
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	for i := range buckets {
		buckets[i].Percent = float64(buckets[i].Count) / float64(total) * 100
	}

	// Gerando recomendações
	recommendedRequest := avgMem * 1.25 // 25% buffer
	if recommendedRequest > current.Current.Memory.Request {
		recommendedRequest = current.Current.Memory.Request
	}

	return types.MemoryAnalysis{
		CurrentUsage:  current.Current.Memory.Usage,
		HistoricalAvg: avgMem,
		Peak:          peakMem,
		Distribution:  buckets,
		Recommendation: &types.Recommendation{
			Current:    current.Current.Memory.Request,
			Suggested:  recommendedRequest,
			Confidence: 0.95,
			Reason:     fmt.Sprintf("Baseado no uso histórico de %.0fMi com pico de %.0fMi", avgMem*1024, peakMem*1024),
			Priority:   1,
		},
		Recommended: recommendedRequest,
	}
}

// analyzePods realiza análise detalhada de pods
func analyzePods(current *types.ResourceMetrics, historical []*types.ResourceMetrics) types.PodAnalysis {
	var totalPods, peakPods int
	for _, m := range historical {
		totalPods += m.Current.Pods.Running
		if m.Current.Pods.Running > peakPods {
			peakPods = m.Current.Pods.Running
		}
	}
	avgPods := float64(totalPods) / float64(len(historical))

	// Gerando recomendações
	recommendedMin := int(math.Ceil(avgPods * 0.7)) // 30% abaixo da média
	if recommendedMin < 2 {
		recommendedMin = 2
	}
	recommendedMax := int(math.Ceil(float64(peakPods) * 1.4)) // 40% acima do pico

	return types.PodAnalysis{
		CurrentRunning: current.Current.Pods.Running,
		HistoricalAvg:  avgPods,
		Peak:           peakPods,
		Recommendation: &types.Recommendation{
			Current:    float64(current.Current.Pods.MinReplicas),
			Suggested:  float64(recommendedMin),
			Confidence: 0.9,
			Reason:     fmt.Sprintf("Baseado no uso histórico de %.1f pods com pico de %d", avgPods, peakPods),
		},
		MaxReplicasRecommendation: &types.Recommendation{
			Current:    float64(current.Current.Pods.MaxReplicas),
			Suggested:  float64(recommendedMax),
			Confidence: 0.85,
			Reason:     fmt.Sprintf("Baseado no pico histórico de %d pods com 40%% de buffer", peakPods),
		},
	}
}

// generateRecommendations gera recomendações baseadas nas análises
func generateRecommendations(cpu types.CPUAnalysis, mem types.MemoryAnalysis, pods types.PodAnalysis) []types.Recommendation {
	recommendations := make([]types.Recommendation, 0, 3)

	// Copia as recomendações de cada análise
	if cpu.Recommendation != nil {
		recommendations = append(recommendations, *cpu.Recommendation)
	}
	if mem.Recommendation != nil {
		recommendations = append(recommendations, *mem.Recommendation)
	}
	if pods.Recommendation != nil {
		recommendations = append(recommendations, *pods.Recommendation)
	}

	// Ordena por prioridade
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority < recommendations[j].Priority
	})

	return recommendations
}

// calculateDistribution calcula a distribuição geral de recursos
func calculateDistribution(current *types.ResourceMetrics, historical []*types.ResourceMetrics) types.ResourceDistribution {
	return types.ResourceDistribution{
		CPU:    calculateCPUDistribution(current, historical),
		Memory: calculateMemoryDistribution(current, historical),
	}
}

// calculateCPUDistribution calcula a distribuição de CPU
func calculateCPUDistribution(current *types.ResourceMetrics, historical []*types.ResourceMetrics) []types.DistributionBucket {
	buckets := make([]types.DistributionBucket, 5)
	bucketSize := 0.2 // 200m

	for i := range buckets {
		start := float64(i) * bucketSize
		end := start + bucketSize
		buckets[i] = types.DistributionBucket{
			Range:    fmt.Sprintf("%.0f-%.0fm", start*1000, end*1000),
			StartVal: start,
			EndVal:   end,
		}
	}

	return buckets
}

// calculateMemoryDistribution calcula a distribuição de memória
func calculateMemoryDistribution(current *types.ResourceMetrics, historical []*types.ResourceMetrics) []types.DistributionBucket {
	buckets := make([]types.DistributionBucket, 5)
	bucketSize := current.Current.Memory.Limit / 5

	for i := range buckets {
		start := float64(i) * bucketSize
		end := start + bucketSize
		buckets[i] = types.DistributionBucket{
			Range:    fmt.Sprintf("%.0f-%.0fMi", start*1024, end*1024),
			StartVal: start,
			EndVal:   end,
		}
	}

	return buckets
}

// identifyOptimizations identifica oportunidades de otimização
func identifyOptimizations(cpu types.CPUAnalysis, mem types.MemoryAnalysis, pods types.PodAnalysis) []types.OptimizationOp {
	ops := make([]types.OptimizationOp, 0)

	// Adiciona otimizações baseadas nas recomendações
	if cpu.Recommendation != nil && cpu.CurrentUsage > cpu.Recommended {
		ops = append(ops, types.OptimizationOp{
			Resource:    "cpu",
			Type:        "request",
			Current:     cpu.CurrentUsage,
			Recommended: cpu.Recommended,
			Savings:     cpu.CurrentUsage - cpu.Recommended,
			Reason:      cpu.Recommendation.Reason,
		})
	}

	if mem.Recommendation != nil && mem.CurrentUsage > mem.Recommended {
		ops = append(ops, types.OptimizationOp{
			Resource:    "memory",
			Type:        "request",
			Current:     mem.CurrentUsage,
			Recommended: mem.Recommended,
			Savings:     mem.CurrentUsage - mem.Recommended,
			Reason:      mem.Recommendation.Reason,
		})
	}

	if pods.Recommendation != nil && float64(pods.CurrentRunning) > float64(pods.Recommended) {
		ops = append(ops, types.OptimizationOp{
			Resource:    "pods",
			Type:        "replicas",
			Current:     float64(pods.CurrentRunning),
			Recommended: float64(pods.Recommended),
			Savings:     float64(pods.CurrentRunning - pods.Recommended),
			Reason:      pods.Recommendation.Reason,
		})
	}

	return ops
}

// GetMetrics retorna as métricas atuais e históricas de um deployment
func (s *Service) GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error) {
	params := types.MetricsParams{
		Namespace:  namespace,
		Deployment: deployment,
		Period:     period,
	}

	// Obtém métricas atuais
	current, err := s.GetCurrentMetrics(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter métricas atuais: %w", err)
	}

	// Obtém métricas históricas
	historical, err := s.GetHistoricalMetrics(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter métricas históricas: %w", err)
	}

	// Converte métricas históricas para o formato esperado
	minLen := len(historical.CPU)
	if len(historical.Memory) < minLen {
		minLen = len(historical.Memory)
	}
	if len(historical.Pods) < minLen {
		minLen = len(historical.Pods)
	}

	historicalMetrics := make([]*types.ResourceMetrics, minLen)
	for i := 0; i < minLen; i++ {
		historicalMetrics[i] = &types.ResourceMetrics{
			Current: &types.CurrentMetrics{
				CPU:    historical.CPU[i],
				Memory: historical.Memory[i],
				Pods:   historical.Pods[i],
			},
		}
	}

	return &types.MetricsResponse{
		Current:    current,
		Historical: historicalMetrics,
		Metadata:   current.Metadata,
	}, nil
}
