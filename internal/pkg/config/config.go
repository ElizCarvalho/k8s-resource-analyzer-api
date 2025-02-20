package config

import (
	"os"
	"strconv"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/errors"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
	"github.com/joho/godotenv"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
	Mimir   MimirConfig
	K8s     K8sConfig
	Pricing PricingConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type MimirConfig struct {
	URL            string
	ServiceName    string
	Namespace      string
	LocalPort      string
	ServicePort    string
	OrgID          string
	RetryMax       int
	RetryBackoff   time.Duration
	MaxBackoff     time.Duration
	TimeoutQuery   time.Duration
	TimeoutRange   time.Duration
	TimeoutConnect time.Duration
	CBMaxFailures  int
	CBResetTimeout time.Duration
	CBHalfOpenMax  int
}

type K8sConfig struct {
	KubeconfigPath string
	InCluster      bool
}

type PricingConfig struct {
	ExchangeURL string
	Timeout     time.Duration
}

// LoadConfig carrega e valida todas as configurações
func LoadConfig() (*Config, error) {
	// Carrega o ambiente correto
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Tenta carregar o arquivo .env específico do ambiente
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		// Se não encontrar, tenta carregar o .env padrão
		if err := godotenv.Load(); err != nil {
			logger.Warn("Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
		}
	} else {
		logger.Info("Configurações carregadas do arquivo",
			logger.NewField("file", envFile),
		)
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnvOrDefault("PORT", "9000"),
			GinMode: getEnvOrDefault("GIN_MODE", "debug"),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
		},
		Mimir: MimirConfig{
			URL:            getEnvOrDefault("MIMIR_URL", "http://localhost:8080"),
			ServiceName:    getEnvOrDefault("MIMIR_SERVICE_NAME", "lgtm-mimir-query-frontend"),
			Namespace:      getEnvOrDefault("MIMIR_NAMESPACE", "monitoring"),
			LocalPort:      getEnvOrDefault("MIMIR_LOCAL_PORT", "8080"),
			ServicePort:    getEnvOrDefault("MIMIR_SERVICE_PORT", "8080"),
			OrgID:          getEnvOrDefault("MIMIR_ORG_ID", "anonymous"),
			RetryMax:       getEnvAsIntOrDefault("MIMIR_RETRY_MAX", 3),
			RetryBackoff:   getEnvAsDurationOrDefault("MIMIR_RETRY_INITIAL_BACKOFF", time.Second),
			MaxBackoff:     getEnvAsDurationOrDefault("MIMIR_RETRY_MAX_BACKOFF", 10*time.Second),
			TimeoutQuery:   getEnvAsDurationOrDefault("MIMIR_TIMEOUT_QUERY", 10*time.Second),
			TimeoutRange:   getEnvAsDurationOrDefault("MIMIR_TIMEOUT_QUERY_RANGE", 30*time.Second),
			TimeoutConnect: getEnvAsDurationOrDefault("MIMIR_TIMEOUT_CONNECT", 5*time.Second),
			CBMaxFailures:  getEnvAsIntOrDefault("MIMIR_CB_MAX_FAILURES", 5),
			CBResetTimeout: getEnvAsDurationOrDefault("MIMIR_CB_RESET_TIMEOUT", 60*time.Second),
			CBHalfOpenMax:  getEnvAsIntOrDefault("MIMIR_CB_HALF_OPEN_MAX", 2),
		},
		K8s: K8sConfig{
			KubeconfigPath: getKubeconfigPath(),
			InCluster:      getEnvOrDefault("IN_CLUSTER", "false") == "true",
		},
		Pricing: PricingConfig{
			ExchangeURL: getEnvOrDefault("EXCHANGE_URL", "https://api.exchangerate.host"),
			Timeout:     30 * time.Second,
		},
	}

	// Valida a configuração
	if err := config.validate(); err != nil {
		return nil, err
	}

	// Loga as configurações carregadas
	config.logConfig()

	return config, nil
}

func (c *Config) validate() error {
	// Validações básicas
	if c.Server.Port == "" {
		return errors.NewInvalidConfigurationError("port", "PORT é obrigatório")
	}

	if c.Mimir.URL == "" {
		return errors.NewInvalidConfigurationError("mimir_url", "MIMIR_URL é obrigatório")
	}

	return nil
}

func (c *Config) logConfig() {
	logger.Info("Configurações carregadas",
		logger.NewField("port", c.Server.Port),
		logger.NewField("gin_mode", c.Server.GinMode),
		logger.NewField("log_level", c.Logging.Level),
		logger.NewField("log_format", c.Logging.Format),
		logger.NewField("mimir_url", c.Mimir.URL),
		logger.NewField("in_cluster", c.K8s.InCluster),
	)
}

// Funções auxiliares
func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsIntOrDefault(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvAsDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func getKubeconfigPath() string {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			kubeconfigPath = homeDir + "/.kube/config"
		}
	}
	return kubeconfigPath
}
