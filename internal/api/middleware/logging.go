// Package middleware provides HTTP middleware for the API.
// This package implements common functionalities such as logging,
// authentication, rate limiting and other request interceptions.
package middleware

import (
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs HTTP request information
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request start time
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Prepare log fields
		fields := []logger.Field{
			logger.NewField("method", c.Request.Method),
			logger.NewField("path", c.Request.URL.Path),
			logger.NewField("status", c.Writer.Status()),
			logger.NewField("duration", duration.String()),
			logger.NewField("client_ip", c.ClientIP()),
			logger.NewField("user_agent", c.Request.UserAgent()),
		}

		// Add query params if they exist
		if len(c.Request.URL.RawQuery) > 0 {
			fields = append(fields, logger.NewField("query", c.Request.URL.RawQuery))
		}

		// Add errors if they exist
		if len(c.Errors) > 0 {
			fields = append(fields, logger.NewField("errors", c.Errors.String()))
		}

		// Log request with appropriate level
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

// ErrorLogger is a middleware that logs unhandled errors
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log each error if they exist
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

// RecoveryLogger is a middleware that recovers from panics and logs the error
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
