package mimir

import (
	"os"
	"strconv"
	"time"
)

// RetryConfig contém as configurações de retry
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// TimeoutConfig contém as configurações de timeout por tipo de operação
type TimeoutConfig struct {
	Query      time.Duration
	QueryRange time.Duration
	Connect    time.Duration
}

// CircuitBreakerConfig contém as configurações do circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures      int
	ResetTimeout     time.Duration
	HalfOpenMaxCalls int
}

// Config contém as configurações para conexão com o Mimir
type Config struct {
	// URL base do Mimir (ex: http://localhost:8080)
	BaseURL string

	// Namespace onde o Mimir está instalado
	Namespace string

	// Nome do serviço do Mimir
	ServiceName string

	// Porta local para port-forward
	LocalPort string

	// Porta do serviço Mimir
	ServicePort string

	// ID da organização para autenticação no Mimir
	OrgID string

	// Configurações de retry
	Retry RetryConfig

	// Configurações de timeout
	Timeouts TimeoutConfig

	// Configurações do circuit breaker
	CircuitBreaker CircuitBreakerConfig
}

// NewConfig cria uma nova configuração do Mimir com valores padrão
// que podem ser sobrescritos por variáveis de ambiente
func NewConfig() *Config {
	return &Config{
		BaseURL:     getEnv("MIMIR_URL", "http://localhost:8080"),
		Namespace:   getEnv("MIMIR_NAMESPACE", "monitoring"),
		ServiceName: getEnv("MIMIR_SERVICE", "lgtm-mimir-query-frontend"),
		LocalPort:   getEnv("MIMIR_LOCAL_PORT", "8080"),
		ServicePort: getEnv("MIMIR_SERVICE_PORT", "8080"),
		OrgID:       getEnv("MIMIR_ORG_ID", "anonymous"),
		Retry: RetryConfig{
			MaxRetries:     getEnvAsInt("MIMIR_RETRY_MAX", 3),
			InitialBackoff: getEnvAsDuration("MIMIR_RETRY_INITIAL_BACKOFF", time.Second),
			MaxBackoff:     getEnvAsDuration("MIMIR_RETRY_MAX_BACKOFF", 10*time.Second),
		},
		Timeouts: TimeoutConfig{
			Query:      getEnvAsDuration("MIMIR_TIMEOUT_QUERY", 10*time.Second),
			QueryRange: getEnvAsDuration("MIMIR_TIMEOUT_QUERY_RANGE", 30*time.Second),
			Connect:    getEnvAsDuration("MIMIR_TIMEOUT_CONNECT", 5*time.Second),
		},
		CircuitBreaker: CircuitBreakerConfig{
			MaxFailures:      getEnvAsInt("MIMIR_CB_MAX_FAILURES", 5),
			ResetTimeout:     getEnvAsDuration("MIMIR_CB_RESET_TIMEOUT", 60*time.Second),
			HalfOpenMaxCalls: getEnvAsInt("MIMIR_CB_HALF_OPEN_MAX", 2),
		},
	}
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
