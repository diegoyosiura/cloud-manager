package authentication

import (
	"cloud-manager/internal/utils"
	"testing"
)

func TestGCPAuthIntegration(t *testing.T) {
	auth := &GCPAuth{
		ProjectID: utils.GetEnvWithValidation("GCP_KEY_ID"),
		AuthJSON:  utils.GetEnvWithValidation("GCP_JSON_INFO"),
	}

	if err := TestGCPAuth(auth); err != nil {
		t.Fatalf("GCP authentication test failed: %v", err)
	}
}
