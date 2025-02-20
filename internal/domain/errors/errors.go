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
	ErrResourceNotFound = errors.New("recurso não encontrado")

	// ErrInvalidMetrics indica que as métricas obtidas são inválidas
	ErrInvalidMetrics = errors.New("métricas inválidas")

	// ErrInvalidConfiguration indica configuração inválida de recursos
	ErrInvalidConfiguration = errors.New("configuração inválida")

	// ErrUnavailableMetrics indica que as métricas não estão disponíveis
	ErrUnavailableMetrics = errors.New("métricas indisponíveis")

	// ErrInvalidPeriod indica período de análise inválido
	ErrInvalidPeriod = errors.New("período inválido")
)

// ResourceError representa um erro relacionado a recursos Kubernetes
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

// Unwrap retorna o erro interno
func (e *ResourceError) Unwrap() error {
	return e.Err
}

// NewResourceNotFoundError cria um novo erro de recurso não encontrado
func NewResourceNotFoundError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrResourceNotFound,
	}
}

// NewInvalidMetricsError cria um novo erro de métricas inválidas
func NewInvalidMetricsError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrInvalidMetrics,
	}
}

// NewInvalidConfigurationError cria um novo erro de configuração inválida
func NewInvalidConfigurationError(resource, message string) error {
	return &ResourceError{
		Resource: resource,
		Message:  message,
		Err:      ErrInvalidConfiguration,
	}
}

// IsResourceNotFound verifica se o erro é do tipo ErrResourceNotFound
func IsResourceNotFound(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}

// IsInvalidMetrics verifica se o erro é do tipo ErrInvalidMetrics
func IsInvalidMetrics(err error) bool {
	return errors.Is(err, ErrInvalidMetrics)
}

// IsInvalidConfiguration verifica se o erro é do tipo ErrInvalidConfiguration
func IsInvalidConfiguration(err error) bool {
	return errors.Is(err, ErrInvalidConfiguration)
}

// IsUnavailableMetrics verifica se o erro é do tipo ErrUnavailableMetrics
func IsUnavailableMetrics(err error) bool {
	return errors.Is(err, ErrUnavailableMetrics)
}
