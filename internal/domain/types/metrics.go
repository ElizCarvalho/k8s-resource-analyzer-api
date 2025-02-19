package types

import "time"

// MetricsParams representa os parâmetros para consulta de métricas
type MetricsParams struct {
	Namespace  string
	Deployment string
	Period     time.Duration
}

// QueryResult representa o resultado de uma query no Mimir
type QueryResult struct {
	Value     float64
	Timestamp time.Time
}

// QueryRangeResult representa o resultado de uma query de intervalo no Mimir
type QueryRangeResult struct {
	Values    []QueryResult
	StartTime time.Time
	EndTime   time.Time
}

// K8sMetrics representa as métricas obtidas do Kubernetes
type K8sMetrics struct {
	CPU struct {
		Average float64
		Peak    float64
		Usage   float64
	}
	Memory struct {
		Average float64
		Peak    float64
		Usage   float64
	}
	Pods struct {
		Running int
	}
}

// K8sDeploymentConfig representa a configuração do deployment no Kubernetes
type K8sDeploymentConfig struct {
	CPU struct {
		Request float64
		Limit   float64
	}
	Memory struct {
		Request float64
		Limit   float64
	}
	Pods struct {
		MinReplicas int
		MaxReplicas int
		Replicas    int
		TargetCPU   int
	}
}

// CPUMetrics representa as métricas de CPU
type CPUMetrics struct {
	Average     float64
	Peak        float64
	Usage       float64
	Request     float64
	Limit       float64
	Utilization float64
}

// MemoryMetrics representa as métricas de memória
type MemoryMetrics struct {
	Average     float64
	Peak        float64
	Usage       float64
	Request     float64
	Limit       float64
	Utilization float64
}

// PodMetrics representa as métricas de pods
type PodMetrics struct {
	Running     int
	Replicas    int
	MinReplicas int
	MaxReplicas int
}

// CurrentMetrics representa as métricas atuais
type CurrentMetrics struct {
	CPU    CPUMetrics
	Memory MemoryMetrics
	Pods   PodMetrics
}

// MetricsMetadata representa os metadados das métricas
type MetricsMetadata struct {
	CollectedAt    time.Time
	TimeWindow     time.Duration
	ClusterName    string
	DeploymentName string
	DataSources    DataSourcesInfo
}

// DataSourcesInfo representa informações sobre as fontes de dados
type DataSourcesInfo struct {
	MimirAvailable bool
	K8sAvailable   bool
	CPUSamples     float64
	MemorySamples  float64
	PodSamples     float64
}

// ResourceMetrics representa todas as métricas de um recurso
type ResourceMetrics struct {
	Current  *CurrentMetrics  `json:"current"`
	CPU      *CPUMetrics      `json:"cpu"`
	Memory   *MemoryMetrics   `json:"memory"`
	Pods     *PodMetrics      `json:"pods"`
	Metadata *MetricsMetadata `json:"metadata"`
}

// HistoricalMetrics representa o histórico de métricas
type HistoricalMetrics struct {
	CPU    []CPUMetrics
	Memory []MemoryMetrics
	Pods   []PodMetrics
	Period time.Duration
}

// TrendMetrics representa as tendências de uso
type TrendMetrics struct {
	CPU    TrendData
	Memory TrendData
	Pods   TrendData
}

// TrendData representa os dados de tendência
type TrendData struct {
	Trend      float64
	Confidence float64
	Period     time.Duration
	Direction  string // "increase", "decrease", "stable"
}

// CostAnalysis representa a análise de custos
type CostAnalysis struct {
	Current      CostData      `json:"current"`
	Recommended  CostData      `json:"recommended"`
	Savings      ResourceCosts `json:"savings"`
	SavingsPerc  float64       `json:"savingsPerc"`
	Currency     string        `json:"currency"`
	ExchangeRate float64       `json:"exchangeRate"`
}

// CostData representa os dados de custo
type CostData struct {
	Monthly    float64       `json:"monthly"`
	MonthlyBRL float64       `json:"monthlyBRL"`
	Resources  ResourceCosts `json:"resources"`
	TotalPods  int           `json:"totalPods"`
	TotalCores float64       `json:"totalCores"`
	TotalMemGB float64       `json:"totalMemGB"`
}

// ResourceCosts representa os custos por tipo de recurso
type ResourceCosts struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Total  float64 `json:"total"`
}

// ResourceAnalysis representa a análise detalhada dos recursos
type ResourceAnalysis struct {
	CPU             CPUAnalysis
	Memory          MemoryAnalysis
	Pods            PodAnalysis
	Recommendations []Recommendation
	Distribution    ResourceDistribution
	OptimizationOps []OptimizationOp
}

// CPUAnalysis representa a análise de CPU
type CPUAnalysis struct {
	CurrentUsage    float64              `json:"currentUsage"`
	HistoricalAvg   float64              `json:"historicalAvg"`
	Peak            float64              `json:"peak"`
	Distribution    []DistributionBucket `json:"distribution"`
	Recommendation  *Recommendation      `json:"recommendation"`
	Recommendations []Recommendation     `json:"recommendations"`
	Recommended     float64              `json:"recommended"`
}

// MemoryAnalysis representa a análise de memória
type MemoryAnalysis struct {
	CurrentUsage    float64              `json:"currentUsage"`
	HistoricalAvg   float64              `json:"historicalAvg"`
	Peak            float64              `json:"peak"`
	Distribution    []DistributionBucket `json:"distribution"`
	Recommendation  *Recommendation      `json:"recommendation"`
	Recommendations []Recommendation     `json:"recommendations"`
	Recommended     float64              `json:"recommended"`
}

// PodAnalysis representa a análise de pods
type PodAnalysis struct {
	CurrentRunning            int              `json:"currentRunning"`
	HistoricalAvg             float64          `json:"historicalAvg"`
	Peak                      int              `json:"peak"`
	Recommendation            *Recommendation  `json:"recommendation"`
	MaxReplicasRecommendation *Recommendation  `json:"maxReplicasRecommendation"`
	Recommendations           []Recommendation `json:"recommendations"`
	Recommended               int              `json:"recommended"`
}

// DistributionBucket representa um bucket na distribuição de recursos
type DistributionBucket struct {
	Range    string
	Count    int
	Percent  float64
	StartVal float64
	EndVal   float64
}

// ResourceDistribution representa a distribuição de recursos
type ResourceDistribution struct {
	CPU    []DistributionBucket
	Memory []DistributionBucket
}

// Recommendation representa uma recomendação de recursos
type Recommendation struct {
	Current    float64 `json:"current"`
	Suggested  float64 `json:"suggested"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
	Priority   int     `json:"priority"`
}

// OptimizationOp representa uma oportunidade de otimização
type OptimizationOp struct {
	Resource    string
	Type        string // "request", "limit", "replicas"
	Current     float64
	Recommended float64
	Savings     float64
	Reason      string
}

// Alert representa um alerta sobre recursos
type Alert struct {
	Type        string
	Severity    string // "critical", "warning", "info"
	Message     string
	Resource    string
	CurrentVal  float64
	Threshold   float64
	Occurrences int
}

// MetricsResponse representa a resposta com métricas
type MetricsResponse struct {
	Current    *ResourceMetrics   `json:"current"`
	Historical []*ResourceMetrics `json:"historical"`
	Metadata   *MetricsMetadata   `json:"metadata"`
}

// TrendsResponse representa a resposta com tendências
type TrendsResponse struct {
	CPU    *TrendMetrics `json:"cpu"`
	Memory *TrendMetrics `json:"memory"`
	Pods   *TrendMetrics `json:"pods"`
}
