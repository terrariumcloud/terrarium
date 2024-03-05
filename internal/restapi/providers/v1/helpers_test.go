package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Test_GetProviderNameFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/hashicorp/random/versions", nil)
	req = mux.SetURLVars(req, map[string]string{
		"organization_name": "hashicorp",
		"name":              "random",
	})

	providerName := GetProviderNameFromRequest(req)

	expectedProviderName := "hashicorp/random"
	if providerName != expectedProviderName {
		t.Errorf("Expected provider name to be %s, but got %s", expectedProviderName, providerName)
	}
}

func Test_GetProviderInputsFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/providers/v1/hashicorp/random/2.0.0/download/linux/amd64", nil)
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
