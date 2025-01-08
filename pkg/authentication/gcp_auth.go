package authentication

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"sync"
)

// GCPAuth is a structure that encapsulates the configuration and state
// for authenticating with Google Cloud Platform (GCP) services, specifically
// the Google Cloud Storage API.
type GCPAuth struct {
	ProjectID     string          // The GCP project ID for the service account.
	AuthJSON      string          // JSON-encoded credentials for the service account.
	Authenticated bool            // Tracks if authentication was successful.
	EmailHost     string          // SMTP Host
	EmailPort     string          // SMTP Port
	EmailUser     string          // SMTP User
	EmailPassword string          // SMTP PWD
	Client        *storage.Client // GCP Storage Client instance for interacting with resources.

	mu sync.Mutex // Mutex to ensure thread-safe access to the struct.
}

// NewGCPAuthFromAuth creates a new GCPAuth instance, initializing it with fields
// extracted from a map[string]string and validating its configuration.
//
// Parameters:
//   - fields: A map containing required GCP fields, such as "gcp_project_id"
//     and "gcp_auth_json".
//
// Returns:
// - A pointer to the initialized GCPAuth instance or an error if validation fails.
func NewGCPAuthFromAuth(fields map[string]string) (*GCPAuth, error) {
	config := &GCPAuth{
		mu:            sync.Mutex{},             // Initialize mutex for thread safety.
		Authenticated: false,                    // Set default authentication state to false.
		ProjectID:     fields["gcp_project_id"], // Read the project ID from input.
		AuthJSON:      fields["gcp_auth_json"],  // Read the credentials JSON from input.
		EmailHost:     fields["email_host"],     // SMTP User
		EmailPort:     fields["email_port"],     // SMTP User
		EmailUser:     fields["email_user"],     // SMTP User
		EmailPassword: fields["email_password"], // SMTP PWD
	}
	// Validate the configuration immediately after initialization.
	return config, config.Validate()
}

// Validate checks that the required fields (ProjectID and AuthJSON) are populated.
//
// Returns:
// - nil if the validation is successful.
// - An error if any mandatory field is missing.
func (g *GCPAuth) Validate() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check for missing fields.
	if g.ProjectID == "" || g.AuthJSON == "" {
		return fmt.Errorf("missing required GCP authentication fields")
	}
	return nil // All fields are valid.
}

// Authenticate performs the actual authentication process by creating a GCP Storage client
// and validating access by listing buckets in the specified project.
//
// Returns:
// - nil if authentication is successful.
// - An error if validation fails, credentials are invalid, or access to GCP resources is denied.
func (g *GCPAuth) Authenticate() error {
	g.mu.Lock()

	// Skip reauthentication if already authenticated.
	if g.Authenticated {
		return nil
	}

	var err error

	g.mu.Unlock()
	// Validate the configuration before attempting authentication.
	if err := g.Validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	// Create a GCP Storage client using the provided JSON credentials.
	ctx := context.Background() // Use a background context for client creation.
	g.Client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(g.AuthJSON)))
	if err != nil {
		// Return an error if the client creation fails.
		return fmt.Errorf("failed to create GCP client: %v", err)
	}

	// Perform a simple resource access test by listing buckets in the given project.
	it := g.Client.Buckets(ctx, g.ProjectID)
	_, err = it.Next()
	if err != nil {
		// Handle specific error cases, such as missing permissions or no buckets found.
		if errors.Is(err, storage.ErrBucketNotExist) {
			return fmt.Errorf("bucket does not exist or no access to buckets: %v", err)
		}
		// General error for failed bucket listing.
		return fmt.Errorf("failed to list buckets: %v", err)
	}

	// Mark the authentication as successful upon completion.
	g.Authenticated = true
	return nil
}

// TestGCPAuth is a utility function to validate and test the authentication process.
//
// Parameters:
// - auth: A GCPAuth instance to validate and authenticate.
//
// Returns:
// - nil if the validation and authentication succeed.
// - An error if validation or authentication fails.
func TestGCPAuth(auth *GCPAuth) error {
	// Step 1: Validate the provided configuration.
	if err := auth.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err) // Return error if validation fails.
	}

	// Step 2: Attempt to authenticate using the provided configuration.
	if err := auth.Authenticate(); err != nil {
		return fmt.Errorf("authentication test failed: %w", err) // Return error if authentication fails.
	}

	// Return nil if both validation and authentication are successful.
	return nil
}
