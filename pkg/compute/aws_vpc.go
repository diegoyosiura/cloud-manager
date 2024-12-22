package compute

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWSInstanceToVPC converts an AWS EC2 Instance object into a generic VPC structure.
// This function maps AWS-specific instance properties, such as CPU, memory, and networking,
// into a unified VPC structure usable within the application logic.
// Parameters:
//   - instance: A pointer to an AWS EC2 instance object.
//
// Returns:
//   - A VPC object populated with details from the AWS EC2 instance.
func AWSInstanceToVPC(instance *ec2.Instance) VPC {

	// Initialize variables to hold private and public IPs
	privateIP := ""
	publicIP := ""

	// Extract private IP address if available
	if instance.PrivateIpAddress != nil {
		privateIP = *instance.PrivateIpAddress
	}

	// Extract public IP address if available
	if instance.PublicIpAddress != nil {
		publicIP = *instance.PublicIpAddress
	}

	// Constructing the VPC object
	vpc := VPC{
		ID:          *instance.InstanceId,                 // Instance ID
		Name:        *instance.KeyName,                    // Key name (possibly representing the instance)
		Region:      *instance.Placement.AvailabilityZone, // The availability zone of the instance
		Provider:    "aws",                                // Static value "aws" for provider
		Description: *instance.InstanceType,               // Instance type for its description

		CPUCount: *instance.CpuOptions.CoreCount, // Number of CPU cores
		VirtualCPUCount: *instance.CpuOptions.CoreCount * // Total virtual CPUs based on cores and threads per core
			*instance.CpuOptions.ThreadsPerCore,
		CPUDescription: *instance.Hypervisor, // Hypervisor description (e.g., "xen" or "nitro")
		GPUCount:       0,                    // Placeholder for GPU count (not extracted in this implementation)
		GPUDescription: "",                   // Placeholder for GPU details
		MemoryGB:       0,                    // Placeholder for memory size in GB

		PrivateIP: privateIP, // Resolved private IP address
		PublicIP:  publicIP,  // Resolved public IP address

		ProviderSpecific: instance,                                         // Store the original AWS Instance object
		State:            mapInstanceStateToVPCState(*instance.State.Name), // Map AWS instance state to VPC state
	}

	return vpc
}

// mapInstanceStateToVPCState maps the state of an AWS EC2 instance to a generic VPC state.
// Parameters:
//   - state: A string representing the state of the AWS instance (e.g., "running", "stopped").
//
// Returns:
//   - A VPCStateEnum value that represents the equivalent state in the application's domain model.
func mapInstanceStateToVPCState(state string) VPCStateEnum {
	switch state {
	case "pending", "shutting-down", "stopping":
		return VPCStateModifying // States related to transitioning or modification
	case "running":
		return VPCStateAvailable // Instance is active and available
	case "terminated":
		return VPCStateDeleted // Instance is permanently deleted
	case "stopped":
		return VPCStateUnavailable // Instance is stopped and unavailable
	default:
		return VPCStateUnavailable // Default to unavailable for unknown states
	}
}
