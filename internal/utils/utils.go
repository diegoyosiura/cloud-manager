package utils

import (
	"log"
	"os"
	"strings"
)

// GetEnvWithValidation retrieves an environment variable and ensures it is not empty.
func GetEnvWithValidation(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

// GetOptionalEnv retrieves an environment variable or returns the provided default value if the variable is not set.
func GetOptionalEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func IsValidArn(arn string) bool {
	return strings.HasPrefix(arn, "arn:aws:")
}
