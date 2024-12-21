package authentication

import (
	"testing"
)

// TestNewAuthConfig_ValidAWS verifies that a valid AWS configuration is correctly initialized.
func TestNewAuthConfig_ValidAWS(t *testing.T) {
	fields := map[string]string{
		"aws_access_key_id":     "testAccessKey",
		"aws_secret_access_key": "testSecretKey",
		"aws_region":            "us-east-1",
	}

	config, err := NewAuthConfig("aws", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.ProviderName != "aws" {
		t.Errorf("expected provider to be 'aws', got '%s'", config.ProviderName)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("unexpected validation error for valid AWS config: %v", err)
	}
}

// TestNewAuthConfig_InvalidProvider ensures that an unsupported provider returns an appropriate error.
func TestNewAuthConfig_InvalidProvider(t *testing.T) {
	fields := map[string]string{
		"irrelevant_field": "value",
	}

	_, err := NewAuthConfig("unsupported", fields)
	if err == nil {
		t.Fatalf("expected error for unsupported provider, got nil")
	}

	expectedErr := "unsupported provider: unsupported"
	if err.Error() != expectedErr {
		t.Fatalf("expected error message: %s, got: %v", expectedErr, err)
	}
}

// TestAuthConfig_ValidateAWS_MissingFields ensures that missing required fields in AWS configuration produce an error.
func TestAuthConfig_ValidateAWS_MissingFields(t *testing.T) {
	fields := map[string]string{
		"aws_secret_access_key": "testSecretKey",
		"aws_region":            "us-east-1",
	}

	config, err := NewAuthConfig("aws", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = config.Validate()
	if err == nil {
		t.Fatalf("expected error for missing AWSAccessKeyID, got nil")
	}

	expectedErr := "missing AWSAccessKeyID"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: expected %s, got: %v", expectedErr, err)
	}
}

// TestNewAuthConfig_ValidateGCP_Valid checks if a valid GCP configuration passes validation successfully.
func TestNewAuthConfig_ValidateGCP_Valid(t *testing.T) {
	fields := map[string]string{
		"gcp_project_id": "testGCPProject",
		"gcp_auth_json":  "{}",
	}

	config, err := NewAuthConfig("gcp", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("unexpected validation error for valid GCP config: %v", err)
	}
}

// TestAuthConfig_ValidateGCP_MissingFields ensures appropriate error is returned for missing required GCP fields.
func TestAuthConfig_ValidateGCP_MissingFields(t *testing.T) {
	fields := map[string]string{
		"gcp_auth_json": "{}",
	}

	config, err := NewAuthConfig("gcp", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = config.Validate()
	if err == nil {
		t.Fatalf("expected error for missing GCPProjectID, got nil")
	}

	expectedErr := "missing GCPProjectID"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: expected %s, got: %v", expectedErr, err)
	}
}

// TestAuthConfig_ValidateOCI_Valid verifies that a valid OCI configuration is correctly validated.
func TestAuthConfig_ValidateOCI_Valid(t *testing.T) {
	fields := map[string]string{
		"oci_tenancy_id":     "testTenancyID",
		"oci_user_id":        "testUserID",
		"oci_region":         "us-ashburn-1",
		"oci_private_key":    "testPrivateKey",
		"oci_fingerprint":    "testFingerprint",
		"oci_key_passphrase": "testPassphrase",
	}

	config, err := NewAuthConfig("oci", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("unexpected validation error for valid OCI config: %v", err)
	}
}

// TestAuthConfig_ValidateOCI_MissingFields checks if appropriate errors are triggered for missing OCI fields.
func TestAuthConfig_ValidateOCI_MissingFields(t *testing.T) {
	fields := map[string]string{
		"oci_user_id":        "testUserID",
		"oci_region":         "us-ashburn-1",
		"oci_private_key":    "testPrivateKey",
		"oci_fingerprint":    "testFingerprint",
		"oci_key_passphrase": "testPassphrase",
	}

	config, err := NewAuthConfig("oci", fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = config.Validate()
	if err == nil {
		t.Fatalf("expected error for missing OCITenancyID, got nil")
	}

	expectedErr := "missing OCITenancyID"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: expected %s, got: %v", expectedErr, err)
	}
}

// TestAuthConfig_Authenticate_VerifyIntegration tests the integration of the Authenticate method for supported providers.
// Simulates success/failure based on field completion.
func TestAuthConfig_Authenticate_VerifyIntegration(t *testing.T) {
	fields := map[string]string{
		"aws_access_key_id":     "testAccessKey",
		"aws_secret_access_key": "testSecretKey",
		"aws_region":            "us-east-1",
	}

	config, err := NewAuthConfig("aws", fields)
	if err != nil {
		t.Fatalf("unexpected error during initialization: %v", err)
	}

	err = config.Authenticate()
	if err != nil {
		t.Errorf("unexpected authentication error for valid AWS config: %v", err)
	}
}
