package authentication

// Provider is an interface that defines the contract for provider-specific authentication configurations.
// Each provider must implement its own Validate and Authenticate logic.
type Provider interface {
	Validate() error     // Ensures all required fields are properly set for the provider.
	Authenticate() error // Handles the provider-specific authentication logic.
}
