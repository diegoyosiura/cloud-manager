package compute

import (
	"cloud-manager/pkg/authentication"
	"fmt"
)

// Manager is a generic interface for managing VPCs across cloud providers.
// It includes methods for listing, creating, and deleting VPCs in various states.
type Manager interface {
	ListRunningVPCs(map[string]interface{}) ([]VPC, error)  // Lists VPCs in "Running" state.
	ListStartingVPCs(map[string]interface{}) ([]VPC, error) // Lists VPCs in "Starting" state.
	ListStoppingVPCs(map[string]interface{}) ([]VPC, error) // Lists VPCs in "Stopping" state.
	ListStoppedVPCs(map[string]interface{}) ([]VPC, error)  // Lists VPCs in "Stopped" state.
	ListCreatingVPCs(map[string]interface{}) ([]VPC, error) // Lists VPCs in "Creating" state.
	ListDeletingVPCs(map[string]interface{}) ([]VPC, error) // Lists VPCs in "Deleting" state.
	ListDeletedVPCs(map[string]interface{}) ([]VPC, error)  // Lists VPCs in "Deleted" state.
	ListAllVPCs(map[string]interface{}) ([]VPC, error)      // Lists VPCs across all states.
	CreateVPC(name, cidr string) (VPC, error)               // Creates a new VPC.
	DeleteVPC(id string) error                              // Deletes a VPC by ID.
	GetVPC(id string) (VPC, error)                          // Retrieves a specific VPC by ID.
}

// NewVPCManager is a factory function that returns a Manager implementation based on the cloud provider.
func NewVPCManager(authConfig *authentication.AuthConfig) (Manager, error) {
	// Realiza autenticação.
	if err := authConfig.Authenticate(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Caso a autenticação for com OCI, inicializa o cliente da OCI.
	switch authConfig.ProviderName {
	case "oci":
		// Returns an OCI-specific manager implementation.
		ociConfig, ok := authConfig.Config.(*authentication.OCIAuth)
		if !ok {
			return nil, fmt.Errorf("invalid OCI authentication config")
		}
		return &OCIManager{Auth: ociConfig}, nil

	default:
		// Returns an error if the cloud provider is unsupported.
		return nil, fmt.Errorf("unsupported provider: %s", authConfig.ProviderName)
	}
}
