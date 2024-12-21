package authentication

import (
	"testing"
)

// TestNewAWSAuthFromAuth_Valid verifica se a inicialização de AWSAuth com entradas válidas ocorre sem erros.
func TestNewAWSAuthFromAuth_Valid(t *testing.T) {
	fields := map[string]string{
		"aws_access_key_id":     "test-access-key-id",
		"aws_secret_access_key": "test-secret-access-key",
		"aws_region":            "us-east-1",
	}

	auth, err := NewAWSAuthFromAuth(fields)
	if err != nil {
		t.Fatalf("erro inesperado ao criar AWSAuth: %v", err)
	}

	if string(auth.AccessKeyID) != "test-access-key-id" {
		t.Errorf("esperado AccessKeyID 'test-access-key-id', recebido '%s'", auth.AccessKeyID)
	}

	if string(auth.SecretAccessKey) != "test-secret-access-key" {
		t.Errorf("esperado SecretAccessKey 'test-secret-access-key', recebido '%s'", auth.SecretAccessKey)
	}

	if auth.Region != "us-east-1" {
		t.Errorf("esperado Region 'us-east-1', recebido '%s'", auth.Region)
	}
}

// TestNewAWSAuthFromAuth_Invalid verifica erros ao inicializar AWSAuth com campos ausentes.
func TestNewAWSAuthFromAuth_Invalid(t *testing.T) {
	fields := map[string]string{
		"aws_access_key_id":     "",
		"aws_secret_access_key": "test-secret-access-key",
		"aws_region":            "",
	}

	_, err := NewAWSAuthFromAuth(fields)
	if err == nil {
		t.Fatalf("esperado erro ao usar campos inválidos, mas nenhum erro foi retornado")
	}

	expectedErr := "missing required AWS authentication fields: [AccessKeyID Region]"
	if err.Error() != expectedErr {
		t.Errorf("mensagem incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestAWSAuth_Validate_ValidFields verifica a validação com dados válidos.
func TestAWSAuth_Validate_ValidFields(t *testing.T) {
	auth := &AWSAuth{
		AccessKeyID:     []byte("test-access-key-id"),
		SecretAccessKey: []byte("test-secret-access-key"),
		Region:          "us-east-1",
	}

	if err := auth.Validate(); err != nil {
		t.Errorf("erro inesperado na validação de dados válidos: %v", err)
	}
}

// TestAWSAuth_Validate_MissingFields verifica se a validação falha com campos ausentes.
func TestAWSAuth_Validate_MissingFields(t *testing.T) {
	auth := &AWSAuth{
		AccessKeyID:     nil,
		SecretAccessKey: nil,
		Region:          "",
	}

	err := auth.Validate()
	if err == nil {
		t.Fatalf("esperado erro para validação com campos inválidos, mas retornou nil")
	}

	expectedErr := "missing required AWS authentication fields: [AccessKeyID SecretAccessKey Region]"
	if err.Error() != expectedErr {
		t.Errorf("mensagem incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestAWSAuth_Authenticate_AlreadyAuthenticated verifica autenticação quando já autenticado.
func TestAWSAuth_Authenticate_AlreadyAuthenticated(t *testing.T) {
	auth := &AWSAuth{
		Authenticated: true,
	}

	err := auth.Authenticate()
	if err != nil {
		t.Errorf("erro inesperado ao autenticar já autenticado: %v", err)
	}
}

// TestAWSAuth_Authenticate_InvalidConfig verifica falha ao autenticar com configuração inválida.
func TestAWSAuth_Authenticate_InvalidConfig(t *testing.T) {
	auth := &AWSAuth{
		AccessKeyID:     nil,
		SecretAccessKey: nil,
		Region:          "",
	}

	err := auth.Authenticate()
	if err == nil {
		t.Fatalf("esperado erro, mas retornado nil com configuração inválida")
	}

	if err.Error() != "validation failed: missing required AWS authentication fields: [AccessKeyID SecretAccessKey Region]" {
		t.Errorf("mensagem inesperada de erro: '%v'", err)
	}
}
