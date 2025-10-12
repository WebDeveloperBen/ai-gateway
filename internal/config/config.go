// Package config defines the gateway's configuration loading and parsing
// logic, including environment, file, and default values.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/joho/godotenv"
)

var Envs = loadConfig()

func Reload() { Envs = loadConfig() }

type Config struct {
	ProxyPort                   string
	AdminPort                   string
	AuthSecret                  string
	DBConnectionString          string
	DBManagedIdentityConfig     string
	JWTExpiration               time.Duration
	IsProd                      bool
	ApplicationName             string
	Version                     string
	AzureOpenAiAPIKey           string
	KVBackend                   string
	DBBackend                   model.RepositoryBackend
	RedisAddr                   string
	RedisPW                     string
	AppRegistrationClientID     string
	AppRegistrationClientSecret string
	AppRegistrationTenantID     string
	AppRegistrationRedirectURL  string
	EnableRedisCircuitBreaker   bool
}

// Loads all environment variables from the .env file
func loadConfig() Config {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	return Config{
		ProxyPort:                   getEnv("PROXY_PORT", ":8000"),
		AdminPort:                   getEnv("ADMIN_PORT", ":8080"),
		JWTExpiration:               getEnvAsDuration("JWT_EXPIRATION_IN_SECONDS", time.Hour),
		IsProd:                      getEnvAsBoolean("IS_PROD", false),
		DBConnectionString:          getEnv("POSTGRES_DSN", ""),
		DBManagedIdentityConfig:     getEnv("AZURE_DB_CONFIG", ""),
		AuthSecret:                  getEnv("AUTH_SECRET", "superSecretNeedsToBeChanged"),
		ApplicationName:             getEnv("APPLICATION_NAME", "LLM Gateway"),
		Version:                     getEnv("VERSION", "v1.0.0"),
		AzureOpenAiAPIKey:           getEnv("AZURE_OPENAI_API_KEY", ""),
		KVBackend:                   getEnv("KV_BACKEND", "memory"),
		RedisAddr:                   getEnv("REDIS_ADDR", ""),
		RedisPW:                     getEnv("REDIS_PASSWORD", ""),
		DBBackend:                   model.RepositoryBackend(getEnv("DB_BACKEND", "postgres")),
		AppRegistrationClientID:     getEnv("AZURE_APP_REGISTRATION_CLIENT_ID", "dummy-client-id"),
		AppRegistrationClientSecret: getEnv("AZURE_APP_REGISTRATION_CLIENT_SECRET", "dummy-client-secret"),
		AppRegistrationTenantID:     getEnv("AZURE_APP_REGISTRATION_TENANT_ID", "dummy-tenant-id"),
		AppRegistrationRedirectURL:  getEnv("AZURE_APP_REGISTRATION_REDIRECT_URL", "http://localhost:3000/auth/callback"),
		EnableRedisCircuitBreaker:   getEnvAsBoolean("REDIS_CIRCUIT_BREAKER_ENABLED", true),
	}
}

// Gets an environment variable and provides a default value if not set
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvAsInt64 Gets an environment variable as an int64 and provides a default value if not set
func GetEnvAsInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return fallback
}

// Modified function to return time.Duration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(intValue) * time.Second
		}
	}
	return fallback
}

// Gets an environment variable as a boolean from a string, provides the fallback value if not set
func getEnvAsBoolean(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// GetDatabaseConnection returns the appropriate database connection string
// Prioritizes managed identity config over connection string for Azure deployments
func (c *Config) GetDatabaseConnection() string {
	if c.DBManagedIdentityConfig != "" {
		return c.DBManagedIdentityConfig
	}
	return c.DBConnectionString
}
