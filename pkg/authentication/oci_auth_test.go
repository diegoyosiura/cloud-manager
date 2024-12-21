package authentication

import (
	"cloud-manager/internal/utils"
	"os"
	"testing"
)

var testConfig OCIAuth

func TestMain(m *testing.M) {
	// Initialize shared test configuration
	testConfig = OCIAuth{
		TenancyID:     utils.GetEnvWithValidation("ORACLE_API_TENANCY"),
		UserID:        utils.GetEnvWithValidation("ORACLE_API_USER"),
		Region:        utils.GetEnvWithValidation("ORACLE_API_REGION"),
		PrivateKey:    utils.GetEnvWithValidation("ORACLE_API_PRIVATE_KEY"),
		Fingerprint:   utils.GetEnvWithValidation("ORACLE_API_FINGERPRINT"),
		KeyPassphrase: utils.GetOptionalEnv("ORACLE_API_KEY_PASSPHRASE", ""),
	}

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

// TestNewOCIAuth validates the creation of an OCIAuth configuration and checks that it initializes without errors.
func TestNewOCIAuth(t *testing.T) {
	// Attempt to create a new OCIAuth instance
	_, err := NewOCIAuth(testConfig)
	if err != nil {
		t.Fatalf("failed to create OCI client: %v", err) // Fail the test if there's an error
	}
}

// TestTestOCIAuth checks if the provided OCIAuth configuration is valid by running a test authentication.
func TestTestOCIAuth(t *testing.T) {
	// Attempt to test the OCIAuth configuration
	err := TestOCIAuth(testConfig)
	if err != nil {
		t.Fatalf("failed to test OCI authentication: %v", err) // Fail the test if there's an error
	}
}
