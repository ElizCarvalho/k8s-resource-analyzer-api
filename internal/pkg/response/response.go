package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response é a estrutura base para todas as respostas da API
type Response struct {
	Success   bool        `json:"success"`         // Indica se a requisição foi bem sucedida
	Message   string      `json:"message"`         // Mensagem descritiva
	Data      interface{} `json:"data,omitempty"`  // Dados da resposta (opcional)
	Error     string      `json:"error,omitempty"` // Mensagem de erro (opcional)
	Timestamp time.Time   `json:"timestamp"`       // Timestamp da resposta
	RequestID string      `json:"request_id"`      // ID único da requisição
}

// Success envia uma resposta de sucesso
func Success(c *gin.Context, message string, data interface{}) {
	resp := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: c.GetString("RequestID"),
	}

	c.JSON(http.StatusOK, resp)
}

// Error envia uma resposta de erro
func Error(c *gin.Context, statusCode int, message string, err error) {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	resp := Response{
		Success:   false,
		Message:   message,
		Error:     errorMessage,
		Timestamp: time.Now(),
		RequestID: c.GetString("RequestID"),
	}

	c.JSON(statusCode, resp)
}

// BadRequest envia uma resposta de erro 400
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message, err)
}

// NotFound envia uma resposta de erro 404
func NotFound(c *gin.Context, message string, err error) {
	Error(c, http.StatusNotFound, message, err)
}

// InternalServerError envia uma resposta de erro 500
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}
