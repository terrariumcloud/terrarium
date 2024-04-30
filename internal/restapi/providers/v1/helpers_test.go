package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Test_GetProviderNameFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/test-org/test-provider/versions", nil)
	req = mux.SetURLVars(req, map[string]string{
		"organization_name": "test-org",
		"name":              "test-provider",
	})

	providerName := GetProviderNameFromRequest(req)

	expectedProviderName := "test-org/test-provider"
	if providerName != expectedProviderName {
		t.Errorf("Expected provider name to be %s, but got %s", expectedProviderName, providerName)
	}
}

func Test_GetProviderInputsFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/test-org/test-provider/2.0.0/download/linux/amd64", nil)
	req = mux.SetURLVars(req, map[string]string{
		"version": "2.0.0",
		"os":      "linux",
		"arch":    "amd64",
	})

	version, os, arch := GetProviderInputsFromRequest(req)

	expectedVersion := "2.0.0"
	if version != expectedVersion {
		t.Errorf("Expected version to be %s, but got %s", expectedVersion, version)
	}

	expectedOS := "linux"
	if os != expectedOS {
		t.Errorf("Expected OS to be %s, but got %s", expectedOS, os)
	}

	expectedArch := "amd64"
	if arch != expectedArch {
		t.Errorf("Expected arch. to be %s, but got %s", expectedArch, arch)
	}
}

func Test_GetProviderLocationFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/test-org/test-provider/2.0.0/linux/amd64/terraform-provider-test-provider_2.0.0_linux_amd64.zip", nil)
	req = mux.SetURLVars(req, map[string]string{
		"organization_name": "test-org",
		"name":              "test-provider",
		"version":           "2.0.0",
		"os":                "linux",
		"arch":              "amd64",
	})

	provider := GetProviderLocationFromRequest(req)

	expectedProviderName := "test-org/test-provider"
	if provider.Name != expectedProviderName {
		t.Errorf("Expected provider name to be %s, but got %s", expectedProviderName, provider.Name)
	}

	expectedVersion := "2.0.0"
	if provider.Version != expectedVersion {
		t.Errorf("Expected version to be %s, but got %s", expectedVersion, provider.Version)
	}

	expectedOS := "linux"
	if provider.Os != expectedOS {
		t.Errorf("Expected OS to be %s, but got %s", expectedOS, provider.Os)
	}

	expectedArch := "amd64"
	if provider.Arch != expectedArch {
		t.Errorf("Expected arch. to be %s, but got %s", expectedArch, provider.Arch)
	}
}

func Test_GetVersionedProviderFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/test-org/test-provider/2.0.0/terraform-provider-test-provider_2.0.0_SHA256SUMS", nil)
	req = mux.SetURLVars(req, map[string]string{
		"organization_name": "test-org",
		"name":              "test-provider",
		"version":           "2.0.0",
	})

	provider := GetVersionedProviderFromRequest(req)

	expectedProviderName := "test-org/test-provider"
	if provider.Name != expectedProviderName {
		t.Errorf("Expected provider name to be %s, but got %s", expectedProviderName, provider.Name)
	}

	expectedVersion := "2.0.0"
	if provider.Version != expectedVersion {
		t.Errorf("Expected version to be %s, but got %s", expectedVersion, provider.Version)
	}
}
