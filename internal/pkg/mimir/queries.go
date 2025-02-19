package mimir

import "fmt"

// Queries PromQL para métricas de recursos
const (
	// Formato da query para uso de CPU por pod (taxa de 5 minutos)
	podCPUUsageQuery = `rate(container_cpu_usage_seconds_total{container!="POD",container!="",pod=~"%s.*"}[5m])`

	// Formato da query para uso de memória por pod
	podMemoryUsageQuery = `container_memory_working_set_bytes{container!="POD",container!="",pod=~"%s.*"}`

	// Formato da query para uso de CPU por deployment (taxa de 5 minutos)
	deploymentCPUUsageQuery = `sum(rate(container_cpu_usage_seconds_total{container!="POD",container!="",namespace="%s",pod=~"%s.*"}[5m])) by (pod)`

	// Formato da query para uso de memória por deployment
	deploymentMemoryUsageQuery = `sum(container_memory_working_set_bytes{container!="POD",container!="",namespace="%s",pod=~"%s.*"}) by (pod)`
)

// GetPodCPUUsageQuery retorna a query para uso de CPU de um pod
func GetPodCPUUsageQuery(podNamePrefix string) string {
	return fmt.Sprintf(podCPUUsageQuery, podNamePrefix)
}

// GetPodMemoryUsageQuery retorna a query para uso de memória de um pod
func GetPodMemoryUsageQuery(podNamePrefix string) string {
	return fmt.Sprintf(podMemoryUsageQuery, podNamePrefix)
}

// GetDeploymentCPUUsageQuery retorna a query para uso de CPU de um deployment
func GetDeploymentCPUUsageQuery(namespace, deploymentName string) string {
	return fmt.Sprintf(deploymentCPUUsageQuery, namespace, deploymentName)
}

// GetDeploymentMemoryUsageQuery retorna a query para uso de memória de um deployment
func GetDeploymentMemoryUsageQuery(namespace, deploymentName string) string {
	return fmt.Sprintf(deploymentMemoryUsageQuery, namespace, deploymentName)
}
