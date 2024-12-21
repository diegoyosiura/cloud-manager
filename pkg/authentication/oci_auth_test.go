package authentication

import (
	"errors"
	"os"
	"testing"
)

// Test cases for configuration validation
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		fields  *OCIAuth
		wantErr error
	}{
		{
			name: "Missing TenancyID",
			fields: &OCIAuth{
				UserID:      "ocid1.user.oc1...",
				Region:      "us-ashburn-1",
				PrivateKey:  "some-private-key",
				Fingerprint: "some-fingerprint",
			},
			wantErr: errors.New("tenancy ID is required"),
		},
		{
			name: "All Fields Valid",
			fields: &OCIAuth{
				TenancyID:   "ocid1.tenancy.oc1...",
				UserID:      "ocid1.user.oc1...",
				Region:      "us-ashburn-1",
				PrivateKey:  "some-private-key",
				Fingerprint: "some-fingerprint",
			},
			wantErr: nil,
		},
		{
			name: "Missing PrivateKey",
			fields: &OCIAuth{
				TenancyID:   "ocid1.tenancy.oc1...",
				UserID:      "ocid1.user.oc1...",
				Region:      "us-ashburn-1",
				Fingerprint: "some-fingerprint",
			},
			wantErr: errors.New("private key is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.Validate()
			if err == nil && tt.wantErr == nil {
				return // Test passed
			}
			if err == nil || (err != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestNewOCIAuth validates the creation of an OCIAuth configuration and checks that it initializes without errors.
func TestNewOCIAuth(t *testing.T) {
	fields := map[string]string{
		"oci_tenancy_id":     "ocid1.tenancy.oc1...",
		"oci_user_id":        "ocid1.user.oc1...",
		"oci_region":         "us-ashburn-1",
		"oci_private_key":    "some-private-key",
		"oci_fingerprint":    "some-fingerprint",
		"oci_key_passphrase": "",
	}

	auth, err := NewOCIAuthFromAuth(fields)
	if err != nil {
		t.Fatalf("failed to create OCIAuth: %v", err)
	}

	if err := auth.Validate(); err != nil {
		t.Errorf("Validate() unexpectedly failed: %v", err)
	}
}

// TestTestOCIAuth checks if the provided OCIAuth configuration is valid by running a test authentication.
func TestTestOCIAuth(t *testing.T) {
	auth := &OCIAuth{
		TenancyID:     os.Getenv("ORACLE_API_TENANCY"),
		UserID:        os.Getenv("ORACLE_API_USER"),
		Region:        os.Getenv("ORACLE_API_REGION"),
		PrivateKey:    os.Getenv("ORACLE_API_PRIVATE_KEY"),
		Fingerprint:   os.Getenv("ORACLE_API_FINGERPRINT"),
		KeyPassphrase: os.Getenv("ORACLE_API_KEY_PASSPHRASE"),
	}
	
	// Mock the TestOCIAuth behavior if needed (e.g., mock OCI SDK calls)
	if err := TestOCIAuth(auth); err != nil {
		t.Fatalf("failed to test OCI authentication: %v", err)
	}
}