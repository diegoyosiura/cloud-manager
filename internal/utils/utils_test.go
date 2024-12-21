package utils

import (
	"os"
	"testing"
)

// TestGetEnvWithValidation tests the GetEnvWithValidation function to ensure that it correctly retrieves
// the value of a required environment variable. If the variable is missing or empty, it is expected to
// fail the execution.
func TestGetEnvWithValidation(t *testing.T) {
	key := "REQUIRED_ENV"    // Name of the required environment variable
	expectedValue := "value" // Expected value to test against

	// Set a temporary environment variable for testing
	_ = os.Setenv(key, expectedValue)
	defer func() { _ = os.Unsetenv(key) }() // Ensure the variable is unset after the test

	// Call the GetEnvWithValidation function
	value := GetEnvWithValidation(key)
	if value != expectedValue {
		t.Fatalf("Expected %s, got %s", expectedValue, value) // Fail the test if values do not match
	}
}

// TestGetOptionalEnv tests the GetOptionalEnv function to ensure that it correctly retrieves
// the value of an optional environment variable. If the variable is not set, the default value
// provided should be returned instead.
func TestGetOptionalEnv(t *testing.T) {
	key := "OPTIONAL_ENV"     // Name of the optional environment variable
	defaultValue := "default" // Default value to return if the variable is not set

	// Set a temporary environment variable for testing
	_ = os.Setenv(key, "set-value")
	defer func() { _ = os.Unsetenv(key) }() // Ensure the variable is unset after the test

	// Call the GetOptionalEnv function when the variable is set
	value := GetOptionalEnv(key, defaultValue)
	if value != "set-value" {
		t.Fatalf("Expected %s, got %s", "set-value", value) // Fail the test if values do not match
	}

	// Unset the environment variable to test the default behavior
	_ = os.Unsetenv(key)
	value = GetOptionalEnv(key, defaultValue)
	if value != defaultValue {
		t.Fatalf("Expected %s, got %s", defaultValue, value) // Fail the test if the default value is incorrect
	}
}
