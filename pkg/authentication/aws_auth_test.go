package authentication

import (
	"cloud-manager/internal/utils"
	"testing"
)

func TestAWSAuthIntegration(t *testing.T) {
	auth := AWSAuth{
		AccessKeyID:     []byte(utils.GetEnvWithValidation("AWS_KEY")),
		SecretAccessKey: []byte(utils.GetEnvWithValidation("AWS_SECRETE")),
		Region:          utils.GetEnvWithValidation("AWS_REGION"),
	}

	if err := TestAWSAuth(auth); err != nil {
		t.Fatalf("AWS authentication test failed: %v", err)
	}
}
