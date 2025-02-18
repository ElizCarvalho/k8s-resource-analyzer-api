package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PingResponse representa a resposta do endpoint /ping
type PingResponse struct {
	Message   string    `json:"message" example:"pong"`
	Timestamp time.Time `json:"timestamp" example:"2024-02-18T00:00:00Z"`
	Status    string    `json:"status" example:"ok"`
}

// @Summary Endpoint de health check
// @Description Retorna pong se a API estiver funcionando
// @Tags health
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	timestamp := time.Now()
	log.Printf("Recebida requisição ping em: %v", timestamp)

	response := PingResponse{
		Message:   "pong",
		Timestamp: timestamp,
		Status:    "ok",
	}

	c.JSON(http.StatusOK, response)
}
