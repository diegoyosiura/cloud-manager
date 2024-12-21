package main

import (
	"cloud-manager/internal/utils"
	"cloud-manager/pkg/authentication"
	"fmt"
	"os"
)

func main() {
	// Validate number of arguments; ensure user provides a command.
	if len(os.Args) < 2 {
		fmt.Println("Usage: cloud-manager <authenticate-<provider>>")
		fmt.Println("Available providers: aws, azure, gcp, oci")
		os.Exit(1)
	}

	// Retrieve the command (e.g., authenticate-aws, authenticate-azure).
	command := os.Args[1]

	// Extract provider name from the command by removing the "authenticate-" prefix.
	// Example: command "authenticate-aws" -> provider "aws".
	provider := extractProviderFromCommand(command)
	if provider == "" {
		fmt.Printf("Invalid command. Usage: cloud-manager <authenticate-<provider>>\n")
		os.Exit(1)
	}

	// Load environment variables into a generic map of fields.
	fields := loadEnvVariables(provider)

	// Initialize an AuthConfig instance based on the provider and environment variables.
	authConfig, err := authentication.NewAuthConfig(provider, fields)
	if err != nil {
		fmt.Printf("Failed to initialize authentication for provider '%s': %v\n", provider, err)
		os.Exit(1)
	}

	// Validate the configuration.
	if err := authConfig.Validate(); err != nil {
		fmt.Printf("Validation failed for provider '%s': %v\n", provider, err)
		os.Exit(1)
	}

	// Authenticate using the configuration.
	if err := authConfig.Authenticate(); err != nil {
		fmt.Printf("Authentication failed for provider '%s': %v\n", provider, err)
		os.Exit(1)
	}

	// If successful, print a success message.
	fmt.Printf("Authentication successful for provider '%s'.\n", provider)
}

// extractProviderFromCommand extracts the provider name from the command.
// Example: "authenticate-aws" -> "aws".
func extractProviderFromCommand(command string) string {
	if len(command) > len("authenticate-") && command[:len("authenticate-")] == "authenticate-" {
		return command[len("authenticate-"):]
	}
	return ""
}

// loadEnvVariables loads environment variables into a map based on the provider.
// It retrieves variables specific to each cloud provider as required.
func loadEnvVariables(provider string) map[string]string {
	envVars := map[string]string{}

	switch provider {
	case "aws":
		envVars["aws_access_key_id"] = utils.GetEnvWithValidation("AWS_KEY")         // Access Key ID.
		envVars["aws_secret_access_key"] = utils.GetEnvWithValidation("AWS_SECRETE") // Secret Access Key.
		envVars["aws_region"] = utils.GetEnvWithValidation("AWS_REGION")             // Region.
	case "azure":
		envVars["azure_client_id"] = utils.GetEnvWithValidation("AZURE_CLIENT_KEY")         // Client ID.
		envVars["azure_client_secret"] = utils.GetEnvWithValidation("AZURE_CLIENT_SECRETE") // Client Secret.
		envVars["azure_tenant_id"] = utils.GetEnvWithValidation("AZURE_DIRECTORY_ID")       // Tenant ID.
		envVars["azure_subscription_id"] = utils.GetEnvWithValidation("AZURE_OBJECT_ID")    // Subscription ID.
	case "gcp":
		envVars["gcp_project_id"] = utils.GetEnvWithValidation("GCP_KEY_ID")   // Project ID.
		envVars["gcp_auth_json"] = utils.GetEnvWithValidation("GCP_JSON_INFO") // JSON Credentials.
	case "oci":
		envVars["oci_tenancy_id"] = os.Getenv("ORACLE_API_TENANCY")            // Tenancy ID.
		envVars["oci_user_id"] = os.Getenv("ORACLE_API_USER")                  // User ID.
		envVars["oci_region"] = os.Getenv("ORACLE_API_REGION")                 // Region.
		envVars["oci_private_key"] = os.Getenv("ORACLE_API_PRIVATE_KEY")       // Private Key.
		envVars["oci_fingerprint"] = os.Getenv("ORACLE_API_FINGERPRINT")       // Fingerprint.
		envVars["oci_key_passphrase"] = os.Getenv("ORACLE_API_KEY_PASSPHRASE") // Private Key Passphrase (optional).
	default:
		// Handle unsupported providers by returning an empty map.
		fmt.Printf("Unsupported provider: %s\n", provider)
		os.Exit(1)
	}

	return envVars
}
