
# Cloud Manager
![Build Status](https://github.com/diegoyosiura/cloud-manager/actions/workflows/go.yml/badge.svg)
![Build Status](https://github.com/diegoyosiura/cloud-manager/actions/workflows/codeql.yml/badge.svg)
![Build Status](https://github.com/diegoyosiura/cloud-manager/actions/workflows/dependency-review.yml/badge.svg)


**Cloud Manager** is a Go-based application designed to manage AWS VPCs and EC2 instances efficiently.

## Features

- **AWSManager**: A dedicated manager for handling AWS resources.
    - List all VPCs in a specified region.
    - List all stopped EC2 instances.
    - Attach EC2 instances to specific VPCs.

## Example Usage

The following example demonstrates how to use the `AWSManager` to interact with AWS resources.

```go
func main() {
    awsManager, err := NewAWSManager("us-east-1")
    if err != nil {
        fmt.Printf("Error creating AWSManager: %v\n", err)
        return
    }

    // List stopped instances
    stoppedInstances, err := awsManager.ListStoppedInstances()
    if err != nil {
        fmt.Printf("Error listing stopped instances: %v\n", err)
        return
    }

    fmt.Println("Stopped instances:")
    for _, instance := range stoppedInstances {
        fmt.Printf("- Instance ID: %s\n", aws.StringValue(instance.InstanceId))
    }

    // List VPCs
    vpcs, err := awsManager.ListVPCs()
    if err != nil {
        fmt.Printf("Error listing VPCs: %v\n", err)
        return
    }

    fmt.Println("Available VPCs:")
    for _, vpc := range vpcs {
        fmt.Printf("- VPC ID: %s, CIDR: %s\n", aws.StringValue(vpc.VpcId), aws.StringValue(vpc.CidrBlock))
    }

    // Attach an instance to a VPC
    err = awsManager.AttachInstanceToVPC("i-1234567890abcdef0", "subnet-12345678")
    if err != nil {
        fmt.Printf("Error attaching instance to VPC: %v\n", err)
    }
}
```

## Getting Started

### Prerequisites

- Install Go (1.19 or later).
- Configure AWS credentials and set the `AWS_REGION` environment variable.

### Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:diegoyosiura/cloud-manager.git
   cd cloud-manager
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

Run the `main.go` to interact with AWS resources:

```bash
go run main.go
```

### Testing

To run tests for the application, use:

```bash
go test ./... -v
```

## License

This project is licensed under the MIT License.
