package types

// ResourceAnalysis representa a análise completa de recursos
type ResourceAnalysis struct {
	CPU    *ResourceTypeAnalysis `json:"cpu"`
	Memory *ResourceTypeAnalysis `json:"memory"`
	Pods   *PodAnalysis          `json:"pods"`
	Status string                `json:"status"`
}

// ResourceTypeAnalysis representa a análise de um tipo de recurso (CPU ou Memória)
type ResourceTypeAnalysis struct {
	CurrentUsage     float64            `json:"currentUsage"`  // em milicores para CPU, Mi para memória
	HistoricalAvg    float64            `json:"historicalAvg"` // em milicores para CPU, Mi para memória
	Peak             float64            `json:"peak"`          // em milicores para CPU, Mi para memória
	Distribution     map[string]float64 `json:"distribution"`  // percentual em cada faixa
	Utilization      float64            `json:"utilization"`   // em percentual
	UtilizationTrend *UtilizationTrend  `json:"utilizationTrend"`
	Recommendation   *Recommendation    `json:"recommendation"`
}

// PodAnalysis representa a análise de pods
type PodAnalysis struct {
	CurrentRunning    int               `json:"currentRunning"`
	HistoricalAvg     float64           `json:"historicalAvg"`
	Peak              int               `json:"peak"`
	ScalingEfficiency float64           `json:"scalingEfficiency"`
	UtilizationTrend  *UtilizationTrend `json:"utilizationTrend"`
	Recommendations   struct {
		Current     *Recommendation `json:"current"`
		MinReplicas *Recommendation `json:"minReplicas"`
		MaxReplicas *Recommendation `json:"maxReplicas"`
	} `json:"recommendations"`
}

// UtilizationTrend representa a tendência de utilização
type UtilizationTrend struct {
	Current     float64      `json:"current"`    // em percentual
	Historical  float64      `json:"historical"` // em percentual
	Trend       float64      `json:"trend"`      // variação percentual
	Pattern     string       `json:"pattern"`    // "increasing", "decreasing", "stable", "fluctuating"
	Seasonality *Seasonality `json:"seasonality"`
}

// Seasonality representa informações de sazonalidade
type Seasonality struct {
	Pattern string `json:"pattern"` // "daily", "weekly", "monthly", "none"
	Period  string `json:"period"`  // duração do ciclo (ex: "24h", "7d", "30d")
}

// Recommendation representa uma recomendação de recurso
type Recommendation struct {
	Current   float64 `json:"current"`   // em milicores para CPU, Mi para memória
	Suggested float64 `json:"suggested"` // em milicores para CPU, Mi para memória
	Reason    string  `json:"reason"`    // "overprovisioned", "underprovisioned", "optimal"
}
