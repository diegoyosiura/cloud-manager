package compute

// VPC is a generic and extensible representation of a Virtual Private Cloud (VPC) instance.
// It allows uniform representation of VPCs across different cloud providers.
type VPC struct {
	ID              string       `json:"id"`                // Unique identifier for the VPC.
	Name            string       `json:"name"`              // Display name of the VPC.
	Region          string       `json:"region"`            // Region where the VPC resides.
	Provider        string       `json:"provider"`          // Cloud provider (e.g., "oci", "aws", etc.).
	Description     string       `json:"description"`       // Detailed description of the VPC (e.g., shape or configuration).
	CidrBlock       string       `json:"cidr_block"`        // CIDR block associated with the VPC.
	PublicIP        string       `json:"public_ip"`         // CIDR block associated with the VPC.
	PrivateIP       string       `json:"private_ip"`        // CIDR block associated with the VPC.
	State           VPCStateEnum `json:"state"`             // Current state of the VPC (e.g., "available", "creating", "deleting").
	CPUCount        int64        `json:"cpu_count"`         // Number of physical CPUs (if applicable).
	VirtualCPUCount int64        `json:"virtual_cpu_count"` // Number of virtual CPUs.
	CPUDescription  string       `json:"cpu_description"`   // Description of the CPU type.
	GPUCount        int64        `json:"gpu_count"`         // Number of GPUs (if applicable).
	GPUDescription  string       `json:"gpu_description"`   // Description of the GPU type.
	MemoryGB        int64        `json:"memory_gb"`         // Total memory in GB.

	// ProviderSpecific holds provider-specific details about the VPC.
	// For OCI, use the OCIInstance; for other providers, use respective implementations.
	ProviderSpecific interface{} `json:"providerSpecific"`
}
