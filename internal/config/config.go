// Package config defines the gateway's configuration loading and parsing
// logic, including environment, file, and default values.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var Envs = loadConfig()

type Config struct {
	AppPort            string
	AuthSecret         string
	DBConnectionString string
	JWTExpiration      time.Duration
	IsProd             bool
}

// Loads all environment variables from the .env file
func loadConfig() Config {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	return Config{
		AppPort:            getEnv("PORT", ":8000"),
		JWTExpiration:      getEnvAsDuration("JWT_EXPIRATION_IN_SECONDS", time.Hour),
		IsProd:             getEnvAsBoolean("IS_PROD", true),
		DBConnectionString: getEnv("POSTGRES_DNS", ""),
		AuthSecret:         getEnv("AUTH_SECRET", "superSecretNeedsToBeChanged"),
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
	return time.Duration(fallback) * time.Second
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
