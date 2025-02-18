package main

import (
	"os"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/middleware"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/routes"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ElizCarvalho/k8s-resource-analyzer-api/docs" // Importa os docs gerados pelo Swagger
)

// @title K8s Resource Analyzer API
// @version 1.0
// @description API para análise e otimização de recursos Kubernetes com foco em FinOps. Fornece métricas de utilização, recomendações de custos e análise de eficiência dos recursos em clusters Kubernetes.
// @host localhost:9000
// @BasePath /api/v1
// @schemes http https
// @contact.name Elizabeth Carvalho
// @contact.url https://github.com/ElizCarvalho/k8s-resource-analyzer-api
// @contact.email elizabethcarvalh0@yahoo.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @tag.name health
// @tag.description Endpoints para monitoramento da saúde da API

func main() {
	// Inicializa o logger
	logger.Setup()
	log := logger.Logger

	// Configurar modo de execução
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// Inicializar router
	r := gin.Default()

	// Adiciona middleware de RequestID
	r.Use(middleware.RequestID())

	// Configurar rotas
	routes.SetupRoutes(r)

	// Configurar Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:9000/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	// Iniciar servidor
	port := getEnv("PORT", "9000")
	log.Info().
		Str("port", port).
		Msg("🚀 Servidor iniciando...")

	log.Info().
		Str("url", "http://localhost:"+port+"/swagger/index.html").
		Msg("📚 Documentação Swagger disponível")

	if err := r.Run(":" + port); err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Erro ao iniciar servidor")
	}
}

// Utilitário para obter variáveis de ambiente
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
