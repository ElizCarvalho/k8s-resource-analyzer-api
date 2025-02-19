package routes

import (
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/handler"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/metrics"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/k8s"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/mimir"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(r *gin.Engine, k8sClient *k8s.Client, mimirClient *mimir.Client, metricsService metrics.MetricsProvider) {
	// Cria os handlers
	healthHandler := handler.NewHealthHandler(k8sClient, mimirClient)
	analyzerHandler := handler.NewAnalyzerHandler(metricsService)

	// Rotas de health check
	r.GET("/health", healthHandler.Check)

	// Grupo de rotas v1
	api := r.Group("/api/v1")
	{
		// Rotas de recursos
		api.GET("/resources/analyze", analyzerHandler.AnalyzeResources)
	}
}
