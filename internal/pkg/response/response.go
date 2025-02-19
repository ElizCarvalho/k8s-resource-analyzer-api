package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response é a estrutura base para todas as respostas da API
type Response struct {
	Success   bool        `json:"success"`              // Indica se a requisição foi bem sucedida
	Message   string      `json:"message"`              // Mensagem descritiva
	Data      interface{} `json:"data,omitempty"`       // Dados da resposta (opcional)
	Error     string      `json:"error,omitempty"`      // Mensagem de erro (opcional)
	Timestamp time.Time   `json:"timestamp"`            // Timestamp da resposta
	RequestID string      `json:"request_id,omitempty"` // ID único da requisição (opcional)
}

// Success retorna uma resposta de sucesso
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// SuccessWithRequestID retorna uma resposta de sucesso com request ID
func SuccessWithRequestID(c *gin.Context, message string, data interface{}, requestID string) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: requestID,
	})
}

// Error retorna uma resposta de erro
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

// ErrorWithRequestID retorna uma resposta de erro com request ID
func ErrorWithRequestID(c *gin.Context, code int, message string, requestID string) {
	c.JSON(code, Response{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	})
}

// BadRequest envia uma resposta de erro 400
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message)
}

// NotFound envia uma resposta de erro 404
func NotFound(c *gin.Context, message string, err error) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError envia uma resposta de erro 500
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message)
}
