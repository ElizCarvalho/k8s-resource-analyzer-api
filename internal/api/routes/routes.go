// Package routes configura as rotas da API HTTP.
// Este pacote define os endpoints disponíveis, seus handlers e middlewares.
package routes

import (
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/handler"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/middleware"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/analyzer"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/k8s"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/mimir"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, k8sClient *k8s.Client, mimirClient *mimir.Client, analyzerService analyzer.ResourceAnalyzer) {
	// Configura middlewares globais
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorLogger())
	router.Use(middleware.RecoveryLogger())

	// Configura os handlers
	analyzerHandler := handler.NewAnalyzerHandler(analyzerService)

	// Grupo de rotas v1
	v1 := router.Group("/api/v1")
	{
		// Endpoints de recursos
		resources := v1.Group("/resources")
		{
			// Análise de recursos
			resources.GET("/:deployment/analysis", analyzerHandler.AnalyzeResources)
		}

		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
	}
}
