// Package handler implementa os handlers HTTP da API.
// Este pacote é responsável por receber as requisições HTTP,
// validar os inputs, chamar os serviços apropriados e formatar as respostas.
package handler

import (
	"net/http"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/errors"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/analyzer"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
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
	// Extrai e valida parâmetros
	var req GetMetricsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Parâmetros inválidos", err,
			logger.NewField("namespace", req.Namespace),
			logger.NewField("period", req.Period),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parâmetros inválidos: " + err.Error(),
		})
		return
	}

	deployment := c.Param("deployment")
	if deployment == "" {
		err := errors.NewInvalidConfigurationError("deployment", "nome não especificado")
		logger.Error("Deployment não especificado", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Info("Requisição recebida",
		logger.NewField("namespace", req.Namespace),
		logger.NewField("deployment", deployment),
		logger.NewField("period", req.Period),
	)

	// Converte o período para time.Duration
	period, err := time.ParseDuration(req.Period)
	if err != nil {
		err = errors.NewInvalidConfigurationError("period", "período inválido")
		logger.Error("Período inválido", err,
			logger.NewField("period", req.Period),
			logger.NewField("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Obtém métricas
	logger.Info("Obtendo métricas",
		logger.NewField("namespace", req.Namespace),
		logger.NewField("deployment", deployment),
		logger.NewField("period", period),
	)

	metricsResponse, err := h.resourceAnalyzer.GetMetrics(c.Request.Context(), req.Namespace, deployment, period)
	if err != nil {
		logger.Error("Erro ao obter métricas", err,
			logger.NewField("namespace", req.Namespace),
			logger.NewField("deployment", deployment),
		)
		status := http.StatusInternalServerError
		if errors.IsResourceNotFound(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Obtém tendências
	logger.Info("Obtendo tendências")
	trendsResponse, err := h.resourceAnalyzer.GetTrends(c.Request.Context(), req.Namespace, deployment, period)
	if err != nil {
		logger.Error("Erro ao obter tendências", err,
			logger.NewField("namespace", req.Namespace),
			logger.NewField("deployment", deployment),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Analisa recursos
	logger.Info("Analisando recursos")
	resourceAnalysis := h.resourceAnalyzer.AnalyzeResources(metricsResponse.Current, metricsResponse.Historical)

	// Calcula custos
	logger.Info("Calculando custos")
	costAnalysis, err := h.resourceAnalyzer.CalculateCosts(c.Request.Context(), metricsResponse.Current, metricsResponse.Analysis)
	if err != nil {
		logger.Error("Erro ao calcular custos", err,
			logger.NewField("namespace", req.Namespace),
			logger.NewField("deployment", deployment),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Gera alertas
	logger.Info("Gerando alertas")
	alerts := h.resourceAnalyzer.GenerateAlerts(metricsResponse.Current, metricsResponse.Historical)

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

	logger.Info("Enviando resposta",
		logger.NewField("namespace", req.Namespace),
		logger.NewField("deployment", deployment),
		logger.NewField("alerts_count", len(alerts)),
	)

	c.JSON(http.StatusOK, response)
}
