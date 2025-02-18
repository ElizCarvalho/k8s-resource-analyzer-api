package handlers

import (
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// PingResponse representa a resposta do endpoint /ping
type PingResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2024-02-18T00:00:00Z"`
}

// @Summary Endpoint de health check
// @Description Retorna pong se a API estiver funcionando
// @Tags health
// @Produce json
// @Success 200 {object} response.Response
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	timestamp := time.Now()

	logger.Info().
		Time("timestamp", timestamp).
		Str("handler", "ping").
		Str("method", "GET").
		Str("path", "/ping").
		Str("ip", c.ClientIP()).
		Msg("Recebida requisição ping")

	data := PingResponse{
		Status:    "ok",
		Timestamp: timestamp,
	}

	response.Success(c, "pong", data)
}
