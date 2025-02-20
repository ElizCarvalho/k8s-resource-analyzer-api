package routes

import (
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/handler"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/analyzer"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/k8s"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/mimir"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, k8sClient *k8s.Client, mimirClient *mimir.Client, analyzerService analyzer.ResourceAnalyzer) {
	// Configura os handlers
	analyzerHandler := handler.NewAnalyzerHandler(analyzerService)

	// Configura as rotas
	v1 := router.Group("/api/v1")
	{
		resources := v1.Group("/resources")
		{
			resources.GET("/:deployment/analysis", analyzerHandler.AnalyzeResources)
		}
	}
}
