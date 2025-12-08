package config

import (
	"fmt"
	"os"
	"strconv"
)

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the defaultValue.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// mustGetEnv retrieves the value of the environment variable named by the key.
// If the variable is not present or empty, it panics with an error message.
// Use this for critical configuration that must be set in production.
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// getEnvAsInt retrieves the value of the environment variable named by the key
// and converts it to an integer. If the variable is not present or cannot be
// converted to an integer, it returns the defaultValue.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(valueStr); err == nil {
		return intValue
	}
	return defaultValue
}
