# **OCI Authentication Library**

This Go library facilitates authentication with **Oracle Cloud Infrastructure (OCI)** services. It provides an easy-to-use interface to validate and authenticate users using OCI credentials, along with utility methods to ensure the configuration is correct.

---

## **Features**

- Validates OCI authentication parameters, like `TenancyID`, `UserID`, `Region`, `PrivateKey`, and `Fingerprint`.
- Performs the authentication process using OCI's Identity Client.
- Thread-safe implementation using `sync.Mutex` for safe concurrent access.
- Includes testing functionality to verify successful authentication via listing OCI regions.

---

## **Getting Started**

### **Prerequisites**

1. **Go Environment**: Ensure that Go 1.21+ is installed on your system.
2. **OCI SDK**: This library depends on Oracle's official Go SDK.

To install the OCI SDK:
```bash
go get github.com/oracle/oci-go-sdk/v65
```

---

### **Installation**

You can clone this repository or integrate the library by importing its package file(s) into your project:
```bash
git clone <repository_url>
cd <repository_folder>
```

---

### **Usage**

This library provides utility functions to perform validation, authentication, and testing. Below is an example usage of the library:

```go
package main

import (
	"log"
	"authentication" // Replace with the actual package path
)

func main() {
	// Define your OCI credentials
	authFields := map[string]string{
		"oci_tenancy_id":    "ocid1.tenancy.oc1..<unique_ID>",
		"oci_user_id":       "ocid1.user.oc1..<unique_ID>",
		"oci_region":        "us-ashburn-1",
		"oci_private_key":   "<your-private-key>",
		"oci_fingerprint":   "<your-fingerprint>",
		"oci_key_passphrase": "<your-key-passphrase>", // Optional
	}

	// Initialize the OCI Auth configuration
	auth, err := authentication.NewOCIAuthFromAuth(authFields)
	if err != nil {
		log.Fatalf("Failed to initialize OCI authentication: %v", err)
	}

	// Test Authentication
	if err := authentication.TestOCIAuth(auth); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	log.Println("Authentication successful!")
}
```

---

## **Configuration Fields**

The following fields must be provided to initialize and authenticate with OCI:

| Field Name          | Description                                                                                     | Required | Example           |
|---------------------|-------------------------------------------------------------------------------------------------|----------|-------------------|
| `oci_tenancy_id`    | The tenancy ID for your OCI account.                                                           | Yes      | `"ocid1.tenancy.oc1..<unique_ID>"` |
| `oci_user_id`       | The user ID for authentication.                                                                | Yes      | `"ocid1.user.oc1..<unique_ID>"` |
| `oci_region`        | The OCI region for your requests.                                                              | Yes      | `"us-ashburn-1"`  |
| `oci_private_key`   | The private key used for authentication.                                                       | Yes      | `"-----BEGIN PRIVATE KEY-----\n<key>\n-----END PRIVATE KEY-----"` |
| `oci_fingerprint`   | The fingerprint of your private key.                                                           | Yes      | `"01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF"` |
| `oci_key_passphrase`| The passphrase for your private key (only if required).                                         | No       | `"your-passphrase"` |

---

## **Key Methods**

### **NewOCIAuthFromAuth(fields map[string]string) (*OCIAuth, error)**

- **Description**: Initializes a new instance of `OCIAuth` with the given configuration fields and validates them.
- **Parameters**:
    - `fields`: A map containing the OCI authentication details.
- **Returns**:
    - An `OCIAuth` instance or a validation error if required fields are missing.

### **Validate() error**

- **Description**: Ensures all mandatory fields are correctly provided in the configuration.
- **Returns**:
    - Returns `nil` if validation succeeds or an error indicating which field is missing.

### **Authenticate() error**

- **Description**: Authenticates with OCI by creating an Identity Client using the configured credentials. Performs a basic test by listing OCI regions to ensure the configuration is correct.
- **Returns**:
    - Returns `nil` if authentication is successful or an error if any step fails.

### **TestOCIAuth(auth *OCIAuth) error**

- **Description**: Combines validation and authentication to ensure the provided credentials work as expected.
- **Parameters**:
    - `auth`: The `OCIAuth` instance to be tested.
- **Returns**:
    - Returns `nil` on success or an error on failure.

---

## **Error Handling**

The library provides meaningful error messages throughout the validation and authentication processes. For example:
- If a required field is missing, `Validate()` will indicate which field is missing.
- Network errors during authentication and region listing will return detailed error messages to facilitate troubleshooting.

---

## **Logging**

The library uses Go's standard `log` package to display logs to the console. To integrate advanced logging, you can replace `log` with structured logging libraries like `logrus` or `zap`.

Example with `logrus`:
```go
import "github.com/sirupsen/logrus"

logrus.SetFormatter(&logrus.JSONFormatter{})
logrus.Info("Logged with logrus")
```

---

## **Testing**

### **Unit Tests**
You can write unit tests to verify various functionalities, such as:
1. Validation of input fields.
2. Authenticating with mock credentials.
3. Simulating OCI API calls (using mocks).

### **Manual Testing**
To manually test the authentication process, ensure your OCI credentials are correctly set and run the following:
```bash
go run main.go
```

---

## **Dependencies**

This library requires the following dependencies:
1. **OCI SDK for Go**: [github.com/oracle/oci-go-sdk](https://github.com/oracle/oci-go-sdk)
2. **Standard Packages**:
    - `context`
    - `sync`
    - `strings`
    - `log`

---

## **Roadmap**

Future improvements for the library could include:
- **Improved error messages** with additional details for failed OCI API requests.
- **Timeout configurations** for network requests during authentication.
- **Sensitive data cleanup** to remove private key/passphrase from memory after use.
- **Structured logging** using libraries like `logrus` or `zap`.

---

## **Contributing**

We welcome contributions to improve this library! To contribute:
1. Fork this repository.
2. Create a feature branch: `git checkout -b feature-name`.
3. Submit a pull request with a detailed description of your changes.

---

## **License**

This library is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

## **Contact**

For questions or suggestions, please contact the maintainer or open an issue on the repository.