package authentication

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"sync"
)

// AzureAuth represents the configuration and state for authenticating
// with Microsoft Azure using the Azure SDK for Go.
type AzureAuth struct {
	ClientID       string // Azure Client ID (Application ID) used for authentication.
	ClientSecret   string // Azure Client Secret used for authentication.
	TenantID       string // Azure Tenant ID that the application belongs to.
	SubscriptionID string // Azure Subscription ID to operate within.
	EmailHost      string // SMTP Host
	EmailPort      string // SMTP Port
	EmailUser      string // SMTP User
	EmailPassword  string // SMTP PWD

	Authenticated bool                               // Tracks whether authentication was performed successfully.
	Credential    *azidentity.ClientSecretCredential // Credential object used for authorization with Azure.
	Client        *armresources.Client               // Azure Resource Manager client for interacting with Azure resources.

	mu sync.Mutex
}

// NewAzureAuthFromAuth initializes a new AzureAuth object using a map of fields.
// The function populates the struct with values taken from the fields map and validates it.
func NewAzureAuthFromAuth(fields map[string]string) (*AzureAuth, error) {
	config := &AzureAuth{
		mu:             sync.Mutex{},
		Authenticated:  false,                           // Start with unauthenticated state.
		ClientID:       fields["azure_client_id"],       // Extract Azure Client ID from fields.
		ClientSecret:   fields["azure_client_secret"],   // Extract Azure Client Secret from fields.
		TenantID:       fields["azure_tenant_id"],       // Extract Azure Tenant ID from fields.
		SubscriptionID: fields["azure_subscription_id"], // Extract Azure Subscription ID from fields.
		EmailHost:      fields["email_host"],            // SMTP User
		EmailPort:      fields["email_port"],            // SMTP User
		EmailUser:      fields["email_user"],            // SMTP User
		EmailPassword:  fields["email_password"],        // SMTP PWD
	}
	// Return the initialized AzureAuth structure and validate the configuration.
	return config, config.Validate()
}

// Validate checks if all required Azure authentication fields in the struct are populated.
// It returns an error if any mandatory fields are missing.
func (a *AzureAuth) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Check for empty mandatory fields: ClientID, ClientSecret, TenantID, and SubscriptionID.
	if a.ClientID == "" || a.ClientSecret == "" || a.TenantID == "" || a.SubscriptionID == "" {
		return fmt.Errorf("missing required Azure authentication fields")
	}
	// Return nil if all fields are valid (no missing fields).
	return nil
}

// Authenticate performs Azure authentication using the ClientSecretCredential.
// If authentication is successful, it also initializes a resource manager client for further operations.
func (a *AzureAuth) Authenticate() error {
	a.mu.Lock()
	// Avoid reauthentication if already authenticated.
	if a.Authenticated {
		return nil
	}

	var err error

	a.mu.Unlock()
	// Validate the AzureAuth struct to ensure all required fields are set.
	if err := a.Validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create an Azure client credential object for authentication using ClientID, ClientSecret, and TenantID.
	a.Credential, err = azidentity.NewClientSecretCredential(a.TenantID, a.ClientID, a.ClientSecret, nil)
	if err != nil {
		// Return an error if the credential creation fails, providing more context.
		return fmt.Errorf("failed to create Azure credentials: %v. Check TenantID, ClientID, ClientSecret", err)
	}

	// Initialize a new Resource Manager client for the specified subscription using the created credentials.
	a.Client, err = armresources.NewClient(a.SubscriptionID, a.Credential, nil)
	if err != nil {
		// Return an error if the resource manager client cannot be created, with possible reasons.
		return fmt.Errorf("failed to create Azure client: %v. Check SubscriptionID or permissions", err)
	}

	// Set authentication state to true to indicate successful authentication.
	a.Authenticated = true
	return nil
}

// TestAzureAuth tests the AzureAuth configuration by validating the input and attempting authentication.
// It ensures both validation and authentication complete without errors.
func TestAzureAuth(auth *AzureAuth) error {
	// Step 1: Validate the configuration to check for missing or invalid fields.
	if err := auth.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err) // Return error if validation fails.
	}

	// Step 2: Attempt to authenticate using the provided credentials and configuration.
	if err := auth.Authenticate(); err != nil {
		return fmt.Errorf("authentication test failed: %w", err) // Return error if authentication fails.
	}

	// Both validation and authentication have succeeded.
	return nil
}
