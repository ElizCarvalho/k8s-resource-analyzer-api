package types

// K8sMetrics representa métricas do Kubernetes
type K8sMetrics struct {
	CPU struct {
		Usage       float64 `json:"usage"`       // em milicores
		Average     float64 `json:"average"`     // em milicores
		Peak        float64 `json:"peak"`        // em milicores
		Utilization float64 `json:"utilization"` // em percentual
	} `json:"cpu"`
	Memory struct {
		Usage       float64 `json:"usage"`       // em Mi
		Average     float64 `json:"average"`     // em Mi
		Peak        float64 `json:"peak"`        // em Mi
		Utilization float64 `json:"utilization"` // em percentual
	} `json:"memory"`
	Pods struct {
		Running     int     `json:"running"`
		Utilization float64 `json:"utilization"` // em percentual
	} `json:"pods"`
}

// K8sDeploymentConfig representa configuração de um deployment
type K8sDeploymentConfig struct {
	CPU struct {
		Request float64 `json:"request"` // em milicores
		Limit   float64 `json:"limit"`   // em milicores
	} `json:"cpu"`
	Memory struct {
		Request float64 `json:"request"` // em Mi
		Limit   float64 `json:"limit"`   // em Mi
	} `json:"memory"`
	Pods struct {
		Replicas    int     `json:"replicas"`
		MinReplicas int     `json:"minReplicas"`
		MaxReplicas int     `json:"maxReplicas"`
		TargetCPU   float64 `json:"targetCPU"` // em percentual
	} `json:"pods"`
	ClusterName string `json:"clusterName"`
}
