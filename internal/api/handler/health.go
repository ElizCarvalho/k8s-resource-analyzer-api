package handler

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/response"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/version"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var startTime = time.Now()

// HealthResponse representa a resposta do health check
type HealthResponse struct {
	Status       string            `json:"status" example:"healthy"`
	Timestamp    time.Time         `json:"timestamp" example:"2024-02-18T00:00:00Z"`
	Version      string            `json:"version" example:"1.0.0"`
	Environment  string            `json:"environment" example:"development"`
	Uptime       string            `json:"uptime" example:"24h0m0s"`
	System       SystemInfo        `json:"system"`
	Dependencies map[string]Status `json:"dependencies"`
}

// SystemInfo representa informações do sistema
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
}

// Status representa o status de uma dependência
type Status struct {
	Status  string `json:"status" example:"healthy"`
	Message string `json:"message,omitempty" example:"conectado com sucesso"`
	Error   string `json:"error,omitempty" example:"timeout ao conectar"`
}

// K8sClient interface para o cliente Kubernetes
type K8sClient interface {
	CheckConnection(ctx context.Context) error
}

// MimirClient interface para o cliente Mimir
type MimirClient interface {
	CheckConnection(ctx context.Context) error
}

// HealthHandler é o handler para health check
type HealthHandler struct {
	k8sClient   K8sClient
	mimirClient MimirClient
}

// NewHealthHandler cria uma nova instância do HealthHandler
func NewHealthHandler(k8sClient K8sClient, mimirClient MimirClient) *HealthHandler {
	return &HealthHandler{
		k8sClient:   k8sClient,
		mimirClient: mimirClient,
	}
}

// Check verifica a saúde da API
// @Summary Verifica a saúde da API
// @Description Endpoint para verificar se a API está funcionando corretamente
// @Tags health
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	requestID := uuid.New().String()

	// Cria um novo logger com o request ID
	log := logger.NewLogger().With("request_id", requestID)

	// Adiciona o logger ao contexto
	ctx := logger.WithContext(c.Request.Context(), log)
	c.Request = c.Request.WithContext(ctx)

	log.Info("iniciando health check")

	// Verifica as dependências
	dependencies := h.checkDependencies(ctx)

	// Determina o status geral baseado nas dependências
	status := "healthy"
	for _, dep := range dependencies {
		if dep.Status != "healthy" {
			status = "degraded"
			break
		}
	}

	healthResponse := HealthResponse{
		Status:      status,
		Timestamp:   time.Now(),
		Version:     version.GetVersion(),
		Environment: os.Getenv("GIN_MODE"),
		Uptime:      time.Since(startTime).Round(time.Second).String(),
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
		},
		Dependencies: dependencies,
	}

	log.Info("health check concluído com sucesso")
	response.SuccessWithRequestID(c, "API funcionando normalmente", healthResponse, requestID)
}

func (h *HealthHandler) checkDependencies(ctx context.Context) map[string]Status {
	dependencies := make(map[string]Status)

	// Verifica Kubernetes
	if err := h.k8sClient.CheckConnection(ctx); err != nil {
		dependencies["kubernetes"] = Status{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	} else {
		dependencies["kubernetes"] = Status{
			Status:  "healthy",
			Message: "conectado ao cluster",
		}
	}

	// Verifica Mimir
	if err := h.mimirClient.CheckConnection(ctx); err != nil {
		dependencies["mimir"] = Status{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	} else {
		dependencies["mimir"] = Status{
			Status:  "healthy",
			Message: "conectado ao serviço",
		}
	}

	return dependencies
}
