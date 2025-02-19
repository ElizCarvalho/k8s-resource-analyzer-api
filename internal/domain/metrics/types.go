package metrics

import (
	"time"
)

type MetricsResponse struct {
	Current    *CurrentMetrics    `json:"current"`
	Historical *HistoricalMetrics `json:"historical"`
	Metadata   *MetricsMetadata   `json:"metadata"`
}

type TrendsResponse struct {
	CPU    *TrendMetrics `json:"cpu"`
	Memory *TrendMetrics `json:"memory"`
	Pods   *TrendMetrics `json:"pods"`
}

type ResourceMetrics struct {
	Average     float64 `json:"average"`
	Peak        float64 `json:"peak"`
	Usage       float64 `json:"usage"`
	Request     float64 `json:"request"`
	Limit       float64 `json:"limit"`
	Utilization float64 `json:"utilization"`
}

type PodMetrics struct {
	Running     int `json:"running"`
	Replicas    int `json:"replicas"`
	MinReplicas int `json:"minReplicas"`
	MaxReplicas int `json:"maxReplicas"`
}

type CurrentMetrics struct {
	CPU    *ResourceMetrics `json:"cpu"`
	Memory *ResourceMetrics `json:"memory"`
	Pods   *PodMetrics      `json:"pods"`
}

type HistoricalMetrics struct {
	CPU    []*ResourceMetrics `json:"cpu"`
	Memory []*ResourceMetrics `json:"memory"`
	Pods   []*PodMetrics      `json:"pods"`
	Period string             `json:"period"`
}

type MetricsMetadata struct {
	CollectedAt time.Time `json:"collectedAt"`
	TimeWindow  string    `json:"timeWindow"`
}

type TrendMetrics struct {
	Trend      float64 `json:"trend"`
	Confidence float64 `json:"confidence"`
	Period     string  `json:"period"`
}

type CPUAnalysis struct {
	CurrentUsage   float64         `json:"currentUsage"`
	HistoricalAvg  float64         `json:"historicalAvg"`
	Peak           float64         `json:"peak"`
	Distribution   map[string]int  `json:"distribution"`
	Recommendation *Recommendation `json:"recommendation"`
}

type MemoryAnalysis struct {
	CurrentUsage   float64         `json:"currentUsage"`
	HistoricalAvg  float64         `json:"historicalAvg"`
	Peak           float64         `json:"peak"`
	Distribution   map[string]int  `json:"distribution"`
	Recommendation *Recommendation `json:"recommendation"`
}

type PodAnalysis struct {
	CurrentRunning            int             `json:"currentRunning"`
	HistoricalAvg             float64         `json:"historicalAvg"`
	Peak                      int             `json:"peak"`
	Recommendation            *Recommendation `json:"recommendation"`
	MaxReplicasRecommendation *Recommendation `json:"maxReplicasRecommendation"`
}

type ResourceAnalysis struct {
	CPU    *CPUAnalysis    `json:"cpu"`
	Memory *MemoryAnalysis `json:"memory"`
	Pods   *PodAnalysis    `json:"pods"`
}

type Recommendation struct {
	Current    float64 `json:"current"`
	Suggested  float64 `json:"suggested"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type CostAnalysis struct {
	Current     *CostData      `json:"current"`
	Recommended *CostData      `json:"recommended"`
	Savings     *ResourceCosts `json:"savings"`
}

type CostData struct {
	Monthly *ResourceCosts `json:"monthly"`
}

type ResourceCosts struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Total  float64 `json:"total"`
}

type Alert struct {
	Type      string  `json:"type"`
	Severity  string  `json:"severity"`
	Message   string  `json:"message"`
	Current   float64 `json:"current"`
	Threshold float64 `json:"threshold"`
}
