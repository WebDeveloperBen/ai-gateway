package config

import (
	"os"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	t.Run("Returns value when env var exists", func(t *testing.T) {
		key := "TEST_VAR_EXISTS"
		expected := "test-value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		result := getEnv(key, "fallback")
		assert.Equal(t, expected, result)
	})

	t.Run("Returns fallback when env var does not exist", func(t *testing.T) {
		result := getEnv("TEST_VAR_NOT_EXISTS", "fallback")
		assert.Equal(t, "fallback", result)
	})

	t.Run("Returns empty string when env var exists but is empty", func(t *testing.T) {
		key := "TEST_VAR_EMPTY"
		os.Setenv(key, "")
		defer os.Unsetenv(key)

		result := getEnv(key, "fallback")
		assert.Equal(t, "", result)
	})
}

func TestGetEnvAsInt64(t *testing.T) {
	t.Run("Returns value when env var is valid int64", func(t *testing.T) {
		key := "TEST_INT_VAR"
		os.Setenv(key, "12345")
		defer os.Unsetenv(key)

		result := GetEnvAsInt64(key, 999)
		assert.Equal(t, int64(12345), result)
	})

	t.Run("Returns fallback when env var does not exist", func(t *testing.T) {
		result := GetEnvAsInt64("TEST_INT_NOT_EXISTS", 999)
		assert.Equal(t, int64(999), result)
	})

	t.Run("Returns fallback when env var is not a valid int64", func(t *testing.T) {
		key := "TEST_INT_INVALID"
		os.Setenv(key, "not-a-number")
		defer os.Unsetenv(key)

		result := GetEnvAsInt64(key, 999)
		assert.Equal(t, int64(999), result)
	})

	t.Run("Handles negative numbers", func(t *testing.T) {
		key := "TEST_INT_NEGATIVE"
		os.Setenv(key, "-500")
		defer os.Unsetenv(key)

		result := GetEnvAsInt64(key, 0)
		assert.Equal(t, int64(-500), result)
	})
}

func TestGetEnvAsDuration(t *testing.T) {
	t.Run("Returns duration when env var is valid int", func(t *testing.T) {
		key := "TEST_DURATION_VAR"
		os.Setenv(key, "300")
		defer os.Unsetenv(key)

		result := getEnvAsDuration(key, time.Hour)
		assert.Equal(t, 300*time.Second, result)
	})

	t.Run("Returns fallback when env var does not exist", func(t *testing.T) {
		result := getEnvAsDuration("TEST_DURATION_NOT_EXISTS", time.Minute)
		assert.Equal(t, time.Minute, result)
	})

	t.Run("Returns fallback when env var is not a valid int", func(t *testing.T) {
		key := "TEST_DURATION_INVALID"
		os.Setenv(key, "not-a-number")
		defer os.Unsetenv(key)

		result := getEnvAsDuration(key, time.Minute)
		assert.Equal(t, time.Minute, result)
	})

	t.Run("Handles zero duration", func(t *testing.T) {
		key := "TEST_DURATION_ZERO"
		os.Setenv(key, "0")
		defer os.Unsetenv(key)

		result := getEnvAsDuration(key, time.Hour)
		assert.Equal(t, time.Duration(0), result)
	})
}

