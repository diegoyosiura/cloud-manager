package compute

import (
	"context"
	"errors"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// OCIManager manages VPC-related operations in Oracle Cloud Infrastructure (OCI).
// It interacts with the OCI SDK for tasks like listing, creating, and deleting VPCs.
type OCIManager struct {
	Auth   *authentication.OCIAuth // OCI authentication details.
	Client *core.ComputeClient     // OCI Compute Client for interacting with OCI services.
}

// ListVPCs filters VPCs based on a lifecycle state and additional fields.
// Parameters:
// - fields: A generic map where keys (e.g., "oci_compartment_id") provide filtering options.
// - enum: The lifecycle state to filter VPCs (e.g., Running, Stopped).
// Returns: A list of filtered VPCs or an error if the request fails.
func (m *OCIManager) ListVPCs(fields map[string]interface{}, enum core.InstanceLifecycleStateEnum) ([]VPC, error) {
	if m.Client == nil {
		cl, err := core.NewComputeClientWithConfigurationProvider(m.Auth.GetConfigurationProvider())
		if err != nil {
			return nil, err
		}
		m.Client = &cl
	}

	request := convertMapInstanceRequest(fields)
	request.LifecycleState = enum

	if value, ok := fields["oci_compartment_id"]; !ok {
		return nil, errors.New("oci_compartment_id is required")
	} else {
		request.CompartmentId = common.String(value.(string))
	}

	resp, err := m.Client.ListInstances(context.Background(), request)

	if err != nil {
		return nil, err
	}

	var response []VPC
	for _, vpc := range resp.Items {
		response = append(response, OCIInstanceToVPC(vpc))
	}
	return response, nil
}

// convertMapInstanceRequest converts the "fields" map into an OCI ListInstancesRequest.
// Default values are used if the "oci_instance_request" field is not provided.
func convertMapInstanceRequest(fields map[string]interface{}) core.ListInstancesRequest {
	if value, ok := fields["oci_instance_request"]; !ok {
		return core.ListInstancesRequest{
			Limit:     common.Int(100),
			SortOrder: core.ListInstancesSortOrderDesc,
			SortBy:    core.ListInstancesSortByTimecreated,
		}
	} else {
		return value.(core.ListInstancesRequest)
	}
}

// Various List functions specialize in filtering VPCs by lifecycle state.
// These include:
// - ListRunningVPCs: Lists VPCs in the "Running" state.
// - ListStoppingVPCs: Lists VPCs in the "Stopping" state.
// - ListStoppedVPCs: Lists VPCs in the "Stopped" state.
// - ListCreatingVPCs: Lists VPCs in the "Creating" state.
// - ListDeletingVPCs: Lists VPCs in the "Deleting" state.
// - ListDeletedVPCs: Lists VPCs in the "Deleted" state.
// - ListAllVPCs: Aggregates all VPCs from any lifecycle state.

func (m *OCIManager) ListRunningVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateRunning)
}

func (m *OCIManager) ListStartingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateStarting)
}

func (m *OCIManager) ListStoppingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateStopping)
}
func (m *OCIManager) ListStoppedVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateStopped)
}

func (m *OCIManager) ListCreatingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateProvisioning)
}

func (m *OCIManager) ListDeletingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateTerminating)
}

func (m *OCIManager) ListDeletedVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, core.InstanceLifecycleStateTerminated)
}

func (m *OCIManager) ListAllVPCs(fields map[string]interface{}) ([]VPC, error) {
	var response []VPC

	for _, v := range core.GetInstanceLifecycleStateEnumStringValues() {
		state, _ := core.GetMappingInstanceLifecycleStateEnum(v)
		vpcs, err := m.ListVPCs(fields, state)

		if err != nil {
			return nil, err
		}

		response = append(response, vpcs...)
	}
	return response, nil
}
func (m *OCIManager) CreateVPC(name, cidr string) (VPC, error) {
	return VPC{}, nil
}

func (m *OCIManager) DeleteVPC(id string) error {
	return nil
}

func (m *OCIManager) GetVPC(id string) (VPC, error) {
	return VPC{}, nil
}
