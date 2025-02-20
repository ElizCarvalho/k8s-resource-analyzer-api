package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/analyzer"
	"github.com/gin-gonic/gin"
)

// AnalyzerHandler é o handler para análise de recursos
type AnalyzerHandler struct {
	resourceAnalyzer analyzer.ResourceAnalyzer
}

// NewAnalyzerHandler cria uma nova instância do AnalyzerHandler
func NewAnalyzerHandler(resourceAnalyzer analyzer.ResourceAnalyzer) *AnalyzerHandler {
	return &AnalyzerHandler{
		resourceAnalyzer: resourceAnalyzer,
	}
}

// GetMetricsRequest representa o request para obter métricas
type GetMetricsRequest struct {
	Namespace string `form:"namespace" binding:"required"`
	Period    string `form:"period" binding:"required"`
}

// AnalyzeResources analisa os recursos de um deployment
func (h *AnalyzerHandler) AnalyzeResources(c *gin.Context) {
	var req GetMetricsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parâmetros inválidos: " + err.Error(),
		})
		return
	}

	deployment := c.Param("deployment")
	if deployment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Deployment não especificado",
		})
		return
	}

	fmt.Printf("Requisição recebida: namespace=%s, deployment=%s, period=%s\n", req.Namespace, deployment, req.Period)

	// Converte o período para time.Duration
	period, err := time.ParseDuration(req.Period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Período inválido: " + err.Error(),
		})
		return
	}

	// Obtém métricas
	fmt.Printf("Obtendo métricas para %s/%s (período: %s)...\n", req.Namespace, deployment, period)
	metricsResponse, err := h.resourceAnalyzer.GetMetrics(c.Request.Context(), req.Namespace, deployment, period)
	if err != nil {
		fmt.Printf("Erro ao obter métricas: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao obter métricas: " + err.Error(),
		})
		return
	}
	fmt.Printf("Métricas obtidas com sucesso\n")

	// Obtém tendências
	fmt.Printf("Obtendo tendências...\n")
	trendsResponse, err := h.resourceAnalyzer.GetTrends(c.Request.Context(), req.Namespace, deployment, period)
	if err != nil {
		fmt.Printf("Erro ao obter tendências: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao obter tendências: " + err.Error(),
		})
		return
	}
	fmt.Printf("Tendências obtidas com sucesso\n")

	// Analisa recursos
	fmt.Printf("Analisando recursos...\n")
	resourceAnalysis := h.resourceAnalyzer.AnalyzeResources(metricsResponse.Current, metricsResponse.Historical)
	fmt.Printf("Análise de recursos concluída\n")

	// Calcula custos
	fmt.Printf("Calculando custos...\n")
	costAnalysis, err := h.resourceAnalyzer.CalculateCosts(c.Request.Context(), metricsResponse.Current, metricsResponse.Analysis)
	if err != nil {
		fmt.Printf("Erro ao calcular custos: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao calcular custos: " + err.Error(),
		})
		return
	}
	fmt.Printf("Custos calculados com sucesso\n")

	// Gera alertas
	fmt.Printf("Gerando alertas...\n")
	alerts := h.resourceAnalyzer.GenerateAlerts(metricsResponse.Current, metricsResponse.Historical)
	fmt.Printf("Alertas gerados com sucesso\n")

	// Monta a resposta
	response := gin.H{
		"current":    metricsResponse.Current,
		"historical": metricsResponse.Historical,
		"metadata":   metricsResponse.Metadata,
		"trends":     trendsResponse,
		"analysis":   resourceAnalysis,
		"costs":      costAnalysis,
		"alerts":     alerts,
	}

	fmt.Printf("Enviando resposta...\n")
	c.JSON(http.StatusOK, response)
}
