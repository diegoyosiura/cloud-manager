package authentication

// Imports required packages:
// - context: for managing execution contexts.
// - fmt: for formatting and creating error or log messages.
// - github.com/oracle/oci-go-sdk/v65: Oracle Cloud Infrastructure (OCI) SDK library for interacting with OCI resources.
// - log: for logging messages to the console.
// - strings: for string operations, such as replacing patterns or characters.
// - sync: for thread-safe operations using a sync.Mutex.
import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"strings"
	"sync"
)

// OCIAuth is a struct that encapsulates the configuration and state required
// to authenticate with Oracle Cloud Infrastructure (OCI) services.
type OCIAuth struct {
	Namespace     string // The Namespace of the account.
	CompartmentID string // The Compartment ID of the account (mandatory).
	TenancyID     string // The tenancy ID of the account (mandatory).
	UserID        string // The user ID in the tenancy (mandatory).
	Region        string // The OCI region where services will be used (mandatory).
	PrivateKey    string // The private key for authentication (mandatory).
	Fingerprint   string // Fingerprint of the private key (mandatory).
	KeyPassphrase string // The passphrase for the private key (optional if the private key doesn't require it).
	SMTPSecret    string // The passphrase for SMTP Authentication.
	EmailHost     string // SMTP Host
	EmailPort     string // SMTP Port
	EmailUser     string // SMTP User
	EmailPassword string // SMTP PWD

	Authenticated bool                    // Tracks whether the user is successfully authenticated.
	Client        identity.IdentityClient // The client used to interact with the OCI identity service.

	privateKeyProvider common.ConfigurationProvider

	mu sync.Mutex // A mutex used to ensure thread safety when accessing the struct.
}

// NewOCIAuthFromAuth creates a new instance of OCIAuth based on the provided fields.
//
// Parameters:
// - fields: A map[string]string containing the required authentication fields (e.g., tenancy ID, user ID, etc.).
//
// Returns:
// - A pointer to the initialized OCIAuth instance.
// - An error if the configuration is invalid based on the Validate method.
func NewOCIAuthFromAuth(fields map[string]string) (*OCIAuth, error) {
	config := &OCIAuth{
		mu:            sync.Mutex{},                 // Initializes the mutex for thread safety.
		Authenticated: false,                        // Authentication is set to "false" by default.
		Namespace:     fields["oci_namespace"],      // Reads the namespace from the input fields.
		CompartmentID: fields["oci_compartment_id"], // Reads the compartment ID from the input fields.
		TenancyID:     fields["oci_tenancy_id"],     // Reads the tenancy ID from the input fields.
		UserID:        fields["oci_user_id"],        // Reads the user ID from the input fields.
		Region:        fields["oci_region"],         // Reads the region from the input fields.
		PrivateKey:    fields["oci_private_key"],    // Reads the private key from the input fields.
		Fingerprint:   fields["oci_fingerprint"],    // Reads the fingerprint from the input fields.
		KeyPassphrase: fields["oci_key_passphrase"], // Reads the private key passphrase from the input fields.
		EmailHost:     fields["email_host"],         // SMTP User
		EmailPort:     fields["email_port"],         // SMTP User
		EmailUser:     fields["email_user"],         // SMTP User
		EmailPassword: fields["email_password"],     // SMTP PWD
	}
	// Validates the populated configuration to ensure all necessary fields are set.
	return config, config.Validate()
}

// Validate ensures that the OCIAuth struct contains all mandatory fields.
//
// Returns:
// - nil if all required fields are populated.
// - An error if any of the required fields (TenancyID, UserID, Region, PrivateKey, or Fingerprint) is missing.
func (o *OCIAuth) Validate() error {
	// Locks the mutex to ensure thread safety during validation.
	o.mu.Lock()
	defer o.mu.Unlock() // Unlocks the mutex after validation is complete.

	// Checks for missing fields and returns errors for each unfulfilled requirement.
	if o.CompartmentID == "" {
		return fmt.Errorf("compartment ID is required")
	}
	if o.TenancyID == "" {
		return fmt.Errorf("tenancy ID is required")
	}
	if o.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if o.Region == "" {
		return fmt.Errorf("region is required")
	}
	if o.PrivateKey == "" {
		return fmt.Errorf("private key is required")
	}
	if o.Fingerprint == "" {
		return fmt.Errorf("fingerprint is required")
	}
	return nil // Validation is successful if all fields are populated.
}

// Authenticate performs the authentication process by configuring an OCI Identity client
// with the provided credentials. It validates the configuration, replaces newline placeholders
// in the private key, and tests the connection by listing available OCI regions.
//
// Returns:
// - nil if authentication succeeds.
// - An error if validation fails, client creation fails, or the test action fails.
func (o *OCIAuth) Authenticate() error {
	// Locks the struct to ensure authentication is thread-safe.
	o.mu.Lock()

	// If the user is already authenticated, skip the process.
	if o.Authenticated {
		o.mu.Unlock()
		return nil
	}
	o.mu.Unlock() // Unlock early to allow validation outside of the lock.

	// Validates the configuration to ensure all required fields are set.
	if err := o.Validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	o.mu.Lock()         // Lock again for setup within the struct.
	defer o.mu.Unlock() // Ensures the mutex is unlocked even if an error occurs.

	// Replace any "\\n" placeholders in the private key with actual newlines ("\n") for proper formatting.
	o.PrivateKey = strings.Replace(o.PrivateKey, "\\n", "\n", -1)

	// Creates a new RawConfigurationProvider with the necessary credentials for OCI services.
	o.privateKeyProvider = common.NewRawConfigurationProvider(
		o.TenancyID,      // The tenancy ID.
		o.UserID,         // The user ID.
		o.Region,         // The OCI region.
		o.Fingerprint,    // The private key's fingerprint.
		o.PrivateKey,     // The private key itself.
		&o.KeyPassphrase, // The private key's passphrase.
	)

	// Uses the configuration provider to create an OCI identity client.
	var err error
	o.Client, err = identity.NewIdentityClientWithConfigurationProvider(o.privateKeyProvider)
	if err != nil {
		// Returns an error if the client cannot be created.
		return fmt.Errorf("unable to create OCI Identity Client: %v", err)
	}

	// Uses the client to retrieve a list of available regions in OCI as a basic test action.
	response, err := o.Client.ListRegions(context.Background())
	if err != nil {
		// Returns an error if the API call to list regions fails.
		return fmt.Errorf("error occurred while listing regions: %v", err)
	}

	// Checks if the list of regions is empty.
	if len(response.Items) > 0 {
		o.Authenticated = true // Sets the authentication status to true on success.
		return nil             // Returns nil to indicate successful authentication.
	}

	return fmt.Errorf("authentication failed: no regions retrieved")
}

func (o *OCIAuth) GetAllRegions() ([]string, error) {
	response, err := o.Client.ListRegions(context.Background())
	if err != nil {
		// Returns an error if the API call to list regions fails.
		return nil, err
	}

	var regions []string

	for i := 0; i < len(response.Items); i++ {
		regions = append(regions, *response.Items[i].Name)
	}

	return regions, nil
}

func (o *OCIAuth) GetAllCompartments(request identity.ListCompartmentsRequest) ([]string, error) {
	response, err := o.Client.ListCompartments(context.Background(), request)
	if err != nil {
		// Returns an error if the API call to list regions fails.
		return nil, err
	}

	var regions []string

	for i := 0; i < len(response.Items); i++ {
		regions = append(regions, *response.Items[i].Name)
	}

	return regions, nil
}

func (o *OCIAuth) GetConfigurationProvider() common.ConfigurationProvider {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.privateKeyProvider
}
