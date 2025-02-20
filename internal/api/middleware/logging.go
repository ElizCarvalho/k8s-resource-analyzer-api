// Package middleware fornece middlewares para a API HTTP.
// Este pacote implementa funcionalidades comuns como logging,
// autenticação, rate limiting e outras interceptações de requisições.
package middleware

import (
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger é um middleware que loga informações sobre as requisições HTTP
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tempo de início da requisição
		start := time.Now()

		// Processa a requisição
		c.Next()

		// Calcula a duração
		duration := time.Since(start)

		// Prepara os campos do log
		fields := []logger.Field{
			logger.NewField("method", c.Request.Method),
			logger.NewField("path", c.Request.URL.Path),
			logger.NewField("status", c.Writer.Status()),
			logger.NewField("duration", duration.String()),
			logger.NewField("client_ip", c.ClientIP()),
			logger.NewField("user_agent", c.Request.UserAgent()),
		}

		// Adiciona query params se existirem
		if len(c.Request.URL.RawQuery) > 0 {
			fields = append(fields, logger.NewField("query", c.Request.URL.RawQuery))
		}

		// Adiciona erros se existirem
		if len(c.Errors) > 0 {
			fields = append(fields, logger.NewField("errors", c.Errors.String()))
		}

		// Loga a requisição com o nível apropriado
		msg := "Request processed"
		if c.Writer.Status() >= 500 {
			logger.Error(msg, nil, fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn(msg, fields...)
		} else {
			logger.Info(msg, fields...)
		}
	}
}

// ErrorLogger é um middleware que loga erros não tratados
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Se houver erros, loga cada um
		for _, err := range c.Errors {
			logger.Error("Unhandled error", err.Err,
				logger.NewField("method", c.Request.Method),
				logger.NewField("path", c.Request.URL.Path),
				logger.NewField("error_type", err.Type),
				logger.NewField("meta", err.Meta),
			)
		}
	}
}

// RecoveryLogger é um middleware que recupera de panics e loga o erro
func RecoveryLogger() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered", nil,
			logger.NewField("error", recovered),
			logger.NewField("method", c.Request.Method),
			logger.NewField("path", c.Request.URL.Path),
			logger.NewField("client_ip", c.ClientIP()),
		)

		c.AbortWithStatus(500)
	})
}
