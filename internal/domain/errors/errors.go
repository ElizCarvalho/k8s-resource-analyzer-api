// Package errors fornece tipos de erro específicos do domínio.
// Este pacote centraliza a definição de erros, facilitando o tratamento
// e a padronização de mensagens de erro em toda a aplicação.
package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrResourceNotFound indica que um recurso não foi encontrado
	ErrResourceNotFound = errors.New("resource not found")

	// ErrInvalidMetrics indica que as métricas obtidas são inválidas
	ErrInvalidMetrics = errors.New("invalid metrics")

	// ErrInvalidConfiguration indica configuração de recurso inválida
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrUnavailableMetrics indica que as métricas não estão disponíveis
	ErrUnavailableMetrics = errors.New("metrics unavailable")

	// ErrInvalidPeriod indica período de análise inválido
	ErrInvalidPeriod = errors.New("invalid period")
)

// ResourceError representa um erro relacionado a recursos do Kubernetes
type ResourceError struct {
	Resource string
	Message  string
	Err      error
}

// Error implementa a interface error
func (e *ResourceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Resource, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Resource, e.Message)
}

// Unwrap returns the internal error
func (e *ResourceError) Unwrap() error {
	return e.Err
}

// NewResourceNotFoundError creates a new resource not found error
func NewResourceNotFoundError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrResourceNotFound,
	}
}

// NewInvalidMetricsError creates a new invalid metrics error
func NewInvalidMetricsError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrInvalidMetrics,
	}
}

// NewInvalidConfigurationError creates a new invalid configuration error
func NewInvalidConfigurationError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrInvalidConfiguration,
	}
}

// IsResourceNotFound checks if the error is of type ErrResourceNotFound
func IsResourceNotFound(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}

// IsInvalidMetrics checks if the error is of type ErrInvalidMetrics
func IsInvalidMetrics(err error) bool {
	return errors.Is(err, ErrInvalidMetrics)
}

// IsInvalidConfiguration checks if the error is of type ErrInvalidConfiguration
func IsInvalidConfiguration(err error) bool {
	return errors.Is(err, ErrInvalidConfiguration)
}

// IsUnavailableMetrics checks if the error is of type ErrUnavailableMetrics
func IsUnavailableMetrics(err error) bool {
	return errors.Is(err, ErrUnavailableMetrics)
}
