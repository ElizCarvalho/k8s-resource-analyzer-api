package routes

import (
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(r *gin.Engine) {
	// Grupo de rotas v1
	api := r.Group("/api/v1")
	{
		// Rotas de health check
		api.GET("/ping", handlers.PingHandler)
	}
}
