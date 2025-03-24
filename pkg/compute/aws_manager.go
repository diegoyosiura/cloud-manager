package compute

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"github.com/oracle/oci-go-sdk/v65/common"
)

// AWSManager provides functionality for managing AWS VPCs and their lifecycle states.
// It abstracts AWS SDK interactions, enabling listing, creating, deleting, and retrieving VPCs.
type AWSManager struct {
	Auth   *authentication.AWSAuth // Stores AWS authentication and session configurations.
	Ec2Svc *ec2.EC2                // AWS EC2 Service client for managing VPCs.
}

// ListVPCs retrieves a list of VPCs filtered by lifecycle state and additional custom parameters.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//   - instanceStateCode: A string representing the lifecycle state of instances (e.g., "running", "stopped").
//
// Returns:
//   - A slice of `VPC` objects that match the inputs.
//   - An error if the operation fails.
func (m *AWSManager) ListVPCs(fields map[string]interface{}, instanceStateCode string) ([]VPC, error) {
	// Lazily initialize Ec2Svc if not already set
	if m.Ec2Svc == nil {
		m.Ec2Svc = ec2.New(m.Auth.Session)
	}

	// Convert the fields map to AWS DescribeInstancesInput
	input := convertMapDescribeInstancesInput(fields)

	// Add lifecycle state filter, if specified
	if instanceStateCode != "" {
		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{common.String(instanceStateCode)},
		})
	}

	// Describe instances through AWS SDK
	result, err := m.Ec2Svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	// Convert AWS instance data into custom VPC objects
	var response []VPC
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			response = append(response, AWSInstanceToVPC(instance))
		}
	}
	return response, nil
}

// convertMapDescribeInstancesInput converts a map of filter fields into an AWS SDK DescribeInstancesInput object.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing request attributes.
//
// Returns:
//   - A pointer to an AWS DescribeInstancesInput object.
func convertMapDescribeInstancesInput(fields map[string]interface{}) *ec2.DescribeInstancesInput {
	if value, ok := fields["aws_describe_instances_input"]; ok {
		if input, valid := value.(*ec2.DescribeInstancesInput); valid {
			return input
		}
	}
	// Default to an empty DescribeInstancesInput object
	return &ec2.DescribeInstancesInput{}
}

// ListRunningVPCs retrieves a list of VPCs with instances in the "running" state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListRunningVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "running")
}

// ListStartingVPCs retrieves a list of VPCs with instances in the "pending" (starting) state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListStartingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "pending")
}

// ListStoppingVPCs retrieves a list of VPCs with instances in the "stopping" state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListStoppingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "stopping")
}

// ListStoppedVPCs retrieves a list of VPCs with instances in the "stopped" state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListStoppedVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "stopped")
}

// ListCreatingVPCs retrieves a list of VPCs with instances in the "pending" (creating) state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListCreatingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "pending")
}

// ListDeletingVPCs retrieves a list of VPCs with instances in the "pending" (deleting) state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListDeletingVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "pending")
}

// ListDeletedVPCs retrieves a list of VPCs with instances in the "terminated" (deleted) state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListDeletedVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "terminated")
}

// ListAllVPCs retrieves a list of all VPCs, regardless of lifecycle state.
// Parameters:
//   - fields: A map (`map[string]interface{}`) containing optional filters for the request.
//
// Returns:
//   - A slice of `VPC` objects.
//   - An error if the operation fails.
func (m *AWSManager) ListAllVPCs(fields map[string]interface{}) ([]VPC, error) {
	return m.ListVPCs(fields, "")
}

// CreateVPC creates a new VPC with the specified name and CIDR block.
// Parameters:
//   - name: The name of the VPC to create.
//   - cidr: The CIDR block for the new VPC.
//
// Returns:
//   - A `VPC` object representing the created VPC (placeholder).
//   - An error if the operation fails.
func (m *AWSManager) CreateVPC(name, cidr string) (*VPC, error) {
	return &VPC{}, nil
}

// DeleteVPC deletes a VPC with the specified ID.
// Parameters:
//   - id: The ID of the VPC to delete.
//
// Returns:
//   - An error if the operation fails (placeholder implementation).
func (m *AWSManager) DeleteVPC(id string) error {
	return nil
}

// GetVPC retrieves the details of a VPC with the specified ID.
// Parameters:
//   - id: The ID of the VPC to retrieve.
//
// Returns:
//   - A `VPC` object representing the retrieved VPC (placeholder).
//   - An error if the operation fails.
func (m *AWSManager) GetVPC(id string) (*VPC, error) {
	// Lazily initialize Ec2Svc if not already set
	if m.Ec2Svc == nil {
		m.Ec2Svc = ec2.New(m.Auth.Session)
	}

	result, err := m.Ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{&id}})

	if err != nil {
		return nil, err
	}

	var response []VPC
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			response = append(response, AWSInstanceToVPC(instance))
		}
	}

	if len(response) != 1 {
		return nil, errors.New("invalid instance count")
	}
	return &response[0], nil
}
func (m *AWSManager) Start(id string) (*VPC, error) {
	request, _ := m.Ec2Svc.StartInstancesRequest(&ec2.StartInstancesInput{InstanceIds: []*string{&id}})
	err := request.Send()
	if err != nil {
		return nil, err
	}
	return m.GetVPC(id)
}

func (m *AWSManager) Stop(id string) (*VPC, error) {
	request, _ := m.Ec2Svc.StopInstancesRequest(&ec2.StopInstancesInput{InstanceIds: []*string{&id}})
	err := request.Send()
	if err != nil {
		return nil, err
	}
	return m.GetVPC(id)
}

func (m *AWSManager) Restart(id string) (*VPC, error) {
	request, _ := m.Ec2Svc.RebootInstancesRequest(&ec2.RebootInstancesInput{InstanceIds: []*string{&id}})
	err := request.Send()
	if err != nil {
		return nil, err
	}
	return m.GetVPC(id)
}
