package authentication

import (
	"testing"
)

// TestNewGCPAuthFromAuth_Valid verifica se a inicialização de GCPAuth com entradas válidas ocorre corretamente.
func TestNewGCPAuthFromAuth_Valid(t *testing.T) {
	fields := map[string]string{
		"gcp_project_id": "test-project-id",
		"gcp_auth_json":  `{"type": "service_account", "project_id": "test-project-id"}`,
	}

	auth, err := NewGCPAuthFromAuth(fields)
	if err != nil {
		t.Fatalf("erro inesperado ao criar GCPAuth: %v", err)
	}

	if auth.ProjectID != "test-project-id" {
		t.Errorf("esperado ProjectID 'test-project-id', mas foi recebido '%s'", auth.ProjectID)
	}

	if auth.AuthJSON != `{"type": "service_account", "project_id": "test-project-id"}` {
		t.Errorf("esperado AuthJSON correspondente, mas foi recebido '%s'", auth.AuthJSON)
	}
}

// TestNewGCPAuthFromAuth_Invalid verifica se a inicialização de GCPAuth com entradas inválidas retorna um erro.
func TestNewGCPAuthFromAuth_Invalid(t *testing.T) {
	fields := map[string]string{
		"gcp_project_id": "", // Campo vazio para simular entrada inválida.
		"gcp_auth_json":  "",
	}

	_, err := NewGCPAuthFromAuth(fields)
	if err == nil {
		t.Fatalf("esperado erro para campos obrigatórios ausentes, mas nenhum erro foi retornado")
	}

	expectedErr := "missing required GCP authentication fields"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestGCPAuth_Validate_ValidFields verifica se a validação ocorre corretamente com campos válidos.
func TestGCPAuth_Validate_ValidFields(t *testing.T) {
	auth := &GCPAuth{
		ProjectID: "test-project-id",
		AuthJSON:  `{"type": "service_account", "project_id": "test-project-id"}`,
	}

	if err := auth.Validate(); err != nil {
		t.Errorf("erro inesperado na validação com campos válidos: %v", err)
	}
}

// TestGCPAuth_Validate_MissingFields verifica se a validação falha quando campos obrigatórios estão ausentes.
func TestGCPAuth_Validate_MissingFields(t *testing.T) {
	auth := &GCPAuth{
		ProjectID: "",
		AuthJSON:  "",
	}

	err := auth.Validate()
	if err == nil {
		t.Fatalf("esperado erro para campos obrigatórios ausentes, mas foi retornado nil")
	}

	expectedErr := "missing required GCP authentication fields"
	if err.Error() != expectedErr {
		t.Errorf("mensagem de erro incorreta: esperado '%s', recebido '%s'", expectedErr, err.Error())
	}
}

// TestGCPAuth_Authenticate_AlreadyAuthenticated verifica o comportamento quando já autenticado.
func TestGCPAuth_Authenticate_AlreadyAuthenticated(t *testing.T) {
	auth := &GCPAuth{
		Authenticated: true, // Simula a autenticação já realizada
	}

	err := auth.Authenticate()
	if err != nil {
		t.Errorf("erro inesperado para autenticação já realizada: %v", err)
	}
}

// TestGCPAuth_Authenticate_InvalidConfig verifica a autenticação com uma configuração inválida.
func TestGCPAuth_Authenticate_InvalidConfig(t *testing.T) {
	auth := &GCPAuth{
		ProjectID: "", // Configuração inválida.
		AuthJSON:  "",
	}

	err := auth.Authenticate()
	if err == nil {
		t.Fatalf("esperado erro para configuração inválida, mas foi retornado nil")
	}

	if err.Error() != "validation failed: missing required GCP authentication fields" {
		t.Errorf("mensagem de erro inesperada: '%v'", err)
	}
}
