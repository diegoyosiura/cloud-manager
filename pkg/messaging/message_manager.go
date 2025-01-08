package messaging

import (
	"cloud-manager/pkg/authentication"
	"fmt"
	"sync"
)

type MessageManager interface {
	AddMessage(m Message)
	AddMessages(m []Message)
	setup() (bool, error)
	CancelSend() (bool, error)
	Send() (chan Message, bool, error)
	SendStatus() (float64, error)
}

func NewMessageManager(authConfig *authentication.AuthConfig) (MessageManager, error) {
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
		return &OciManager{Auth: ociConfig, MessagesMT: &sync.RWMutex{}}, nil
	case "aws":
		// Returns an AWS-specific manager implementation.
		awsConfig, ok := authConfig.Config.(*authentication.AWSAuth)
		if !ok {
			return nil, fmt.Errorf("invalid AWS authentication config")
		}
		return &AWSManager{Auth: awsConfig, MessagesMT: &sync.RWMutex{}}, nil

	default:
		// Returns an error if the cloud provider is unsupported.
		return nil, fmt.Errorf("unsupported provider: %s", authConfig.ProviderName)
	}
}