func TestGetEnvAsBoolean(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fallback bool
		expected bool
	}{
		{"true string", "true", false, true},
		{"false string", "false", true, false},
		{"1 is true", "1", false, true},
		{"0 is false", "0", true, false},
		{"TRUE is true", "TRUE", false, true},
		{"FALSE is false", "FALSE", true, false},
		{"t is true", "t", false, true},
		{"f is false", "f", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_BOOL_VAR"
			os.Setenv(key, tt.value)
			defer os.Unsetenv(key)

			result := getEnvAsBoolean(key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("Returns fallback when env var does not exist", func(t *testing.T) {
		result := getEnvAsBoolean("TEST_BOOL_NOT_EXISTS", true)
		assert.Equal(t, true, result)
	})

	t.Run("Returns fallback when env var is not a valid boolean", func(t *testing.T) {
		key := "TEST_BOOL_INVALID"
		os.Setenv(key, "maybe")
		defer os.Unsetenv(key)

		result := getEnvAsBoolean(key, true)
		assert.Equal(t, true, result)
	})
}

func TestConfigGetDatabaseConnection(t *testing.T) {
	t.Run("Returns managed identity config when set", func(t *testing.T) {
		cfg := &Config{
			DBManagedIdentityConfig: "managed-identity-connection",
			DBConnectionString:      "regular-connection",
		}

		result := cfg.GetDatabaseConnection()
		assert.Equal(t, "managed-identity-connection", result)
	})

	t.Run("Returns connection string when managed identity is empty", func(t *testing.T) {
		cfg := &Config{
			DBManagedIdentityConfig: "",
			DBConnectionString:      "regular-connection",
		}

		result := cfg.GetDatabaseConnection()
		assert.Equal(t, "regular-connection", result)
	})

	t.Run("Prioritizes managed identity over connection string", func(t *testing.T) {
		cfg := &Config{
			DBManagedIdentityConfig: "managed",
			DBConnectionString:      "connection",
		}

		result := cfg.GetDatabaseConnection()
		assert.Equal(t, "managed", result)
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("Loads config with default values", func(t *testing.T) {
		cfg := loadConfig()

		assert.NotNil(t, cfg)
		assert.NotEmpty(t, cfg.ProxyPort)
		assert.NotEmpty(t, cfg.AdminPort)
		assert.NotEmpty(t, cfg.AuthSecret)
		assert.NotEmpty(t, cfg.ApplicationName)
		assert.NotEmpty(t, cfg.Version)
	})

	t.Run("Config has expected default values", func(t *testing.T) {
		cfg := loadConfig()

		assert.Equal(t, ":8000", cfg.ProxyPort)
		assert.Equal(t, ":8080", cfg.AdminPort)
		assert.Equal(t, "LLM Gateway", cfg.ApplicationName)
		assert.Equal(t, "v1.0.0", cfg.Version)
		assert.Equal(t, "memory", cfg.KVBackend)
		assert.Equal(t, model.RepositoryBackend("postgres"), cfg.DBBackend)
		assert.True(t, cfg.EnableRedisCircuitBreaker)
	})

	t.Run("Respects environment variables", func(t *testing.T) {
		os.Setenv("PROXY_PORT", ":9000")
		os.Setenv("ADMIN_PORT", ":9090")
		os.Setenv("APPLICATION_NAME", "Test Gateway")
		os.Setenv("IS_PROD", "true")
		defer func() {
			os.Unsetenv("PROXY_PORT")
			os.Unsetenv("ADMIN_PORT")
			os.Unsetenv("APPLICATION_NAME")
			os.Unsetenv("IS_PROD")
		}()

		cfg := loadConfig()

		assert.Equal(t, ":9000", cfg.ProxyPort)
		assert.Equal(t, ":9090", cfg.AdminPort)
		assert.Equal(t, "Test Gateway", cfg.ApplicationName)
		assert.True(t, cfg.IsProd)
	})
}

func TestReload(t *testing.T) {
	t.Run("Reload updates global Envs variable", func(t *testing.T) {
		originalPort := Envs.ProxyPort

		os.Setenv("PROXY_PORT", ":7777")
		defer os.Unsetenv("PROXY_PORT")

		Reload()

		assert.NotEqual(t, originalPort, Envs.ProxyPort)
		assert.Equal(t, ":7777", Envs.ProxyPort)

		Reload()
	})
}

func TestConfigStructure(t *testing.T) {
	t.Run("Config has all expected fields", func(t *testing.T) {
		cfg := Config{}

		cfg.ProxyPort = ":8000"
		cfg.AdminPort = ":8080"
		cfg.AuthSecret = "secret"
		cfg.DBConnectionString = "connection"
		cfg.DBManagedIdentityConfig = "managed"
		cfg.JWTExpiration = time.Hour
		cfg.IsProd = true
		cfg.ApplicationName = "App"
		cfg.Version = "v1.0"
		cfg.AzureOpenAiAPIKey = "key"
		cfg.KVBackend = "redis"
		cfg.DBBackend = model.RepositoryBackend("postgres")
		cfg.RedisAddr = "localhost:6379"
		cfg.RedisPW = "password"
		cfg.AppRegistrationClientID = "client"
		cfg.AppRegistrationClientSecret = "secret"
		cfg.AppRegistrationTenantID = "tenant"
		cfg.AppRegistrationRedirectURL = "url"
		cfg.EnableRedisCircuitBreaker = true

		require.NotNil(t, cfg)
		assert.Equal(t, ":8000", cfg.ProxyPort)
	})
}
