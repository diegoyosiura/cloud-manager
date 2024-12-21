package compute

// VPCStateEnum defines the possible states for a Virtual Private Cloud (VPC) lifecycle.
type VPCStateEnum string

// Constants representing the various states of a VPC.
const (
	// VPCStateAvailable The VPC is currently available and active.
	VPCStateAvailable VPCStateEnum = "AVAILABLE"
	// VPCStateUnavailable The VPC is currently unavailable and inactive.
	VPCStateUnavailable VPCStateEnum = "UNAVAILABLE"
	// VPCStateCreating The VPC is being created.
	VPCStateCreating VPCStateEnum = "CREATING"
	// VPCStateModifying The VPC is actively being updated or modified.
	VPCStateModifying VPCStateEnum = "MODIFYING"
	// VPCStateDeleting The VPC is in the process of being deleted.
	VPCStateDeleting VPCStateEnum = "DELETING"
	// VPCStateFailed The VPC has failed creation or encountered an error during modification.
	VPCStateFailed VPCStateEnum = "FAILED"
	// VPCStateDeleted The VPC has been successfully deleted and is no longer present.
	VPCStateDeleted VPCStateEnum = "DELETED"
)
