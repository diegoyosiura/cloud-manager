package bucket

import (
	"fmt"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"os"
)

type BucketManager interface {
	List(name string) (r []BucketObject, err error)
	Create(name string, waitCreate bool) error
	Delete(name string) error
	Upload(bucket string, objectName string, f *os.File, partSize int64, threads int) error
	DownloadLink(bucketName string, objectName string, expires int64) (string, error)
	Update(bucket string, objectName string, f *os.File, partSize int64, threads int) error
	DeleteObject(bucketName string, objectName string) error
}

// NewBucketManager
func NewBucketManager(authConfig *authentication.AuthConfig) (BucketManager, error) {
	// Realiza autenticação.
	if err := authConfig.Authenticate(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Caso a autenticação for com OCI, inicializa o cliente da OCI.
	switch authConfig.ProviderName {
	case "oci":
		// Returns an OCI-specific manager implementation.
		ociConfig, ok := authConfig.Config.(*authentication.OCIAuth)
		if !ok {
			return nil, fmt.Errorf("invalid OCI authentication config")
		}
		return &OCIManager{Auth: ociConfig}, nil
	case "aws":
		// Returns an AWS-specific manager implementation.
		awsConfig, ok := authConfig.Config.(*authentication.AWSAuth)
		if !ok {
			return nil, fmt.Errorf("invalid OCI authentication config")
		}
		return &AWSManager{Auth: awsConfig}, nil

	default:
		// Returns an error if the cloud provider is unsupported.
		return nil, fmt.Errorf("unsupported provider: %s", authConfig.ProviderName)
	}
}
