package authentication

import (
	"testing"
)

// TestNewAuthConfig_ValidAWS verifica se uma configuração válida da AWS é corretamente inicializada.
func TestNewAuthConfig_ValidAWS(t *testing.T) {
	fields := map[string]string{
		"aws_access_key_id":     "testAccessKey",
		"aws_secret_access_key": "testSecretKey",
		"aws_region":            "us-east-1",
	}

	config, err := NewAuthConfig("aws", fields)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if config.ProviderName != "aws" {
		t.Errorf("esperado provider 'aws', mas foi recebido '%s'", config.ProviderName)
	}
}

// TestNewAuthConfig_InvalidProvider garante que um provedor não suportado retorna o erro apropriado.
func TestNewAuthConfig_InvalidProvider(t *testing.T) {
	fields := map[string]string{
		"irrelevant_field": "value",
	}

	_, err := NewAuthConfig("unsupported", fields)
	if err == nil {
		t.Fatalf("esperado erro para provedor não suportado, mas foi recebido nil")
	}

	expectedErr := "unsupported provider: unsupported"
	if err.Error() != expectedErr {
		t.Fatalf("mensagem de erro esperada: %s, mas foi recebido: %v", expectedErr, err)
	}
}

// TestNewAuthConfig_MissingAWSFields verifica se campos ausentes na configuração AWS disparam erros no construtor.
func TestNewAuthConfig_MissingAWSFields(t *testing.T) {
	fields := map[string]string{
		"aws_secret_access_key": "testSecretKey",
		"aws_region":            "us-east-1",
	}

	_, err := NewAuthConfig("aws", fields)
	if err == nil {
		t.Fatalf("esperado erro para campo AWSAccessKeyID ausente, mas foi recebido nil")
	}

	expectedErr := "missing required AWS authentication fields: [AccessKeyID]"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro inesperada: esperado %s, mas recebido: %v", expectedErr, err)
	}
}

// TestNewAuthConfig_ValidOCI verifica se uma configuração válida da OCI é validada corretamente.
func TestNewAuthConfig_ValidOCI(t *testing.T) {
	fields := map[string]string{
		"oci_tenancy_id":     "testTenancyID",
		"oci_user_id":        "testUserID",
		"oci_region":         "us-ashburn-1",
		"oci_private_key":    "testPrivateKey",
		"oci_fingerprint":    "testFingerprint",
		"oci_key_passphrase": "testPassphrase",
	}

	config, err := NewAuthConfig("oci", fields)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if config.ProviderName != "oci" {
		t.Errorf("esperado provider 'oci', mas foi recebido '%s'", config.ProviderName)
	}
}

// TestNewAuthConfig_MissingOCIFields verifica se campos ausentes na configuração OCI disparam erros no construtor.
func TestNewAuthConfig_MissingOCIFields(t *testing.T) {
	fields := map[string]string{
		"oci_user_id":     "testUserID",
		"oci_region":      "us-ashburn-1",
		"oci_private_key": "testPrivateKey",
		"oci_fingerprint": "testFingerprint",
	}

	_, err := NewAuthConfig("oci", fields)
	if err == nil {
		t.Fatalf("esperado erro para campo OCITenancyID ausente, mas foi recebido nil")
	}

	expectedErr := "tenancy ID is required"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro inesperada: esperado %s, mas recebido: %v", expectedErr, err)
	}
}
