package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Backup das variáveis de ambiente originais
	originalEnv := map[string]string{
		"PORT":               os.Getenv("PORT"),
		"GIN_MODE":           os.Getenv("GIN_MODE"),
		"LOG_LEVEL":          os.Getenv("LOG_LEVEL"),
		"LOG_FORMAT":         os.Getenv("LOG_FORMAT"),
		"MIMIR_URL":          os.Getenv("MIMIR_URL"),
		"MIMIR_NAMESPACE":    os.Getenv("MIMIR_NAMESPACE"),
		"MIMIR_SERVICE_NAME": os.Getenv("MIMIR_SERVICE_NAME"),
		"MIMIR_ORG_ID":       os.Getenv("MIMIR_ORG_ID"),
		"IN_CLUSTER":         os.Getenv("IN_CLUSTER"),
		"KUBECONFIG":         os.Getenv("KUBECONFIG"),
	}

	// Restaura as variáveis de ambiente originais ao final
	defer func() {
		for k, v := range originalEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Configura variáveis de ambiente para teste
	testEnv := map[string]string{
		"PORT":               "8080",
		"GIN_MODE":           "release",
		"LOG_LEVEL":          "debug",
		"LOG_FORMAT":         "json",
		"MIMIR_URL":          "http://mimir:9090",
		"MIMIR_NAMESPACE":    "monitoring",
		"MIMIR_SERVICE_NAME": "mimir",
		"MIMIR_ORG_ID":       "test",
		"IN_CLUSTER":         "true",
	}

	for k, v := range testEnv {
		os.Setenv(k, v)
	}

	// Testa carregamento da configuração
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() retornou erro: %v", err)
	}

	// Verifica valores carregados
	if cfg.Server.Port != testEnv["PORT"] {
		t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, testEnv["PORT"])
	}
	if cfg.Server.GinMode != testEnv["GIN_MODE"] {
		t.Errorf("Server.GinMode = %v, want %v", cfg.Server.GinMode, testEnv["GIN_MODE"])
	}
	if cfg.Logging.Level != testEnv["LOG_LEVEL"] {
		t.Errorf("Logging.Level = %v, want %v", cfg.Logging.Level, testEnv["LOG_LEVEL"])
	}
	if cfg.Logging.Format != testEnv["LOG_FORMAT"] {
		t.Errorf("Logging.Format = %v, want %v", cfg.Logging.Format, testEnv["LOG_FORMAT"])
	}
	if cfg.Mimir.URL != testEnv["MIMIR_URL"] {
		t.Errorf("Mimir.URL = %v, want %v", cfg.Mimir.URL, testEnv["MIMIR_URL"])
	}
	if cfg.K8s.InCluster != true {
		t.Error("K8s.InCluster = false, want true")
	}
}

func TestConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "configuração válida",
			config: &Config{
				Server: ServerConfig{
					Port: "8080",
				},
				Mimir: MimirConfig{
					URL: "http://mimir:9090",
				},
			},
			wantErr: false,
		},
		{
			name: "porta vazia",
			config: &Config{
				Server: ServerConfig{
					Port: "",
				},
				Mimir: MimirConfig{
					URL: "http://mimir:9090",
				},
			},
			wantErr: true,
		},
		{
			name: "URL do Mimir vazia",
			config: &Config{
				Server: ServerConfig{
					Port: "8080",
				},
				Mimir: MimirConfig{
					URL: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Backup e restauração da variável de ambiente
	key := "TEST_ENV_VAR"
	originalValue := os.Getenv(key)
	defer os.Setenv(key, originalValue)

	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "variável definida",
			key:      key,
			value:    "test_value",
			fallback: "default_value",
			want:     "test_value",
		},
		{
			name:     "variável não definida",
			key:      key,
			value:    "",
			fallback: "default_value",
			want:     "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}

			if got := getEnvOrDefault(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsIntOrDefault(t *testing.T) {
	key := "TEST_INT_VAR"
	originalValue := os.Getenv(key)
	defer os.Setenv(key, originalValue)

	tests := []struct {
		name     string
		key      string
		value    string
		fallback int
		want     int
	}{
		{
			name:     "valor inteiro válido",
			key:      key,
			value:    "42",
			fallback: 0,
			want:     42,
		},
		{
			name:     "valor inválido",
			key:      key,
			value:    "not_a_number",
			fallback: 10,
			want:     10,
		},
		{
			name:     "variável não definida",
			key:      key,
			value:    "",
			fallback: 5,
			want:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}

			if got := getEnvAsIntOrDefault(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnvAsIntOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsDurationOrDefault(t *testing.T) {
	key := "TEST_DURATION_VAR"
	originalValue := os.Getenv(key)
	defer os.Setenv(key, originalValue)

	tests := []struct {
		name     string
		key      string
		value    string
		fallback time.Duration
		want     time.Duration
	}{
		{
			name:     "duração válida",
			key:      key,
			value:    "5s",
			fallback: time.Second,
			want:     5 * time.Second,
		},
		{
			name:     "duração inválida",
			key:      key,
			value:    "invalid",
			fallback: 2 * time.Second,
			want:     2 * time.Second,
		},
		{
			name:     "variável não definida",
			key:      key,
			value:    "",
			fallback: 3 * time.Second,
			want:     3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}

			if got := getEnvAsDurationOrDefault(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnvAsDurationOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKubeconfigPath(t *testing.T) {
	// Backup das variáveis de ambiente
	originalKubeconfig := os.Getenv("KUBECONFIG")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("KUBECONFIG", originalKubeconfig)
		os.Setenv("HOME", originalHome)
	}()

	tests := []struct {
		name         string
		kubeconfig   string
		home         string
		wantContains string
		wantNotEmpty bool
	}{
		{
			name:         "KUBECONFIG definido",
			kubeconfig:   "/custom/path/config",
			home:         "/home/user",
			wantContains: "/custom/path/config",
			wantNotEmpty: true,
		},
		{
			name:         "KUBECONFIG não definido, usa HOME",
			kubeconfig:   "",
			home:         "/home/user",
			wantContains: "/.kube/config",
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.kubeconfig != "" {
				os.Setenv("KUBECONFIG", tt.kubeconfig)
			} else {
				os.Unsetenv("KUBECONFIG")
			}
			if tt.home != "" {
				os.Setenv("HOME", tt.home)
			}

			got := getKubeconfigPath()

			if tt.wantNotEmpty && got == "" {
				t.Error("getKubeconfigPath() retornou string vazia")
			}
			if tt.wantContains != "" && got != "" && !contains(got, tt.wantContains) {
				t.Errorf("getKubeconfigPath() = %v, deve conter %v", got, tt.wantContains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return s == substr || len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
