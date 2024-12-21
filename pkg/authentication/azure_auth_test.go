package authentication

import (
	"testing"
)

// TestNewAzureAuthFromAuth_Valid verifica se a inicialização de AzureAuth com entradas válidas ocorre corretamente.
func TestNewAzureAuthFromAuth_Valid(t *testing.T) {
	fields := map[string]string{
		"azure_client_id":       "test-client-id",
		"azure_client_secret":   "test-client-secret",
		"azure_tenant_id":       "test-tenant-id",
		"azure_subscription_id": "test-subscription-id",
	}

	auth, err := NewAzureAuthFromAuth(fields)
	if err != nil {
		t.Fatalf("erro inesperado ao criar AzureAuth: %v", err)
	}

	if auth.ClientID != "test-client-id" {
		t.Errorf("esperado ClientID 'test-client-id', recebido '%s'", auth.ClientID)
	}

	if auth.ClientSecret != "test-client-secret" {
		t.Errorf("esperado ClientSecret 'test-client-secret', recebido '%s'", auth.ClientSecret)
	}

	if auth.TenantID != "test-tenant-id" {
		t.Errorf("esperado TenantID 'test-tenant-id', recebido '%s'", auth.TenantID)
	}

	if auth.SubscriptionID != "test-subscription-id" {
		t.Errorf("esperado SubscriptionID 'test-subscription-id', recebido '%s'", auth.SubscriptionID)
	}
}

// TestNewAzureAuthFromAuth_Invalid verifica se a inicialização de AzureAuth com entradas inválidas retorna erros.
func TestNewAzureAuthFromAuth_Invalid(t *testing.T) {
	fields := map[string]string{
		"azure_client_id":       "test-client-id",
		"azure_client_secret":   "",
		"azure_tenant_id":       "test-tenant-id",
		"azure_subscription_id": "",
	}

	_, err := NewAzureAuthFromAuth(fields)
	if err == nil {
		t.Fatalf("esperado erro para campos obrigatórios ausentes, mas nenhum erro foi retornado")
	}

	expectedErr := "missing required Azure authentication fields"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestAzureAuth_Validate_ValidFields verifica se a validação ocorre corretamente com campos válidos.
func TestAzureAuth_Validate_ValidFields(t *testing.T) {
	auth := &AzureAuth{
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		TenantID:       "test-tenant-id",
		SubscriptionID: "test-subscription-id",
	}

	if err := auth.Validate(); err != nil {
		t.Errorf("erro inesperado na validação com campos válidos: %v", err)
	}
}

// TestAzureAuth_Validate_MissingFields verifica se a validação falha ao faltar campos obrigatórios.
func TestAzureAuth_Validate_MissingFields(t *testing.T) {
	auth := &AzureAuth{
		ClientID:       "",
		ClientSecret:   "",
		TenantID:       "test-tenant-id",
		SubscriptionID: "",
	}

	err := auth.Validate()
	if err == nil {
		t.Fatalf("esperado erro para campos obrigatórios ausentes, mas foi retornado nil")
	}

	expectedErr := "missing required Azure authentication fields"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestAzureAuth_Authenticate_AlreadyAuthenticated testa autenticação quando já autenticado.
func TestAzureAuth_Authenticate_AlreadyAuthenticated(t *testing.T) {
	auth := &AzureAuth{
		Authenticated: true,
	}

	err := auth.Authenticate()
	if err != nil {
		t.Errorf("erro inesperado para autenticação já realizada: %v", err)
	}
}

// TestAzureAuth_Authenticate_InvalidFields verifica a autenticação com uma configuração inválida.
func TestAzureAuth_Authenticate_InvalidFields(t *testing.T) {
	auth := &AzureAuth{
		ClientID:       "",
		ClientSecret:   "",
		TenantID:       "",
		SubscriptionID: "",
	}

	err := auth.Authenticate()
	if err == nil {
		t.Fatalf("esperado erro para configuração inválida, mas foi retornado nil")
	}

	if err.Error() != "validation failed: missing required Azure authentication fields" {
		t.Errorf("mensagem de erro inesperada: '%v'", err)
	}
}

// Exemplo de teste de autenticação mockada para evitar dependência real de Azure SDK.
func TestAzureAuth_Authenticate_Simulated(t *testing.T) {
	auth := &AzureAuth{
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		TenantID:       "test-tenant-id",
		SubscriptionID: "test-subscription-id",
	}

	err := auth.Authenticate()
	if err != nil {
		t.Errorf("erro inesperado ao autenticar com configuração simulada: %v", err)
	}
}
