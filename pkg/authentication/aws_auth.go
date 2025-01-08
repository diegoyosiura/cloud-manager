package authentication

import (
	"cloud-manager/internal/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"sync"
)

type AWSAuth struct {
	AccessKeyID     []byte // AWS Access Key ID stored as a byte slice for security
	SecretAccessKey []byte // AWS Secret Access Key stored as a byte slice for security
	EmailHost       string // SMTP Host
	EmailPort       string // SMTP Port
	EmailUser       []byte // SMTP User
	EmailPassword   []byte // SMTP PWD
	Region          string // AWS Region for resource operations

	Authenticated bool             // Tracks if authentication was successful
	Session       *session.Session // AWS Session instance for API interactions

	mu sync.Mutex
}

// NewAWSAuthFromAuth initializes an AWSAuth configuration from a map of fields.
// This function maps input fields into the AWSAuth struct and validates them.
func NewAWSAuthFromAuth(fields map[string]string) (*AWSAuth, error) {
	config := &AWSAuth{
		mu:              sync.Mutex{},
		Authenticated:   false,                                   // Authentication starts as false
		AccessKeyID:     []byte(fields["aws_access_key_id"]),     // Convert key ID to byte slice for security
		SecretAccessKey: []byte(fields["aws_secret_access_key"]), // Convert secret key to byte slice for security
		Region:          fields["aws_region"],                    // Set the region value
		EmailHost:       fields["email_host"],                    // SMTP User
		EmailPort:       fields["email_port"],                    // SMTP User
		EmailUser:       []byte(fields["email_user"]),            // SMTP User
		EmailPassword:   []byte(fields["email_password"]),        // SMTP PWD
	}

	// Validate the configuration to ensure all required fields are present
	if err := config.Validate(); err != nil {
		return nil, err // Return an error if validation fails
	}

	return config, nil // Return the valid AWSAuth instance
}

// Validate ensures that all required AWS authentication fields are provided and non-empty.
// Returns an error if any essential fields are missing or invalid.
func (a *AWSAuth) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var missingFields []string // Slice to accumulate missing fields

	// Check if each required field is empty and add to the missingFields slice
	if len(a.AccessKeyID) == 0 {
		missingFields = append(missingFields, "AccessKeyID")
	}
	if len(a.SecretAccessKey) == 0 {
		missingFields = append(missingFields, "SecretAccessKey")
	}
	if a.Region == "" {
		missingFields = append(missingFields, "Region")
	}

	// If there are any missing fields, return a detailed error message
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required AWS authentication fields: %v", missingFields)
	}

	return nil // Return nil if all fields are valid
}

// InitializeSession sets up the AWS session if it is not already initialized.
// Uses the stored AccessKeyID, SecretAccessKey, and Region for session configuration.
func (a *AWSAuth) initializeSession() error {
	// Check if the session is already initialized to avoid duplication
	if a.Session == nil {
		// Create a configuration using the provided credentials and region
		sessionConfig := &aws.Config{
			Region:      aws.String(a.Region),
			Credentials: credentials.NewStaticCredentials(string(a.AccessKeyID), string(a.SecretAccessKey), ""), // Static credentials
		}

		// Attempt to create a new AWS session
		sess, err := session.NewSession(sessionConfig)
		if err != nil {
			return fmt.Errorf("failed to create AWS session: %w", err) // Return an error if session initialization fails
		}

		// Store the session and mark authentication status as false
		a.Session = sess
		a.Authenticated = false
	}
	return nil // Session initialized successfully
}

// Authenticate establishes a connection to AWS services and validates credentials via STS API.
// Ensures that the authentication is only performed once unless reauthentication is required.
func (a *AWSAuth) Authenticate() error {
	a.mu.Lock()
	// Skip reauthentication if already authenticated
	if a.Authenticated {
		return nil
	}
	a.mu.Unlock()
	if err := a.Validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	// Attempt to initialize the session
	err := a.initializeSession()
	if err != nil {
		return err // Return error if session initialization fails
	}

	// Create an STS (Security Token Service) client using the session
	stsSvc := sts.New(a.Session)

	// Perform a GetCallerIdentity API call to validate credentials
	identityData, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil || identityData == nil {
		return fmt.Errorf("failed to authenticate with AWS STS: %w", err) // Return error if authentication fails
	}

	// Validate the ARN (Amazon Resource Name) returned by STS
	if !utils.IsValidArn(*identityData.Arn) {
		return fmt.Errorf("invalid ARN returned from STS: %s", aws.StringValue(identityData.Arn))
	}

	// If validation is successful, mark as authenticated
	a.Authenticated = true
	return nil
}

// TestAWSAuth validates the AWSAuth configuration and performs an authentication test.
// Ensures both validation and authentication logic function correctly.
func TestAWSAuth(auth *AWSAuth) error {
	// Step 1: Validate the configuration before attempting authentication
	if err := auth.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err) // Return error if validation fails
	}

	// Step 2: Perform authentication test with AWS services
	if err := auth.Authenticate(); err != nil {
		return fmt.Errorf("authentication test failed: %w", err) // Return error if authentication fails
	}

	return nil // Validation and authentication succeeded
}
