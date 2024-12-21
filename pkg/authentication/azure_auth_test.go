package authentication

import (
	"testing"
)

func TestAzureAuthIntegration(t *testing.T) {
	auth := &AzureAuth{
		ClientID:       "your-client-id",
		ClientSecret:   "your-client-secret",
		TenantID:       "your-tenant-id",
		SubscriptionID: "your-subscription-id",
	}

	if err := TestAzureAuth(auth); err != nil {
		t.Fatalf("Azure authentication test failed: %v", err)
	}
}
