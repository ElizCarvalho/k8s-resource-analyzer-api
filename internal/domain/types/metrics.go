package types

import "time"

// ===== Tipos Base de Métricas =====

// ResourceMetrics representa métricas de um recurso (CPU ou Memória)
type ResourceMetrics struct {
	Usage        float64        `json:"usage"`
	Request      float64        `json:"request"`
	Limit        float64        `json:"limit"`
	Average      float64        `json:"average"`
	Peak         float64        `json:"peak"`
	Utilization  float64        `json:"utilization"`
	Distribution map[string]int `json:"distribution"`
	Timestamp    int64          `json:"timestamp,omitempty"`
}

// PodMetrics representa métricas de pods
type PodMetrics struct {
	Running     int     `json:"running"`
	Total       int     `json:"total,omitempty"`
	Replicas    int     `json:"replicas"`
	MinReplicas int     `json:"minReplicas"`
	MaxReplicas int     `json:"maxReplicas"`
	Utilization float64 `json:"utilization"`
	Timestamp   int64   `json:"timestamp,omitempty"`
}

// ===== Tipos de Resposta =====

// CurrentMetrics representa métricas atuais
type CurrentMetrics struct {
	Deployment struct {
		Config struct {
			CPU struct {
				Request float64 `json:"request"` // em milicores
				Limit   float64 `json:"limit"`   // em milicores
			} `json:"cpu"`
			Memory struct {
				Request float64 `json:"request"` // em Mi
				Limit   float64 `json:"limit"`   // em Mi
			} `json:"memory"`
			HPA struct {
				MinReplicas int     `json:"minReplicas"`
				MaxReplicas int     `json:"maxReplicas"`
				TargetCPU   float64 `json:"targetCPU"` // em percentual
			} `json:"hpa"`
		} `json:"config"`
	} `json:"deployment"`
	Analysis struct {
		CPU struct {
			Distribution map[string]float64 `json:"distribution"` // faixas em milicores
			Alerts       struct {
				HighCPU    int `json:"highCPU"`    // pods com CPU > 999m
				NearLimit  int `json:"nearLimit"`  // pods com CPU > 900m
				HighMemory int `json:"highMemory"` // pods com memória > 800Mi
			} `json:"alerts"`
			Usage struct {
				Current struct {
					Average float64 `json:"average"` // em milicores
					Peak    float64 `json:"peak"`    // em milicores
				} `json:"current"`
				Historical struct {
					Average float64 `json:"average"` // em milicores
					Peak    float64 `json:"peak"`    // em milicores
				} `json:"historical"`
			} `json:"usage"`
		} `json:"cpu"`
		Memory struct {
			Usage struct {
				Current struct {
					Average float64 `json:"average"` // em Mi
					Peak    float64 `json:"peak"`    // em Mi
				} `json:"current"`
				Historical struct {
					Average float64 `json:"average"` // em Mi
					Peak    float64 `json:"peak"`    // em Mi
				} `json:"historical"`
			} `json:"usage"`
		} `json:"memory"`
	} `json:"analysis"`
	CPU    *ResourceMetrics `json:"cpu"`
	Memory *ResourceMetrics `json:"memory"`
	Pods   *PodMetrics      `json:"pods"`
}

// HistoricalMetrics representa métricas históricas
type HistoricalMetrics struct {
	CPU    []*ResourceMetrics `json:"cpu"`
	Memory []*ResourceMetrics `json:"memory"`
	Pods   []*PodMetrics      `json:"pods"`
}

// TrendsResponse representa a resposta com tendências
type TrendsResponse struct {
	CPU    *TrendMetrics `json:"cpu"`
	Memory *TrendMetrics `json:"memory"`
	Pods   *TrendMetrics `json:"pods"`
}

// MetricsResponse representa a resposta com métricas de um deployment
type MetricsResponse struct {
	Current    *CurrentMetrics                 `json:"current"`
	Historical *HistoricalMetrics              `json:"historical"`
	Analysis   *ResourceRecommendationAnalysis `json:"analysis"`
	Costs      *CostAnalysis                   `json:"costs"`
	Metadata   struct {
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
	} `json:"metadata"`
}

// ===== Tipos de Metadados e Parâmetros =====

// MetricsParams representa os parâmetros para consulta de métricas
type MetricsParams struct {
	Namespace  string        `json:"namespace"`
	Deployment string        `json:"deployment"`
	Period     time.Duration `json:"period"`
}

// MetricsMetadata representa metadados das métricas
type MetricsMetadata struct {
	CollectedAt time.Time `json:"collectedAt"`
	TimeWindow  string    `json:"timeWindow"`
	Cluster     string    `json:"cluster"`
	Sources     []string  `json:"sources"`
	Reliability struct {
		CPU    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
		Pods   float64 `json:"pods"`
	} `json:"reliability"`
}

// ===== Tipos de Análise =====

// TrendMetrics representa métricas de tendência
type TrendMetrics struct {
	Trend      float64 `json:"trend"`
	Confidence float64 `json:"confidence"`
	Period     string  `json:"period"`
}

// DistributionBucket representa um bucket na distribuição de recursos
type DistributionBucket struct {
	Range    string  `json:"range"`
	Count    int     `json:"count"`
	Percent  float64 `json:"percent"`
	StartVal float64 `json:"startVal"`
	EndVal   float64 `json:"endVal"`
}

// Alert representa um alerta sobre recursos
type Alert struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"` // "critical", "warning", "info"
	Message     string  `json:"message"`
	Resource    string  `json:"resource"`
	CurrentVal  float64 `json:"currentVal"`
	Threshold   float64 `json:"threshold"`
	Occurrences int     `json:"occurrences"`
}

// ResourceRecommendation representa uma recomendação para um recurso
type ResourceRecommendation struct {
	Status         string              `json:"status"`
	Recommendation *ResourceSuggestion `json:"recommendation,omitempty"`
}

// PodRecommendation representa uma recomendação para pods
type PodRecommendation struct {
	Status         string    `json:"status"`
	Recommendation *PodCount `json:"recommendation,omitempty"`
}

// PodCount representa a contagem de pods
type PodCount struct {
	Current   int `json:"current"`
	Suggested int `json:"suggested"`
	Min       int `json:"min"`
	Max       int `json:"max"`
}

// ResourceSuggestion representa uma sugestão de recurso
type ResourceSuggestion struct {
	Current   float64 `json:"current"`
	Suggested float64 `json:"suggested"`
	Action    string  `json:"action"`
}

// ResourceRecommendationAnalysis representa a análise completa dos recursos
type ResourceRecommendationAnalysis struct {
	CPU    *ResourceRecommendation `json:"cpu"`
	Memory *ResourceRecommendation `json:"memory"`
	Pods   *PodRecommendation      `json:"pods"`
}
