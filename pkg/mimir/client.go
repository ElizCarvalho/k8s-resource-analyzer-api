package mimir

import (
	"os"
	"strconv"
)

type Config struct {
	Host  string
	Port  int
	OrgID string
}

func NewConfig() *Config {
	return &Config{
		Host:  getEnvOrDefault("MIMIR_HOST", "localhost"),
		Port:  getEnvAsIntOrDefault("MIMIR_PORT", 8080),
		OrgID: getEnvOrDefault("MIMIR_ORG_ID", "anonymous"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
