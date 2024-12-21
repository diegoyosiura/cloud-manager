
# Cloud Manager

**Cloud Manager** is a Go-based application designed to manage cloud resources across multiple providers, including AWS, Azure, Google Cloud Platform (GCP), and Oracle Cloud Infrastructure (OCI). The project emphasizes best practices and leverages the latest features of the Go programming language to ensure efficient and optimized cloud resource management.

## Project Structure

The project is organized into several directories, each serving a specific purpose:

```plaintext
cloud-manager/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   └── utils/
│       ├── utils.go
│       └── utils_test.go
├── pkg/
│   └── authentication/
│       ├── auth.go
│       ├── auth_test.go
│       ├── aws_auth.go
│       ├── aws_auth_test.go
│       ├── azure_auth.go
│       ├── azure_auth_test.go
│       ├── gcp_auth.go
│       ├── gcp_auth_test.go
│       ├── oci_auth.go
│       ├── oci_auth_test.go
│       └── debug.go
├── .gitignore
└── go.mod
```

### Directories and Their Roles

- **cmd/**: Contains the entry point of the application. The `main.go` file initializes and starts the application.

- **internal/**: Houses internal packages that are not intended for external use.
    - **config/**: Manages configuration settings for the application.
        - `config.go`: Handles loading and parsing of configuration files.
        - `config_test.go`: Contains tests for configuration management.
    - **utils/**: Provides utility functions used throughout the application.
        - `utils.go`: Includes helper functions for tasks like environment variable management.
        - `utils_test.go`: Contains tests for utility functions.

- **pkg/**: Contains packages that can be imported by other applications or services.
    - **authentication/**: Manages authentication with various cloud providers.
        - `auth.go`: Defines common authentication interfaces and structures.
        - `auth_test.go`: Contains tests for the common authentication functionalities.
        - `aws_auth.go`: Implements authentication methods for AWS.
        - `aws_auth_test.go`: Contains tests for AWS authentication.
        - `azure_auth.go`: Implements authentication methods for Azure.
        - `azure_auth_test.go`: Contains tests for Azure authentication.
        - `gcp_auth.go`: Implements authentication methods for GCP.
        - `gcp_auth_test.go`: Contains tests for GCP authentication.
        - `oci_auth.go`: Implements authentication methods for OCI.
        - `oci_auth_test.go`: Contains tests for OCI authentication.
        - `debug.go`: Provides debugging utilities for authentication processes.

## Authentication Implementations

The application provides authentication mechanisms for multiple cloud providers, each implemented in its respective file within the `pkg/authentication/` directory. Below is an overview of the authentication structures and methods for each provider:

### AWS Authentication (`aws_auth.go`)

- **AWSAuth Structure**: Holds AWS credentials and region information.
  ```go
  type AWSAuth struct {
      AccessKeyID     string
      SecretAccessKey string
      Region          string
  }
  ```

- **Validate Method**: Ensures all necessary fields are populated.
  ```go
  func (a *AWSAuth) Validate() error {
      // Validation logic
  }
  ```

- **TestAWSAuth Function**: Tests the AWS authentication by creating a session and performing a simple operation.
  ```go
  func TestAWSAuth(auth AWSAuth) error {
      // Authentication testing logic
  }
  ```

### Azure Authentication (`azure_auth.go`)

- **AzureAuth Structure**: Holds Azure credentials and subscription information.
  ```go
  type AzureAuth struct {
      ClientID       string
      ClientSecret   string
      TenantID       string
      SubscriptionID string
  }
  ```

- **Validate Method**: Ensures all necessary fields are populated.
  ```go
  func (a *AzureAuth) Validate() error {
      // Validation logic
  }
  ```

- **TestAzureAuth Function**: Tests the Azure authentication by creating a resource groups client and performing a simple operation.
  ```go
  func TestAzureAuth(auth AzureAuth) error {
      // Authentication testing logic
  }
  ```

### GCP Authentication (`gcp_auth.go`)

- **GCPAuth Structure**: Holds GCP project ID and authentication JSON.
  ```go
  type GCPAuth struct {
      ProjectID string
      AuthJSON  string
  }
  ```

- **Validate Method**: Ensures all necessary fields are populated.
  ```go
  func (g *GCPAuth) Validate() error {
      // Validation logic
  }
  ```

- **TestGCPAuth Function**: Tests the GCP authentication by creating a storage client and performing a simple operation.
  ```go
  func TestGCPAuth(auth GCPAuth) error {
      // Authentication testing logic
  }
  ```

### OCI Authentication (`oci_auth.go`)

- **OCIAuth Structure**: Holds OCI credentials and tenancy information.
  ```go
  type OCIAuth struct {
      TenancyID     string
      UserID        string
      PrivateKey    string
      Fingerprint   string
      Region        string
      PrivateKeyPassphrase string
  }
  ```

- **Validate Method**: Ensures all necessary fields are populated.
  ```go
  func (o *OCIAuth) Validate() error {
      // Validation logic
  }
  ```

- **TestOCIAuth Function**: Tests the OCI authentication by creating a compute client and performing a simple operation.
  ```go
  func TestOCIAuth(auth OCIAuth) error {
      // Authentication testing logic
  }
  ```

## Configuration Management

The `internal/config/config.go` file manages the application's configuration settings. It includes functions to load and parse configuration files, ensuring that all necessary settings are correctly initialized before the application runs.

## Utility Functions

The `internal/utils/utils.go` file provides utility functions used throughout the application. For example:

- **GetEnvWithValidation Function**: Retrieves an environment variable and ensures it is not empty.
  ```go
  func GetEnvWithValidation(key string) string {
      // Function logic
  }
  ```

- **GetOptionalEnv Function**: Retrieves an environment variable or returns a default value if the variable is not set.
  ```go
  func GetOptionalEnv(key, defaultValue string) string {
      // Function logic
  }
  ```

## Getting Started

To build and run the application, ensure you have Go installed and properly configured. Clone the repository, navigate to the project directory, and execute the following commands:

```bash
go mod tidy
go build -o cloud-manager ./cmd
./cloud-manager
```

## Testing

The project includes tests for various components, located alongside their respective implementations. To run the tests, use the `go test` command:

```bash
go test ./...
```

This command will execute all tests in the project, ensuring that each component functions as expected.

## Contributing