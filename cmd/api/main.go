package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/api/routes"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/analyzer"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/collector"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/k8s"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/mimir"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/config"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/pricing"
	"github.com/gin-gonic/gin"
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
	// Carrega as configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Configura o logger
	logger.Setup()

	// Configura o cliente Kubernetes
	k8sClient, err := k8s.NewClient(&k8s.ClientConfig{
		KubeconfigPath: cfg.K8s.KubeconfigPath,
		InCluster:      cfg.K8s.InCluster,
	})
	if err != nil {
		log.Fatalf("Erro ao criar cliente Kubernetes: %v", err)
	}

	// Configura o cliente Mimir
	mimirClient := mimir.NewClient(&mimir.ClientConfig{
		BaseURL:     cfg.Mimir.URL,
		ServiceName: cfg.Mimir.ServiceName,
		Namespace:   cfg.Mimir.Namespace,
		LocalPort:   cfg.Mimir.LocalPort,
		ServicePort: cfg.Mimir.ServicePort,
		OrgID:       cfg.Mimir.OrgID,
	})

	// Configura o cliente de preços
	pricingClient := pricing.NewClient(&pricing.Config{
		ExchangeURL: cfg.Pricing.ExchangeURL,
		Timeout:     cfg.Pricing.Timeout,
	})

	// Configura o modo do Gin
	gin.SetMode(cfg.Server.GinMode)

	// Cria o coletor de métricas
	metricsCollector := collector.NewK8sMimirCollector(k8sClient, mimirClient)

	// Cria o serviço de análise
	analyzerService := analyzer.NewService(metricsCollector, pricingClient)

	// Configura o router
	router := gin.Default()

	// Configura as rotas
	routes.SetupRoutes(router, k8sClient, mimirClient, analyzerService)

	// Configura o servidor
	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// Inicia o servidor em uma goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Configura o canal para sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando servidor...")

	// Contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tenta desligar o servidor graciosamente
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor desligado com sucesso")
}
