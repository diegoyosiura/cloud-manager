package authentication

import "errors"

// AuthConfig is a general configuration structure that holds the provider name and its associated configuration.
// It uses the Provider interface to abstract provider-specific behavior.
type AuthConfig struct {
	ProviderName string   // Name of the cloud provider (e.g., "aws", "azure", "gcp", "oci").
	Config       Provider // The configuration object for the specific provider implementing the Provider interface.
}

// NewAuthConfig initializes a new instance of AuthConfig based on the given provider name and input fields.
// The function delegates the creation of provider-specific configurations to their respective constructors.
func NewAuthConfig(provider string, fields map[string]string) (*AuthConfig, error) {
	var config Provider
	var err error

	// Determine the provider and create its associated configuration.
	switch provider {
	case "aws":
		config, err = NewAWSAuthFromAuth(fields) // Initializes AWS-specific configuration.
	case "azure":
		config, err = NewAzureAuthFromAuth(fields) // Initializes Azure-specific configuration.
	case "gcp":
		config, err = NewGCPAuthFromAuth(fields) // Initializes GCP-specific configuration.
	case "oci":
		config, err = NewOCIAuthFromAuth(fields) // Initializes OCI-specific configuration.
	default:
		// Return an error if the provider is unsupported.
		return nil, errors.New("unsupported provider: " + provider)
	}

	// Return an error if provider-specific initialization failed.
	if err != nil {
		return nil, err
	}

	// Return a new AuthConfig object containing the provider name and its configuration.
	return &AuthConfig{
		ProviderName: provider,
		Config:       config,
	}, nil
}

// Validate checks if the associated provider's configuration is valid by calling its Validate method.
// It ensures that all required fields are correctly set for the specific provider.
func (a *AuthConfig) Validate() error {
	if a.Config == nil {
		// Return an error if no configuration has been provided for the specified provider.
		return errors.New("no configuration provided for provider: " + a.ProviderName)
	}
	return a.Config.Validate()
}

// Authenticate delegates the authentication logic to the specific provider's Authenticate method.
// It returns an error if the provider's configuration is missing or the authentication fails.
func (a *AuthConfig) Authenticate() error {
	if a.Config == nil {
		// Return an error if no configuration has been provided for the specified provider.
		return errors.New("no configuration provided for provider: " + a.ProviderName)
	}
	return a.Config.Authenticate()
}
