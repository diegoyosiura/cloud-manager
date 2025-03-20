package compute

import (
	"github.com/oracle/oci-go-sdk/v65/core"
	"math"
)

// OCIInstanceToVPC converts an OCI Instance object into a generic VPC structure.
// It extracts fields like CPU, GPU, memory, and other details from the instance shape configuration.
func OCIInstanceToVPC(instance core.Instance) VPC {
	// Calculates values such as CPU count, GPU count, memory (GB), and descriptions from the ShapeConfig within the instance.
	CPUCount := int64(0)
	VirtualCPUCount := int64(0)
	GPUCount := int64(0)
	MemoryGB := int64(0)
	CPUDescription := ""
	GPUDescription := ""

	if instance.ShapeConfig.Ocpus != nil {
		CPUCount = int64(math.Round(float64(*instance.ShapeConfig.Ocpus)))
	}
	if instance.ShapeConfig.Vcpus != nil {
		VirtualCPUCount = int64(math.Round(float64(*instance.ShapeConfig.Vcpus)))
	}
	if instance.ShapeConfig.Gpus != nil {
		GPUCount = int64(math.Round(float64(*instance.ShapeConfig.Gpus)))
	}
	if instance.ShapeConfig.MemoryInGBs != nil {
		MemoryGB = int64(math.Round(float64(*instance.ShapeConfig.MemoryInGBs)))
	}
	if instance.ShapeConfig.ProcessorDescription != nil {
		CPUDescription = *instance.ShapeConfig.ProcessorDescription
	}
	if instance.ShapeConfig.GpuDescription != nil {
		GPUDescription = *instance.ShapeConfig.GpuDescription
	}

	vpc := VPC{
		ID:          *instance.Id,
		Name:        *instance.DisplayName,
		Region:      *instance.AvailabilityDomain,
		Provider:    "oci",
		Description: *instance.Shape,

		CPUCount:        CPUCount,
		VirtualCPUCount: VirtualCPUCount,
		CPUDescription:  CPUDescription,
		GPUCount:        GPUCount,
		GPUDescription:  GPUDescription,
		MemoryGB:        MemoryGB,

		ProviderSpecific: instance,
	}

	switch instance.LifecycleState {
	case core.InstanceLifecycleStateMoving:
		vpc.State = VPCStateCreating
		break
	case core.InstanceLifecycleStateProvisioning:
		vpc.State = VPCStateCreating
		break
	case core.InstanceLifecycleStateRunning:
		vpc.State = VPCStateAvailable
		break
	case core.InstanceLifecycleStateStarting:
		vpc.State = VPCStateModifying
		break
	case core.InstanceLifecycleStateStopping:
		vpc.State = VPCStateModifying
		break
	case core.InstanceLifecycleStateStopped:
		vpc.State = VPCStateUnavailable
		break
	case core.InstanceLifecycleStateCreatingImage:
		vpc.State = VPCStateCreating
		break
	case core.InstanceLifecycleStateTerminating:
		vpc.State = VPCStateDeleting
		break
	case core.InstanceLifecycleStateTerminated:
		vpc.State = VPCStateDeleted
		break
	}
	return vpc
}
